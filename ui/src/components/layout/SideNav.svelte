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
    
    function logout() {
        // Clear PocketBase auth store
        pb.authStore.clear();
        // Clear the auth cookie
        document.cookie = 'pb_auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax';
        // Clear localStorage
        localStorage.removeItem('pocketbase_auth');
        localStorage.removeItem('authToken');
        // Redirect to login
        window.location.href = '/login';
    }
    
    // Get user information if available
    const user = pb.authStore.model;
    const userName = user ? user.name || user.username || user.email || 'User' : 'User';
    
    // Application version
    const appVersion = 'v0.0.1';
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
        
        <div class="mt-auto p-3 border-t border-gray-200 dark:border-gray-700 text-sm">
            <!-- Single row with user info and GitHub link -->
            <div class="flex items-center justify-between">
                <!-- User profile info -->
                <a href="/account" class="flex items-center group text-gray-700 dark:text-gray-300 hover:text-indigo-600 dark:hover:text-indigo-400">
                    <div class="bg-gray-200 dark:bg-gray-700 rounded-full p-1.5 mr-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                        </svg>
                    </div>
                    <span class="truncate max-w-[100px] group-hover:text-indigo-600 dark:group-hover:text-indigo-400">{userName}</span>
                </a>
                
                <!-- GitHub and version -->
                <div class="flex items-center text-xs text-gray-400 dark:text-gray-500">
                    <a 
                        href="https://github.com/busybytelab/glimmer" 
                        target="_blank" 
                        rel="noopener noreferrer"
                        class="flex items-center hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
                        title="View Glimmer on GitHub"
                    >
                        <span class="text-xs mr-1">Visit</span>
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="currentColor">
                            <path fill-rule="evenodd" clip-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" />
                        </svg>
                    </a>
                    <span class="mx-1.5">|</span>
                    <span>{appVersion}</span>
                </div>
            </div>
        </div>
    </div>
</aside> 