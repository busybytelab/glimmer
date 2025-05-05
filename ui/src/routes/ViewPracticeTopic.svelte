<script lang="ts">
    import { onMount } from 'svelte';
    import pb from '$lib/pocketbase';
    import type { PracticeTopic, PracticeSession } from '$lib/types';
    import AppLayout from '../components/layout/AppLayout.svelte';
    import FormButton from '../components/common/FormButton.svelte';

    let topic: PracticeTopic | null = null;
    let pastPractices: PracticeSession[] = [];
    let loading = true;
    let error: string | null = null;
    let sidebarOpen = true;
    let topicId: string | null = null;

    onMount(() => {
        // Extract the topic ID from the URL path
        const path = window.location.pathname;
        console.log('Current path:', path);
        
        if (path.includes('/practice-topic/')) {
            // Extract ID from path
            topicId = path.split('/practice-topic/')[1];
            console.log('Extracted topicId:', topicId);
            
            // Load data if we have an ID
            if (topicId) {
                loadTopic(topicId);
                loadPastPractices(topicId);
            } else {
                error = 'Invalid topic ID';
                loading = false;
            }
        } else {
            console.error('Invalid URL format, expected /practice-topic/ID');
            error = 'Invalid URL format';
            loading = false;
        }
    });

    async function loadTopic(id: string) {
        try {
            console.log('Loading topic with ID:', id);
            loading = true;
            error = null;
            
            const result = await pb.collection('practice_topics').getOne(id);
            console.log('Topic data loaded:', result);
            
            // Format tags if needed
            if (result.tags && !Array.isArray(result.tags)) {
                try {
                    if (typeof result.tags === 'string' && (result.tags as string).trim().startsWith('[')) {
                        result.tags = JSON.parse(result.tags as string);
                    } else if (typeof result.tags === 'string') {
                        result.tags = (result.tags as string).split(',').map((tag: string) => tag.trim()).filter(Boolean);
                    }
                } catch (e) {
                    console.error('Error parsing tags:', e);
                    result.tags = [];
                }
            } else if (!result.tags) {
                result.tags = [];
            }
            
            topic = result as unknown as PracticeTopic;
            console.log('Topic assigned:', topic);
        } catch (err) {
            console.error('Failed to load topic:', err);
            error = 'Failed to load practice topic';
        } finally {
            loading = false;
        }
    }

    async function loadPastPractices(topicId: string) {
        try {
            console.log('Loading past practices for topic:', topicId);
            // Update field name from "topic" to "practice_topic"
            const result = await pb.collection('practice_sessions').getList(1, 10, {
                filter: `practice_topic="${topicId}"`,
                sort: '-created',
                expand: 'learner,practice_topic'
            });
            
            console.log('Past practices loaded:', result);
            pastPractices = result.items as unknown as PracticeSession[];
        } catch (err) {
            console.error('Failed to load past practices:', err);
            // Don't set error here as it would override topic loading error
        }
    }

    function goBack() {
        // Navigate using path-based routing
        (window as any).navigate('/practice-topics');
    }

    function goToEdit() {
        if (!topic) return;
        
        // Navigate to practice topics page with edit query parameter
        const url = new URL(window.location.origin + '/practice-topics');
        url.searchParams.set('edit', topic.id);
        
        console.log('Navigating to edit:', url.toString());
        (window as any).navigate(url.pathname + url.search);
    }

    async function startNewPractice() {
        if (!topic) return;
        try {
            // Get current user info
            const authData = pb.authStore.model;
            if (!authData) {
                console.error('User not authenticated');
                error = 'You must be logged in to start a practice';
                return;
            }
            
            console.log('Auth data:', authData);
            
            // First, determine if the current user is an instructor or a learner
            try {
                // Check if user is an instructor
                const instructorRecord = await pb.collection('instructors').getFirstListItem(`user="${authData.id}"`);
                if (instructorRecord) {
                    console.log('User is an instructor', instructorRecord);
                    
                    // If user is an instructor, redirect to the create practice session page
                    // where they can select a learner
                    (window as any).navigate(`/create-practice/${topic.id}`);
                    return;
                }
            } catch (err) {
                // Not an instructor, continue checking if learner
                console.log('User is not an instructor');
            }
            
            // If we get here, check if user is a learner
            try {
                const learnerRecord = await pb.collection('learners').getFirstListItem(`user="${authData.id}"`);
                console.log('Found learner record:', learnerRecord);
                
                // Create the session with all required fields
                const sessionData = {
                    practice_topic: topic.id,
                    learner: learnerRecord.id,
                    practice_items: "[]", // Send as string since PocketBase expects JSON
                    assigned_at: new Date().toISOString(),
                    status: 'InProgress',
                    name: `Practice: ${topic.name}`
                };
                
                console.log('Creating session with data:', sessionData);
                
                const newSession = await pb.collection('practice_sessions').create(sessionData);
                
                console.log('Created session:', newSession);
                
                // Redirect to practice session page using path-based routing
                (window as any).navigate(`/practice-session/${newSession.id}`);
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
</script>

<AppLayout bind:sidebarOpen>
    <div class="container mx-auto px-4 py-8">
        <div class="flex justify-between items-center mb-6">
            <h1 class="text-2xl font-bold text-gray-900">Practice Topic</h1>
            <div class="flex space-x-2">
                <FormButton
                    type="button"
                    variant="secondary"
                    on:click={goBack}
                >
                    Back to Topics
                </FormButton>
                <FormButton
                    type="button"
                    variant="primary"
                    on:click={goToEdit}
                >
                    Edit Topic
                </FormButton>
            </div>
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
                        <h3 class="text-sm font-medium text-gray-700 mb-2">Tags:</h3>
                        <div class="flex flex-wrap gap-1">
                            {#each topic.tags as tag}
                                <span class="bg-gray-100 text-gray-800 text-xs font-medium px-2 py-0.5 rounded">
                                    {tag}
                                </span>
                            {/each}
                        </div>
                    </div>
                {/if}
                
                {#if topic.learning_goals && Array.isArray(topic.learning_goals) && topic.learning_goals.length > 0}
                    <div class="mb-4">
                        <h3 class="text-sm font-medium text-gray-700 mb-2">Learning Goals:</h3>
                        <ul class="list-disc list-inside">
                            {#each topic.learning_goals as goal}
                                <li class="text-gray-600">{goal}</li>
                            {/each}
                        </ul>
                    </div>
                {/if}
                
                <div class="mt-6">
                    <FormButton
                        type="button"
                        variant="primary"
                        on:click={startNewPractice}
                    >
                        Start New Practice
                    </FormButton>
                </div>
            </div>
            
            {#if pastPractices.length > 0}
                <div class="bg-white shadow-md rounded-lg p-6">
                    <h2 class="text-xl font-semibold text-gray-900 mb-4">Past Practice Sessions</h2>
                    
                    <div class="overflow-x-auto">
                        <table class="min-w-full divide-y divide-gray-200">
                            <thead class="bg-gray-50">
                                <tr>
                                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Score</th>
                                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                                </tr>
                            </thead>
                            <tbody class="bg-white divide-y divide-gray-200">
                                {#each pastPractices as practice}
                                    <tr>
                                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                            {new Date(practice.assigned_at).toLocaleDateString()}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                            {practice.name || 'Unnamed practice'}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap text-sm">
                                            <span class={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full 
                                                ${practice.status === 'Completed' ? 'bg-green-100 text-green-800' : 
                                                practice.status === 'InProgress' ? 'bg-blue-100 text-blue-800' : 
                                                'bg-gray-100 text-gray-800'}`}>
                                                {practice.status}
                                            </span>
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                            {practice.score ? `${practice.score}%` : '-'}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                            <a 
                                                href={`#practice-session/${practice.id}`}
                                                class="text-indigo-600 hover:text-indigo-900"
                                            >
                                                {practice.status === 'InProgress' ? 'Continue' : 'View'}
                                            </a>
                                        </td>
                                    </tr>
                                {/each}
                            </tbody>
                        </table>
                    </div>
                </div>
            {:else}
                <div class="bg-gray-50 border border-gray-200 p-4 rounded-md">
                    <p class="text-gray-600">No practice sessions yet. Start your first one!</p>
                </div>
            {/if}
        {/if}
    </div>
</AppLayout> 