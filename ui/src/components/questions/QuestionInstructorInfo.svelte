<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    
    /**
     * The practice item containing instructor information
     */
    export let item: PracticeItem;
    
    /**
     * Whether this is displayed in the instructor view mode
     */
    export let showInstructorInfo: boolean = false;
</script>

{#if showInstructorInfo}
    <div class="mt-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-md">
        <h5 class="text-sm font-medium text-gray-900 dark:text-white mb-2">Correct Answer:</h5>
        <p class="text-green-600 dark:text-green-400 font-medium">{item.correct_answer}</p>
        
        {#if item.explanation}
            <h5 class="text-sm font-medium text-gray-900 dark:text-white mt-4 mb-2">Explanation:</h5>
            <p class="text-gray-700 dark:text-gray-300">{item.explanation}</p>
        {/if}
        
        {#if item.user_answer && String(item.user_answer) !== String(item.correct_answer) && item.explanation_for_incorrect && item.explanation_for_incorrect[item.user_answer]}
            <h5 class="text-sm font-medium text-gray-900 dark:text-white mt-4 mb-2">Explanation for "{item.user_answer}":</h5>
            <p class="text-gray-700 dark:text-gray-300">{item.explanation_for_incorrect[item.user_answer]}</p>
        {/if}
        
        {#if item.hints && Array.isArray(item.hints) && item.hints.length > 0}
            <h5 class="text-sm font-medium text-gray-900 dark:text-white mt-4 mb-2">Available Hints:</h5>
            <ul class="list-disc list-inside text-gray-700 dark:text-gray-300 space-y-1">
                {#each item.hints as hint}
                    <li>{hint}</li>
                {/each}
            </ul>
        {/if}
        
        {#if item.hint_level_reached !== undefined && item.hint_level_reached > 0}
            <h5 class="text-sm font-medium text-gray-900 dark:text-white mt-4 mb-2">Learner Used Hints:</h5>
            <p class="text-gray-700 dark:text-gray-300">
                Reached hint level {item.hint_level_reached} of {Array.isArray(item.hints) ? item.hints.length : 0}
            </p>
        {/if}
    </div>
{/if} 