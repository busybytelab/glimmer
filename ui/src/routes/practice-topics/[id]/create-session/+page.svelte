<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import type { PracticeTopic, Learner } from '$lib/types';
    import pb from '$lib/pocketbase';
    import LearnersList from '../../../../components/learners/LearnersList.svelte';
    import { practiceService } from '$lib/services/practice';
    import ActionToolbar from '../../../../components/common/ActionToolbar.svelte';
    import Breadcrumbs from '../../../../components/common/Breadcrumbs.svelte';

    // Define the breadcrumb item type
    type BreadcrumbItem = {
        label: string;
        href?: string;
        icon?: string;
    };

    let topic: PracticeTopic | null = null;
    let learners: Learner[] = [];
    let loading = true;
    let loadingLearners = true;
    let creatingSession = false;
    let selectedLearner: Learner | null = null;
    let error: string | null = null;
    let topicId: string | null = null;
    let breadcrumbItems: BreadcrumbItem[] = [];

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
                sort: 'nickname',
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
        // Navigate back to the topic page
        goto(`/practice-topics/${topicId}`);
    }

    async function createPracticeSession(learner: Learner) {
        if (!topic) return;
        try {
            creatingSession = true;
            error = null;
            selectedLearner = learner;
            
            // Create the session using the practice service
            const newSession = await practiceService.createSession({
                learnerId: learner.id,
                practiceTopicId: topic.id
            });
            
            // Redirect to practice session page
            goto(`/practice-sessions/${newSession.id}`);
        } catch (err) {
            console.error('Failed to create practice session:', err);
            error = err instanceof Error ? err.message : 'An unexpected error occurred';
        } finally {
            creatingSession = false;
            selectedLearner = null;
        }
    }

    function updateBreadcrumbs() {
        if (!topic) return;
        
        breadcrumbItems = [
            {
                label: 'Topics',
                href: '/practice-topics',
                icon: 'topic'
            },
            {
                label: topic.name,
                href: `/practice-topics/${topic.id}`,
                icon: 'topic'
            },
            {
                label: 'Create Session',
                icon: 'create'
            }
        ];
    }

    // Back action for the toolbar
    $: backAction = {
        id: 'back',
        label: 'Back',
        icon: 'back',
        variant: 'secondary' as const,
        onClick: goBack
    };
</script>

<div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
        <div>
            <Breadcrumbs items={breadcrumbItems} />
        </div>
        <ActionToolbar actions={[backAction]} />
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary"></div>
        </div>
    {:else if error}
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
            <strong class="font-bold">Error!</strong>
            <span class="block sm:inline"> {error}</span>
        </div>
    {:else if topic}
        <div class="bg-white shadow-md rounded-lg p-6 mb-6">
            <h2 class="text-xl font-semibold text-gray-900 mb-2">{topic.name}</h2>
            
            {#if topic.description}
                <p class="text-gray-600 mb-4">{topic.description}</p>
            {/if}
            
            <div class="flex flex-wrap gap-2 mb-4">
                <span class="bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded">
                    {topic.subject}
                </span>
                
                {#if topic.target_age_range}
                    <span class="bg-green-100 text-green-800 text-xs font-medium px-2.5 py-0.5 rounded">
                        Age: {topic.target_age_range}
                    </span>
                {/if}
                
                {#if topic.target_grade_level}
                    <span class="bg-purple-100 text-purple-800 text-xs font-medium px-2.5 py-0.5 rounded">
                        Grade: {topic.target_grade_level}
                    </span>
                {/if}

                {#if topic.difficulty_level}
                    <span class="bg-orange-100 text-orange-800 text-xs font-medium px-2.5 py-0.5 rounded">
                        Difficulty: {topic.difficulty_level}
                    </span>
                {/if}
            </div>
        </div>
        
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Select a Learner</h2>
        
        <LearnersList
            learners={learners}
            loading={loadingLearners}
            emptyMessage="No learners found. Please add learners to your account first."
            onClick={createPracticeSession}
            cardActions={[
                {
                    label: 'Create Practice',
                    color: 'primary',
                    onClick: createPracticeSession
                }
            ]}
        />

        {#if creatingSession}
            <div class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center z-50">
                <div class="bg-white rounded-lg p-8 max-w-md w-full mx-4">
                    <div class="flex flex-col items-center">
                        <div class="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-primary mb-4"></div>
                        <h3 class="text-xl font-semibold text-gray-900 mb-2">Creating Practice Session</h3>
                        <p class="text-gray-600 text-center mb-4">
                            We're generating personalized practice items for {selectedLearner?.nickname || 'the learner'}.
                        </p>
                        <p class="text-gray-600 text-center mb-4">
                            This may take a few minutes. Please don't refresh the page.
                        </p>
                        <div class="w-full bg-gray-200 rounded-full h-2.5">
                            <div class="bg-primary h-2.5 rounded-full animate-pulse"></div>
                        </div>
                    </div>
                </div>
            </div>
        {/if}
    {/if}
</div> 