<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import type { PracticeTopic } from '$lib/types';
    import pb from '$lib/pocketbase';
    import ActionToolbar from '../../../components/common/ActionToolbar.svelte';
    import Breadcrumbs from '../../../components/common/Breadcrumbs.svelte';
    import LoadingSpinner from '../../../components/common/LoadingSpinner.svelte';
    import ErrorAlert from '../../../components/common/ErrorAlert.svelte';
    import { rolesService } from '$lib/services/roles';

    // Define the breadcrumb item type
    type BreadcrumbItem = {
        label: string;
        href?: string;
        icon?: string;
    };

    let topic: PracticeTopic | null = null;
    let pastPractices: any[] = [];
    let loading = true;
    let error: string | null = null;
    let topicId: string | null = null;
    let breadcrumbItems: BreadcrumbItem[] = [];
    let isInstructor = false;

    onMount(async () => {
        try {
            // Get the topic ID from the page store
            topicId = $page.params.id;
            
            // Load data if we have an ID
            if (topicId) {
                await loadTopic(topicId);
                await loadPastPractices(topicId);
                isInstructor = await rolesService.isInstructor();
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

    async function loadPastPractices(id: string) {
        try {
            const result = await pb.collection('practice_sessions').getList(1, 10, {
                filter: `practice_topic="${id}"`,
                sort: '-created',
                expand: 'learner,practice_topic'
            });
            pastPractices = result.items;
        } catch (err) {
            console.error('Failed to load past practices:', err);
        }
    }

    function goBack() {
        goto('/practice-topics');
    }

    function editTopic() {
        if (!topic) return;
        goto(`/practice-topics/${topic.id}/edit`);
    }

    async function startPractice() {
        if (!topic) return;

        try {
            const authData = pb.authStore.model;
            if (!authData) {
                console.error('User not authenticated');
                error = 'You must be logged in to start a practice';
                return;
            }

            if (!isInstructor) {
                error = 'Only instructors can create practice sessions';
                return;
            }

            goto(`/practice-topics/${topic.id}/create-session`);
        } catch (err) {
            console.error('Failed to create practice session:', err);
            error = 'Failed to start new practice: ' + (err instanceof Error ? err.message : String(err));
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
                icon: 'topic'
            }
        ];
    }

    // Actions for the toolbar
    $: topicActions = [
        {
            id: 'back',
            label: 'Back',
            icon: 'back',
            variant: 'secondary' as const,
            onClick: goBack
        },
        ...(isInstructor ? [{
            id: 'edit',
            label: 'Edit',
            icon: 'edit',
            variant: 'primary' as const,
            onClick: editTopic
        }] : [])
    ];
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
    <div class="flex justify-between items-center mb-6">
        <div>
            <Breadcrumbs items={breadcrumbItems} />
        </div>
        <ActionToolbar actions={topicActions} />
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <LoadingSpinner size="md" color="primary" />
        </div>
    {:else if error}
        <ErrorAlert message={error} />
    {:else if topic}
        <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6 mb-6">
            <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">{topic.name}</h2>
            
            {#if topic.description}
                <p class="text-gray-600 dark:text-gray-300 mb-4">{topic.description}</p>
            {/if}
            
            <div class="flex flex-wrap gap-2 mb-4">
                <span class="bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-300 text-xs font-medium px-2.5 py-0.5 rounded">
                    {topic.subject}
                </span>
                
                {#if topic.target_age_range}
                    <span class="bg-green-100 dark:bg-green-900/40 text-green-800 dark:text-green-300 text-xs font-medium px-2.5 py-0.5 rounded">
                        Age: {topic.target_age_range}
                    </span>
                {/if}
                
                {#if topic.target_grade_level}
                    <span class="bg-purple-100 dark:bg-purple-900/40 text-purple-800 dark:text-purple-300 text-xs font-medium px-2.5 py-0.5 rounded">
                        Grade: {topic.target_grade_level}
                    </span>
                {/if}
            </div>
            
            {#if topic.tags && Array.isArray(topic.tags) && topic.tags.length > 0}
                <div class="mb-4">
                    <div class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">Tags:</div>
                    <div class="flex flex-wrap gap-1">
                        {#each topic.tags as tag}
                            <span class="bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300 text-xs font-medium px-2 py-0.5 rounded">
                                {tag}
                            </span>
                        {/each}
                    </div>
                </div>
            {/if}
                
            <div class="mb-4">
                <div class="grid grid-cols-2 gap-4">
                    <div class="bg-gray-50 dark:bg-gray-700 p-3 rounded">
                        <div class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">Difficulty Level:</div>
                        <div class="dark:text-gray-200">{topic.difficulty_level ?? 'Not specified'}</div>
                    </div>
                    
                    <div class="bg-gray-50 dark:bg-gray-700 p-3 rounded">
                        <div class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">Target Ages:</div>
                        <div class="dark:text-gray-200">{topic.target_age_range || 'Not specified'}</div>
                    </div>
                </div>
            </div>
            
            {#if isInstructor}
                <div class="mt-8 flex justify-center">
                    <button
                        type="button"
                        class="inline-flex items-center justify-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 dark:bg-indigo-700 dark:hover:bg-indigo-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-offset-gray-800 transition-all duration-200"
                        on:click={startPractice}
                    >
                        <svg class="h-6 w-6 mr-2" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
                        </svg>
                        Create Practice Session
                    </button>
                </div>
            {/if}
        </div>
            
        <!-- Past Practice Sessions Section -->
        {#if pastPractices.length > 0}
            <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Past Practice Sessions</h3>
                
                <div class="overflow-x-auto">
                    <table class="min-w-full bg-white dark:bg-gray-800">
                        <thead>
                            <tr class="border-b border-gray-200 dark:border-gray-700">
                                <th class="px-4 py-2 text-left text-sm font-semibold text-gray-700 dark:text-gray-300">Session Name</th>
                                <th class="px-4 py-2 text-left text-sm font-semibold text-gray-700 dark:text-gray-300">Learner</th>
                                <th class="px-4 py-2 text-left text-sm font-semibold text-gray-700 dark:text-gray-300">Date</th>
                                <th class="px-4 py-2 text-left text-sm font-semibold text-gray-700 dark:text-gray-300">Status</th>
                                <th class="px-4 py-2 text-left text-sm font-semibold text-gray-700 dark:text-gray-300">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each pastPractices as session}
                                <tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700">
                                    <td class="px-4 py-2 text-sm text-gray-700 dark:text-gray-300">{session.name || 'Practice Session'}</td>
                                    <td class="px-4 py-2 text-sm text-gray-700 dark:text-gray-300">
                                        {session.expand?.learner?.nickname || 'Unknown Learner'}
                                    </td>
                                    <td class="px-4 py-2 text-sm text-gray-700 dark:text-gray-300">
                                        {new Date(session.created).toLocaleDateString()}
                                    </td>
                                    <td class="px-4 py-2 text-sm">
                                        <span class={`px-2 py-1 rounded-full text-xs font-medium 
                                            ${session.status === 'Completed' ? 'bg-green-100 dark:bg-green-900/40 text-green-800 dark:text-green-300' : 
                                              session.status === 'InProgress' ? 'bg-yellow-100 dark:bg-yellow-900/40 text-yellow-800 dark:text-yellow-300' : 
                                              'bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-300'}`}>
                                            {session.status}
                                        </span>
                                    </td>
                                    <td class="px-4 py-2 text-sm text-gray-700 dark:text-gray-300">
                                        <a 
                                            href={`/practice-sessions/${session.id}`} 
                                            class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 hover:underline"
                                        >
                                            View
                                        </a>
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
            </div>
        {/if}
    {/if}
</div> 