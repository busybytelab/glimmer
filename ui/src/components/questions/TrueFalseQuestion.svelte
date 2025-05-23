<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    import QuestionHeader from './QuestionHeader.svelte';
    import QuestionAnswerInfo from './QuestionAnswerInfo.svelte';
    import QuestionInstructorInfo from './QuestionInstructorInfo.svelte';
    import SaveButton from './SaveButton.svelte';

    export let item: PracticeItem;
    export let index: number;
    export let disabled = false;
    export let showAnswer = false;
    export let showInstructorInfo = false;
    export let onAnswerChange: ((answer: string) => void) | undefined = undefined;
    export let printMode = false;

    let selectedAnswer: string | null = null;

    function handleAnswerChange(event: Event) {
        const target = event.target as HTMLInputElement;
        selectedAnswer = target.value;
    }

    function handleSave() {
        if (selectedAnswer && onAnswerChange) {
            onAnswerChange(selectedAnswer);
        }
    }
    
    // Helper for immediate feedback
    $: isAnswered = !!item.user_answer;
    $: showFeedback = isAnswered && item.is_correct !== undefined && !showAnswer;
</script>

<div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 bg-white dark:bg-gray-800
            {showFeedback ? (item.is_correct ? 'border-green-300 dark:border-green-600' : 'border-red-300 dark:border-red-600') : ''}">
    <QuestionHeader {item} {index} />
    <p class="text-gray-700 dark:text-gray-300 mb-4">{item.question_text}</p>
    
    <div class="space-y-2">
        {#each ['True', 'False'] as option}
            <div class="flex items-center">
                {#if printMode}
                    <div class="w-4 h-4 border border-gray-400 dark:border-gray-500 rounded mr-2"></div>
                {:else}
                    <input
                        type="radio"
                        id={`question-${index}-option-${option}-${item.id}`}
                        name={`question-${index}-${item.id}`}
                        value={option}
                        class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 dark:border-gray-600 dark:bg-gray-700"
                        {disabled}
                        checked={selectedAnswer === option || item.user_answer === option}
                        on:change={handleAnswerChange}
                    />
                {/if}
                <label 
                    for={printMode ? undefined : `question-${index}-option-${option}-${item.id}`}
                    class="ml-3 {showInstructorInfo && String(option) === String(item.correct_answer) 
                        ? 'text-green-600 dark:text-green-400 font-medium' 
                        : 'text-gray-700 dark:text-gray-300'} {item.user_answer === option 
                        ? 'font-medium border-b border-indigo-300 dark:border-indigo-500' 
                        : ''}"
                >
                    {option} {showInstructorInfo && String(option) === String(item.correct_answer) ? '✓' : ''}
                </label>
            </div>
        {/each}
    </div>
    
    {#if showFeedback}
        <div class="mt-4 text-sm">
            {#if item.is_correct}
                <p class="text-green-600 dark:text-green-400 font-medium">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 inline mr-1" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                    </svg>
                    Correct!
                </p>
            {:else}
                <p class="text-red-600 dark:text-red-400 font-medium">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 inline mr-1" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                    </svg>
                    Try again
                </p>
            {/if}
        </div>
    {/if}

    {#if !disabled && !showAnswer && !printMode}
        <SaveButton
            disabled={!selectedAnswer || selectedAnswer === item.user_answer}
            onClick={handleSave}
        />
    {/if}

    <QuestionAnswerInfo {item} {showAnswer} {showInstructorInfo} />
    <QuestionInstructorInfo {item} {showInstructorInfo} />
</div> 