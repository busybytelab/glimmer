<!-- ChatsSubMenu.svelte -->
<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { chatListStore, loadChats, searchChats } from '$lib/stores/chatStore';
    import { goto } from '$app/navigation';
    import { error as errorStore } from '$lib/stores';
    
    // Active chat handling (for highlighting)
    export let activeChatId: string | null = null;
    
    // Search filter from parent component
    export let searchFilter: string = '';
    
    // Access to the toggle sidebar function from parent
    export let toggleSidebar: (() => void) | undefined = undefined;
    
    // Search state
    let isSearching = false;
    let searchTimeout: ReturnType<typeof setTimeout> | null = null;
    
    // Tooltip state
    let tooltipText = '';
    let tooltipVisible = false;
    let tooltipX = 0;
    let tooltipY = 0;
    
    // Format relative time (like "5m", "2h", "3d")
    function formatRelativeTime(date: Date): string {
        const now = new Date();
        const diffMs = now.getTime() - date.getTime();
        const diffMins = Math.floor(diffMs / (1000 * 60));
        
        if (diffMins < 60) {
            return `${diffMins}m`;
        }
        
        const diffHours = Math.floor(diffMins / 60);
        if (diffHours < 24) {
            return `${diffHours}h`;
        }
        
        const diffDays = Math.floor(diffHours / 24);
        return `${diffDays}d`;
    }

    // Format full date for tooltip
    function formatFullDate(date: Date): string {
        return date.toLocaleString(undefined, {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    // Show tooltip
    function showTooltip(event: MouseEvent, date: Date) {
        tooltipText = formatFullDate(date);
        tooltipVisible = true;
        tooltipX = event.clientX;
        tooltipY = event.clientY + 20; // Offset below cursor
    }

    // Hide tooltip
    function hideTooltip() {
        tooltipVisible = false;
    }
    
    // Watch for changes to searchFilter from parent
    $: {
        if (searchFilter && searchFilter.trim()) {
            performSearch(searchFilter);
        } else if (searchFilter === '') {
            // If search is cleared, load all chats
            loadChats().catch(err => {
                console.error('Failed to load chats:', err);
                errorStore.set(err instanceof Error ? err.message : 'Failed to load chats');
            });
        }
    }
    
    // Function to perform search with debounce
    function performSearch(query: string) {
        if (searchTimeout) {
            clearTimeout(searchTimeout);
        }
        
        searchTimeout = setTimeout(() => {
            if (query.trim()) {
                isSearching = true;
                searchChats(query)
                    .catch(err => {
                        console.error('Search failed:', err);
                        errorStore.set(err instanceof Error ? err.message : 'Failed to search chats');
                    })
                    .finally(() => {
                        isSearching = false;
                    });
            }
        }, 300); // 300ms debounce
    }
    
    onMount(() => {
        // Load chats when component mounts
        loadChats().catch(err => {
            console.error('Failed to load chats:', err);
            errorStore.set(err instanceof Error ? err.message : 'Failed to load chats');
        });
    });
    
    onDestroy(() => {
        // Clean up timeout on component destroy
        if (searchTimeout) {
            clearTimeout(searchTimeout);
        }
    });
    
    async function selectChat(id: string) {
        // Navigate to the selected chat
        await goto(`/chat/${id}`);
        
        // Check screen size to determine if we should close the sidebar
        const isMobile = window.innerWidth < 768;
        if (isMobile && toggleSidebar) {
            // Add a slight delay to ensure navigation completes
            setTimeout(() => toggleSidebar(), 100);
        }
    }
</script>

<div class="pl-6 mt-2 mb-4">
    <!-- Chat list -->
    <div class="max-h-48 overflow-y-auto pr-1">
        {#if $chatListStore.loading && !isSearching}
            <div class="flex justify-center items-center py-2">
                <div class="animate-spin rounded-full h-4 w-4 border-t-2 border-indigo-500"></div>
            </div>
        {:else if $chatListStore.error}
            <div class="text-center text-red-500 dark:text-red-400 text-xs py-2">
                {$chatListStore.error}
                <button
                    on:click={() => loadChats()}
                    class="mt-1 text-indigo-600 dark:text-indigo-400 underline text-xs"
                >
                    Retry
                </button>
            </div>
        {:else if $chatListStore.chats.length === 0}
            {#if searchFilter.trim()}
                <p class="text-center text-gray-500 dark:text-gray-400 text-xs py-2">No chats matching "{searchFilter}"</p>
            {:else}
                <p class="text-center text-gray-500 dark:text-gray-400 text-xs py-2">No previous chats</p>
            {/if}
        {:else}
            <div class="space-y-1">
                {#each $chatListStore.chats as chat}
                    <div 
                        class="flex items-center justify-between p-1.5 rounded cursor-pointer text-sm truncate
                            {activeChatId === chat.id ? 
                            'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300' : 
                            'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300'}"
                        on:click={() => selectChat(chat.id)}
                        on:keydown={(e) => e.key === 'Enter' && selectChat(chat.id)}
                        tabindex="0"
                        role="button"
                        aria-label="Open chat: {chat.title || 'Untitled Chat'}"
                    >
                        <!-- Chat icon and title -->
                        <div class="flex items-center min-w-0">
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 mr-1.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M18 10c0 3.866-3.582 7-8 7a8.841 8.841 0 01-4.083-.98L2 17l1.338-3.123C2.493 12.767 2 11.434 2 10c0-3.866 3.582-7 8-7s8 3.134 8 7zM7 9H5v2h2V9zm8 0h-2v2h2V9zM9 9h2v2H9V9z" clip-rule="evenodd" />
                            </svg>
                            <span class="truncate text-xs">{chat.title || 'Untitled Chat'}</span>
                        </div>
                        <!-- Relative time with custom tooltip -->
                        <span 
                            class="text-xs text-gray-500 dark:text-gray-400 ml-2 flex-shrink-0 cursor-help"
                            on:mouseenter={(e) => showTooltip(e, chat.updatedAt)}
                            on:mouseleave={hideTooltip}
                            role="tooltip"
                            aria-label="Last updated: {formatFullDate(chat.updatedAt)}"
                        >
                            {formatRelativeTime(chat.updatedAt)}
                        </span>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>

<!-- Custom tooltip -->
{#if tooltipVisible}
    <div 
        class="fixed z-50 px-2 py-1 text-xs bg-gray-900 text-white rounded shadow-lg pointer-events-none"
        style="left: {tooltipX}px; top: {tooltipY}px; transform: translateX(-50%);"
    >
        {tooltipText}
    </div>
{/if} 