<script lang="ts">
    import pb from '$lib/pocketbase';
    import { theme } from '$lib/stores';
    // Use the public URL instead of importing the asset
    const glimmerLogoUrl = '/glimmer.svg';
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import ChatsSubMenu from './ChatsSubMenu.svelte';
    import { createNewChat as createNewChatAction } from '$lib/stores/chatStore';
    
    export let toggleSidebar: () => void;
    
    // Tracking search state
    let showChatSearch = false;
    let searchQuery = '';
    
    const navItems = [
        { href: '/dashboard', label: 'Dashboard', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6' },
        { href: '/practice-topics', label: 'Practice Topics', icon: 'M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253' },
        { href: '/learners', label: 'Learners', icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' },
        { href: '/chat', label: 'Chat', icon: 'M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z' },
    ];
    
    // Check if we're in the chat route
    $: isChatRoute = $page?.url?.pathname?.startsWith('/chat');
    
    // Get the chat ID if we're in a specific chat
    $: chatId = isChatRoute && $page?.params?.id ? $page.params.id : null;
    
    function isActive(href: string): boolean {
        if ($page) {
            const path = $page.url.pathname;
            if (href === '/dashboard' && path === '/dashboard') {
                return true;
            }
            return path.startsWith(href);
        }
        return false;
    }
    
    // Handle navigation with better mobile support
    async function handleNavClick(event: MouseEvent, href: string) {
        // Always prevent default to handle navigation manually
        event.preventDefault();
        event.stopPropagation();
        
        // Check screen size to determine if we should close the sidebar
        const isMobile = window.innerWidth < 768;
        
        try {
            // Navigate first
            await goto(href);
            
            // If mobile, close the sidebar after navigation
            if (isMobile && toggleSidebar) {
                // Add a slight delay to ensure navigation completes
                setTimeout(() => toggleSidebar(), 100);
            }
        } catch (error) {
            console.error('Navigation failed:', error);
        }
    }
    
    // Toggle search input
    function toggleSearch(event: MouseEvent) {
        event.preventDefault();
        event.stopPropagation();
        showChatSearch = !showChatSearch;
        
        if (showChatSearch) {
            // Focus the search input after it's shown
            setTimeout(() => {
                const searchInput = document.getElementById('chatSearchInput');
                if (searchInput) {
                    searchInput.focus();
                }
            }, 10);
        } else {
            // Clear search when hiding
            searchQuery = '';
        }
    }
    
    // Handle search input click to prevent navigation
    function handleSearchInputClick(event: MouseEvent) {
        event.preventDefault();
        event.stopPropagation();
    }
    
    // Update the createNewChat function to include required parameters and close sidebar on mobile
    async function createNewChat(event: MouseEvent) {
        event.preventDefault();
        event.stopPropagation();
        
        try {
            // Create a new chat with default system prompt and model
            const defaultSystemPrompt = "You are a helpful learning assistant for kids.";
            const newChatId = await createNewChatAction(defaultSystemPrompt);
            
            // Navigate to the new chat page
            await goto(`/chat/${newChatId}`);
            
            // Check screen size to determine if we should close the sidebar
            const isMobile = window.innerWidth < 768;
            if (isMobile && toggleSidebar) {
                // Add a slight delay to ensure navigation completes
                setTimeout(() => toggleSidebar(), 100);
            }
        } catch (error) {
            console.error('Failed to create new chat:', error);
        }
    }
    
    // Toggle theme function
    function toggleTheme() {
        $theme = $theme === 'dark' ? 'light' : 'dark';
    }
    
    // Get user information if available
    const user = pb.authStore.model;
    const userName = user ? user.name || user.username || user.email || 'User' : 'User';
    
    // Application version
    const appVersion = 'v0.0.1';
    // GitHub repo URL
    const githubUrl = 'https://github.com/busybytelab/glimmer';
</script>

<aside class="bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 w-64 h-full flex-shrink-0">
    <div class="h-full flex flex-col">
        <div class="flex items-center justify-center h-16 border-b border-gray-200 dark:border-gray-700">
            <img src={glimmerLogoUrl} alt="Glimmer Logo" class="h-8 w-8 mr-2" />
            <div class="text-xl font-semibold text-gray-800 dark:text-white">Glimmer</div>
        </div>
        
        <nav class="flex-1 overflow-y-auto p-4">
            <ul class="space-y-2">
                {#each navItems as item}
                    <li>
                        <div class="relative">
                            <a 
                                href={item.href}
                                on:click={(e) => handleNavClick(e, item.href)}
                                class={`flex items-center p-2 rounded-md transition-colors ${
                                    isActive(item.href) 
                                        ? 'bg-gray-100 dark:bg-gray-700 text-indigo-600 dark:text-indigo-400' 
                                        : 'text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700'
                                }`}
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-3 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={item.icon} />
                                </svg>
                                
                                {#if item.href === '/chat' && isActive(item.href) && showChatSearch}
                                    <!-- Search input replaces the label when search is active -->
                                    <input
                                        id="chatSearchInput"
                                        type="text"
                                        placeholder="Search chats..."
                                        bind:value={searchQuery}
                                        on:click={handleSearchInputClick}
                                        class="w-full bg-transparent border-none outline-none text-sm text-gray-700 dark:text-gray-200 placeholder-gray-500 dark:placeholder-gray-400"
                                    />
                                {:else}
                                    <span class="flex-1">{item.label}</span>
                                {/if}
                                
                                {#if item.href === '/chat' && isActive(item.href)}
                                    <!-- Search icon -->
                                    <button
                                        on:click={toggleSearch}
                                        class="p-1 rounded-md text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
                                        title={showChatSearch ? "Close search" : "Search chats"}
                                    >
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            {#if showChatSearch}
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                                            {:else}
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                                            {/if}
                                        </svg>
                                    </button>
                                    
                                    <!-- New chat icon -->
                                    <button
                                        on:click={createNewChat}
                                        class="ml-1 p-1 rounded-md text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
                                        title="New chat"
                                        aria-label="Create new chat"
                                    >
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                                        </svg>
                                    </button>
                                {/if}
                            </a>
                        </div>
                        
                        <!-- Show ChatsSubMenu only when the Chat menu item is active -->
                        {#if item.href === '/chat' && isChatRoute}
                            <ChatsSubMenu activeChatId={chatId} searchFilter={searchQuery} {toggleSidebar} />
                        {/if}
                    </li>
                {/each}
            </ul>
        </nav>
        
        <div class="mt-auto p-3 border-t border-gray-200 dark:border-gray-700 text-sm">
            <!-- Single row with user info on left, theme toggle and version on right -->
            <div class="flex items-center justify-between">
                <!-- User profile info -->
                <a 
                    href="/account" 
                    on:click={(e) => handleNavClick(e, '/account')}
                    class="flex items-center group text-gray-700 dark:text-gray-200 hover:text-indigo-600 dark:hover:text-indigo-400"
                >
                    <div class="bg-gray-200 dark:bg-gray-700 rounded-full p-1.5 mr-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                        </svg>
                    </div>
                    <span class="truncate max-w-[80px] group-hover:text-indigo-600 dark:group-hover:text-indigo-400">{userName}</span>
                </a>
                
                <!-- Theme toggle and version on right side -->
                <div class="flex items-center space-x-2">
                    <!-- Theme toggle switch -->
                    <button 
                        on:click={toggleTheme}
                        class="p-1.5 rounded-md text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-secondary"
                        aria-label="Toggle theme"
                    >
                        {#if $theme === 'dark'}
                            <!-- Sun icon for light mode -->
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                            </svg>
                        {:else}
                            <!-- Moon icon for dark mode -->
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                            </svg>
                        {/if}
                    </button>
                    
                    <!-- Version as GitHub link -->
                    <a 
                        href={githubUrl} 
                        target="_blank" 
                        rel="noopener noreferrer"
                        class="text-xs text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
                        title="View Glimmer on GitHub"
                    >
                        {appVersion}
                    </a>
                </div>
            </div>
        </div>
    </div>
</aside> 