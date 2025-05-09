<script lang="ts">
    import { onMount } from 'svelte';
    import type { PracticeItem } from '$lib/types';
    import AppLayout from '../components/layout/AppLayout.svelte';
    import FormButton from '../components/common/FormButton.svelte';
    import QuestionFactory from '../components/questions/QuestionFactory.svelte';
    import { sessionService, type SessionWithExpandedData } from '$lib/services/session';

    let session: SessionWithExpandedData | null = null;
    let practiceItems: PracticeItem[] = [];
    let loading = true;
    let error: string | null = null;
    let sidebarOpen = true;
    let sessionId: string | null = null;
    let isInstructor = false;
    let printMode = false;

    onMount(async () => {
        try {
            // Extract the session ID from the URL path
            const path = window.location.pathname;
            
            if (path.includes('/practice-session/')) {
                // Extract ID from path
                sessionId = path.split('/practice-session/')[1];
                
                // Load data if we have an ID
                if (sessionId) {
                    await loadSession(sessionId);
                    await checkUserRole();
                } else {
                    error = 'Invalid session ID';
                    loading = false;
                }
            } else {
                console.error('Invalid URL format, expected /practice-session/ID');
                error = 'Invalid URL format';
                loading = false;
            }
        } catch (err) {
            console.error('Error in onMount:', err);
            error = err instanceof Error ? err.message : 'An unexpected error occurred';
            loading = false;
        }
    });

    async function checkUserRole() {
        try {
            isInstructor = await sessionService.checkUserRole();
        } catch (err) {
            console.error('Failed to check user role:', err);
            isInstructor = false;
        }
    }

    async function loadSession(id: string) {
        try {
            loading = true;
            error = null;
            
            if (!id) {
                throw new Error('Session ID is required');
            }

            session = await sessionService.loadSession(id);
            
            if (!session) {
                throw new Error('Session not found');
            }

            practiceItems = sessionService.parsePracticeItems(session);
        } catch (err) {
            console.error('Failed to load session:', err);
            error = err instanceof Error ? err.message : 'Failed to load practice session';
            session = null;
            practiceItems = [];
        } finally {
            loading = false;
        }
    }

    function goBack() {
        // Navigate back to practice topics
        (window as any).navigate('/practice-topics');
    }

    function handleAnswerChange(index: number, answer: string) {
        if (practiceItems[index]) {
            practiceItems[index].user_answer = answer;
        }
    }

    function handlePrint() {
        printMode = true;
        sidebarOpen = false; // Force close sidebar before printing
        setTimeout(() => {
            window.print();
            // Reset after printing
            setTimeout(() => {
                printMode = false;
            }, 500);
        }, 200);
    }
</script>

<style>
    @media print {
        /* Add page breaks between questions */
        .question-container {
            page-break-inside: avoid;
            margin-bottom: 30px;
            break-inside: avoid;
            overflow: visible;
        }
        
        /* Print-specific styles */
        .print-header {
            margin-bottom: 20px;
            page-break-after: avoid;
        }
        
        .print-item {
            margin-bottom: 40px;
            page-break-inside: avoid;
        }
    }
</style>

<AppLayout bind:sidebarOpen>
    {#if !printMode}
        <div class="container mx-auto px-4 py-8 no-print">
            <div class="flex justify-between items-center mb-6">
                <h1 class="text-2xl font-bold text-gray-900">Practice Session</h1>
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
                        on:click={handlePrint}
                    >
                        Print Practice
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
            {:else if session}
                <div class="bg-white shadow-md rounded-lg p-6 mb-6">
                    <h2 class="text-xl font-semibold text-gray-900 mb-2">{session.name || 'Unnamed Practice'}</h2>
                    
                    <div class="flex flex-wrap gap-2 mb-4">
                        <span class={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full 
                            ${session.status === 'Completed' ? 'bg-green-100 text-green-800' : 
                            session.status === 'InProgress' ? 'bg-blue-100 text-blue-800' : 
                            'bg-gray-100 text-gray-800'}`}>
                            {session.status}
                        </span>
                        
                        {#if session.assigned_at}
                            <span class="bg-gray-100 text-gray-800 text-xs font-medium px-2.5 py-0.5 rounded">
                                Assigned: {new Date(session.assigned_at).toLocaleDateString()}
                            </span>
                        {/if}
                        
                        {#if session.completed_at}
                            <span class="bg-gray-100 text-gray-800 text-xs font-medium px-2.5 py-0.5 rounded">
                                Completed: {new Date(session.completed_at).toLocaleDateString()}
                            </span>
                        {/if}
                    </div>

                    {#if session.expand?.learner}
                        <div class="mb-4">
                            <h3 class="text-sm font-medium text-gray-700 mb-2">Learner:</h3>
                            <p class="text-gray-600">{session.expand.learner.name}</p>
                        </div>
                    {/if}

                    {#if session.expand?.practice_topic}
                        <div class="mb-4">
                            <h3 class="text-sm font-medium text-gray-700 mb-2">Topic:</h3>
                            <p class="text-gray-600">{session.expand.practice_topic.name}</p>
                        </div>
                    {/if}

                    {#if practiceItems.length > 0}
                        <div class="mt-6">
                            <h3 class="text-lg font-medium text-gray-900 mb-4">Practice Items</h3>
                            <div class="space-y-6">
                                {#each practiceItems as item, index}
                                    <div class="question-container">
                                        <QuestionFactory
                                            {item}
                                            {index}
                                            viewType={session.status === 'Completed' ? (isInstructor ? 'instructor' : 'answered') : 'learner'}
                                            disabled={session.status === 'Completed'}
                                            onAnswerChange={(answer) => handleAnswerChange(index, answer)}
                                        />
                                    </div>
                                {/each}
                            </div>
                        </div>
                    {:else}
                        <div class="bg-gray-50 border border-gray-200 p-4 rounded-md">
                            <p class="text-gray-600">No practice items available.</p>
                        </div>
                    {/if}
                </div>
            {/if}
        </div>
    {/if}

    <!-- Print-only version that shows when printing -->
    <div class="print-only print-container">
        {#if session}
            <div class="print-header">
                <h1 class="text-3xl font-bold mb-2">{session.name || 'Practice Session'}</h1>
                
                {#if session.expand?.practice_topic}
                    <h2 class="text-xl mb-1">Topic: {session.expand.practice_topic.name}</h2>
                {/if}
                
                {#if session.expand?.learner}
                    <p class="text-lg">Learner: {session.expand.learner.name}</p>
                {/if}
            </div>

            {#if practiceItems.length > 0}
                    {#each practiceItems as item, index}
                        <div class="print-item question-container">
                            <QuestionFactory
                                {item}
                                {index}
                                printMode={true}
                                viewType={session.status === 'Completed' ? (isInstructor ? 'instructor' : 'answered') : 'learner'}
                                disabled={session.status === 'Completed'}
                                onAnswerChange={(answer) => handleAnswerChange(index, answer)}
                            />
                        </div>
                    {/each}
            {:else}
                <p>No practice items available.</p>
            {/if}
        {/if}
    </div>
</AppLayout> 