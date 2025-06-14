<script lang="ts">
    import ProfileCard from '$components/account/ProfileCard.svelte';
    import pb from '$lib/pocketbase';
    
    
    // /**
    //  * Handle LLM settings save event from the LLMSettingsCard component
    //  * @param event Custom event containing the saved settings
    //  */
    // function handleLLMSettingsSave(event: CustomEvent<{
    //     openaiApiKey: string;
    //     ollamaEndpoint: string;
    //     defaultModel: string;
    // }>) {
    //     // You could handle global app state updates here if needed
    //     console.log('LLM settings saved:', event.detail);
    // }
    
    /**
     * Log out the current user and redirect to the login page
     */
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
</script>

<div class="py-6">
    <div class="container max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between mb-6">
            <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">Account</h1>
            
            <!-- Sign out button -->
            <button 
                class="flex items-center px-3 py-2 rounded-md text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/30 border border-red-200 dark:border-red-900/30 transition-colors"
                aria-label="Sign out"
                title="Sign out"
                on:click={logout}
            >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
                Sign Out
            </button>
        </div>
        
        <div class="grid grid-cols-1 gap-6">
            <!-- Profile Card -->
            <ProfileCard />
            
            <!-- LLM Settings disabled since backend is not ready yet -->
            <!-- <LLMSettingsCard 
                isLoading={isLoading.llmSettings}
                on:save={handleLLMSettingsSave}
            /> -->
        </div>
    </div>
</div> 