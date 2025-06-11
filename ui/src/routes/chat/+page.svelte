<script lang="ts">
	import { page } from '$app/stores';
	import { error as errorStore } from '$lib/stores';
	import { llmService } from '$lib/services/llm';
	import type { Usage } from '$lib/services/llm';
	import type { ChatMessage } from '$lib/services/chat';
	import SelectField from '../../components/common/SelectField.svelte';
	import ChatInput from '../../components/chat/ChatInput.svelte';
	import MarkdownRenderer from '../../components/common/MarkdownRenderer.svelte';
	import { onMount } from 'svelte';
	import { 
		activeChatStore, 
		loadChat, 
		createNewChat,
		deleteChat,
		archiveChat,
		updateChatTitle
	} from '$lib/stores/chatStore';
	import { goto } from '$app/navigation';
	import { chatService } from '$lib/services/chat';
	
	type Message = {
		role: 'user' | 'assistant' | 'system';
		content: string;
		timestamp: Date;
		usage?: Usage;
	};
	
	let message = '';
	let messages: Message[] = [];
	let isLoading = false;
	let systemPrompt = "You are a helpful learning assistant for kids.";
	let selectedModel = ""; // Empty means use default model
	let totalTokens = 0;
	let totalPromptTokens = 0;
	let totalCompletionTokens = 0;
	let totalCost = 0;
	let isLoadingModels = false;
	let modelError: string | null = null;
	let showSettings = false;
	let hasInitialized = false;
	
	// Current chat ID from the URL
	$: chatId = $page.params.id || null;
	
	// Get initial prompt from URL if available
	$: initialPrompt = $page.data.initialPrompt;
	
	// Set the message input value when initialPrompt is available
	$: if (initialPrompt && !chatId && !message) {
		message = initialPrompt;
	}
	
	// Settings
	let autoScrollEnabled = true;
	let sendWithEnter = true;
	
	// Available models - will be populated from backend
	let availableModels: { id: string; name: string; isDefault?: boolean }[] = [];
	
	// Chat title editing state
	let isEditingTitle = false;
	let editedTitle = '';
	let titleInputEl: HTMLInputElement;
	
	// Load chat when ID changes
	$: {
		if (chatId) {
			loadChatData(chatId);
		} else {
			// If no chat ID, show empty state
			activeChatStore.set({ chat: null, loading: false, error: null });
			messages = [];
		}
	}
	
	// Convert chat messages from backend to UI format
	function convertChatMessagesToUIFormat(chatMessages: ChatMessage[]): Message[] {
		return chatMessages.map(msg => ({
			role: msg.role,
			content: msg.content,
			timestamp: msg.createdAt
		}));
	}
	
	// Function to scroll to bottom of messages
	function scrollToBottom() {
		if (autoScrollEnabled) {
			setTimeout(() => {
				window.scrollTo({
					top: document.body.scrollHeight,
					behavior: 'smooth'
				});
			}, 50);
		}
	}
	
	// When active chat changes, update local messages and system prompt
	$: if ($activeChatStore.chat) {
		systemPrompt = $activeChatStore.chat.systemPrompt;
		// Always reset messages when switching chats
		if ($activeChatStore.chat.messages) {
			messages = convertChatMessagesToUIFormat($activeChatStore.chat.messages);
		} else {
			messages = [];
		}
		
		// Scroll to bottom when messages are loaded
		if ($activeChatStore.chat.messages && $activeChatStore.chat.messages.length > 0) {
			scrollToBottom();
		}
	}
	
	async function loadChatData(id: string) {
		try {
			await loadChat(id);
		} catch (err) {
			console.error('Error loading chat:', err);
			errorStore.set(err instanceof Error ? err.message : 'Failed to load chat');
		}
	}
	
	async function handleCreateNewChat() {
		try {
			const newChatId = await createNewChat(systemPrompt, selectedModel || undefined);
			goto(`/chat/${newChatId}`);
		} catch (err) {
			console.error('Error creating new chat:', err);
			errorStore.set(err instanceof Error ? err.message : 'Failed to create new chat');
		}
	}
	
	onMount(() => {
		document.addEventListener('keydown', handleKeyboardShortcuts);
		fetchModels();
		hasInitialized = true;
		
		return () => {
			document.removeEventListener('keydown', handleKeyboardShortcuts);
		};
	});
	
	// Fetch available models from backend
	async function fetchModels() {
		isLoadingModels = true;
		modelError = null;
		
		try {
			const data = await llmService.getInfo();
			
			// Transform the models data into the format we need
			availableModels = data.platforms.flatMap(platform => 
				platform.models.map(model => ({
					id: model.name,
					name: `${model.name}${model.isDefault ? ' (Default)' : ''}`,
					isDefault: model.isDefault
				}))
			);
			
			// If no models were found, add a default option
			if (availableModels.length === 0) {
				availableModels = [{ id: "", name: "Default" }];
			}
			
			// Select the default model if available
			const defaultModel = availableModels.find(m => m.isDefault);
			if (defaultModel) {
				selectedModel = defaultModel.id;
			}
		} catch (err) {
			console.error('Error fetching models:', err);
			modelError = err instanceof Error ? err.message : 'Failed to fetch models';
			// Fallback to default model
			availableModels = [{ id: "", name: "Default" }];
		} finally {
			isLoadingModels = false;
		}
	}
	
	// Calculate the total tokens and cost used in the conversation
	$: {
		let prompt = 0;
		let completion = 0;
		let total = 0;
		let cost = 0;
		
		messages.forEach(msg => {
			if (msg.usage) {
				prompt += msg.usage.PromptTokens || 0;
				completion += msg.usage.CompletionTokens || 0;
				total += msg.usage.TotalTokens || 0;
				cost += msg.usage.Cost || 0;
			}
		});
		
		totalPromptTokens = prompt;
		totalCompletionTokens = completion;
		totalTokens = total;
		totalCost = cost;
	}
	
	// Create tooltip content for total tokens
	$: totalTokensTooltip = `Prompt: ${totalPromptTokens} tokens\nCompletion: ${totalCompletionTokens} tokens\nCost: $${totalCost.toFixed(6)}`;
	
	function toggleSettings() {
		showSettings = !showSettings;
	}
	
	function clearConversation() {
		if (confirm('Are you sure you want to clear the entire conversation?')) {
			messages = [];
			totalTokens = 0;
			totalPromptTokens = 0;
			totalCompletionTokens = 0;
			totalCost = 0;
		}
	}
	
	// Handle message submission from the ChatInput component
	async function handleMessageSubmit(event: CustomEvent<{ message: string }>) {
		const userMessageContent = event.detail.message;
		
		if (!chatId) {
			// For a new chat, create it first
			try {
				isLoading = true; // Set loading state
				const newChatId = await createNewChat(systemPrompt, selectedModel || undefined);
				
				// Add user message to local state
				const userMessage: Message = {
					role: 'user',
					content: userMessageContent,
					timestamp: new Date()
				};
				messages = [...messages, userMessage];
				
				// First save the user message to the backend
				await chatService.addMessageToChat(newChatId, userMessageContent, 'user');
				
				// Send message to LLM and get response
				const chatHistory = [userMessageContent]; // Only the first message for new chat
				const data = await llmService.chat(chatHistory.join('\n\n'), systemPrompt, selectedModel || undefined);
				
				// Add assistant response to local state
				const assistantMessage: Message = {
					role: 'assistant',
					content: data.response,
					timestamp: new Date(),
					usage: data.usage
				};
				messages = [...messages, assistantMessage];
				
				// Save the assistant's response
				await chatService.addMessageToChat(newChatId, data.response, 'assistant');
				
				// Generate and update chat title based on the first message
				const truncatedMessage = userMessageContent.length > 50 
					? userMessageContent.substring(0, 47) + '...' 
					: userMessageContent;
				await updateChatTitle(newChatId, truncatedMessage);
				
				// Then navigate to the new chat - both messages will be loaded
				await goto(`/chat/${newChatId}`);
			} catch (err) {
				console.error('Error creating new chat:', err);
				errorStore.set(err instanceof Error ? err.message : 'Failed to create new chat');
				isLoading = false; // Reset loading state on error
			}
		} else {
			// For an existing chat, just send the message
			message = userMessageContent;
			await sendMessage();
		}
	}
	
	// Handle keydown events from ChatInput
	function handleInputKeyDown(_event: CustomEvent<KeyboardEvent>) {
		// This function intentionally left empty for now
		// Previously contained shortcuts that were interfering with typing
	}
	
	// Update sendMessage to not rely on messageInputRef
	async function sendMessage() {
		if (!message.trim() || !chatId) return;
		
		const userMessageContent = message.trim();
		
		// Add user message to local state
		const userMessage: Message = {
			role: 'user',
			content: userMessageContent,
			timestamp: new Date()
		};
		messages = [...messages, userMessage];
		
		// Clear input
		message = '';
		isLoading = true;
		errorStore.set(null);
		
		// Scroll to bottom after adding user message
		scrollToBottom();
		
		try {
			// First, save the user message to the database
			await chatService.addMessageToChat(chatId, userMessageContent, 'user');
			
			// Prepare chat history by concatenating all previous messages
			const chatHistory = messages.map(msg => `${msg.role === 'user' ? 'User' : 'Assistant'}: ${msg.content}`).join('\n\n');
			
			// Send chat request using the LLM service
			const data = await llmService.chat(chatHistory, systemPrompt, selectedModel || undefined);
			
			// Add assistant response to local state
			const assistantMessage: Message = {
				role: 'assistant',
				content: data.response,
				timestamp: new Date(),
				usage: data.usage
			};
			messages = [...messages, assistantMessage];
			
			// Save the assistant response to the database
			await chatService.addMessageToChat(chatId, data.response, 'assistant');
			
			// Scroll to bottom after receiving response
			scrollToBottom();
			
		} catch (err) {
			console.error('Error sending message:', err);
			errorStore.set(err instanceof Error ? err.message : 'Failed to send message');
		} finally {
			isLoading = false;
		}
	}
	
	function handleKeyboardShortcuts(event: KeyboardEvent) {
		// Only handle shortcuts when not in input/textarea
		if (event.target instanceof HTMLInputElement || 
			event.target instanceof HTMLTextAreaElement) {
			return;
		}
		
		// n to create new chat
		if (event.key === 'n' && !event.ctrlKey && !event.metaKey) {
			event.preventDefault();
			handleCreateNewChat();
		}
	}
	
	async function handleDeleteChat() {
		if (!chatId) return;
		
		if (confirm('Are you sure you want to delete this chat? This action cannot be undone.')) {
			try {
				await deleteChat(chatId);
				// Navigate to main chat page
				goto('/chat');
			} catch (error) {
				console.error('Failed to delete chat:', error);
				errorStore.set(error instanceof Error ? error.message : 'Failed to delete chat');
			}
		}
	}
	
	async function handleArchiveChat() {
		if (!chatId) return;
		
		try {
			await archiveChat(chatId, true);
			// Navigate to main chat page after archiving
			goto('/chat');
		} catch (error) {
			console.error('Failed to archive chat:', error);
			errorStore.set(error instanceof Error ? error.message : 'Failed to archive chat');
		}
	}

	// Function to start editing the chat title
	function startEditingTitle() {
		if (!$activeChatStore.chat) return;
		editedTitle = $activeChatStore.chat.title || 'New Chat';
		isEditingTitle = true;
		
		// Focus the input after DOM update
		setTimeout(() => {
			if (titleInputEl) {
				titleInputEl.focus();
				titleInputEl.select();
			}
		}, 10);
	}

	// Function to save edited title
	async function saveChatTitle() {
		if (!chatId || !$activeChatStore.chat) return;
		
		if (editedTitle.trim() === '') {
			editedTitle = 'New Chat';
		}
		
		try {
			if (editedTitle !== $activeChatStore.chat.title) {
				await updateChatTitle(chatId, editedTitle);
			}
		} catch (error) {
			console.error('Failed to update chat title:', error);
			errorStore.set(error instanceof Error ? error.message : 'Failed to update chat title');
		} finally {
			isEditingTitle = false;
		}
	}

	// Function to cancel editing
	function cancelEditingTitle() {
		isEditingTitle = false;
		if ($activeChatStore.chat) {
			editedTitle = $activeChatStore.chat.title || 'New Chat';
		}
	}

	// Handle key events in the title input
	function handleTitleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			event.preventDefault();
			saveChatTitle();
		} else if (event.key === 'Escape') {
			event.preventDefault();
			cancelEditingTitle();
		}
	}

	// Handle double click on chat title
	function handleTitleDoubleClick() {
		if (chatId) {
			startEditingTitle();
		}
	}
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6 flex-1">
	<div class="flex justify-between items-center mb-6">
		<div>
			{#if isEditingTitle && $activeChatStore.chat}
				<!-- Edit mode for title -->
				<div class="flex items-center">
					<input 
						bind:this={titleInputEl}
						bind:value={editedTitle}
						on:keydown={handleTitleKeyDown}
						on:blur={saveChatTitle}
						class="text-2xl font-semibold text-gray-900 dark:text-white bg-transparent border-b border-indigo-300 dark:border-indigo-600 focus:outline-none focus:border-indigo-500 px-1 py-0.5 mr-2"
					/>
				</div>
			{:else}
				<!-- Display mode -->
				<div class="group flex items-center">
					<h1 
						class="text-2xl font-semibold text-gray-900 dark:text-white cursor-pointer"
						on:dblclick={handleTitleDoubleClick}
						title={chatId ? "Double-click to edit title" : undefined}
					>
						{$activeChatStore.chat ? $activeChatStore.chat.title : 'Chat'}
					</h1>
					{#if chatId && $activeChatStore.chat}
						<button 
							on:click={startEditingTitle}
							class="ml-2 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 opacity-0 group-hover:opacity-100 transition-opacity"
							title="Edit title"
							aria-label="Edit title"
						>
							<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
								<path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
							</svg>
						</button>
					{/if}
				</div>
				<p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
					Using {selectedModel ? availableModels.find(m => m.id === selectedModel)?.name : 'Default Model'}
				</p>
			{/if}
		</div>
		
		<div class="flex items-center space-x-4">
			<button 
				type="button" 
				on:click={toggleSettings}
				class="text-sm text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white flex items-center"
				title="Chat Settings"
			>
				<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z" clip-rule="evenodd" />
				</svg>
				Settings
			</button>
		</div>
	</div>
	
	{#if showSettings}
		<div class="mb-6 bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
			<div class="p-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
				<h2 class="text-md font-medium text-gray-900 dark:text-white">Chat Settings</h2>
				<button 
					type="button"
					on:click={toggleSettings}
					class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300"
					aria-label="Close settings"
				>
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
						<path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
					</svg>
				</button>
			</div>
			<div class="p-4 space-y-4">
				<div>
					<h3 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Model Selection</h3>
					<p class="text-xs text-gray-500 dark:text-gray-400 mb-2">Choose the AI model to use for the conversation.</p>
					<SelectField
						id="modelSelect"
						label="Model"
						bind:value={selectedModel}
						disabled={isLoadingModels}
						cols=""
						inline={true}
					>
						{#if isLoadingModels}
							<option value="">Loading models...</option>
						{:else if modelError}
							<option value="">Error loading models</option>
						{:else}
							{#each availableModels as model}
								<option value={model.id}>{model.name}</option>
							{/each}
						{/if}
					</SelectField>
				</div>

				<div>
					<h3 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">System Prompt</h3>
					<p class="text-xs text-gray-500 dark:text-gray-400 mb-2">Define the assistant's behavior and personality.</p>
					<textarea
						bind:value={systemPrompt}
						rows="4"
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-100 dark:placeholder-gray-400 text-sm"
						placeholder="Define how the assistant should behave..."
					></textarea>
				</div>

				<div class="pt-4 border-t border-gray-200 dark:border-gray-700">
					<h3 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-4">Interface Settings</h3>
					
					<div class="space-y-4">
						<div class="flex items-center justify-between">
							<div>
								<h4 class="text-sm font-medium text-gray-700 dark:text-gray-300">Auto-scroll Messages</h4>
								<p class="text-xs text-gray-500 dark:text-gray-400">Automatically scroll to new messages</p>
							</div>
							<button 
								type="button"
								on:click={() => autoScrollEnabled = !autoScrollEnabled}
								class="ml-4 relative inline-flex flex-shrink-0 h-6 w-11 border-2 border-transparent rounded-full cursor-pointer transition-colors ease-in-out duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-offset-gray-800 {autoScrollEnabled ? 'bg-indigo-600' : 'bg-gray-200 dark:bg-gray-700'}"
							>
								<span class="sr-only">Auto-scroll messages</span>
								<span 
									class="pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow transform ring-0 transition ease-in-out duration-200 {autoScrollEnabled ? 'translate-x-5' : 'translate-x-0'}"
								></span>
							</button>
						</div>
						
						<div class="flex items-center justify-between">
							<div>
								<h4 class="text-sm font-medium text-gray-700 dark:text-gray-300">Send with Enter</h4>
							</div>
							<button 
								type="button"
								on:click={() => sendWithEnter = !sendWithEnter}
								class="ml-4 relative inline-flex flex-shrink-0 h-6 w-11 border-2 border-transparent rounded-full cursor-pointer transition-colors ease-in-out duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-offset-gray-800 {sendWithEnter ? 'bg-indigo-600' : 'bg-gray-200 dark:bg-gray-700'}"
							>
								<span class="sr-only">Send with Enter</span>
								<span 
									class="pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow transform ring-0 transition ease-in-out duration-200 {sendWithEnter ? 'translate-x-5' : 'translate-x-0'}"
								></span>
							</button>
						</div>
					</div>
				</div>

				<div class="pt-4 border-t border-gray-200 dark:border-gray-700">
					<button 
						type="button" 
						on:click={clearConversation}
						class="text-sm text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-300 font-medium"
					>
						Clear conversation
					</button>
				</div>

				<div class="pt-4 border-t border-gray-200 dark:border-gray-700">
					<h3 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-4">Chat Management</h3>
					
					<div class="space-y-3">
						<button 
							type="button" 
							on:click={handleArchiveChat}
							class="flex items-center text-sm text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200 font-medium"
							disabled={!chatId}
						>
							<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" viewBox="0 0 20 20" fill="currentColor">
								<path d="M4 3a2 2 0 100 4h12a2 2 0 100-4H4z" />
								<path fill-rule="evenodd" d="M3 8h14v7a2 2 0 01-2 2H5a2 2 0 01-2-2V8zm5 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" clip-rule="evenodd" />
							</svg>
							Archive chat
						</button>
						
						<button 
							type="button" 
							on:click={handleDeleteChat}
							class="flex items-center text-sm text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-300 font-medium"
							disabled={!chatId}
						>
							<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
							</svg>
							Delete chat
						</button>
					</div>
				</div>
			</div>
		</div>
	{/if}
	
	{#if !hasInitialized}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-8 w-8 border-t-2 border-indigo-500"></div>
		</div>
	{:else if !chatId}
		<div class="w-full flex flex-col items-center px-4 py-12 text-center space-y-6">
			<div class="w-full">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-indigo-400 dark:text-indigo-500 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
				</svg>
				<!-- New chat input using the ChatInput component -->
				<div class="w-full mx-auto">
					<div class="relative">
						<ChatInput 
							bind:value={message}
							placeholder="Type here to start a new chat..."
							disabled={isLoading}
							buttonText="Send"
							isLoading={isLoading}
							on:submit={handleMessageSubmit}
							on:keydown={handleInputKeyDown}
							showKeyboardShortcuts={true}
							totalTokensInfo={null}
						/>
					</div>
				</div>
			</div>
		</div>
	{:else if $activeChatStore.loading}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-8 w-8 border-t-2 border-indigo-500"></div>
		</div>
	{:else if $activeChatStore.error}
		<div class="p-6 text-center text-red-500 dark:text-red-400">
			<p>{$activeChatStore.error}</p>
			<button
				on:click={() => chatId && loadChatData(chatId)}
				class="mt-2 text-indigo-600 dark:text-indigo-400 underline text-sm"
			>
				Retry
			</button>
		</div>
	{:else}
		<div class="p-4">
			{#if messages.length === 0}
				<div class="text-center py-12 border border-dashed border-gray-200 dark:border-gray-700 rounded-lg">
					<p class="text-gray-500 dark:text-gray-400">
						No messages yet. Start a conversation below.
					</p>
				</div>
			{:else}
				<div class="space-y-6">
					{#each messages as msg}
						<div class="flex flex-col">
							<div class="flex items-center mb-2">
								<span class="font-medium text-sm text-gray-700 dark:text-gray-300 mr-2">
									{msg.role === 'user' ? 'You' : msg.role === 'assistant' ? 'Assistant' : 'System'}
								</span>
								<span class="text-xs text-gray-500 dark:text-gray-400">
									{msg.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
								</span>
								{#if msg.role === 'assistant' && msg.usage}
									<span class="text-xs text-gray-400 dark:text-gray-500 ml-2 tooltip" 
										title="Model: {msg.usage.LlmModelName || 'Default'}&#10;Cache hit: {msg.usage.CacheHit ? 'Yes' : 'No'}&#10;Prompt: {msg.usage.PromptTokens} tokens&#10;Completion: {msg.usage.CompletionTokens} tokens&#10;Cost: ${(msg.usage.Cost || 0).toFixed(6)}">
										({msg.usage.TotalTokens} tokens)
									</span>
								{/if}
							</div>
							<div class="{msg.role === 'user' ? 'bg-blue-50 dark:bg-blue-900/30 border-blue-100 dark:border-blue-800' : msg.role === 'system' ? 'bg-purple-50 dark:bg-purple-900/30 border-purple-100 dark:border-purple-800' : 'bg-gray-50 dark:bg-gray-900 border-gray-100 dark:border-gray-700'} border rounded-lg p-4">
								{#if msg.role === 'assistant'}
									<MarkdownRenderer content={msg.content} />
								{:else}
									<p class="whitespace-pre-wrap break-words text-left text-gray-800 dark:text-gray-200">{msg.content}</p>
								{/if}
							</div>
						</div>
					{/each}
					{#if isLoading}
						<div class="flex flex-col">
							<div class="flex items-center mb-2">
								<span class="font-medium text-sm text-gray-700 dark:text-gray-300">Assistant</span>
							</div>
							<div class="bg-gray-50 dark:bg-gray-700 border border-gray-100 dark:border-gray-600 rounded-lg p-4">
								<div class="flex space-x-2">
									<div class="w-2 h-2 bg-gray-400 dark:bg-gray-500 rounded-full animate-bounce bounce-delay-1"></div>
									<div class="w-2 h-2 bg-gray-400 dark:bg-gray-500 rounded-full animate-bounce bounce-delay-2"></div>
								</div>
							</div>
						</div>
					{/if}
				</div>
			{/if}
		</div>
		
		<!-- Input area (only show when a chat is selected) -->
		{#if chatId}
			<div class="border-t border-gray-200 dark:border-gray-700 p-4">
				<ChatInput 
					bind:value={message}
					placeholder="Type your message... (Press / to focus)"
					disabled={isLoading || $activeChatStore.loading}
					isLoading={isLoading}
					sendWithEnter={sendWithEnter}
					on:submit={handleMessageSubmit}
					on:keydown={handleInputKeyDown}
					showKeyboardShortcuts={true}
					totalTokensInfo={chatId ? { count: totalTokens, tooltip: totalTokensTooltip } : null}
				/>
			</div>
		{/if}
	{/if}
</div>

<style>
	.tooltip {
		position: relative;
		cursor: help;
	}
	
	.tooltip:hover::after {
		content: attr(title);
		position: absolute;
		bottom: 100%;
		left: 50%;
		transform: translateX(-50%);
		padding: 0.5rem;
		background-color: rgba(0, 0, 0, 0.8);
		color: white;
		border-radius: 0.25rem;
		white-space: pre;
		z-index: 10;
		min-width: 200px;
		text-align: left;
		font-weight: normal;
	}
	
	/* Animation keyframes for the bouncing dots */
	@keyframes bounce {
		0%, 100% {
			transform: translateY(0);
		}
		50% {
			transform: translateY(-4px);
		}
	}
	
	.animate-bounce {
		animation: bounce 1s infinite;
	}
	
	.bounce-delay-1 {
		animation-delay: 0.1s;
	}
	
	.bounce-delay-2 {
		animation-delay: 0.2s;
	}
</style> 