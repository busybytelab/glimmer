<script lang="ts">
    import { onMount } from 'svelte';
    import type { PracticeTopicLibrary, PracticeSessionLibrary, BreadcrumbItem, IconType } from '$lib/types';
    import { libraryService } from '$lib/services/library';
    import { topicsService } from '$lib/services/topics';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    import Breadcrumbs from '$components/common/Breadcrumbs.svelte';
    import ActionToolbar from '$components/common/ActionToolbar.svelte';
    import TopicLibraryCard from '$components/library/TopicLibraryCard.svelte';
    import PopularityBadge from '$components/common/PopularityBadge.svelte';
    import { toast } from '$lib/stores/toast';

    // State management
    let topTopics: PracticeTopicLibrary[] = [];
    let allSessions: PracticeSessionLibrary[] = [];
    let filteredSessions: PracticeSessionLibrary[] = [];
    let selectedTopicId: string | null = null;
    let selectedTopic: PracticeTopicLibrary | null = null;
    
    // UI state
    let loading = true;
    let error: string | null = null;
    let loadingSessions = false;
    let showAllTopics = false;
    let allTopics: PracticeTopicLibrary[] = [];

    

    // Breadcrumbs
    const breadcrumbItems: BreadcrumbItem[] = [
        {
            label: 'Account',
            href: '/account/dashboard',
            icon: 'home' as IconType
        },
        {
            label: 'Library',
            icon: 'library' as IconType
        }
    ];

    onMount(async () => {
        await loadInitialData();
    });

    async function loadInitialData() {
        try {
            loading = true;
            error = null;

            // Load top 3 topics and all sessions in parallel
            const [topTopicsData, allSessionsData] = await Promise.all([
                libraryService.getTopTopicsLibrary(3),
                libraryService.getSessionsLibrary()
            ]);

            topTopics = topTopicsData;
            allSessions = allSessionsData;
            filteredSessions = allSessions;
        } catch (err) {
            console.error('Failed to load library data:', err);
            error = err instanceof Error ? err.message : 'Failed to load library data';
        } finally {
            loading = false;
        }
    }

    async function handleTopicSelection(topicId: string | null, topic: PracticeTopicLibrary | null = null) {
        try {
            loadingSessions = true;
            selectedTopicId = topicId;
            selectedTopic = topic;

            if (topicId) {
                // Filter sessions by selected topic
                filteredSessions = await libraryService.getSessionsLibrary(topicId);
            } else {
                // Show all sessions
                filteredSessions = allSessions;
            }
        } catch (err) {
            console.error('Failed to filter sessions:', err);
            error = err instanceof Error ? err.message : 'Failed to filter sessions';
        } finally {
            loadingSessions = false;
        }
    }

    async function handleShowAllTopics() {
        if (showAllTopics) {
            showAllTopics = false;
            return;
        }

        try {
            loading = true;
            error = null;
            
            allTopics = await libraryService.getTopicsLibrary();
            showAllTopics = true;
        } catch (err) {
            console.error('Failed to load all topics:', err);
            error = err instanceof Error ? err.message : 'Failed to load all topics';
        } finally {
            loading = false;
        }
    }

    function clearTopicSelection() {
        handleTopicSelection(null, null);
    }

    async function handleAdd(libraryTopic: PracticeTopicLibrary) {
        try {
            const newTopic = await topicsService.importFromLibrary(libraryTopic);
            
            toast.success(`Successfully added "${newTopic.name}" to your topics!`);
            
            // Optionally redirect to the new topic or show success message
            console.log('Successfully imported topic:', newTopic);
        } catch (err) {
            console.error('Failed to add topic:', err);
            const errorMessage = err instanceof Error ? err.message : 'Failed to add topic';
            toast.error(errorMessage);
        }
    }

    // Actions for the toolbar
    const libraryActions = [
        {
            id: 'refresh',
            label: 'Refresh',
            icon: 'reset' as IconType,
            variant: 'secondary' as const,
            onClick: loadInitialData
        }
    ];
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
    <!-- Header with breadcrumbs and actions -->
    <div class="flex justify-between items-center mb-6">
        <div>
            <Breadcrumbs items={breadcrumbItems} />
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white mt-2">Library</h1>
            <p class="text-gray-600 dark:text-gray-400 mt-1">
                Browse community-shared topics and practice sheets
            </p>
        </div>
        <div class="hidden sm:block">
            <ActionToolbar actions={libraryActions} />
        </div>
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <LoadingSpinner size="lg" color="primary" />
        </div>
    {:else if error}
        <ErrorAlert message={error} />
    {:else}
        <!-- Topics Section -->
        <div class="mb-8">
            <div class="flex items-center justify-between mb-4">
                <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Popular Topics</h2>
                {#if selectedTopicId}
                    <button
                        on:click={clearTopicSelection}
                        class="text-sm text-indigo-600 dark:text-indigo-400 hover:text-indigo-800 dark:hover:text-indigo-300"
                    >
                        Clear filter
                    </button>
                {/if}
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <!-- Top 3 Topics -->
                {#if !showAllTopics}
                    {#each topTopics as topic}
                        <TopicLibraryCard 
                            {topic}
                            isSelected={selectedTopicId === topic.id}
                            onClick={handleTopicSelection}
                            onAdd={handleAdd}
                        />
                    {/each}

                    <!-- Browse All Topics Card -->
                    <div 
                        class="border border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-4 hover:border-gray-400 dark:hover:border-gray-500 transition-colors cursor-pointer bg-gray-50 dark:bg-gray-700/50 flex flex-col items-center justify-center text-center min-h-[140px]"
                        on:click={handleShowAllTopics}
                        role="button"
                        tabindex="0"
                        on:keydown={(e) => {
                            if (e.key === 'Enter' || e.key === ' ') {
                                e.preventDefault();
                                handleShowAllTopics();
                            }
                        }}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-400 dark:text-gray-500 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                        </svg>
                        <span class="text-sm font-medium text-gray-600 dark:text-gray-400">Browse All Topics</span>
                        <span class="text-xs text-gray-500 dark:text-gray-500 mt-1">View complete library</span>
                    </div>
                {:else}
                    <!-- All Topics Grid -->
                    {#each allTopics as topic}
                        <TopicLibraryCard 
                            {topic}
                            isSelected={selectedTopicId === topic.id}
                            onClick={handleTopicSelection}
                            onAdd={handleAdd}
                        />
                    {/each}

                    <!-- Show Less Card -->
                    <div 
                        class="border border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-4 hover:border-gray-400 dark:hover:border-gray-500 transition-colors cursor-pointer bg-gray-50 dark:bg-gray-700/50 flex flex-col items-center justify-center text-center min-h-[140px]"
                        on:click={() => showAllTopics = false}
                        role="button"
                        tabindex="0"
                        on:keydown={(e) => {
                            if (e.key === 'Enter' || e.key === ' ') {
                                e.preventDefault();
                                showAllTopics = false;
                            }
                        }}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-400 dark:text-gray-500 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
                        </svg>
                        <span class="text-sm font-medium text-gray-600 dark:text-gray-400">Show Less</span>
                        <span class="text-xs text-gray-500 dark:text-gray-500 mt-1">Back to top topics</span>
                    </div>
                {/if}
            </div>
        </div>

        <!-- Sessions Section -->
        <div>
            <div class="flex items-center justify-between mb-4">
                <div>
                    <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
                        {selectedTopic ? `Practice sessions for "${selectedTopic.name}"` : 'Practice sessions'}
                    </h2>
                    {#if selectedTopic}
                        <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                            Showing sessions from the selected topic
                        </p>
                    {/if}
                </div>
                {#if loadingSessions}
                    <LoadingSpinner size="sm" color="primary" />
                {/if}
            </div>

            {#if filteredSessions.length === 0}
                <div class="bg-yellow-50 dark:bg-yellow-900/20 border-l-4 border-yellow-400 dark:border-yellow-600 p-4">
                    <div class="flex">
                        <div class="flex-shrink-0">
                            <svg class="h-5 w-5 text-yellow-400 dark:text-yellow-300" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                            </svg>
                        </div>
                        <div class="ml-3">
                            <p class="text-sm text-yellow-700 dark:text-yellow-200">
                                {selectedTopic ? 'No sessions found for this topic.' : 'No sessions available in the library.'}
                            </p>
                        </div>
                    </div>
                </div>
            {:else}
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {#each filteredSessions as session}
                        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow bg-white dark:bg-gray-800">
                            <div class="flex items-start justify-between mb-2">
                                <h3 class="font-medium text-gray-900 dark:text-white text-sm">{session.name}</h3>
                                {#if session.total_usage}
                                    <PopularityBadge count={session.total_usage} size="sm" />
                                {/if}
                            </div>
                            
                            {#if session.description}
                                <p class="text-xs text-gray-600 dark:text-gray-400 mb-3 line-clamp-2">{session.description}</p>
                            {/if}
                            
                            <div class="flex flex-wrap gap-1 mb-3">
                                {#if session.expand?.practice_topic_library?.category}
                                    <span class="text-xs bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-300 px-2 py-1 rounded">
                                        {session.expand.practice_topic_library.category}
                                    </span>
                                {/if}
                                {#if session.target_year}
                                    <span class="text-xs bg-orange-100 dark:bg-orange-900/30 text-orange-800 dark:text-orange-300 px-2 py-1 rounded">
                                        Year {session.target_year}
                                    </span>
                                {/if}
                            </div>

                            <div class="flex justify-end">
                                <button class="text-xs text-indigo-600 dark:text-indigo-400 hover:text-indigo-800 dark:hover:text-indigo-300 font-medium">
                                    Use Template
                                </button>
                            </div>
                        </div>
                    {/each}
                </div>
            {/if}
        </div>
    {/if}

    <!-- Mobile Action Toolbar -->
    <div class="sm:hidden fixed bottom-4 right-4">
        <ActionToolbar actions={libraryActions} />
    </div>
</div>

 