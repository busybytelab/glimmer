<script lang="ts">
	import AppLayout from '../components/layout/AppLayout.svelte';
	import { error as errorStore } from '$lib/stores';
	import { llmService } from '$lib/services/llm';
	import type { Usage } from '$lib/services/llm';
	import SelectField from '../components/common/SelectField.svelte';
	
	type Message = {
		role: 'user' | 'assistant';
		content: string;
		timestamp: Date;
		usage?: Usage;
	};
	
	let sidebarOpen = true;
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
	
	// Available models - will be populated from backend
	let availableModels: { id: string; name: string; isDefault?: boolean }[] = [];
	
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
	
	// Fetch models when component mounts
	fetchModels();
	
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
			
			// Scroll to bottom
			setTimeout(() => {
				const chatContainer = document.querySelector('.chat-container');
				if (chatContainer) {
					chatContainer.scrollTop = chatContainer.scrollHeight;
				}
			}, 50);
			
		} catch (err) {
			console.error('Error sending message:', err);
			errorStore.set(err instanceof Error ? err.message : 'Failed to send message');
		} finally {
			isLoading = false;
		}
	}
	
	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' && !event.shiftKey) {
			event.preventDefault();
			sendMessage();
		}
	}
</script>

<AppLayout bind:sidebarOpen>
	<div class="max-w-4xl mx-auto px-4 sm:px-6 md:px-8 my-8">
		<div class="flex justify-between items-center mb-6">
			<h1 class="text-2xl font-semibold text-gray-900">Chat with LLM</h1>
			
			<div class="flex items-center">
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
		</div>
		
		<div class="bg-white shadow rounded-lg">
			<!-- Chat messages container with auto-expanding height -->
			<div class="px-4 py-5 sm:p-6 chat-container">
				{#if messages.length === 0}
					<div class="text-center py-12 border border-dashed border-gray-200 rounded-lg">
						<p class="text-gray-500">
							No messages yet. Start a conversation below.
						</p>
					</div>
				{:else}
					<div class="space-y-6">
						{#each messages as msg}
							<div class="flex flex-col">
								<div class="flex items-center mb-2">
									<span class="font-medium text-sm text-gray-700 mr-2">
										{msg.role === 'user' ? 'You' : 'Assistant'}
									</span>
									<span class="text-xs text-gray-500">
										{msg.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
									</span>
									{#if msg.role === 'assistant' && msg.usage}
										<span class="text-xs text-gray-400 ml-2 tooltip" 
										      title="Model: {msg.usage.LlmModelName || 'Default'}&#10;Cache hit: {msg.usage.CacheHit ? 'Yes' : 'No'}&#10;Prompt: {msg.usage.PromptTokens} tokens&#10;Completion: {msg.usage.CompletionTokens} tokens&#10;Cost: ${(msg.usage.Cost || 0).toFixed(6)}">
											({msg.usage.TotalTokens} tokens)
										</span>
									{/if}
								</div>
								<div class="{msg.role === 'user' ? 'bg-blue-50 border-blue-100' : 'bg-gray-50 border-gray-100'} border rounded-lg p-4">
									<p class="whitespace-pre-wrap break-words text-left text-gray-800">{msg.content}</p>
								</div>
							</div>
						{/each}
						{#if isLoading}
							<div class="flex flex-col">
								<div class="flex items-center mb-2">
									<span class="font-medium text-sm text-gray-700">Assistant</span>
								</div>
								<div class="bg-gray-50 border border-gray-100 rounded-lg p-4">
									<div class="flex space-x-2">
										<div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
										<div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
										<div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.4s"></div>
									</div>
								</div>
							</div>
						{/if}
					</div>
				{/if}
			</div>
			
			<!-- Message input -->
			<div class="border-t border-gray-200 px-4 py-5 sm:p-6">
				<div class="flex">
					<textarea
						bind:value={message}
						on:keydown={handleKeyDown}
						placeholder="Type your message..."
						rows="3"
						class="		dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600 dark:placeholder-gray-500
 flex-1 px-4 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-0 focus:ring-secondary focus:border-secondary resize-none"
					></textarea>
					<button 
						on:click={sendMessage}
						disabled={isLoading || !message.trim()}
						class="bg-secondary text-white px-6 py-2 rounded-r-md hover:bg-blue-600 focus:outline-none focus:ring-0 focus:ring-secondary focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed flex-shrink-0"
					>
						{isLoading ? 'Sending...' : 'Send'}
					</button>
				</div>
				<div class="flex justify-between mt-2">
					<p class="text-xs text-gray-500">Press Enter to send, Shift+Enter for new line</p>
					<p class="text-xs text-gray-500 font-medium tooltip" title={totalTokensTooltip}>{totalTokens} tokens</p>
				</div>
			</div>
		</div>
	</div>
</AppLayout>

<style>
	/* Remove fixed height from chat container to allow it to grow naturally */
	.chat-container {
		max-height: none;
		min-height: 400px;
	}
	
	/* Tooltip styles */
	.tooltip {
		position: relative;
		cursor: help;
	}
	
	.tooltip:hover::before {
		content: attr(title);
		position: absolute;
		bottom: 100%;
		right: 0;
		transform: translateY(-5px);
		background-color: #333;
		color: white;
		padding: 0.5rem 0.75rem;
		border-radius: 0.25rem;
		font-size: 0.75rem;
		white-space: pre-line;
		width: max-content;
		max-width: 250px;
		z-index: 10;
		font-weight: normal;
		box-shadow: 0 3px 10px rgba(0, 0, 0, 0.2);
		text-align: left;
	}
</style> 