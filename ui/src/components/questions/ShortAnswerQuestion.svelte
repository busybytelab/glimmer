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
        const target = event.target as HTMLTextAreaElement;
        if (onAnswerChange) {
            onAnswerChange(target.value);
        }
    }
</script>

<div class="border border-gray-200 rounded-lg p-4">
    <QuestionHeader {index} />
    <p class="text-gray-600 mb-4">{item.question_text}</p>
    
    {#if printMode}
        <div class="h-32 border border-gray-300 rounded-md"></div>
    {:else}
        <textarea
            id={`question-${index}-${item.id}`}
            class="w-full h-32 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="Type your answer here..."
            {disabled}
            on:input={handleAnswerChange}
        ></textarea>
    {/if}

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