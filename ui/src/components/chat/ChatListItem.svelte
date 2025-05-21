<!-- DEPRECATED: This component is deprecated and should be removed. Make sure no usage/import -->

<script lang="ts">
    import type { Chat } from '$lib/services/chat';
    import { updateChatTitle } from '$lib/stores/chatStore';
    import { createEventDispatcher } from 'svelte';
    
    export let chat: Chat;
    export let isActive: boolean = false;
    // Allow passing the last message content from parent component
    export let lastMessageContent: string | null = null;
    
    // State for edit mode
    let isEditing = false;
    let editedTitle = chat.title || 'New Chat';
    let inputEl: HTMLInputElement; // Reference to the input element
    
    // State for dropdown menu
    let isMenuOpen = false;
    
    const dispatch = createEventDispatcher();
    
    // Function to start editing
    function startEditing() {
        editedTitle = chat.title || 'New Chat';
        isEditing = true;
        
        // Focus the input after DOM update
        setTimeout(() => {
            if (inputEl) {
                inputEl.focus();
                inputEl.select();
            }
        }, 10);
    }
    
    // Function to save edited title
    async function saveTitle() {
        if (editedTitle.trim() === '') {
            editedTitle = 'New Chat';
        }
        
        try {
            if (editedTitle !== chat.title) {
                await updateChatTitle(chat.id, editedTitle);
            }
        } catch (error) {
            console.error('Failed to update chat title:', error);
        } finally {
            isEditing = false;
        }
    }
    
    // Function to cancel editing
    function cancelEditing() {
        isEditing = false;
        editedTitle = chat.title || 'New Chat';
    }
    
    // Toggle menu open/closed
    function toggleMenu(event: MouseEvent) {
        event.stopPropagation();
        isMenuOpen = !isMenuOpen;
    }
    
    // Keyboard handler for toggle menu
    function handleKeyToggleMenu(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.stopPropagation();
            isMenuOpen = !isMenuOpen;
        }
    }
    
    // Close menu when clicking outside
    function handleClickOutside() {
        if (isMenuOpen) {
            isMenuOpen = false;
        }
    }
    
    // Menu actions
    function handleDeleteChat(event: MouseEvent) {
        event.stopPropagation();
        isMenuOpen = false;
        dispatch('delete', { chatId: chat.id });
    }
    
    // Keyboard handler for delete chat
    function handleKeyDeleteChat(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.stopPropagation();
            isMenuOpen = false;
            dispatch('delete', { chatId: chat.id });
        }
    }
    
    function handleArchiveChat(event: MouseEvent) {
        event.stopPropagation();
        isMenuOpen = false;
        dispatch('archive', { chatId: chat.id });
    }
    
    // Keyboard handler for archive chat
    function handleKeyArchiveChat(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.stopPropagation();
            isMenuOpen = false;
            dispatch('archive', { chatId: chat.id });
        }
    }
    
    // Handle key events in the input
    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.preventDefault();
            saveTitle();
        } else if (event.key === 'Escape') {
            event.preventDefault();
            cancelEditing();
        }
    }
    
    /**
     * Creates a short headline from a message by extracting the first few words
     * @param message The full message text
     * @param wordCount Number of words to include in the headline
     * @returns A short headline or "No messages yet" if message is empty
     */
    function createMessageHeadline(message: string | null, wordCount: number = 4): string {
        if (!message) return "No messages yet";
        
        const words = message.trim().split(/\s+/);
        if (words.length <= wordCount) return message;
        
        return words.slice(0, wordCount).join(' ') + '...';
    }
    
    // Use passed lastMessageContent to create a headline
    $: lastMessageHeadline = createMessageHeadline(lastMessageContent);
    
    // Format relative time (like "10m ago", "1h ago", "2d ago")
    function formatRelativeTime(date: Date): string {
        const now = new Date();
        const diffMs = now.getTime() - date.getTime();
        const diffMins = Math.floor(diffMs / (1000 * 60));
        
        if (diffMins < 60) {
            return `${diffMins}m ago`;
        }
        
        const diffHours = Math.floor(diffMins / 60);
        if (diffHours < 24) {
            return `${diffHours}h ago`;
        }
        
        const diffDays = Math.floor(diffHours / 24);
        return `${diffDays}d ago`;
    }
</script>

<svelte:window on:click={handleClickOutside} />

<button 
    on:click={() => {
        if (!isEditing) {
            dispatch('click');
        }
    }}
    class="w-full text-left p-2 rounded-md transition-colors {isActive ? 'bg-indigo-100 dark:bg-indigo-900/30' : 'hover:bg-gray-100 dark:hover:bg-gray-800'}"
>
    <div class="flex justify-between items-start relative">
        {#if isEditing}
            <!-- Edit mode -->
            <div class="flex-1 flex items-center">
                <input 
                    bind:this={inputEl}
                    bind:value={editedTitle}
                    on:keydown={handleKeyDown}
                    class="w-full font-medium text-sm text-gray-900 dark:text-gray-100 bg-white dark:bg-gray-700 border border-indigo-300 dark:border-indigo-600 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                />
                <div class="flex ml-1">
                    <span 
                        role="button"
                        tabindex="0"
                        on:click|stopPropagation={saveTitle}
                        on:keydown={(e) => e.key === 'Enter' && saveTitle()}
                        class="text-green-600 dark:text-green-400 hover:text-green-700 dark:hover:text-green-300 p-1 cursor-pointer"
                        title="Save"
                        aria-label="Save title"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                        </svg>
                    </span>
                    <span 
                        role="button"
                        tabindex="0"
                        on:click|stopPropagation={cancelEditing}
                        on:keydown={(e) => e.key === 'Enter' && cancelEditing()}
                        class="text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300 p-1 cursor-pointer"
                        title="Cancel"
                        aria-label="Cancel editing"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                        </svg>
                    </span>
                </div>
            </div>
        {:else}
            <!-- Display mode -->
            <div class="group flex items-center max-w-[140px]">
                <span class="font-medium text-sm text-gray-900 dark:text-gray-100 truncate">
                    {chat.title || "[No Title]"}
                </span>
                <span 
                    role="button"
                    tabindex="0"
                    on:click|stopPropagation={startEditing}
                    on:keydown={(e) => e.key === 'Enter' && startEditing()}
                    class="ml-1 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer"
                    title="Edit title"
                    aria-label="Edit title"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                        <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                    </svg>
                </span>
            </div>
            
            <!-- Right side with time and menu -->
            <div class="flex items-center">
                <span class="text-xs text-gray-500 dark:text-gray-400 whitespace-nowrap mr-1.5">
                    {formatRelativeTime(chat.updatedAt)}
                </span>
                
                <!-- Three dots menu button -->
                <div class="relative">
                    <span 
                        role="button"
                        tabindex="0"
                        on:click|stopPropagation={toggleMenu}
                        on:keydown={(e) => e.key === 'Enter' && handleKeyToggleMenu(e)}
                        class="text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 p-1 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer"
                        title="Chat options"
                        aria-label="Chat options menu"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                            <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z" />
                        </svg>
                    </span>
                    
                    <!-- Dropdown menu -->
                    {#if isMenuOpen}
                        <div class="absolute right-0 mt-1 w-36 bg-white dark:bg-gray-800 rounded-md shadow-lg border border-gray-200 dark:border-gray-700 z-10">
                            <div class="py-1">
                                <span
                                    role="button"
                                    tabindex="0"
                                    on:click|stopPropagation={handleArchiveChat}
                                    on:keydown={(e) => e.key === 'Enter' && handleKeyArchiveChat(e)}
                                    class="block w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center cursor-pointer"
                                    aria-label="Archive chat"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" viewBox="0 0 20 20" fill="currentColor">
                                        <path d="M4 3a2 2 0 100 4h12a2 2 0 100-4H4z" />
                                        <path fill-rule="evenodd" d="M3 8h14v7a2 2 0 01-2 2H5a2 2 0 01-2-2V8zm5 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" clip-rule="evenodd" />
                                    </svg>
                                    Archive
                                </span>
                                <span
                                    role="button"
                                    tabindex="0"
                                    on:click|stopPropagation={handleDeleteChat}
                                    on:keydown={(e) => e.key === 'Enter' && handleKeyDeleteChat(e)}
                                    class="block w-full text-left px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center cursor-pointer"
                                    aria-label="Delete chat"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                                    </svg>
                                    Delete
                                </span>
                            </div>
                        </div>
                    {/if}
                </div>
            </div>
        {/if}
    </div>
    <p class="text-xs text-gray-600 dark:text-gray-300 mt-1 truncate">
        {lastMessageHeadline}
    </p>
    <div class="mt-1 flex items-center">
        <span class="text-xs px-1.5 py-0.5 bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 rounded">
            {chat.model || 'Default Model'}
        </span>
    </div>
</button> 