<script lang="ts">
	import { error as errorStore } from '$lib/stores';
	import { llmService } from '$lib/services/llm';
	import type { Usage } from '$lib/services/llm';
	import SelectField from '../../components/common/SelectField.svelte';
	import { onMount } from 'svelte';
	
	type Message = {
		role: 'user' | 'assistant';
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
	
	// Settings
	let autoScrollEnabled = true;
	let sendWithEnter = true;
	
	// Available models - will be populated from backend
	let availableModels: { id: string; name: string; isDefault?: boolean }[] = [];
	
	onMount(async () => {
		await fetchModels();
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
	
	async function sendMessage() {
		if (!message.trim()) return;
		
		// Add user message to chat
		const userMessage: Message = {
			role: 'user',
			content: message,
			timestamp: new Date()
		};
		messages = [...messages, userMessage];
		
		// Clear input
		message = '';
		isLoading = true;
		errorStore.set(null);
		
		try {
			// Prepare chat history by concatenating all previous messages
			const chatHistory = messages.map(msg => `${msg.role === 'user' ? 'User' : 'Assistant'}: ${msg.content}`).join('\n\n');
			
			// Send chat request using the LLM service
			const data = await llmService.chat(chatHistory, systemPrompt, selectedModel || undefined);
			
			// Add assistant response to chat
			const assistantMessage: Message = {
				role: 'assistant',
				content: data.response,
				timestamp: new Date(),
				usage: data.usage
			};
			messages = [...messages, assistantMessage];
			
			// Scroll to bottom if auto-scroll is enabled
			if (autoScrollEnabled) {
				setTimeout(() => {
					const chatContainer = document.querySelector('.chat-container');
					if (chatContainer) {
						chatContainer.scrollTop = chatContainer.scrollHeight;
					}
				}, 50);
			}
			
		} catch (err) {
			console.error('Error sending message:', err);
			errorStore.set(err instanceof Error ? err.message : 'Failed to send message');
		} finally {
			isLoading = false;
		}
	}
	
	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' && !event.shiftKey && sendWithEnter) {
			event.preventDefault();
			sendMessage();
		}
	}
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
	<div class="flex justify-between items-center mb-6">
		<div>
			<h1 class="text-2xl font-semibold text-gray-900 dark:text-white">Chat with LLM</h1>
			<p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
				Using {selectedModel ? availableModels.find(m => m.id === selectedModel)?.name : 'Default Model'}
			</p>
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
								<p class="text-xs text-gray-500 dark:text-gray-400">Press Enter to send message (Shift+Enter for new line)</p>
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
			</div>
		</div>
	{/if}
	
	<div class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
		<!-- Chat messages container without fixed height -->
		<div class="px-4 py-5 sm:p-6 chat-container max-h-[60vh] overflow-y-auto">
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
									{msg.role === 'user' ? 'You' : 'Assistant'}
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
							<div class="{msg.role === 'user' ? 'bg-blue-50 dark:bg-blue-900/30 border-blue-100 dark:border-blue-800' : 'bg-gray-50 dark:bg-gray-700 border-gray-100 dark:border-gray-600'} border rounded-lg p-4">
								<p class="whitespace-pre-wrap break-words text-left text-gray-800 dark:text-gray-200">{msg.content}</p>
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
		
		<!-- Input area -->
		<div class="px-4 py-4 sm:px-6 border-t border-gray-200 dark:border-gray-700">
			<form on:submit|preventDefault={sendMessage} class="flex">
				<textarea
					bind:value={message}
					on:keydown={handleKeyDown}
					placeholder="Type your message..."
					rows="3"
					class="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-l-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 resize-none bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-100 dark:placeholder-gray-400"
				></textarea>
				<button 
					type="submit"
					disabled={isLoading || !message.trim()}
					class="bg-indigo-600 dark:bg-indigo-700 text-white px-6 py-2 rounded-r-md hover:bg-indigo-700 dark:hover:bg-indigo-800 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800 disabled:opacity-50 disabled:cursor-not-allowed flex-shrink-0 transition-colors duration-200"
				>
					{isLoading ? 'Sending...' : 'Send'}
				</button>
			</form>
			<div class="flex justify-between mt-2">
				<p class="text-xs text-gray-500 dark:text-gray-400">Press Enter to send, Shift+Enter for new line</p>
				<p class="text-xs text-gray-500 dark:text-gray-400 font-medium tooltip" title={totalTokensTooltip}>{totalTokens} tokens</p>
			</div>
		</div>
	</div>
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