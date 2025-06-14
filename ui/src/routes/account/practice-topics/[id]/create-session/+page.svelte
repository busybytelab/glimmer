<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import type { PracticeTopic, Learner, BreadcrumbItem, IconType } from '$lib/types';
    import pb from '$lib/pocketbase';
    import LearnersList from '$components/learners/LearnersList.svelte';
    import { practiceService } from '$lib/services/practice';
    import ActionToolbar from '$components/common/ActionToolbar.svelte';
    import Breadcrumbs from '$components/common/Breadcrumbs.svelte';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    import PracticeTopicCard from '$components/practice-topics/PracticeTopicCard.svelte';

    import ExpandableTextArea from '$components/common/ExpandableTextArea.svelte';

    // Define the steps for session creation as a string literal type
    type CreationStep = 'select_learner' | 'edit_prompts';

    let topic: PracticeTopic | null = null;
    let learners: Learner[] = [];
    let loading = true;
    let loadingLearners = true;
    let creatingSession = false;
    let selectedLearner: Learner | null = null;
    let error: string | null = null;
    let topicId: string | null = null;
    let breadcrumbItems: BreadcrumbItem[] = [];
    
    // Added variables for prompts
    let systemPrompt: string = '';
    let basePrompt: string = '';
    let currentStep: CreationStep = 'select_learner';

    onMount(async () => {
        try {
            // Get the topic ID from the page store
            topicId = $page.params.id;
            
            // Load data if we have an ID
            if (topicId) {
                await loadTopic(topicId);
                await loadLearners();
                updateBreadcrumbs();
            } else {
                error = 'Invalid topic ID';
                loading = false;
            }
        } catch (err) {
            console.error('Error in onMount:', err);
            error = err instanceof Error ? err.message : 'An unexpected error occurred';
            loading = false;
        }
    });

    async function loadTopic(id: string) {
        try {
            loading = true;
            error = null;
            
            if (!id) {
                throw new Error('Topic ID is required');
            }

            const result = await pb.collection('practice_topics').getOne<PracticeTopic>(id);

            // Parse tags if they're stored as a string
            if (result.tags) {
                const tagsValue = result.tags as string | string[];
                if (typeof tagsValue === 'string') {
                    try {
                        if (tagsValue.trim().startsWith('[')) {
                            result.tags = JSON.parse(tagsValue);
                        } else {
                            result.tags = tagsValue.split(',').map((tag: string) => tag.trim()).filter(Boolean);
                        }
                    } catch (err) {
                        console.error('Error parsing tags:', err);
                        result.tags = [];
                    }
                }
            } else {
                result.tags = [];
            }

            topic = result;
            
            // Initialize prompts from the topic
            systemPrompt = result.system_prompt || '';
            basePrompt = result.base_prompt || '';
        } catch (err) {
            console.error('Failed to load topic:', err);
            error = 'Failed to load practice topic';
        } finally {
            loading = false;
        }
    }

    async function loadLearners() {
        try {
            loadingLearners = true;

            const result = await pb.collection('learners').getList(1, 50, {
                sort: 'created',
                expand: 'user'
            });

            learners = result.items;
        } catch (err) {
            console.error('Failed to load learners:', err);
            error = 'Failed to load learners';
        } finally {
            loadingLearners = false;
        }
    }

    function goBack() {
        if (currentStep === 'edit_prompts') {
            // Go back to learner selection
            currentStep = 'select_learner';
        } else {
            // Navigate back to the topic page
            goto(`/account/practice-topics/${topicId}`);
        }
    }

    function goToPromptEditing(learner: Learner) {
        selectedLearner = learner;
        currentStep = 'edit_prompts';
    }

    async function createPracticeSession() {
        if (!topic || !selectedLearner) return;
        try {
            creatingSession = true;
            error = null;
            
            // Create the session using the practice service with the edited prompts
            const newSession = await practiceService.createSession({
                learnerId: selectedLearner.id,
                practiceTopicId: topic.id,
                systemPrompt,
                basePrompt
            });
            
            // Redirect to practice session page
            goto(`/account/practice-sessions/${newSession.id}`);
        } catch (err) {
            console.error('Failed to create practice session:', err);
            error = err instanceof Error ? err.message : 'An unexpected error occurred';
            creatingSession = false;
        }
    }

    function updateBreadcrumbs() {
        if (!topic) return;
        
        breadcrumbItems = [
            {
                label: 'Topics',
                href: '/account/practice-topics',
                icon: 'topic' as IconType
            },
            {
                label: topic.name,
                href: `/account/practice-topics/${topic.id}`,
                icon: 'topic' as IconType
            },
            {
                label: 'Create Session',
                icon: 'create' as IconType
            }
        ];
    }

    // Back action for the toolbar
    $: backAction = {
        id: 'back',
        label: 'Back',
        icon: 'back' as IconType,
        variant: 'secondary' as const,
        onClick: goBack
    };

    // Create action for the toolbar
    $: createAction = {
        id: 'create',
        label: 'Create Practice Session',
        icon: 'add' as IconType,
        variant: 'primary' as const,
        onClick: createPracticeSession
    };
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
    <div class="flex justify-between items-center mb-6">
        <div>
            <Breadcrumbs items={breadcrumbItems} />
        </div>
        <ActionToolbar actions={currentStep === 'edit_prompts' ? [backAction, createAction] : [backAction]} />
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <LoadingSpinner size="md" color="primary" />
        </div>
    {:else if error}
        <ErrorAlert message={error} />
    {:else if topic}
        <div class="mb-6">
            <PracticeTopicCard {topic} />
        </div>
        
        {#if currentStep === 'select_learner'}
            <h2 class="text-xl font-semibold text-gray-900 mb-4">Select a Learner</h2>
            
            <LearnersList
                learners={learners}
                loading={loadingLearners}
                emptyMessage="No learners found. Please add learners to your account first."
                onClick={goToPromptEditing}
                cardActions={[
                    {
                        label: 'Select Learner',
                        color: 'primary',
                        onClick: goToPromptEditing
                    }
                ]}
            />
        {:else if currentStep === 'edit_prompts'}
            <h2 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">Review & Edit Prompts</h2>
            
            <div class="mb-6 bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
                <div class="flex items-center">
                    <div class="flex-shrink-0 h-12 w-12 rounded-full bg-primary/20 flex items-center justify-center text-primary">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                        </svg>
                    </div>
                    <div class="ml-4">
                        <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100">Selected Learner</h3>
                        <p class="text-md text-gray-600 dark:text-gray-300">
                            <span class="font-semibold">{selectedLearner?.nickname || 'Unnamed Learner'}</span>
                            {#if selectedLearner?.age}
                                <span class="text-gray-500 dark:text-gray-400 ml-2">• Age: {selectedLearner.age}</span>
                            {/if}
                            {#if selectedLearner?.grade_level}
                                <span class="text-gray-500 dark:text-gray-400 ml-2">• Grade: {selectedLearner.grade_level}</span>
                            {/if}
                        </p>
                    </div>
                </div>
                <p class="mt-3 text-gray-600 dark:text-gray-300">
                    The following prompts will be used to generate personalized practice items for this learner.
                </p>
            </div>
            
            <div class="space-y-6">
                <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
                    <h3 class="font-medium text-lg mb-2 dark:text-gray-200">System Prompt</h3>
                    <p class="text-sm text-gray-500 dark:text-gray-400 mb-3">Instructions that define how the AI should behave when generating practice items.</p>
                    <ExpandableTextArea
                        id="systemPrompt"
                        label=""
                        bind:value={systemPrompt}
                        minRows={5}
                        maxRows={15}
                        language="markdown"
                    />
                </div>
                
                <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
                    <h3 class="font-medium text-lg mb-2 dark:text-gray-200">Base Prompt</h3>
                    <p class="text-sm text-gray-500 dark:text-gray-400 mb-3">The content used to generate practice items for this topic.</p>
                    <ExpandableTextArea
                        id="basePrompt"
                        label=""
                        bind:value={basePrompt}
                        minRows={5}
                        maxRows={15}
                        language="markdown"
                    />
                </div>
                
                <div class="flex justify-end">
                    <button 
                        class="px-4 py-2 bg-primary text-white rounded hover:bg-primary-dark transition-colors"
                        on:click={createPracticeSession}
                    >
                        Create Practice Session
                    </button>
                </div>
            </div>
        {/if}

        {#if creatingSession}
            <div class="fixed inset-0 bg-gray-500 bg-opacity-75 dark:bg-gray-900 dark:bg-opacity-80 flex items-center justify-center z-50">
                <div class="bg-white dark:bg-gray-800 rounded-lg p-8 max-w-md w-full mx-4">
                    <div class="flex flex-col items-center">
                        <LoadingSpinner size="lg" color="primary" />
                        <h3 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">Creating Practice Session</h3>
                        <p class="text-gray-600 dark:text-gray-300 text-center mb-4">
                            We're generating personalized practice items for {selectedLearner?.nickname || 'the learner'}.
                        </p>
                        <p class="text-gray-600 dark:text-gray-300 text-center mb-4">
                            This may take a few minutes. Please don't refresh the page.
                        </p>
                        <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
                            <div class="bg-primary h-2.5 rounded-full animate-pulse"></div>
                        </div>
                    </div>
                </div>
            </div>
        {/if}
    {/if}
</div> 