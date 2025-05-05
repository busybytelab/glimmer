<script lang="ts">
    import { onMount } from 'svelte';
    import pb from '$lib/pocketbase';
    import type { PracticeTopic, Learner } from '$lib/types';
    import AppLayout from '../components/layout/AppLayout.svelte';
    import FormButton from '../components/common/FormButton.svelte';
    import LearnersList from '../components/learners/LearnersList.svelte';

    let topic: PracticeTopic | null = null;
    let learners: Learner[] = [];
    let loading = true;
    let loadingLearners = true;
    let error: string | null = null;
    let sidebarOpen = true;
    let topicId: string | null = null;

    onMount(() => {
        // Extract the topic ID from the URL path
        const path = window.location.pathname;
        console.log('Current path:', path);
        
        if (path.includes('/create-practice/')) {
            // Extract ID from path
            topicId = path.split('/create-practice/')[1];
            console.log('Extracted topicId:', topicId);
            
            // Load data if we have an ID
            if (topicId) {
                loadTopic(topicId);
                loadLearners();
            } else {
                error = 'Invalid topic ID';
                loading = false;
            }
        } else {
            console.error('Invalid URL format, expected /create-practice/ID');
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

    async function loadLearners() {
        try {
            loadingLearners = true;
            
            // Get current user info
            const authData = pb.authStore.model;
            if (!authData) {
                console.error('User not authenticated');
                error = 'You must be logged in to create a practice session';
                return;
            }
            
            // Get instructor record
            const instructorRecord = await pb.collection('instructors').getFirstListItem(`user="${authData.id}"`);
            
            // Get learners from the same account
            const result = await pb.collection('learners').getList(1, 100, {
                filter: `account="${instructorRecord.account}"`,
                sort: 'nickname',
                expand: 'user'
            });
            
            learners = result.items as unknown as Learner[];
            console.log('Loaded learners:', learners);
        } catch (err) {
            console.error('Failed to load learners:', err);
            error = 'Failed to load learners';
        } finally {
            loadingLearners = false;
        }
    }

    function goBack() {
        // Navigate using path-based routing
        if (topicId) {
            (window as any).navigate(`/practice-topic/${topicId}`);
        } else {
            (window as any).navigate('/practice-topics');
        }
    }

    async function createPracticeSession(learner: Learner) {
        if (!topic) return;
        try {
            console.log('Creating practice session for topic:', topic.id, 'and learner:', learner.id);
            
            // Get current user info for instructor ID
            const authData = pb.authStore.model;
            if (!authData) {
                console.error('User not authenticated');
                error = 'You must be logged in to create a practice session';
                return;
            }
            
            // Get instructor record
            const instructorRecord = await pb.collection('instructors').getFirstListItem(`user="${authData.id}"`);
            
            // Create the session with all required fields
            const sessionData = {
                practice_topic: topic.id,
                learner: learner.id,
                instructor: instructorRecord.id,
                practice_items: "[]", // Send as string since PocketBase expects JSON
                assigned_at: new Date().toISOString(),
                status: 'InProgress',
                name: `${learner.nickname} - ${topic.name}`
            };
            
            console.log('Creating session with data:', sessionData);
            
            const newSession = await pb.collection('practice_sessions').create(sessionData);
            
            console.log('Created session:', newSession);
            
            // Redirect to practice session page using path-based routing
            (window as any).navigate(`/practice-session/${newSession.id}`);
        } catch (err) {
            console.error('Failed to create practice session:', err);
            error = 'Failed to create practice session: ' + (err instanceof Error ? err.message : String(err));
        }
    }
</script>

<AppLayout bind:sidebarOpen>
    <div class="container mx-auto px-4 py-8">
        <div class="flex justify-between items-center mb-6">
            <h1 class="text-2xl font-bold text-gray-900">Create New Practice Session</h1>
            <div class="flex space-x-2">
                <FormButton
                    type="button"
                    variant="secondary"
                    on:click={goBack}
                >
                    Back to Topic
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
        {/if}
    </div>
</AppLayout> 