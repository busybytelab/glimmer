<script lang="ts">
    import { onMount } from 'svelte';
    import type { PracticeItem, PracticeResult, BreadcrumbItem } from '$lib/types';
    import { QuestionViewType } from '$lib/types';
    import QuestionFactory from '$components/questions/QuestionFactory.svelte';
    import { sessionService, type SessionWithExpandedData } from '$lib/services/session';
    import { page } from '$app/stores';
    import ActionToolbar from '$components/common/ActionToolbar.svelte';
    import Breadcrumbs from '$components/common/Breadcrumbs.svelte';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    import SessionHeader from '$components/practice-sessions/SessionHeader.svelte';
    import PracticeWizard from '$components/practice-sessions/PracticeWizard.svelte';
    import { updateBreadcrumbs, handlePrint } from '$lib/utils/practice-session';
    import { answersService } from '$lib/services/answers';
    import { resultsService } from '$lib/services/results';

    let session: SessionWithExpandedData | null = null;
    let practiceItems: PracticeItem[] = [];
    let loading = true;
    let error: string | null = null;
    let printMode = false;
    let breadcrumbItems: BreadcrumbItem[] = [];
    let savingItems: Set<number> = new Set();
    let selectedViewType: QuestionViewType = QuestionViewType.LEARNER;
    
    // Smart hint system
    let consecutiveIncorrectAttempts = new Map<string, number>();
    const HINT_THRESHOLD = 2;

    // View mode state
    // TODO: define type for viewMode
    let viewMode: 'all' | 'wizard' = 'all'; // NOTE: wizard mode has lots of bugs, UI need a redesign
    let currentStep = 0;
    // TODO: define type for stepResults
    let stepResults: ('correct' | 'incorrect' | 'pending')[] = [];

    // Reactive declarations for practice items
    $: practiceItemsWithHints = practiceItems.map(item => {
        const attempts = item.id ? (consecutiveIncorrectAttempts.get(item.id) || 0) : 0;
        return {
            item,
            attemptsCount: attempts,
            showHints: attempts >= HINT_THRESHOLD
        };
    });

    onMount(async () => {
        try {
            const sessionId = $page.params.sessionId;
            if (sessionId) {
                await loadSession(sessionId);
                selectedViewType = determineInitialViewType();
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
    
    function determineInitialViewType(): QuestionViewType {
        if (!session) return QuestionViewType.LEARNER;
        return session.status === 'Completed' ? QuestionViewType.ANSWERED : QuestionViewType.LEARNER;
    }

    async function loadSession(id: string) {
        try {
            loading = true;
            error = null;
            
            if (!id) {
                throw new Error('Session ID is required');
            }

            session = await sessionService.loadSessionForLearner(id);
            
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
            if (results.length > 0) {
                const learnerData = session && session.expand?.learner;
                practiceItems = practiceItems.map(item => {
                    const result = results.find((r: PracticeResult) => 
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
                            attempt_number: result.attempt_number,
                            expand: {
                                ...item.expand,
                                learner: learnerData
                            }
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

    // Update step results when practice items change
    $: if (practiceItems.length > 0) {
        stepResults = practiceItems.map(item => {
            if (item.is_correct === undefined) return 'pending';
            return item.is_correct ? 'correct' : 'incorrect';
        });
    }

    function handleStepClick(index: number) {
        // Allow navigation to any completed step or the next available step
        if (index <= currentStep || (index === currentStep + 1 && stepResults[currentStep] === 'correct')) {
            currentStep = index;
        }
    }

    async function handleAnswerChange(index: number, answer: string) {
        if (!session || !practiceItems[index]) return;

        try {
            savingItems.add(index);
            practiceItems[index].user_answer = answer;
            
            const practiceItem = practiceItems[index];
            const now = new Date().toISOString();

            // Call the evaluate answer endpoint
            const { isCorrect } = await answersService.evaluateAnswer(practiceItem.id, answer);
            
            if (isCorrect) {
                consecutiveIncorrectAttempts.set(practiceItem.id, 0);
            } else {
                const currentAttempts = consecutiveIncorrectAttempts.get(practiceItem.id) || 0;
                consecutiveIncorrectAttempts.set(practiceItem.id, currentAttempts + 1);
            }

            const existingResult = await resultsService.getLatestResult(practiceItem.id, session.id);

            if (existingResult) {
                await resultsService.updateResult(existingResult.id, {
                    answer: answer,
                    is_correct: isCorrect,
                    submitted_at: now,
                    attempt_number: (existingResult.attempt_number || 0) + 1
                });
            } else {
                const result = await resultsService.createResult({
                    practice_item: practiceItem.id,
                    practice_session: session.id,
                    learner: session.learner,
                    answer: answer,
                    is_correct: isCorrect,
                    started_at: now,
                    submitted_at: now,
                    attempt_number: 1,
                    hint_level_reached: 0
                });
                console.log('result', result);
            }

            practiceItems[index] = {
                ...practiceItems[index],
                is_correct: isCorrect
            };
            
            practiceItems = [...practiceItems];

            // In wizard mode, move to next step if answer is correct
            if (viewMode === 'wizard' && isCorrect && currentStep < practiceItems.length - 1) {
                currentStep++;
            }
        } catch (err) {
            console.error('Failed to save answer:', err);
            error = 'Failed to save answer: ' + (err instanceof Error ? err.message : String(err));
        } finally {
            savingItems.delete(index);
        }
    }
    
    async function handleHintRequest(index: number, level: number) {
        if (!session || !practiceItems[index]) return;
        
        const practiceItem = practiceItems[index];
        
        try {
            practiceItems[index] = {
                ...practiceItem,
                hint_level_reached: level
            };
            
            practiceItems = [...practiceItems];
            
            const existingResult = await resultsService.getLatestResult(practiceItem.id, session.id);
            
            if (existingResult) {
                await resultsService.updateResult(existingResult.id, {
                    hint_level_reached: level
                });
            }
        } catch (err) {
            console.error('Failed to update hint level:', err);
            error = 'Failed to update hint level: ' + (err instanceof Error ? err.message : String(err));
        }
    }

    // Actions for the toolbar
    $: sessionActions = [
        {
            id: 'print',
            label: 'Print',
            icon: 'print',
            variant: 'primary' as const,
            onClick: handlePrint
        }
    ];

    function toggleViewMode() {
        viewMode = viewMode === 'all' ? 'wizard' : 'all';
    }
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
        <div class="flex items-center space-x-4">
            <button
                class="p-2 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200"
                class:text-primary-500={viewMode === 'wizard'}
                class:text-gray-500={viewMode === 'all'}
                on:click={toggleViewMode}
                title={viewMode === 'all' ? 'Switch to Step by Step view' : 'Switch to All Questions view'}
            >
                {#if viewMode === 'all'}
                    <!-- List view icon -->
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                        <path d="M5 3a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2V5a2 2 0 00-2-2H5zM5 11a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2v-2a2 2 0 00-2-2H5zM11 5a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V5zM11 13a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
                    </svg>
                {:else}
                    <!-- Step by step view icon -->
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M10 2a1 1 0 011 1v1.323l3.954 1.582 1.599-.8a1 1 0 01.894 1.79l-1.233.616 1.738 5.42a1 1 0 01-.285 1.05A3.989 3.989 0 0115 15a3.989 3.989 0 01-2.667-1.019 1 1 0 01-.285-1.05l1.715-5.349L11 4.477V16h2a1 1 0 110 2H7a1 1 0 110-2h2V4.477L6.237 7.582l1.715 5.349a1 1 0 01-.285 1.05A3.989 3.989 0 015 15a3.989 3.989 0 01-2.667-1.019 1 1 0 01-.285-1.05l1.738-5.42-1.233-.616a1 1 0 01.894-1.79l1.599.8L9 4.323V3a1 1 0 011-1z" clip-rule="evenodd" />
                    </svg>
                {/if}
            </button>
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

                {#if practiceItems.length > 0}
                    <div class="mt-6">
                        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Practice Items</h3>
                        
                        {#if viewMode === 'wizard'}
                            <PracticeWizard
                                {practiceItems}
                                {currentStep}
                                {stepResults}
                                {selectedViewType}
                                sessionStatus={session?.status || ''}
                                {savingItems}
                                {consecutiveIncorrectAttempts}
                                {HINT_THRESHOLD}
                                onStepClick={handleStepClick}
                                onAnswerChange={handleAnswerChange}
                                onHintRequest={handleHintRequest}
                            />
                        {:else}
                            <div class="space-y-6">
                                {#each practiceItemsWithHints as { item, showHints }, index}
                                    <div class="question-container">
                                        <QuestionFactory
                                            {item}
                                            {index}
                                            viewType={selectedViewType}
                                            disabled={selectedViewType !== QuestionViewType.LEARNER || session.status === 'Completed' || savingItems.has(index)}
                                            onAnswerChange={(answer: string) => handleAnswerChange(index, answer)}
                                            isInstructor={false}
                                            {showHints}
                                            onHintRequested={(level: number) => handleHintRequest(index, level)}
                                        />
                                    </div>
                                {/each}
                            </div>
                        {/if}
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
                <h2 class="text-xl mb-1">Topic: {session.expand.practice_topic.name}</h2>
            {/if}
            
            {#if session.expand?.learner}
                <p class="text-lg">Learner: {session.expand.learner?.nickname || 'Unknown Learner'}</p>
            {/if}
        </div>

        {#if practiceItemsWithHints.length > 0}
            {#each practiceItemsWithHints as { item, showHints }, index}
                <div class="print-item question-container">
                    <QuestionFactory
                        {item}
                        {index}
                        printMode={true}
                        viewType={selectedViewType}
                        disabled={true}
                        onAnswerChange={(answer: string) => handleAnswerChange(index, answer)}
                        isInstructor={false}
                        {showHints}
                        onHintRequested={(level: number) => handleHintRequest(index, level)}
                    />
                </div>
            {/each}
        {:else}
            <p>No practice items available.</p>
        {/if}
    {/if}
</div> 