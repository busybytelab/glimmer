<script lang="ts">
    /**
     * ChatInput component
     * A reusable input component for chat messages with expandable textarea functionality
     */
    import { createEventDispatcher } from 'svelte';
    import ExpandableTextArea from '../common/ExpandableTextArea.svelte';
    
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
    
    // Handle keydown events - updated to accept both KeyboardEvent and CustomEvent
    function handleKeyDown(event: KeyboardEvent | CustomEvent<KeyboardEvent>) {
        // If it's a CustomEvent, extract the original event
        const keyEvent = event instanceof CustomEvent ? event.detail : event;
        dispatch('keydown', keyEvent);
        
        // Handle Enter key for sending
        if (sendWithEnter && keyEvent.key === 'Enter' && !keyEvent.shiftKey && !disabled) {
            keyEvent.preventDefault();
            handleSubmit();
        }
    }
</script>

<div class="chat-input-container {containerClass}">
    <form on:submit|preventDefault={handleSubmit} class="flex flex-col">
        <div class="flex w-full">
            <div class="flex-1 chat-textarea-wrapper">
                <ExpandableTextArea
                    id="chat-input"
                    label=""
                    bind:value
                    {disabled}
                    {placeholder}
                    minRows={1}
                    maxRows={10}
                    cols=""
                    on:keydown={handleKeyDown}
                />
            </div>
            <button 
                type="submit"
                disabled={disabled || isLoading || !value.trim()}
                class="chat-send-button bg-indigo-600 dark:bg-indigo-700 text-white px-6 py-2 rounded-r-md hover:bg-indigo-700 dark:hover:bg-indigo-800 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800 disabled:opacity-50 disabled:cursor-not-allowed flex-shrink-0 transition-colors duration-200 flex items-center"
            >
                {#if isLoading}
                    <span class="inline-block h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"></span>
                {/if}
                {buttonText}
            </button>
        </div>
    </form>
</div>

<style>
    .chat-input-container {
        width: 100%;
    }
    
    .chat-textarea-wrapper {
        /* Override some ExpandableTextArea styling */
        position: relative;
    }
    
    .chat-textarea-wrapper :global(.editor-container) {
        /* Remove bottom margin */
        margin-bottom: 0;
        /* Fix border radius to match with button */
        border-radius: 0.375rem 0 0 0.375rem;
        /* Remove syntax highlighting (not needed for chat) */
        background-color: white;
    }
    
    /* Dark mode support */
    :global(.dark) .chat-textarea-wrapper :global(.editor-container) {
        background-color: rgb(55, 65, 81); /* dark:bg-gray-700 */
        border-color: rgb(75, 85, 99); /* dark:border-gray-600 */
    }
    
    .chat-textarea-wrapper :global(.syntax-highlight-pre),
    .chat-textarea-wrapper :global(.syntax-highlight-textarea) {
        /* Override font to match normal input */
        font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, 
                     "Helvetica Neue", Arial, sans-serif;
        /* Slightly larger font for better readability */
        font-size: 1rem;
        /* Smaller padding to match original design */
        padding: 0.5rem 0.75rem;
        /* Enable word wrap for chat messages */
        white-space: pre-wrap;
        /* Ensure text color matches theme */
        color: #1f2937;
    }
    
    /* Dark mode text color */
    :global(.dark) .chat-textarea-wrapper :global(.syntax-highlight-textarea) {
        caret-color: #f3f4f6; /* dark:text-gray-100 */
    }
    
    :global(.dark) .chat-textarea-wrapper :global(.syntax-highlight-pre) {
        color: #f3f4f6; /* dark:text-gray-100 */
        background-color: transparent;
    }
    
    /* Fix height for the send button to match the input */
    .chat-send-button {
        align-self: stretch;
        display: flex;
        align-items: center;
        justify-content: center;
    }
</style> 