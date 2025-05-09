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

    // Split the question text by the blank placeholder
    $: parts = item.question_text.split('_____');
</script>

<div class="border border-gray-200 rounded-lg p-4">
    <QuestionHeader {index} />
    
    <div class="text-gray-600 mb-4">
        {#each parts as part, i}
            {part}
            {#if i < parts.length - 1}
                {#if printMode}
                    <div class="inline-block w-32 border-b border-gray-400 mx-1"></div>
                {:else}
                    <input
                        type="text"
                        id={`question-${index}-blank-${i}-${item.id}`}
                        class="inline-block w-32 px-2 py-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                        {disabled}
                        on:input={handleAnswerChange}
                    />
                {/if}
            {/if}
        {/each}
    </div>

    {#if showAnswer}
        <div class="mt-4 p-4 bg-gray-50 rounded-md">
            <h5 class="text-sm font-medium text-gray-900 mb-2">Your Answer:</h5>
            <p class="text-gray-600">{item.user_answer || 'Not answered'}</p>
        </div>
    {/if}

    {#if showInstructorInfo}
        <div class="mt-4 p-4 bg-gray-50 rounded-md">
            <h5 class="text-sm font-medium text-gray-900 mb-2">Correct Answer:</h5>
            <p class="text-gray-600">{item.correct_answer}</p>
            {#if item.explanation}
                <h5 class="text-sm font-medium text-gray-900 mt-2 mb-2">Explanation:</h5>
                <p class="text-gray-600">{item.explanation}</p>
            {/if}
            {#if item.hints_used}
                <h5 class="text-sm font-medium text-gray-900 mt-2 mb-2">Hints Used:</h5>
                <p class="text-gray-600">{item.hints_used}</p>
            {/if}
        </div>
    {/if}
</div> 