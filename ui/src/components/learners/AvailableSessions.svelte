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

  /**
   * Get a friendly message based on session progress and mistakes
   */
  function getProgressMessage(session: PracticeSessionStats): string {
    if (session.wrong_answers_count > 0) {
      return `Let's try again! ${session.wrong_answers_count} question${session.wrong_answers_count === 1 ? '' : 's'} to practice`;
    }
    if (session.answered_items === 0) {
      return 'Ready to start!';
    }
    if (session.answered_items === session.total_items) {
      return 'All done! Great job! 🎉';
    }
    return `${session.answered_items} of ${session.total_items} questions completed`;
  }
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
            
            <p class="text-sm text-gray-600 dark:text-gray-300 mb-2">
              {getProgressMessage(session)}
            </p>

            {#if session.answered_items > 0}
              <div class="w-full bg-gray-200 rounded-full h-2.5 dark:bg-gray-700 mb-2">
                <div 
                  class="h-2.5 rounded-full {session.wrong_answers_count > 0 ? 'bg-orange-400' : 'bg-green-500'}" 
                  style="width: {(session.answered_items / session.total_items) * 100}%"
                ></div>
              </div>
            {/if}
          </div>
          <button
            on:click={() => startPractice(session.id)}
            class="ml-4 px-4 py-2 {session.wrong_answers_count > 0 ? 'bg-orange-500 hover:bg-orange-600' : 'bg-blue-600 hover:bg-blue-700'} text-white rounded-lg transition-colors text-sm font-medium min-w-[100px]"
          >
            {session.answered_items === 0 ? 'Start' : session.wrong_answers_count > 0 ? 'Try Again' : 'Continue'}
          </button>
        </div>
      </div>
    {/each}
  </div>
</div> 