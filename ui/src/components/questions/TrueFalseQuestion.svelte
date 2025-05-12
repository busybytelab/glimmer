<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    import QuestionHeader from './QuestionHeader.svelte';

    export let item: PracticeItem;
    export let index: number;
    export let disabled = false;
    export let showAnswer = false;
    export let showInstructorInfo = false;
    export let onAnswerChange: ((answer: string) => void) | undefined = undefined;
    export let printMode = false;

    function handleAnswerChange(event: Event) {
        const target = event.target as HTMLInputElement;
        if (onAnswerChange) {
            onAnswerChange(target.value);
        }
    }
</script>

<div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 bg-white dark:bg-gray-800">
    <QuestionHeader {item} {index} />
    <p class="text-gray-700 dark:text-gray-300 mb-4">{item.question_text}</p>
    
    <div class="space-y-2">
        <div class="flex items-center">
            {#if printMode}
                <div class="w-4 h-4 border border-gray-400 dark:border-gray-500 rounded mr-2"></div>
            {:else}
                <input
                    type="radio"
                    id={`question-${index}-true-${item.id}`}
                    name={`question-${index}-${item.id}`}
                    value="true"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 dark:border-gray-600 dark:bg-gray-700"
                    {disabled}
                    checked={item.user_answer === 'true'}
                    on:change={handleAnswerChange}
                />
            {/if}
            <label 
                for={printMode ? undefined : `question-${index}-true-${item.id}`}
                class="ml-3 text-gray-700 dark:text-gray-300"
            >
                True
            </label>
        </div>
        <div class="flex items-center">
            {#if printMode}
                <div class="w-4 h-4 border border-gray-400 dark:border-gray-500 rounded mr-2"></div>
            {:else}
                <input
                    type="radio"
                    id={`question-${index}-false-${item.id}`}
                    name={`question-${index}-${item.id}`}
                    value="false"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 dark:border-gray-600 dark:bg-gray-700"
                    {disabled}
                    checked={item.user_answer === 'false'}
                    on:change={handleAnswerChange}
                />
            {/if}
            <label 
                for={printMode ? undefined : `question-${index}-false-${item.id}`}
                class="ml-3 text-gray-700 dark:text-gray-300"
            >
                False
            </label>
        </div>
    </div>

    {#if showAnswer}
        <div class="mt-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-md">
            <h5 class="text-sm font-medium text-gray-900 dark:text-white mb-2">Your Answer:</h5>
            <p class="text-gray-700 dark:text-gray-300">{item.user_answer || 'Not answered'}</p>
        </div>
    {/if}

    {#if showInstructorInfo}
        <div class="mt-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-md">
            <h5 class="text-sm font-medium text-gray-900 dark:text-white mb-2">Correct Answer:</h5>
            <p class="text-gray-700 dark:text-gray-300">{item.correct_answer}</p>
            {#if item.explanation}
                <h5 class="text-sm font-medium text-gray-900 dark:text-white mt-2 mb-2">Explanation:</h5>
                <p class="text-gray-700 dark:text-gray-300">{item.explanation}</p>
            {/if}
        </div>
    {/if}
</div> 