<script lang="ts">
    /**
     * ChatInput component
     * A reusable input component for chat messages with expandable textarea functionality
     */
    import { createEventDispatcher, onMount } from 'svelte';
    
    // Props with proper type definitions
    /**
     * The current message text
     */
    export let value = '';
    
    /**
     * Placeholder text to display when input is empty
     */
    export let placeholder = 'Type your message...';
    
    /**
     * Whether the input is disabled
     */
    export let disabled = false;
    
    /**
     * The button text to display
     */
    export let buttonText = 'Send';
    
    /**
     * Whether the button should show a loading state
     */
    export let isLoading = false;
    
    /**
     * Additional CSS classes for the container
     */
    export let containerClass = '';
    
    /**
     * Whether to send message on Enter key (Shift+Enter for new line)
     */
    export let sendWithEnter = true;

    /**
     * Token usage information if available
     */
    export let tokenUsage: { count?: number; limit?: number } = {};
    
    /**
     * Total tokens info with tooltip
     */
    export let totalTokensInfo: { count: number; tooltip: string } | null = null;
    
    /**
     * Whether to show keyboard shortcuts info
     */
    export let showKeyboardShortcuts = true;
    
    // Reference to the textarea element - used for focusing
    let textareaRef: HTMLTextAreaElement;

    // Set up event dispatcher for component events
    const dispatch = createEventDispatcher<{
        submit: { message: string };
        keydown: KeyboardEvent;
    }>();
    
    // Handle form submission
    function handleSubmit() {
        if (value.trim() && !disabled) {
            dispatch('submit', { message: value.trim() });
        }
    }
    
    // Handle keydown events
    function handleKeyDown(event: KeyboardEvent) {
        dispatch('keydown', event);
        
        // Handle Enter key for sending
        if (sendWithEnter && event.key === 'Enter' && !event.shiftKey && !disabled) {
            event.preventDefault();
            handleSubmit();
        }
    }

    onMount(() => {
        // Focus textarea when component is mounted if needed
        // This can be used later if the component needs to be focused
    });

    // Function to focus the textarea - can be called externally
    export function focus() {
        if (textareaRef) {
            textareaRef.focus();
        }
    }

    // Check if token usage info is available to display
    $: hasTokenInfo = typeof tokenUsage.count === 'number';
</script>

<div class="chat-input-container {containerClass}">
    <form on:submit|preventDefault={handleSubmit} class="flex flex-col w-full">
        <div class="w-full mb-2">
            <textarea
                bind:this={textareaRef}
                bind:value
                {placeholder}
                {disabled}
                class="chat-textarea"
                rows="3"
                on:keydown={handleKeyDown}
            ></textarea>
        </div>
        <div class="flex items-center justify-between mb-1">
            {#if showKeyboardShortcuts}
                <div class="hidden sm:flex items-center text-xs text-gray-500 dark:text-gray-400">
                    <span>Press Enter to send, Shift+Enter for new line</span>
                    <div class="ml-4 flex items-center">
                        <kbd class="mx-1 px-1.5 py-0.5 bg-gray-100 dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded text-xs">/</kbd>
                        <span class="ml-1">to focus</span>
                        <kbd class="mx-1 ml-3 px-1.5 py-0.5 bg-gray-100 dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded text-xs">n</kbd>
                        <span class="ml-1">new chat</span>
                    </div>
                </div>
            {/if}
            <div class="flex items-center ml-auto space-x-4">
                {#if totalTokensInfo}
                    <p class="text-xs text-gray-500 dark:text-gray-400 font-medium tooltip" title={totalTokensInfo.tooltip}>{totalTokensInfo.count} tokens</p>
                {/if}
                <button 
                    type="submit"
                    disabled={disabled || isLoading || !value.trim()}
                    class="chat-send-button"
                >
                    {#if isLoading}
                        <span class="inline-block h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"></span>
                    {/if}
                    {buttonText}
                </button>
            </div>
        </div>
        {#if hasTokenInfo}
            <div class="flex justify-end">
                <p class="text-xs text-gray-500 dark:text-gray-400">
                    {tokenUsage.count} token{tokenUsage.count !== 1 ? 's' : ''}
                    {#if tokenUsage.limit}
                        / {tokenUsage.limit}
                    {/if}
                </p>
            </div>
        {/if}
    </form>
</div>

<style>
    .chat-input-container {
        width: 100%;
    }
    
    .chat-textarea {
        width: 100%;
        min-height: 80px;
        padding: 0.75rem;
        border: 1px solid #d1d5db;
        border-radius: 0.375rem;
        background-color: white;
        color: #1f2937;
        font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, 
                     "Helvetica Neue", Arial, sans-serif;
        font-size: 1rem;
        line-height: 1.5;
        resize: none;
        outline: none;
        overflow-y: auto;
        box-sizing: border-box;
    }
    
    .chat-textarea:focus {
        border-color: #6366f1;
        box-shadow: 0 0 0 1px rgba(99, 102, 241, 0.3);
    }
    
    /* Dark mode support */
    :global(.dark) .chat-textarea {
        background-color: rgb(55, 65, 81); /* dark:bg-gray-700 */
        border-color: rgb(75, 85, 99); /* dark:border-gray-600 */
        color: #f3f4f6; /* dark:text-gray-100 */
    }
    
    :global(.dark) .chat-textarea:focus {
        border-color: #6366f1;
    }
    
    .chat-send-button {
        background-color: #6366f1; /* indigo-600 */
        color: white;
        padding: 0.5rem 1rem;
        border-radius: 0.375rem;
        font-size: 0.875rem;
        font-weight: 500;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: background-color 0.2s;
    }
    
    .chat-send-button:hover:not(:disabled) {
        background-color: #4f46e5; /* indigo-700 */
    }
    
    .chat-send-button:focus {
        outline: none;
        box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.5);
    }
    
    .chat-send-button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
    
    /* Dark mode support for button */
    :global(.dark) .chat-send-button {
        background-color: #4f46e5; /* dark:bg-indigo-700 */
    }
    
    :global(.dark) .chat-send-button:hover:not(:disabled) {
        background-color: #4338ca; /* dark:hover:bg-indigo-800 */
    }
    
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
</style> 