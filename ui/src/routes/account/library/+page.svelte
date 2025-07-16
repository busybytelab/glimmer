<script lang="ts">
    import { onMount } from 'svelte';
    import type { PracticeTopicLibrary, PracticeSessionLibrary, BreadcrumbItem, IconType, Learner } from '$lib/types';
    import { libraryService } from '$lib/services/library';
    import { learnersService } from '$lib/services/learners';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    import Breadcrumbs from '$components/common/Breadcrumbs.svelte';
    import ActionToolbar from '$components/common/ActionToolbar.svelte';
    import SelectField from '$components/common/SelectField.svelte';
    import PracticeSessionLibraryCard from '$components/library/PracticeSessionLibraryCard.svelte';
    import { toast } from '$lib/stores/toast';
    import { goto } from '$app/navigation';

    // State management
    let allSessions: PracticeSessionLibrary[] = [];
    let filteredSessions: PracticeSessionLibrary[] = [];
    let categories: string[] = [];
    let topics: PracticeTopicLibrary[] = [];
    let learners: Learner[] = [];
    let selectedLearnerId: string = '';
    let selectedCategory: string = '';
    let selectedTopicId: string = '';
    let importingSessionId: string | null = null;
    
    // UI state
    let loading = true;
    let error: string | null = null;
    let loadingFilters = false;

    // Get the selected learner and their grade level
    $: selectedLearner = learners.find(l => l.id === selectedLearnerId);
    $: selectedGradeLevel = selectedLearner?.grade_level;
    $: selectedTopic = topics.find(t => t.id === selectedTopicId);

    // Handle learner selection changes - filter content by grade level
    $: if (selectedLearnerId && selectedGradeLevel) {
        handleLearnerChange();
    } else if (selectedLearnerId === '' && learners.length > 1) {
        // Handle "All children" option - reset to show all content
        handleAllChildrenChange();
    }

    // Handle category change - load topics for the category
    $: if (selectedCategory) {
        handleCategoryChange();
    }

    // Handle topic change - filter sessions
    $: if (selectedTopicId) {
        handleTopicChange();
    } else if (!selectedTopicId && selectedCategory) {
        // If no topic selected but category is selected, show all sessions for that category
        handleCategoryOnlyFilter();
    } else if (!selectedTopicId && !selectedCategory) {
        // If no filters, show all sessions
        filteredSessions = allSessions;
    }

    async function handleLearnerChange() {
        try {
            loading = true;
            error = null;

            // Reset filters
            selectedCategory = '';
            selectedTopicId = '';

            // Reload categories and sessions with grade level filter
            const [categoriesData, allSessionsData] = await Promise.all([
                libraryService.getCategoriesLibrary(selectedGradeLevel),
                libraryService.getSessionsLibrary(undefined, selectedGradeLevel)
            ]);

            categories = categoriesData;
            allSessions = allSessionsData;
            filteredSessions = allSessions;
            topics = []; // Reset topics
        } catch (err) {
            console.error('Failed to update content for learner:', err);
            error = err instanceof Error ? err.message : 'Failed to update content for learner';
        } finally {
            loading = false;
        }
    }

    async function handleAllChildrenChange() {
        try {
            loading = true;
            error = null;

            // Reset filters
            selectedCategory = '';
            selectedTopicId = '';

            // Reload categories and sessions without grade level filter
            const [categoriesData, allSessionsData] = await Promise.all([
                libraryService.getCategoriesLibrary(), // No grade level filter
                libraryService.getSessionsLibrary(undefined) // No grade level filter
            ]);

            categories = categoriesData;
            allSessions = allSessionsData;
            filteredSessions = allSessions;
            topics = []; // Reset topics
        } catch (err) {
            console.error('Failed to update content for all children:', err);
            error = err instanceof Error ? err.message : 'Failed to update content for all children';
        } finally {
            loading = false;
        }
    }

    async function handleCategoryChange() {
        if (!selectedCategory) return;
        
        try {
            loadingFilters = true;
            selectedTopicId = ''; // Reset topic selection

            // Load topics for the selected category
            topics = await libraryService.getTopicsLibrary(selectedGradeLevel, selectedCategory);
            
            // Filter sessions by category (through topics)
            const categoryTopicIds = topics.map(t => t.id);
            filteredSessions = allSessions.filter(session => 
                categoryTopicIds.includes(session.practice_topic_library)
            );
        } catch (err) {
            console.error('Failed to load topics for category:', err);
            error = err instanceof Error ? err.message : 'Failed to load topics for category';
        } finally {
            loadingFilters = false;
        }
    }

    async function handleTopicChange() {
        if (!selectedTopicId) return;
        
        try {
            loadingFilters = true;
            // Filter sessions by selected topic
            filteredSessions = await libraryService.getSessionsLibrary(selectedTopicId, selectedGradeLevel);
        } catch (err) {
            console.error('Failed to filter sessions by topic:', err);
            error = err instanceof Error ? err.message : 'Failed to filter sessions by topic';
        } finally {
            loadingFilters = false;
        }
    }

    async function handleCategoryOnlyFilter() {
        if (!selectedCategory) return;
        
        try {
            loadingFilters = true;
            // Filter sessions by category (through topics)
            const categoryTopicIds = topics.map(t => t.id);
            filteredSessions = allSessions.filter(session => 
                categoryTopicIds.includes(session.practice_topic_library)
            );
        } catch (err) {
            console.error('Failed to filter sessions by category:', err);
            error = err instanceof Error ? err.message : 'Failed to filter sessions by category';
        } finally {
            loadingFilters = false;
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

            // Load learners and initial content in parallel
            const [learnersData, categoriesData, allSessionsData] = await Promise.all([
                learnersService.getLearners(),
                libraryService.getCategoriesLibrary(), // Load all categories initially
                libraryService.getSessionsLibrary(undefined) // Load all sessions initially
            ]);

            learners = learnersData;
            categories = categoriesData;
            allSessions = allSessionsData;
            filteredSessions = allSessions;

            // Automatically select the learner if there's only one
            if (learners.length === 1) {
                selectedLearnerId = learners[0].id;
            }
        } catch (err) {
            console.error('Failed to load library data:', err);
            error = err instanceof Error ? err.message : 'Failed to load library data';
        } finally {
            loading = false;
        }
    }

    function clearCategoryFilter() {
        selectedCategory = '';
        selectedTopicId = '';
        topics = [];
        filteredSessions = allSessions;
    }

    function clearTopicFilter() {
        selectedTopicId = '';
        if (selectedCategory) {
            handleCategoryOnlyFilter();
        } else {
            filteredSessions = allSessions;
        }
    }

    function clearAllFilters() {
        selectedLearnerId = learners.length === 1 ? learners[0].id : '';
        selectedCategory = '';
        selectedTopicId = '';
        topics = [];
        if (selectedLearnerId === '') {
            handleAllChildrenChange();
        } else {
            filteredSessions = allSessions;
        }
    }

    async function handleImportSession(session: PracticeSessionLibrary) {
        if (!selectedLearnerId) {
            toast.error('Please select a child');
            return;
        }

        try {
            importingSessionId = session.id;
            const importedSession = await libraryService.importSessionFromLibrary(session, selectedLearnerId);
            
            toast.success(`Got "${session.name}" for ${selectedLearner?.nickname}!`);
            
            // Navigate to the imported session overview
            goto(`/account/practice-sessions/${importedSession.id}/overview`);
        } catch (err) {
            console.error('Failed to import session:', err);
            const errorMessage = err instanceof Error ? err.message : 'Failed to get activity';
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
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">
                Find Activities for Your Child
            </h1>
            <p class="text-gray-600 dark:text-gray-400">
                Discover engaging practice activities from our community library
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
        <!-- Filtering Section -->
        <div class="mb-6 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
            <div class="flex items-center gap-4 mb-4">
                <h3 class="text-lg font-medium text-gray-900 dark:text-white">Filter Activities</h3>
                {#if loadingFilters}
                    <LoadingSpinner size="sm" color="primary" />
                {/if}
            </div>

            <!-- Filter Dropdowns -->
            <div class="grid grid-cols-1 {learners.length > 1 ? 'md:grid-cols-2 lg:grid-cols-3' : 'md:grid-cols-2'} gap-4 mb-4">
                <!-- Child Filter (Only for multiple children) -->
                {#if learners.length > 1}
                    <SelectField
                        id="child-filter"
                        label="Select Child"
                        bind:value={selectedLearnerId}
                        cols="col-span-1"
                    >
                        <option value="">All children</option>
                        {#each learners as learner}
                            <option value={learner.id}>
                                {learner.nickname}
                                {#if learner.grade_level}
                                    (Grade {learner.grade_level})
                                {/if}
                            </option>
                        {/each}
                    </SelectField>
                {/if}
                <!-- Category Filter -->
                <SelectField
                    id="category-filter"
                    label="Subject Area"
                    bind:value={selectedCategory}
                    cols="col-span-1"
                >
                    <option value="">All subjects</option>
                    {#each categories as category}
                        <option value={category}>{category}</option>
                    {/each}
                </SelectField>

                <!-- Topic Filter -->
                <SelectField
                    id="topic-filter"
                    label="Specific Topic"
                    bind:value={selectedTopicId}
                    disabled={!selectedCategory || topics.length === 0}
                    cols="col-span-1"
                >
                    <option value="">All topics in {selectedCategory || 'selected subject'}</option>
                    {#each topics as topic}
                        <option value={topic.id}>{topic.name}</option>
                    {/each}
                </SelectField>
            </div>

            <!-- Active Filters -->
            {#if selectedLearnerId || selectedCategory || selectedTopicId}
                <div class="border-t border-gray-200 dark:border-gray-600 pt-4">
                    <div class="flex items-center gap-2 mb-2">
                        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Active filters:</span>
                        <button
                            on:click={clearAllFilters}
                            class="text-xs text-indigo-600 dark:text-indigo-400 hover:text-indigo-800 dark:hover:text-indigo-300"
                        >
                            Clear all
                        </button>
                    </div>
                    <div class="flex flex-wrap gap-2">
                        {#if selectedLearnerId && selectedLearner}
                            <span class="inline-flex items-center gap-1 px-3 py-1 bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-200 text-sm rounded-full">
                                Child: {selectedLearner.nickname}
                                <button
                                    on:click={() => selectedLearnerId = ''}
                                    class="ml-1 text-purple-600 dark:text-purple-300 hover:text-purple-800 dark:hover:text-purple-100"
                                    aria-label="Remove child filter"
                                >
                                    <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                                        <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                                    </svg>
                                </button>
                            </span>
                        {/if}
                        {#if selectedCategory}
                            <span class="inline-flex items-center gap-1 px-3 py-1 bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-200 text-sm rounded-full">
                                Subject: {selectedCategory}
                                <button
                                    on:click={clearCategoryFilter}
                                    class="ml-1 text-blue-600 dark:text-blue-300 hover:text-blue-800 dark:hover:text-blue-100"
                                    aria-label="Remove subject filter"
                                >
                                    <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                                        <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                                    </svg>
                                </button>
                            </span>
                        {/if}
                        {#if selectedTopicId && selectedTopic}
                            <span class="inline-flex items-center gap-1 px-3 py-1 bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-200 text-sm rounded-full">
                                Topic: {selectedTopic.name}
                                <button
                                    on:click={clearTopicFilter}
                                    class="ml-1 text-green-600 dark:text-green-300 hover:text-green-800 dark:hover:text-green-100"
                                    aria-label="Remove topic filter"
                                >
                                    <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                                        <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                                    </svg>
                                </button>
                            </span>
                        {/if}
                    </div>
                </div>
            {/if}
        </div>

        <!-- Activities Section -->
        <div>
            <div class="flex items-center justify-between mb-4">
                <div>
                    <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
                        Practice Activities
                    </h2>
                    <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                        {#if selectedTopic}
                            Showing activities for "{selectedTopic.name}"
                        {:else if selectedCategory}
                            Showing activities in {selectedCategory}
                        {:else if selectedLearner}
                            Showing activities for {selectedLearner.nickname}
                        {:else}
                            Browse all available activities
                        {/if}
                        ({filteredSessions.length} found)
                    </p>
                </div>
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
                                {#if selectedTopic}
                                    No activities found for "{selectedTopic.name}".
                                {:else if selectedCategory}
                                    No activities found in {selectedCategory}.
                                {:else}
                                    No activities available in the library.
                                {/if}
                                Try adjusting your filters or check back later.
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
                            selectedLearner={selectedLearner}
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

 