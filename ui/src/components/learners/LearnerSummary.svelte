<script lang="ts">
    import type { LearnerProgress } from '$lib/types';

    /**
     * The learner progress data to display
     */
    export let progress: LearnerProgress;

    /**
     * Calculate completion rate for better context
     */
    function getCompletionRate(progress: LearnerProgress): number {
        if (progress.overallProgress.totalSessions === 0) return 0;
        return Math.round((progress.overallProgress.completedSessions / progress.overallProgress.totalSessions) * 100);
    }

    /**
     * Calculate actual average percentage score
     */
    function getAveragePercentage(progress: LearnerProgress): number {
        const completed = progress.recentlyCompleted.length + progress.needsAttention.length;
        if (completed === 0) return 0;
        
        // For this calculation, we assume that completed sessions with 0 wrong answers = 100%
        // and sessions with wrong answers have a score based on correct answers
        const totalPercentage = progress.recentlyCompleted.length * 100 + 
            progress.needsAttention.reduce((sum, session) => {
                const correctAnswers = session.total_items - session.wrong_answers_count;
                return sum + Math.round((correctAnswers / session.total_items) * 100);
            }, 0);
        
        return Math.round(totalPercentage / completed);
    }
</script>

<!-- Learning Summary - Condensed Version -->
<div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm">
    <div class="flex items-center justify-between mb-3">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Learning Summary</h3>
        <span class="text-sm text-gray-500 dark:text-gray-400">
            {progress.overallProgress.totalSessions} session{progress.overallProgress.totalSessions !== 1 ? 's' : ''} completed
        </span>
    </div>
    
    <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div class="text-center">
            <div class="text-2xl font-bold text-green-600 dark:text-green-400">{getCompletionRate(progress)}%</div>
            <div class="text-xs text-gray-600 dark:text-gray-400">Success Rate</div>
        </div>
        
        <div class="text-center">
            <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">{getAveragePercentage(progress)}%</div>
            <div class="text-xs text-gray-600 dark:text-gray-400">Avg. Score</div>
        </div>
        
        <div class="text-center">
            <div class="text-2xl font-bold text-purple-600 dark:text-purple-400">{progress.inProgress.length}</div>
            <div class="text-xs text-gray-600 dark:text-gray-400">In Progress</div>
        </div>
        
        <div class="text-center">
            <div class="text-2xl font-bold text-orange-600 dark:text-orange-400">{progress.needsAttention.length}</div>
            <div class="text-xs text-gray-600 dark:text-gray-400">Needs Review</div>
        </div>
    </div>

    <!-- Strengths and Areas for Improvement - Compact Version -->
    {#if progress.overallProgress.doingWellIn.length > 0 || progress.overallProgress.needsHelpWith.length > 0}
        <div class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                {#if progress.overallProgress.doingWellIn.length > 0}
                    <div>
                        <h4 class="font-medium text-green-800 dark:text-green-200 mb-1 flex items-center">
                            <svg class="w-4 h-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                            </svg>
                            Strengths
                        </h4>
                        <div class="flex flex-wrap gap-1">
                            {#each progress.overallProgress.doingWellIn.slice(0, 3) as topic}
                                <span class="bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300 text-xs px-2 py-1 rounded">
                                    {topic}
                                </span>
                            {/each}
                        </div>
                    </div>
                {/if}

                {#if progress.overallProgress.needsHelpWith.length > 0}
                    <div>
                        <h4 class="font-medium text-orange-800 dark:text-orange-200 mb-1 flex items-center">
                            <svg class="w-4 h-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
                            </svg>
                            Focus Areas
                        </h4>
                        <div class="flex flex-wrap gap-1">
                            {#each progress.overallProgress.needsHelpWith.slice(0, 3) as topic}
                                <span class="bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-300 text-xs px-2 py-1 rounded">
                                    {topic}
                                </span>
                            {/each}
                        </div>
                    </div>
                {/if}
            </div>
        </div>
    {/if}
</div> 