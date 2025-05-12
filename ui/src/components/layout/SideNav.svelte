<script lang="ts">
    import pb from '$lib/pocketbase';
    // Use the public URL instead of importing the asset
    const glimmerLogoUrl = '/glimmer.svg';
    import { page } from '$app/stores';
    
    export let isOpen: boolean = true;
    
    const navItems = [
        { href: '/dashboard', label: 'Dashboard', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6' },
        { href: '/practice-topics', label: 'Practice Topics', icon: 'M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253' },
        { href: '/learners', label: 'Learners', icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' },
        { href: '/chat', label: 'Chat', icon: 'M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z' },
        { href: '/profile', label: 'Profile', icon: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z' },
        { href: '/settings', label: 'Settings', icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z' },
    ];
    
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
</script>

<aside class={`bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 w-64 flex-shrink-0 ${isOpen ? '' : 'hidden'} md:block`}>
    <div class="h-full flex flex-col">
        <div class="flex items-center justify-center h-16 border-b border-gray-200 dark:border-gray-700">
            <img src={glimmerLogoUrl} alt="Glimmer Logo" class="h-8 w-8 mr-2" />
            <div class="text-xl font-semibold text-gray-800 dark:text-white">Glimmer</div>
        </div>
        
        <nav class="flex-1 overflow-y-auto p-4">
            <ul class="space-y-2">
                {#each navItems as item}
                    <li>
                        <a 
                            href={item.href}
                            class={`flex items-center p-2 rounded-md transition-colors ${
                                isActive(item.href) 
                                    ? 'bg-gray-100 dark:bg-gray-700 text-indigo-600 dark:text-indigo-400' 
                                    : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                            }`}
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={item.icon} />
                            </svg>
                            <span>{item.label}</span>
                        </a>
                    </li>
                {/each}
            </ul>
        </nav>
        
        <div class="p-4 border-t border-gray-200 dark:border-gray-700">
            <button 
                class="flex items-center justify-center w-full p-2 bg-red-50 dark:bg-red-900/30 text-red-600 dark:text-red-400 rounded-md hover:bg-red-100 dark:hover:bg-red-900/50 transition-colors"
                on:click={() => {
                    // Clear PocketBase auth store
                    pb.authStore.clear();
                    // Clear the auth cookie
                    document.cookie = 'pb_auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax';
                    // Clear localStorage
                    localStorage.removeItem('pocketbase_auth');
                    localStorage.removeItem('authToken');
                    // Redirect to login
                    window.location.href = '/login';
                }}
            >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
                Sign Out
            </button>
        </div>
    </div>
</aside> 