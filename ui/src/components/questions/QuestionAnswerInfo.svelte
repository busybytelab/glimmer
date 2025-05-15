<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    
    /**
     * The practice item containing answer information
     */
    export let item: PracticeItem;
    
    /**
     * Whether to show the answer information
     */
    export let showAnswer: boolean = false;
    
    /**
     * Whether this is displayed in the instructor view mode
     */
    export let showInstructorInfo: boolean = false;
    
    /**
     * Better check if the answer is correct by comparing strings directly
     * Handle cases where is_correct flag might be incorrect
     */
    $: isCorrect = item.is_correct !== undefined 
        ? item.is_correct 
        : item.user_answer && item.correct_answer && 
          String(item.user_answer).trim().toLowerCase() === String(item.correct_answer).trim().toLowerCase();
    
    /**
     * Get learner name from expand property if it exists
     */
    $: learnerName = 'Learner'; //TODO: fix this, could get from practice_result or practice_session  item.expand?.practice_session?.expand?.learner?.nickname || 'Learner';
</script>

{#if showAnswer}
    <div class="mt-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-md">
        <h5 class="text-sm font-medium text-gray-900 dark:text-white mb-2">
            {showInstructorInfo ? `${learnerName}'s Answer:` : 'Your Answer:'}
        </h5>
        <p class="{isCorrect 
            ? 'text-green-600 dark:text-green-400' 
            : 'text-red-600 dark:text-red-400'} font-medium">
            {item.user_answer || 'Not answered'}
            {#if isCorrect !== undefined}
                {isCorrect ? '✓' : '✗'}
            {/if}
        </p>
    </div>
{/if} 