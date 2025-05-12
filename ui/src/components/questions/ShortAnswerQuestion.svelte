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

    let textarea: HTMLTextAreaElement;
    let lastFocused = false;

    function handleInput() {
        if (onAnswerChange) {
            lastFocused = true;
            onAnswerChange(item.user_answer || '');
        }
    }

    // Restore focus after any updates
    $: if (lastFocused && textarea) {
        textarea.focus();
        lastFocused = false;
    }
</script>

<div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 bg-white dark:bg-gray-800">
    <QuestionHeader {item} {index} />
    <p class="text-gray-700 dark:text-gray-300 mb-4">{item.question_text}</p>
    
    {#if printMode}
        <div class="h-32 border border-gray-300 dark:border-gray-600 rounded-md"></div>
    {:else}
        <textarea
            bind:this={textarea}
            id={`question-${index}-${item.id}`}
            class="w-full h-32 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white"
            placeholder="Type your answer here..."
            {disabled}
            bind:value={item.user_answer}
            on:input={handleInput}
        ></textarea>
    {/if}

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
            {#if item.hint_level_reached}
                <h5 class="text-sm font-medium text-gray-900 dark:text-white mt-2 mb-2">Hints Used:</h5>
                <p class="text-gray-700 dark:text-gray-300">{item.hint_level_reached}</p>
            {/if}
        </div>
    {/if}
</div> 