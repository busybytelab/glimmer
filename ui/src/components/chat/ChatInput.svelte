<script lang="ts">
    /**
     * ChatInput component
     * A reusable input component for chat messages with expandable textarea functionality
     */
    import { createEventDispatcher, onMount } from 'svelte';
    import ExpandableTextArea from '../common/ExpandableTextArea.svelte';
    import type { SvelteComponent } from 'svelte';
    
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

    /**
     * Whether to enable markdown syntax highlighting
     */
    export let enableMarkdown = true;

    /**
     * Whether the textarea should be expanded by default
     */
    export let defaultExpanded = true;
    
    // Reference to the textarea element - used for focusing
    let textareaRef: SvelteComponent;
    let textareaWrapper: HTMLDivElement;

    // Calculate a reasonable maxRows value that effectively keeps the textarea expanded
    // This should be high enough to show most content but not so high that it becomes unwieldy
    const effectiveMaxRows = 20; // Same as the default in ExpandableTextArea

    // Set minRows to a value that provides enough initial space
    const effectiveMinRows = 3;

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
    function handleKeydown(event: CustomEvent<KeyboardEvent>) {
        const keyboardEvent = event.detail;
        dispatch('keydown', keyboardEvent);
        
        // Handle Enter key for sending
        if (sendWithEnter && keyboardEvent.key === 'Enter' && !keyboardEvent.shiftKey && !disabled) {
            keyboardEvent.preventDefault();
            handleSubmit();
        }
    }

    onMount(() => {
        // Ensure the textarea is properly sized for initial value
        if (textareaRef && value) {
            // Force a resize by dispatching an input event
            const event = new Event('input', { bubbles: true });
            const textarea = textareaWrapper?.querySelector('textarea');
            if (textarea) {
                textarea.dispatchEvent(event);
            }
        }
    });

    // Function to focus the textarea - can be called externally
    export function focus() {
        if (textareaRef) {
            const textarea = textareaWrapper?.querySelector('textarea');
            if (textarea) {
                textarea.focus();
            }
        }
    }

    // Check if token usage info is available to display
    $: hasTokenInfo = typeof tokenUsage.count === 'number';
</script>

<div class="chat-input-container {containerClass}">
    <form on:submit|preventDefault={handleSubmit} class="flex flex-col w-full">
        <div class="w-full mb-2 relative">
            <div class="chat-textarea-wrapper" bind:this={textareaWrapper}>
                <ExpandableTextArea
                    bind:this={textareaRef}
                    bind:value
                    id="chat-input"
                    label=""
                    {placeholder}
                    {disabled}
                    minRows={effectiveMinRows}
                    maxRows={effectiveMaxRows}
                    on:keydown={handleKeydown}
                    cols="w-full"
                    language={enableMarkdown ? 'markdown' : ''}
                />
            </div>
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

    .chat-textarea-wrapper {
        position: relative;
        width: 100%;
    }

    .chat-textarea-wrapper :global(.editor-container) {
        border: none !important;
        position: relative !important;
        padding: 0 !important;
    }

    /* Style the pre element that shows the text */
    .chat-textarea-wrapper :global(.syntax-highlight-pre) {
        margin: 0 !important;
        background-color: white !important;
        color: #1f2937 !important;
        padding: 0.75rem !important;
        font-family: inherit !important;
        font-size: 1rem !important;
        line-height: 1.5 !important;
        white-space: pre-wrap !important;
        word-break: break-word !important;
        overflow-wrap: break-word !important;
        tab-size: 4 !important;
        text-align: left !important;
    }

    .chat-textarea-wrapper :global(.syntax-highlight-pre code) {
        display: block !important;
        min-width: 100% !important;
        font-family: inherit !important;
        font-size: inherit !important;
        line-height: inherit !important;
        text-align: left !important;
    }

    .chat-textarea-wrapper :global(textarea) {
        position: absolute !important;
        top: 0 !important;
        left: 0 !important;
        text-align: left !important;
        background-color: transparent !important;
        border: 1px solid #d1d5db;
        border-radius: 0.375rem;
        padding: 0.75rem !important;
        margin: 0 !important;
        width: 100% !important;
        height: 100% !important;
        font-family: inherit !important;
        font-size: 1rem !important;
        line-height: 1.5 !important;
        resize: none;
        transition: border-color 0.2s, box-shadow 0.2s;
        white-space: pre-wrap !important;
        word-break: break-word !important;
        overflow-wrap: break-word !important;
        tab-size: 4 !important;
        caret-color: #1f2937 !important;
        overflow: hidden !important;
        box-sizing: border-box !important;
    }

    .chat-textarea-wrapper :global(textarea:focus) {
        border-color: #6366f1;
        box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.3);
        outline: none;
    }

    /* Dark mode support */
    :global(.dark) .chat-textarea-wrapper :global(.syntax-highlight-pre) {
        background-color: rgb(55, 65, 81) !important; /* dark:bg-gray-700 */
        color: #f3f4f6 !important;
    }

    :global(.dark) .chat-textarea-wrapper :global(textarea) {
        background-color: transparent !important;
        border-color: rgb(75, 85, 99); /* dark:border-gray-600 */
        caret-color: #f3f4f6 !important;
    }

    :global(.dark) .chat-textarea-wrapper :global(textarea:focus) {
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

    /* Add syntax highlighting styles for markdown */
    .chat-textarea-wrapper :global(.token.heading),
    .chat-textarea-wrapper :global(.token.title) {
        color: #2563eb !important; /* blue-600 */
        font-weight: 600 !important;
    }

    .chat-textarea-wrapper :global(.token.bold) {
        font-weight: 600 !important;
    }

    .chat-textarea-wrapper :global(.token.italic) {
        font-style: italic !important;
    }

    .chat-textarea-wrapper :global(.token.list),
    .chat-textarea-wrapper :global(.token.bullet) {
        color: #059669 !important; /* green-600 */
    }

    .chat-textarea-wrapper :global(.token.url),
    .chat-textarea-wrapper :global(.token.link) {
        color: #2563eb !important; /* blue-600 */
        text-decoration: underline !important;
    }

    .chat-textarea-wrapper :global(.token.blockquote) {
        color: #4b5563 !important; /* gray-600 */
        font-style: italic !important;
    }

    .chat-textarea-wrapper :global(.token.code),
    .chat-textarea-wrapper :global(.token.codespan) {
        color: #dc2626 !important; /* red-600 */
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace !important;
    }

    /* Dark mode syntax highlighting */
    :global(.dark) .chat-textarea-wrapper :global(.token.heading),
    :global(.dark) .chat-textarea-wrapper :global(.token.title) {
        color: #3b82f6 !important; /* blue-500 */
    }

    :global(.dark) .chat-textarea-wrapper :global(.token.list),
    :global(.dark) .chat-textarea-wrapper :global(.token.bullet) {
        color: #10b981 !important; /* green-500 */
    }

    :global(.dark) .chat-textarea-wrapper :global(.token.url),
    :global(.dark) .chat-textarea-wrapper :global(.token.link) {
        color: #3b82f6 !important; /* blue-500 */
    }

    :global(.dark) .chat-textarea-wrapper :global(.token.blockquote) {
        color: #9ca3af !important; /* gray-400 */
    }

    :global(.dark) .chat-textarea-wrapper :global(.token.code),
    :global(.dark) .chat-textarea-wrapper :global(.token.codespan) {
        color: #ef4444 !important; /* red-500 */
    }
</style> 