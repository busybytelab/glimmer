<script lang="ts">
    import { onMount } from 'svelte';
    import type { PracticeTopicLibrary, PracticeSessionLibrary, BreadcrumbItem, IconType, Learner } from '$lib/types';
    import { libraryService } from '$lib/services/library';
    import { topicsService } from '$lib/services/topics';
    import { learnersService } from '$lib/services/learners';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    import Breadcrumbs from '$components/common/Breadcrumbs.svelte';
    import ActionToolbar from '$components/common/ActionToolbar.svelte';
    import TopicLibraryCard from '$components/library/TopicLibraryCard.svelte';
    import PracticeSessionLibraryCard from '$components/library/PracticeSessionLibraryCard.svelte';
    import { toast } from '$lib/stores/toast';
    import { goto } from '$app/navigation';
    import LearnersList from '$components/learners/LearnersList.svelte';

    // State management
    let topTopics: PracticeTopicLibrary[] = [];
    let allSessions: PracticeSessionLibrary[] = [];
    let filteredSessions: PracticeSessionLibrary[] = [];
    let selectedTopicId: string | null = null;
    let selectedTopic: PracticeTopicLibrary | null = null;
    let learners: Learner[] = [];
    let selectedLearnerId: string = '';
    let importingSessionId: string | null = null;
    
    // UI state
    let loading = true;
    let error: string | null = null;
    let loadingSessions = false;
    let showAllTopics = false;
    let allTopics: PracticeTopicLibrary[] = [];

    // Get the selected learner and their grade level
    $: selectedLearner = learners.find(l => l.id === selectedLearnerId);
    $: selectedGradeLevel = selectedLearner?.grade_level;

    // Handle learner selection changes
    $: if (selectedLearnerId) {
        handleLearnerChange();
    }

    async function handleLearnerChange() {
        try {
            loading = true;
            error = null;

            // Reset topic selection and view state
            selectedTopicId = null;
            selectedTopic = null;
            showAllTopics = false;
            allTopics = []; // Reset all topics

            // Reload topics and sessions with grade level filter
            const [topTopicsData, allSessionsData] = await Promise.all([
                libraryService.getTopTopicsLibrary(3, selectedGradeLevel),
                libraryService.getSessionsLibrary(undefined, selectedGradeLevel)
            ]);

            topTopics = topTopicsData;
            allSessions = allSessionsData;
            filteredSessions = allSessions;
        } catch (err) {
            console.error('Failed to update content for learner:', err);
            error = err instanceof Error ? err.message : 'Failed to update content for learner';
        } finally {
            loading = false;
        }
    }

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

            // Load learners first to get grade level
            const learnersData = await learnersService.getLearners();
            learners = learnersData;

            // Automatically select the learner if there's only one
            if (learners.length === 1) {
                selectedLearnerId = learners[0].id;
            }

            // Load top 3 topics and all sessions in parallel with grade level filter
            const [topTopicsData, allSessionsData] = await Promise.all([
                libraryService.getTopTopicsLibrary(3, selectedGradeLevel),
                libraryService.getSessionsLibrary(undefined, selectedGradeLevel)
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
                // Filter sessions by selected topic and grade level
                filteredSessions = await libraryService.getSessionsLibrary(topicId, selectedGradeLevel);
            } else {
                // Show all sessions for the grade level
                filteredSessions = await libraryService.getSessionsLibrary(undefined, selectedGradeLevel);
            }
        } catch (err) {
            console.error('Failed to filter sessions:', err);
            error = err instanceof Error ? err.message : 'Failed to filter sessions';
        } finally {
            loadingSessions = false;
        }
    }

    async function handleShowAllTopics() {
        try {
            loading = true;
            error = null;
            showAllTopics = !showAllTopics;

            if (!showAllTopics) {
                // Show top topics
                topTopics = await libraryService.getTopTopicsLibrary(3, selectedGradeLevel);
            } else {
                // Show all topics
                allTopics = await libraryService.getTopicsLibrary(selectedGradeLevel);
            }
        } catch (err) {
            // Revert the toggle if there's an error
            showAllTopics = !showAllTopics;
            console.error('Failed to load topics:', err);
            error = err instanceof Error ? err.message : 'Failed to load topics';
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

    async function handleImportSession(session: PracticeSessionLibrary) {
        if (!selectedLearnerId) {
            toast.error('Please select a child first');
            return;
        }

        try {
            importingSessionId = session.id;
            const importedSession = await libraryService.importSessionFromLibrary(session, selectedLearnerId);
            
            toast.success(`Successfully imported "${session.name}" for ${selectedLearner?.nickname}!`);
            
            // Navigate to the imported session overview
            goto(`/account/practice-sessions/${importedSession.id}/overview`);
        } catch (err) {
            console.error('Failed to import session:', err);
            const errorMessage = err instanceof Error ? err.message : 'Failed to import session';
            toast.error(errorMessage);
        } finally {
            importingSessionId = null;
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
        <!-- Getting Started Guide -->
        {#if !selectedLearnerId}
            <div class="bg-blue-50 dark:bg-blue-900/20 border-l-4 border-blue-400 dark:border-blue-600 p-4 mb-6">
                <div class="flex">
                    <div class="flex-shrink-0">
                        <svg class="h-5 w-5 text-blue-400 dark:text-blue-300" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
                        </svg>
                    </div>
                    <div class="ml-3">
                        <h3 class="text-sm font-medium text-blue-800 dark:text-blue-200">Tips</h3>
                        <div class="mt-2 text-sm text-blue-700 dark:text-blue-300">
                            <p>Follow these steps to use the library:</p>
                            <ol class="list-decimal list-inside mt-2 space-y-1">
                                <li>Select a child to see content tailored to their level</li>
                                <li>Browse topics or use the search to find relevant content</li>
                                <li>Click on a topic to see its practice sessions</li>
                                <li>Click "Add" to import a session for your child</li>
                            </ol>
                        </div>
                    </div>
                </div>
            </div>
        {/if}

        <!-- Child Selection -->
        <div class="mb-8">
            <h2 id="learners-list-label" class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
                Select Child
            </h2>
            <div aria-labelledby="learners-list-label">
                <LearnersList
                    {learners}
                    loading={false}
                    emptyMessage="No children found. Please add a child first."
                    gridCols="grid-cols-1 sm:grid-cols-2 lg:grid-cols-3"
                    showPreferences={false}
                    onClick={(learner) => selectedLearnerId = learner.id}
                    {selectedLearnerId}
                />
            </div>
        </div>

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
                    {#if topTopics.length > 0}
                        {#each topTopics as topic}
                            <TopicLibraryCard 
                                {topic}
                                isSelected={selectedTopicId === topic.id}
                                onClick={handleTopicSelection}
                                onAdd={handleAdd}
                            />
                        {/each}
                    {/if}

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
                    {#if allTopics.length > 0}
                        {#each allTopics as topic}
                            <TopicLibraryCard 
                                {topic}
                                isSelected={selectedTopicId === topic.id}
                                onClick={handleTopicSelection}
                                onAdd={handleAdd}
                            />
                        {/each}
                    {/if}

                    <!-- Show Less Card -->
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
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
                        </svg>
                        <span class="text-sm font-medium text-gray-600 dark:text-gray-400">Show Less</span>
                        <span class="text-xs text-gray-500 dark:text-gray-500 mt-1">Back to top topics</span>
                    </div>
                {/if}

                <!-- No Topics Message -->
                {#if (!showAllTopics && topTopics.length === 0) || (showAllTopics && allTopics.length === 0)}
                    <div class="col-span-full">
                        <div class="bg-yellow-50 dark:bg-yellow-900/20 border-l-4 border-yellow-400 dark:border-yellow-600 p-4">
                            <div class="flex">
                                <div class="flex-shrink-0">
                                    <svg class="h-5 w-5 text-yellow-400 dark:text-yellow-300" viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                                    </svg>
                                </div>
                                <div class="ml-3">
                                    <p class="text-sm text-yellow-700 dark:text-yellow-200">
                                        No topics available for the selected grade level.
                                    </p>
                                </div>
                            </div>
                        </div>
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
                        <PracticeSessionLibraryCard
                            {session}
                            isImporting={importingSessionId === session.id}
                            disabled={!selectedLearnerId}
                            onAdd={handleImportSession}
                        />
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

 