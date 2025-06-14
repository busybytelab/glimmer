<script lang="ts">
    import type { PracticeSessionStats } from '$lib/types';
    import { sessionService } from '$lib/services/session';
    import LoadingSpinner from '../common/LoadingSpinner.svelte';
    import ErrorAlert from '../common/ErrorAlert.svelte';

    export let learnerId: string;

    let loading = true;
    let error: string | null = null;
    let progress: {
        needsAttention: PracticeSessionStats[];
        inProgress: PracticeSessionStats[];
        recentlyCompleted: PracticeSessionStats[];
        overallProgress: {
            totalSessions: number;
            completedSessions: number;
            averageScore: number;
            needsHelpWith: string[];
            doingWellIn: string[];
        };
    } | null = null;

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
        <!-- Overall Progress -->
        <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Overall Progress</h3>
            
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
                    <div class="text-sm text-blue-600 dark:text-blue-300">Total Practice Sessions</div>
                    <div class="text-2xl font-bold text-blue-700 dark:text-blue-200">{progress.overallProgress.totalSessions}</div>
                </div>
                
                <div class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4">
                    <div class="text-sm text-green-600 dark:text-green-300">Completed Successfully</div>
                    <div class="text-2xl font-bold text-green-700 dark:text-green-200">{progress.overallProgress.completedSessions}</div>
                </div>
                
                <div class="bg-purple-50 dark:bg-purple-900/20 rounded-lg p-4">
                    <div class="text-sm text-purple-600 dark:text-purple-300">Average Score</div>
                    <div class="text-2xl font-bold text-purple-700 dark:text-purple-200">{Math.round(progress.overallProgress.averageScore)}%</div>
                </div>
            </div>

            <!-- Strengths and Areas for Improvement -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                {#if progress.overallProgress.doingWellIn.length > 0}
                    <div class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4">
                        <h4 class="font-medium text-green-800 dark:text-green-200 mb-2">Doing Great In:</h4>
                        <ul class="space-y-1">
                            {#each progress.overallProgress.doingWellIn as topic}
                                <li class="text-green-600 dark:text-green-300 text-sm flex items-center">
                                    <svg class="w-4 h-4 mr-2" viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                                    </svg>
                                    {topic}
                                </li>
                            {/each}
                        </ul>
                    </div>
                {/if}

                {#if progress.overallProgress.needsHelpWith.length > 0}
                    <div class="bg-orange-50 dark:bg-orange-900/20 rounded-lg p-4">
                        <h4 class="font-medium text-orange-800 dark:text-orange-200 mb-2">Needs More Practice In:</h4>
                        <ul class="space-y-1">
                            {#each progress.overallProgress.needsHelpWith as topic}
                                <li class="text-orange-600 dark:text-orange-300 text-sm flex items-center">
                                    <svg class="w-4 h-4 mr-2" viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V7z" clip-rule="evenodd" />
                                    </svg>
                                    {topic}
                                </li>
                            {/each}
                        </ul>
                    </div>
                {/if}
            </div>
        </div>

        <!-- Current Practice Sessions -->
        {#if progress.inProgress.length > 0}
            <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Current Practice</h3>
                <div class="space-y-4">
                    {#each progress.inProgress as session}
                        <div class="border dark:border-gray-700 rounded-lg p-4">
                            <div class="flex justify-between items-start">
                                <div>
                                    <h4 class="font-medium text-gray-900 dark:text-white">{session.topic_name}</h4>
                                    <p class="text-sm text-gray-600 dark:text-gray-300 mt-1">
                                        {session.answered_items} of {session.total_items} questions completed
                                    </p>
                                </div>
                                <div class="text-sm text-gray-500 dark:text-gray-400">
                                    Started {formatDate(session.last_answer_time)}
                                </div>
                            </div>
                            <div class="mt-2 w-full bg-gray-200 rounded-full h-2.5 dark:bg-gray-700">
                                <div class="bg-blue-600 h-2.5 rounded-full" style="width: {(session.answered_items / session.total_items) * 100}%" />
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        {/if}

        <!-- Sessions Needing Attention -->
        {#if progress.needsAttention.length > 0}
            <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Needs Review</h3>
                <div class="space-y-4">
                    {#each progress.needsAttention as session}
                        <div class="border-l-4 border-orange-400 dark:border-orange-500 bg-orange-50 dark:bg-orange-900/20 rounded-r-lg p-4">
                            <div class="space-y-2">
                                <div class="flex justify-between items-start">
                                    <div>
                                        <h4 class="font-medium text-gray-900 dark:text-white">{session.topic_name}</h4>
                                        <p class="text-sm text-orange-600 dark:text-orange-300 mt-1">
                                            {session.wrong_answers_count} question{session.wrong_answers_count === 1 ? '' : 's'} to practice again
                                        </p>
                                    </div>
                                    <div class="text-sm text-gray-500 dark:text-gray-400 ml-4">
                                        {formatDate(session.last_answer_time)}
                                    </div>
                                </div>
                                <div class="pt-2 border-t border-orange-200 dark:border-orange-700/30">
                                    <a
                                        href={`/account/practice-sessions/${session.id}/instructor`}
                                        class="inline-flex items-center text-sm text-orange-700 dark:text-orange-300 hover:text-orange-800 dark:hover:text-orange-200"
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

        <!-- Recently Completed -->
        {#if progress.recentlyCompleted.length > 0}
            <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Recently Completed</h3>
                <div class="space-y-4">
                    {#each progress.recentlyCompleted as session}
                        <div class="border-l-4 border-green-400 dark:border-green-500 bg-green-50 dark:bg-green-900/20 rounded-r-lg p-4">
                            <div class="space-y-2">
                                <div class="flex justify-between items-start">
                                    <div>
                                        <h4 class="font-medium text-gray-900 dark:text-white">{session.topic_name}</h4>
                                        <p class="text-sm text-green-600 dark:text-green-300 mt-1">
                                            Completed all {session.total_items} questions successfully!
                                        </p>
                                    </div>
                                    <div class="text-sm text-gray-500 dark:text-gray-400 ml-4">
                                        {formatDate(session.last_answer_time)}
                                    </div>
                                </div>
                                <div class="pt-2 border-t border-green-200 dark:border-green-700/30">
                                    <a
                                        href={`/account/practice-sessions/${session.id}/instructor`}
                                        class="inline-flex items-center text-sm text-green-700 dark:text-green-300 hover:text-green-800 dark:hover:text-green-200"
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