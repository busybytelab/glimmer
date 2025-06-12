<script lang="ts">
  import { goto } from '$app/navigation';
  import type { PracticeSessionStats } from '$lib/types';

  /**
   * Array of available sessions for the learner
   */
  export let sessions: PracticeSessionStats[];
  /**
   * Learner ID for routing
   */
  export let learnerId: string;

  /**
   * Navigate to the session's practice page for the learner
   * @param sessionId The session's unique ID
   */
  const startPractice = (sessionId: string) => {
    goto(`/learners/${learnerId}/practice-sessions/${sessionId}`);
  };
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6 mb-6">
  <div class="flex items-center justify-between mb-4">
    <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
      Available Practice Sessions
    </h2>
  </div>

  <div class="space-y-4">
    {#each sessions as session (session.id)}
      <div class="border dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <h3 class="font-medium text-gray-900 dark:text-white mb-2">
              {session.session_name} <span class="text-xs text-gray-500 ml-2">({session.topic_name})</span>
            </h3>
            {#if session.answered_items > 0}
              <p class="text-sm text-gray-600 dark:text-gray-300">
                {session.answered_items} of {session.total_items} questions answered
                {#if session.total_score > 0}
                  &nbsp;|&nbsp;Score: {Math.round(session.total_score / session.total_items)}
                {/if}
              </p>
            {/if}
          </div>
          <button
            on:click={() => startPractice(session.id)}
            class="ml-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm font-medium min-w-[100px]"
          >
            {session.answered_items === 0 ? 'Start' : 'Continue'}
          </button>
        </div>
      </div>
    {/each}
  </div>
</div> 