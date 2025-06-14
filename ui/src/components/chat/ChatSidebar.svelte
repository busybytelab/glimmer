<!-- DEPRECATED: This component is deprecated and should be removed. Make sure no usage/import -->

<!-- ChatSidebar.svelte -->
<script lang="ts">
    import { onMount } from 'svelte';
    import ChatListItem from './ChatListItem.svelte';
    import { chatListStore, loadChats, searchChats, deleteChat, archiveChat } from '$lib/stores/chatStore';
    import { goto } from '$app/navigation';
    import { error as errorStore } from '$lib/stores';
    
    // Active chat handling (for highlighting)
    export let activeChatId: string | null = null;
    
    // Search functionality
    let searchQuery = '';
    let isSearching = false;
    let searchTimeout: ReturnType<typeof setTimeout> | null = null;
    
    // Watch for changes to searchQuery
    $: {
        // When search query changes, debounce the search
        if (searchTimeout) {
            clearTimeout(searchTimeout);
        }
        
        searchTimeout = setTimeout(() => {
            // Only search if query has content
            if (searchQuery.trim()) {
                isSearching = true;
                searchChats(searchQuery)
                    .catch(err => {
                        console.error('Search failed:', err);
                        errorStore.set(err instanceof Error ? err.message : 'Failed to search chats');
                    })
                    .finally(() => {
                        isSearching = false;
                    });
            } else {
                // If query is empty, load all chats
                loadChats().catch(err => {
                    console.error('Failed to load chats:', err);
                    errorStore.set(err instanceof Error ? err.message : 'Failed to load chats');
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
    
    function selectChat(id: string) {
        goto(`/account/chat/${id}`);
    }

    // Handle delete chat
    async function handleDeleteChat(chatId: string) {
        try {
            await deleteChat(chatId);
            // If we deleted the active chat, navigate to the home page
            if (activeChatId === chatId) {
                goto('/account/chat');
            }
        } catch (error) {
            console.error('Failed to delete chat:', error);
            // Could add a toast notification here
        }
    }

    // Handle archive chat
    async function handleArchiveChat(chatId: string) {
        try {
            await archiveChat(chatId, true);
            // No need to navigate away as the chat will just be hidden
        } catch (error) {
            console.error('Failed to archive chat:', error);
            // Could add a toast notification here
        }
    }
</script>

<div class="h-full flex flex-col bg-gray-50 dark:bg-gray-900 border-r border-gray-200 dark:border-gray-700 w-64">
    <div class="p-2 pt-4 border-0 border-gray-200 dark:border-gray-700">
        <div class="relative w-full">
            <input
                type="text"
                placeholder="Search chats..."
                bind:value={searchQuery}
                class="w-full px-3 py-2 pl-9 bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md text-sm text-gray-800 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            />
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-500 dark:text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
                </svg>
            </div>
            {#if isSearching}
                <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                    <div class="animate-spin rounded-full h-4 w-4 border-t-2 border-indigo-500"></div>
                </div>
            {/if}
        </div>
    </div>
    
    <div class="flex-1 overflow-y-auto p-2">
        {#if $chatListStore.loading && !isSearching}
            <div class="flex justify-center items-center py-4">
                <div class="animate-spin rounded-full h-6 w-6 border-t-2 border-indigo-500"></div>
            </div>
        {:else if $chatListStore.error}
            <div class="text-center text-red-500 dark:text-red-400 p-4">
                {$chatListStore.error}
                <button
                    on:click={() => loadChats()}
                    class="mt-2 text-indigo-600 dark:text-indigo-400 underline text-sm"
                >
                    Retry
                </button>
            </div>
        {:else if $chatListStore.chats.length === 0}
            {#if searchQuery.trim()}
                <p class="text-center text-gray-500 dark:text-gray-400 p-4">No chats matching "{searchQuery}"</p>
            {:else}
                <p class="text-center text-gray-500 dark:text-gray-400 p-4">No previous chats</p>
            {/if}
        {:else}
            <div class="space-y-1">
                {#each $chatListStore.chats as chat}
                    <ChatListItem 
                        {chat}
                        isActive={activeChatId === chat.id}
                        lastMessageContent={chat.lastMessage?.content}
                        on:click={() => selectChat(chat.id)}
                        on:delete={event => handleDeleteChat(event.detail.chatId)}
                        on:archive={event => handleArchiveChat(event.detail.chatId)}
                    />
                {/each}
            </div>
        {/if}
    </div>
</div> 