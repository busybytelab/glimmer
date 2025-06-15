<script lang="ts">
    import type { SessionWithExpandedData } from '$lib/services/session';
    import type { PracticeSessionStats } from '$lib/types';
    import { IconTypeMap } from '$lib/types';
    import SessionStatsDonut from './SessionStatsDonut.svelte';

    /**
     * The session data to display in the header
     */
    export let session: SessionWithExpandedData;

    /**
     * The session statistics
     */
    export let stats: PracticeSessionStats | null = null;

    function getProgressDescription(stats: PracticeSessionStats): string {
        if (stats.answered_items === 0) {
            return 'Not started yet';
        } else if (stats.answered_items === stats.total_items) {
            return `Completed with score ${stats.total_score}%`;
        } else {
            return `${stats.answered_items} of ${stats.total_items} questions answered`;
        }
    }
</script>

<div class="mb-6">
    <div class="flex items-start gap-6">
        {#if session.expand?.learner}
            <div class="flex flex-col items-center bg-gray-50 dark:bg-gray-700/40 rounded-lg p-3">
                <div class="bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300 p-3 rounded-full mb-2">
                    {#if session.expand.learner.avatar}
                        <img 
                            src={session.expand.learner.avatar} 
                            alt={session.expand.learner.nickname || 'Unknown'} 
                            class="h-12 w-12 rounded-full object-cover"
                        />
                    {:else}
                        <svg class="h-12 w-12" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                            <path stroke-linecap="round" stroke-linejoin="round" d={IconTypeMap.user} />
                        </svg>
                    {/if}
                </div>
                <span class="text-gray-900 dark:text-white font-medium text-center text-sm">
                    {session.expand.learner.nickname || 'Unknown'}
                </span>
            </div>
        {/if}

        <div class="flex-1 flex items-start justify-between">
            <div>
                <div class="flex items-center gap-3 mb-2">
                    {#if session.expand?.practice_topic}
                        <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
                            {session.expand.practice_topic.name}
                        </h2>
                        <span class="text-lg text-gray-600 dark:text-gray-400">
                            {session.name || 'Practice Session'}
                        </span>
                    {:else}
                        <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
                            {session.name || 'Practice Session'}
                        </h2>
                    {/if}
                </div>
                
                {#if session.expand?.practice_topic?.description}
                    <p class="text-gray-600 dark:text-gray-400 text-sm">
                        {session.expand.practice_topic.description}
                    </p>
                {/if}

                {#if stats}
                    <p class="text-gray-600 dark:text-gray-400 text-sm mt-2">
                        {getProgressDescription(stats)}
                    </p>
                {/if}
            </div>

            {#if stats}
                <SessionStatsDonut {stats} showLegend={false} />
            {/if}
        </div>
    </div>
</div> 