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

<div class="border border-gray-200 rounded-lg p-4">
    <QuestionHeader {index} />
    <p class="text-gray-600 mb-4">{item.question_text}</p>
    
    <div class="space-y-2">
        {#each (typeof item.options === 'string' ? JSON.parse(item.options) : item.options) as option, optionIndex}
            <div class="flex items-center">
                {#if printMode}
                    <div class="w-4 h-4 border border-gray-400 rounded mr-2"></div>
                {:else}
                    <input
                        type="radio"
                        id={`question-${index}-option-${optionIndex}-${item.id}`}
                        name={`question-${index}-${item.id}`}
                        value={option}
                        class="h-4 w-4 text-indigo-600 focus:ring-indigo-500"
                        {disabled}
                        on:change={handleAnswerChange}
                    />
                {/if}
                <label 
                    for={printMode ? undefined : `question-${index}-option-${optionIndex}-${item.id}`}
                    class="ml-3 text-gray-700"
                >
                    {option}
                </label>
            </div>
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