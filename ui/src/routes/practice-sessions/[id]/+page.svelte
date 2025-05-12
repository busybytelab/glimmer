<script lang="ts">
    import { onMount } from 'svelte';
    import type { PracticeItem } from '$lib/types';
    import QuestionFactory from '../../../components/questions/QuestionFactory.svelte';
    import { sessionService, type SessionWithExpandedData } from '$lib/services/session';
    import { page } from '$app/stores';
    import { goto, beforeNavigate, afterNavigate } from '$app/navigation';
    import pb from '$lib/pocketbase';
    import ActionToolbar from '../../../components/common/ActionToolbar.svelte';
    import Breadcrumbs from '../../../components/common/Breadcrumbs.svelte';
    import LoadingSpinner from '../../../components/common/LoadingSpinner.svelte';
    import ErrorAlert from '../../../components/common/ErrorAlert.svelte';
    import { debounce } from '$lib/utils/debounce';
    import type { RecordModel } from 'pocketbase';

    // Define practice result interface to work with PocketBase RecordModel
    interface PracticeResult extends RecordModel {
        practice_item: string;
        answer: string;
        is_correct: boolean;
        score: number;
        feedback: string;
        hint_level_reached: number;
        attempt_number: number;
    }

    // Define the breadcrumb item type
    type BreadcrumbItem = {
        label: string;
        href?: string;
        icon?: string;
    };

    let session: SessionWithExpandedData | null = null;
    let practiceItems: PracticeItem[] = [];
    let loading = true;
    let error: string | null = null;
    let isInstructor = false;
    let printMode = false;
    let breadcrumbItems: BreadcrumbItem[] = [];
    let savingItems: Set<number> = new Set();

    // Track the currently focused element
    let focusedElement: HTMLElement | null = null;

    // Store the element before navigation
    beforeNavigate(() => {
        focusedElement = document.activeElement as HTMLElement;
    });

    // Restore focus after navigation
    afterNavigate(() => {
        if (focusedElement) {
            focusedElement.focus();
        }
    });

    onMount(async () => {
        try {
            const sessionId = $page.params.id;
            
            if (sessionId) {
                await loadSession(sessionId);
                await checkUserRole();
                updateBreadcrumbs();
            } else {
                error = 'Invalid session ID';
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

            // Fetch existing practice results for this session and learner
            const results = await pb.collection('practice_results').getList(1, 100, {
                filter: `practice_session = "${session.id}" && learner = "${session.learner}"`,
                expand: 'practice_item,learner',
                sort: '-created'
            });

            // Map results to practice items
            if (results.items.length > 0) {
                practiceItems = practiceItems.map(item => {
                    const result = results.items.find((r): r is PracticeResult => 
                        'practice_item' in r && r.practice_item === item.id
                    );
                    if (result) {
                        return {
                            ...item,
                            user_answer: result.answer,
                            is_correct: result.is_correct,
                            score: result.score,
                            feedback: result.feedback,
                            hint_level_reached: result.hint_level_reached,
                            attempt_number: result.attempt_number
                        };
                    }
                    return item;
                });
            }
        } catch (err) {
            console.error('Failed to load session:', err);
            error = err instanceof Error ? err.message : 'Failed to load practice session';
            session = null;
            practiceItems = [];
        } finally {
            loading = false;
        }
    }

    function editSession() {
        if (!session) return;
        goto(`/practice-sessions/edit/${session.id}`);
    }

    async function deleteSession() {
        if (!session) return;
        
        if (!confirm('Are you sure you want to delete this practice session?')) {
            return;
        }
        
        try {
            await pb.collection('practice_sessions').delete(session.id);
            
            // Navigate back to practice topics or dashboard
            if (session.expand?.practice_topic) {
                goto(`/practice-topics/${session.expand.practice_topic.id}`);
            } else {
                goto('/dashboard');
            }
        } catch (err) {
            console.error('Failed to delete practice session:', err);
            error = 'Failed to delete practice session: ' + (err instanceof Error ? err.message : String(err));
        }
    }

    // Create a debounced version of handleAnswerChange
    const debouncedSaveAnswer = debounce(async (index: number, answer: string) => {
        if (!session || !practiceItems[index]) return;

        try {
            // Add this index to the set of saving items
            savingItems.add(index);
            
            // We've already updated the answer in the handleAnswerChange function,
            // so we don't need to update it again here
            
            // Create or update practice result
            const practiceItem = practiceItems[index];
            const now = new Date().toISOString();

            // Check if a result already exists for this item
            const existingResults = await pb.collection('practice_results').getList(1, 1, {
                filter: `practice_item = "${practiceItem.id}" && practice_session = "${session.id}"`,
                sort: '-created'
            });

            let result;
            if (existingResults.items.length > 0) {
                // Update existing result
                result = await pb.collection('practice_results').update(existingResults.items[0].id, {
                    answer: answer,
                    submitted_at: now,
                    attempt_number: (existingResults.items[0].attempt_number || 0) + 1
                });
            } else {
                // Create new result
                result = await pb.collection('practice_results').create({
                    practice_item: practiceItem.id,
                    practice_session: session.id,
                    learner: session.learner,
                    answer: answer,
                    started_at: now,
                    submitted_at: now,
                    attempt_number: 1
                });
            }

            // No need to update the DOM here as we've already done it in handleAnswerChange
        } catch (err) {
            console.error('Failed to save answer:', err);
            error = 'Failed to save answer: ' + (err instanceof Error ? err.message : String(err));
        } finally {
            // Remove this index from the set of saving items
            savingItems.delete(index);
        }
    }, 1000); // Wait 1 second after the last change before saving

    async function handleAnswerChange(index: number, answer: string) {
        // Update the answer immediately for a responsive UI
        practiceItems[index].user_answer = answer;
        
        // Pass to the debounced save function
        debouncedSaveAnswer(index, answer);
    }

    function handlePrint() {
        printMode = true;
        setTimeout(() => {
            window.print();
            // Reset after printing
            setTimeout(() => {
                printMode = false;
            }, 500);
        }, 200);
    }

    function updateBreadcrumbs() {
        if (!session) return;
        
        const items: BreadcrumbItem[] = [
            {
                label: 'Topics',
                href: '/practice-topics',
                icon: 'topic'
            }
        ];
        
        if (session.expand?.practice_topic) {
            items.push({
                label: session.expand.practice_topic.name,
                href: `/practice-topics/${session.expand.practice_topic.id}`,
                icon: 'topic'
            });
        }
        
        // Use the session name or a fallback without date parsing
        const sessionName = session.name || 'Practice Session';
        
        items.push({
            label: sessionName,
            icon: 'session'
        });
        
        breadcrumbItems = items;
    }

    // Actions for the toolbar
    $: sessionActions = [
        {
            id: 'print',
            label: 'Print',
            icon: 'print',
            variant: 'primary' as const,
            onClick: handlePrint
        },
        {
            id: 'edit',
            label: 'Edit',
            icon: 'edit',
            variant: 'secondary' as const,
            onClick: editSession
        },
        {
            id: 'delete',
            label: 'Delete',
            icon: 'delete',
            variant: 'danger' as const,
            onClick: deleteSession
        }
    ];
</script>

<style lang="postcss">
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

{#if !printMode}
<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6 no-print">
    <div class="flex justify-between items-center mb-6">
        <div>
            <Breadcrumbs items={breadcrumbItems} />
        </div>
        <ActionToolbar actions={sessionActions} />
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <LoadingSpinner size="md" color="primary" />
        </div>
    {:else if error}
        <ErrorAlert message={error} />
    {:else if session}
        <div class="bg-white dark:bg-gray-800 shadow rounded-lg">
            <div class="px-4 py-5 sm:p-6">
                <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-4">{session.name || 'Practice Session'}</h2>
            
                {#if session.expand?.practice_topic}
                    <p class="text-gray-600 dark:text-gray-300 mb-4">Topic: {session.expand.practice_topic.name}</p>
                {/if}

            {#if session.expand?.learner}
                    <p class="text-gray-600 dark:text-gray-300 mb-4">Learner: {session.expand.learner.name}</p>
            {/if}

            {#if practiceItems.length > 0}
                <div class="mt-6">
                    <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Practice Items</h3>
                    <div class="space-y-6">
                        {#each practiceItems as item, index}
                            <div class="question-container">
                                <QuestionFactory
                                    {item}
                                    {index}
                                    viewType={session.status === 'Completed' ? (isInstructor ? 'instructor' : 'answered') : 'learner'}
                                    disabled={session.status === 'Completed' || savingItems.has(index)}
                                    onAnswerChange={(answer) => handleAnswerChange(index, answer)}
                                />
                            </div>
                        {/each}
                    </div>
                </div>
            {:else}
                <div class="bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 p-4 rounded-md">
                    <p class="text-gray-600 dark:text-gray-300">No practice items available.</p>
                </div>
            {/if}
            </div>
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
                            disabled={session.status === 'Completed' || savingItems.has(index)}
                            onAnswerChange={(answer) => handleAnswerChange(index, answer)}
                        />
                    </div>
                {/each}
        {:else}
            <p>No practice items available.</p>
        {/if}
    {/if}
</div> 