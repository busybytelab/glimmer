<script lang="ts">
    import { onMount } from 'svelte';
    import type { PracticeItem, BreadcrumbItem, ReviewStatus } from '$lib/types';
    import { QuestionViewType } from '$lib/types';
    import QuestionFactory from '$components/questions/QuestionFactory.svelte';
    import ViewSelector from '$components/questions/ViewSelector.svelte';
    import { sessionService, type SessionWithExpandedData } from '$lib/services/session';
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import { resultsService } from '$lib/services/results';
    import ActionToolbar from '$components/common/ActionToolbar.svelte';
    import Breadcrumbs from '$components/common/Breadcrumbs.svelte';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    import SessionHeader from '$components/practice-sessions/SessionHeader.svelte';
    import ExportSessionButton from '$components/practice-session/ExportSessionButton.svelte';
    import { updateBreadcrumbs, handlePrint } from '$lib/utils/practice-session';

    let session: SessionWithExpandedData | null = null;
    let practiceItems: PracticeItem[] = [];
    let loading = true;
    let error: string | null = null;
    let printMode = false;
    let breadcrumbItems: BreadcrumbItem[] = [];
    let selectedViewType: QuestionViewType = QuestionViewType.PARENT;

    onMount(async () => {
        try {
            const sessionId = $page.params.id;
            if (sessionId) {
                await loadSession(sessionId);
                breadcrumbItems = updateBreadcrumbs(session);
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
            
            // Attach learner data to each practice item
            if (session && session.expand?.learner) {
                const learnerData = session.expand.learner;
                practiceItems = practiceItems.map(item => ({
                    ...item,
                    expand: {
                        ...item.expand,
                        learner: learnerData
                    }
                }));
            }

            // Fetch existing practice results
            const results = await resultsService.getResults(session.id, session.learner);

            // Map results to practice items
            practiceItems = practiceItems.map(item => {
                const result = results.find(r => r.practice_item === item.id);
                return {
                    ...item,
                    result: result || null
                };
            });
        } catch (err) {
            console.error('Failed to load session:', err);
            error = err instanceof Error ? err.message : 'Failed to load practice session';
            session = null;
            practiceItems = [];
        } finally {
            loading = false;
        }
    }

    function handleViewChange(newViewType: QuestionViewType) {
        selectedViewType = newViewType;
    }

    function editSession() {
        if (!session) return;
        goto(`/practice-sessions/${session.id}/edit`);
    }

    async function deleteSession() {
        if (!session) return;
        
        if (!confirm('Are you sure you want to delete this practice session?')) {
            return;
        }
        
        try {
            await sessionService.deleteSession(session.id);
            
            // Navigate back to practice topics or home
            goto('/home');
        } catch (err) {
            console.error('Failed to delete practice session:', err);
            error = 'Failed to delete practice session: ' + (err instanceof Error ? err.message : String(err));
        }
    }

    async function handleReviewStatusChange(itemId: string, status: ReviewStatus) {
        // Update the practice item in the local state
        practiceItems = practiceItems.map(item => {
            if (item.id === itemId) {
                return {
                    ...item,
                    review_status: status,
                    review_date: new Date().toISOString()
                };
            }
            return item;
        });
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
        .question-container {
            page-break-inside: avoid;
            margin-bottom: 30px;
            break-inside: avoid;
            overflow: visible;
        }
        
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
        <div class="flex items-center gap-2">
            <ExportSessionButton sessionId={$page.params.id} disabled={loading} />
            <ActionToolbar actions={sessionActions} />
        </div>
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
                <SessionHeader {session} />
                
                <ViewSelector 
                    viewType={selectedViewType}
                    onViewChange={handleViewChange}
                    isInstructor={true}
                />

                {#if practiceItems.length > 0}
                    <div class="mt-6">
                        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Practice Items</h3>
                        <div class="space-y-6">
                            {#each practiceItems as item, index}
                                <div class="question-container">
                                    <QuestionFactory
                                        {item}
                                        {index}
                                        viewType={selectedViewType}
                                        disabled={true}
                                        isInstructor={true}
                                        onReviewStatusChange={handleReviewStatusChange}
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

<!-- Print-only version -->
<div class="print-only print-container">
    {#if session}
        <div class="print-header">
            <h1 class="text-3xl font-bold mb-2">{session.name || 'Practice Session'}</h1>
            
            {#if session.expand?.practice_topic}
                <h2 class="text-xl mb-1">Topic: {session.expand?.practice_topic.name}</h2>
            {/if}
            
            {#if session.expand?.learner}
                <p class="text-lg">Learner: {session.expand?.learner?.nickname || 'Unknown Learner'}</p>
            {/if}
        </div>

        {#if practiceItems.length > 0}
            {#each practiceItems as item, index}
                <div class="print-item question-container">
                    <QuestionFactory
                        {item}
                        {index}
                        printMode={true}
                        viewType={selectedViewType}
                        disabled={true}
                        isInstructor={true}
                    />
                </div>
            {/each}
        {:else}
            <p>No practice items available.</p>
        {/if}
    {/if}
</div> 