<script lang="ts">
    import type { LearnerProgress } from '$lib/types';
    import { sessionService } from '$lib/services/session';
    import LoadingSpinner from '../common/LoadingSpinner.svelte';
    import ErrorAlert from '../common/ErrorAlert.svelte';
    import LearnerSummary from './LearnerSummary.svelte';

    export let learnerId: string;
    export let learnerName: string = '';

    let loading = true;
    let error: string | null = null;
    let progress: LearnerProgress | null = null;

    async function loadProgress() {
        try {
            loading = true;
            error = null;
            progress = await sessionService.getLearnerProgressForParent(learnerId);
        } catch (err) {
            console.error('Failed to load learner progress:', err);
            error = err instanceof Error ? err.message : 'Failed to load progress';
        } finally {
            loading = false;
        }
    }

    $: if (learnerId) {
        loadProgress();
    }

    function formatDate(dateStr: string): string {
        return new Date(dateStr).toLocaleDateString();
    }
</script>

{#if loading}
    <div class="flex justify-center items-center h-32">
        <LoadingSpinner size="sm" color="primary" />
    </div>
{:else if error}
    <ErrorAlert message={error} />
{:else if progress}
    <div class="space-y-6">
        <!-- Learning Summary -->
        <LearnerSummary {progress} />

        <!-- Sessions Needing Attention - Priority Section -->
        {#if progress.needsAttention.length > 0}
            <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2 flex items-center">
                    <svg class="w-5 h-5 mr-2 text-orange-500" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
                    </svg>
                    Needs Review
                </h3>
                <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
                    Help {learnerName || 'your child'} review these questions for better understanding
                </p>
                <div class="space-y-3">
                    {#each progress.needsAttention as session}
                        <div class="border-l-4 border-orange-400 dark:border-orange-500 bg-orange-50 dark:bg-orange-900/20 rounded-r-lg p-4">
                            <div class="flex justify-between items-start">
                                <div class="flex-1">
                                    <h4 class="font-medium text-gray-900 dark:text-white">{session.topic_name}</h4>
                                    <p class="text-sm text-orange-600 dark:text-orange-300 mt-1">
                                        {session.wrong_answers_count} question{session.wrong_answers_count === 1 ? '' : 's'} to review
                                    </p>
                                </div>
                                <div class="text-right ml-4">
                                    <div class="text-sm text-gray-500 dark:text-gray-400">
                                        {formatDate(session.last_answer_time)}
                                    </div>
                                    <a
                                        href={`/account/practice-sessions/${session.id}/overview`}
                                        class="inline-flex items-center text-sm text-orange-700 dark:text-orange-300 hover:text-orange-800 dark:hover:text-orange-200 mt-1"
                                    >
                                        Review
                                        <svg class="w-4 h-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
                                            <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                                        </svg>
                                    </a>
                                </div>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        {/if}

        <!-- Current Practice Sessions -->
        {#if progress.inProgress.length > 0}
            <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center">
                    <svg class="w-5 h-5 mr-2 text-blue-500" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-8.293l-3-3a1 1 0 00-1.414 0l-3 3a1 1 0 001.414 1.414L9 9.414V13a1 1 0 102 0V9.414l1.293 1.293a1 1 0 001.414-1.414z" clip-rule="evenodd" />
                    </svg>
                    In Progress
                </h3>
                <div class="space-y-4">
                    {#each progress.inProgress as session}
                        <div class="border dark:border-gray-700 rounded-lg p-4">
                            <div class="flex justify-between items-start mb-2">
                                <div class="flex-1">
                                    <h4 class="font-medium text-gray-900 dark:text-white">{session.topic_name}</h4>
                                    <p class="text-sm text-gray-600 dark:text-gray-300 mt-1">
                                        {session.answered_items} of {session.total_items} questions completed
                                    </p>
                                </div>
                                <div class="text-right ml-4">
                                    <div class="text-sm text-gray-500 dark:text-gray-400">
                                        {formatDate(session.last_answer_time)}
                                    </div>
                                    <a
                                        href={`/account/practice-sessions/${session.id}/overview`}
                                        class="inline-flex items-center text-sm text-blue-700 dark:text-blue-300 hover:text-blue-800 dark:hover:text-blue-200 mt-1"
                                    >
                                        Continue
                                        <svg class="w-4 h-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
                                            <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                                        </svg>
                                    </a>
                                </div>
                            </div>
                            <div class="w-full bg-gray-200 rounded-full h-2 dark:bg-gray-700">
                                <div class="bg-blue-600 h-2 rounded-full transition-all duration-300 ease-in-out" style="width: {(session.answered_items / session.total_items) * 100}%"></div>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        {/if}

        <!-- Recently Completed -->
        {#if progress.recentlyCompleted.length > 0}
            <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center">
                    <svg class="w-5 h-5 mr-2 text-green-500" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                    </svg>
                    Recently Completed
                </h3>
                <div class="space-y-3">
                    {#each progress.recentlyCompleted as session}
                        <div class="border-l-4 border-green-400 dark:border-green-500 bg-green-50 dark:bg-green-900/20 rounded-r-lg p-4">
                            <div class="flex justify-between items-start">
                                <div class="flex-1">
                                    <h4 class="font-medium text-gray-900 dark:text-white">{session.topic_name}</h4>
                                    <p class="text-sm text-green-600 dark:text-green-300 mt-1">
                                        Perfect! All {session.total_items} questions correct
                                    </p>
                                </div>
                                <div class="text-right ml-4">
                                    <div class="text-sm text-gray-500 dark:text-gray-400">
                                        {formatDate(session.last_answer_time)}
                                    </div>
                                    <a
                                        href={`/account/practice-sessions/${session.id}/overview`}
                                        class="inline-flex items-center text-sm text-green-700 dark:text-green-300 hover:text-green-800 dark:hover:text-green-200 mt-1"
                                    >
                                        View Details
                                        <svg class="w-4 h-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
                                            <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                                        </svg>
                                    </a>
                                </div>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        {/if}
    </div>
{/if} 