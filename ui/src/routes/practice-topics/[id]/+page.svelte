<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import type { PracticeTopic } from '$lib/types';
    import pb from '$lib/pocketbase';
    import ActionToolbar from '../../../components/common/ActionToolbar.svelte';
    import Breadcrumbs from '../../../components/common/Breadcrumbs.svelte';

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

    onMount(async () => {
        try {
            // Get the topic ID from the page store
            topicId = $page.params.id;
            
            // Load data if we have an ID
            if (topicId) {
                await loadTopic(topicId);
                await loadPastPractices(topicId);
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
        goto(`/practice-topics/edit/${topic.id}`);
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

            try {
                const instructorRecord = await pb.collection('instructors').getFirstListItem(`user="${authData.id}"`);
                if (instructorRecord) {
                    goto(`/practice-topics/${topic.id}/create-session`);
                    return;
                }
            } catch {
                // Not an instructor, continue as learner
            }

            try {
                const learnerRecord = await pb.collection('learners').getFirstListItem(`user="${authData.id}"`);
                
                // Create the session with all required fields
                const sessionData = {
                    practice_topic: topic.id,
                    learner: learnerRecord.id,
                    practice_items: "[]", // Send as string since PocketBase expects JSON
                    assigned_at: new Date().toISOString(),
                    status: 'InProgress',
                    name: `Practice: ${topic.name}`
                };
                
                const newSession = await pb.collection('practice_sessions').create(sessionData);
                
                // Redirect to practice session page using SvelteKit routing
                goto(`/practice-sessions/${newSession.id}`);
            } catch (err) {
                console.error('Failed to find learner record:', err);
                error = 'Could not find your learner profile. Please contact support.';
                return;
            }
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
        {
            id: 'edit',
            label: 'Edit',
            icon: 'edit',
            variant: 'primary' as const,
            onClick: editTopic
        }
    ];
</script>

<div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
        <div>
            <Breadcrumbs items={breadcrumbItems} />
        </div>
        <ActionToolbar actions={topicActions} />
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
            </div>
            
            {#if topic.tags && Array.isArray(topic.tags) && topic.tags.length > 0}
                <div class="mb-4">
                    <div class="text-sm font-semibold text-gray-700 mb-1">Tags:</div>
                    <div class="flex flex-wrap gap-1">
                        {#each topic.tags as tag}
                            <span class="bg-gray-100 text-gray-800 text-xs font-medium px-2 py-0.5 rounded">
                                {tag}
                            </span>
                        {/each}
                    </div>
                </div>
            {/if}
                
            <div class="mb-4">
                <div class="grid grid-cols-2 gap-4">
                    <div class="bg-gray-50 p-3 rounded">
                        <div class="text-sm font-semibold text-gray-700 mb-1">Difficulty Level:</div>
                        <div>{topic.difficulty_level ?? 'Not specified'}</div>
                    </div>
                    
                    <div class="bg-gray-50 p-3 rounded">
                        <div class="text-sm font-semibold text-gray-700 mb-1">Target Ages:</div>
                        <div>{topic.target_age_range || 'Not specified'}</div>
                    </div>
                </div>
            </div>
            
            <div class="mt-8 flex justify-center">
                <button
                    type="button"
                    class="inline-flex items-center justify-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200"
                    on:click={startPractice}
                >
                    <svg class="h-6 w-6 mr-2" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
                    </svg>
                    Create Practice Session
                </button>
            </div>
        </div>
            
        <!-- Past Practice Sessions Section -->
        {#if pastPractices.length > 0}
            <div class="bg-white shadow-md rounded-lg p-6">
                <h3 class="text-lg font-semibold text-gray-900 mb-4">Past Practice Sessions</h3>
                
                <div class="overflow-x-auto">
                    <table class="min-w-full bg-white">
                        <thead>
                            <tr class="border-b border-gray-200">
                                <th class="px-4 py-2 text-left text-sm font-semibold text-gray-700">Session Name</th>
                                <th class="px-4 py-2 text-left text-sm font-semibold text-gray-700">Learner</th>
                                <th class="px-4 py-2 text-left text-sm font-semibold text-gray-700">Date</th>
                                <th class="px-4 py-2 text-left text-sm font-semibold text-gray-700">Status</th>
                                <th class="px-4 py-2 text-left text-sm font-semibold text-gray-700">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each pastPractices as session}
                                <tr class="border-b border-gray-200 hover:bg-gray-50">
                                    <td class="px-4 py-2 text-sm text-gray-700">{session.name || 'Practice Session'}</td>
                                    <td class="px-4 py-2 text-sm text-gray-700">
                                        {session.expand?.learner?.nickname || 'Unknown Learner'}
                                    </td>
                                    <td class="px-4 py-2 text-sm text-gray-700">
                                        {new Date(session.created).toLocaleDateString()}
                                    </td>
                                    <td class="px-4 py-2 text-sm">
                                        <span class={`px-2 py-1 rounded-full text-xs font-medium 
                                            ${session.status === 'Completed' ? 'bg-green-100 text-green-800' : 
                                              session.status === 'InProgress' ? 'bg-yellow-100 text-yellow-800' : 
                                              'bg-blue-100 text-blue-800'}`}>
                                            {session.status}
                                        </span>
                                    </td>
                                    <td class="px-4 py-2 text-sm text-gray-700">
                                        <a 
                                            href={`/practice-sessions/${session.id}`} 
                                            class="text-blue-600 hover:text-blue-800 hover:underline"
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