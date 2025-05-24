<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    import { fade } from 'svelte/transition';

    export let item: PracticeItem;

    // Determine if we should show the explanation
    $: showExplanation = item.explanation || 
        (item.explanation_for_incorrect && 
         item.user_answer && 
         item.explanation_for_incorrect[item.user_answer]);

    // Get the appropriate explanation text
    $: explanationText = item.explanation || 
        (item.explanation_for_incorrect && 
         item.user_answer && 
         item.explanation_for_incorrect[item.user_answer]);

    // Determine if the answer is correct
    $: isCorrect = item.is_correct !== undefined 
        ? item.is_correct 
        : item.user_answer && item.correct_answer && 
          String(item.user_answer).trim().toLowerCase() === String(item.correct_answer).trim().toLowerCase();
</script>

{#if showExplanation}
    <div class="mt-2" transition:fade={{ duration: 300 }}>

        <div class="mt-3 p-4 bg-blue-50 dark:bg-blue-900/30 rounded-md border border-blue-100 dark:border-blue-800" transition:fade={{ duration: 300 }}>
            <div class="flex items-center mb-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1 text-blue-600 dark:text-blue-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                </svg>
                <h5 class="text-sm font-medium text-blue-900 dark:text-blue-200">
                    {isCorrect ? 'Correct Answer Explanation:' : 'Explanation for Your Answer:'}
                </h5>
            </div>
            <p class="text-sm text-blue-700 dark:text-blue-300">{explanationText}</p>
        </div>
    </div>
{/if} 