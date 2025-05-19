<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    import QuestionHeader from './QuestionHeader.svelte';
    import QuestionAnswerInfo from './QuestionAnswerInfo.svelte';
    import QuestionInstructorInfo from './QuestionInstructorInfo.svelte';

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
    
    // Helper for immediate feedback
    $: isAnswered = !!item.user_answer;
    $: showFeedback = isAnswered && item.is_correct !== undefined && !showAnswer;

    // Process question text to find where to place the input
    $: {
        const text = item.question_text;
        if (text.includes('[BLANK]')) {
            questionParts = text.split('[BLANK]');
            isFullTextInput = false;
        } else {
            // Look for common blank indicators
            const blankPattern = /_{3,}|\.{3,}|\s+\.{3,}\s+|\s+\.\.\.\s+/;
            const match = text.match(blankPattern);
            if (match) {
                const parts = text.split(match[0]);
                questionParts = [parts[0], parts[1]];
                isFullTextInput = false;
            } else {
                // For questions without explicit blank indicators (like punctuation questions),
                // use a full text input
                questionParts = [text];
                isFullTextInput = true;
            }
        }
    }

    let questionParts: string[] = [];
    let isFullTextInput = false;
</script>

<div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 bg-white dark:bg-gray-800
            {showFeedback ? (item.is_correct ? 'border-green-300 dark:border-green-600' : 'border-red-300 dark:border-red-600') : ''}">
    <QuestionHeader {item} {index} />
    
    <div class="text-gray-700 dark:text-gray-300 mb-4">
        {#if printMode}
            <p>{item.question_text.replace(/\[BLANK\]|_{3,}|\.{3,}|\s+\.{3,}\s+|\s+\.\.\.\s+/, '___________')}</p>
        {:else}
            {#if isFullTextInput}
                <div class="mb-2">{item.question_text}</div>
                <input
                    type="text"
                    class="w-full px-2 py-1 border border-gray-300 dark:border-gray-600 rounded focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white"
                    value={item.user_answer || ''}
                    {disabled}
                    on:input={handleAnswerChange}
                    placeholder="Type your complete answer here..."
                />
            {:else}
                <div>
                    {#each questionParts as part, i}
                        {part}
                        {#if i < questionParts.length - 1}
                            <input
                                type="text"
                                class="inline-block w-32 mx-1 px-2 py-1 border border-gray-300 dark:border-gray-600 rounded focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white"
                                value={item.user_answer || ''}
                                {disabled}
                                on:input={handleAnswerChange}
                            />
                        {/if}
                    {/each}
                </div>
            {/if}
        {/if}
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

    <QuestionAnswerInfo {item} {showAnswer} {showInstructorInfo} />
    <QuestionInstructorInfo {item} {showInstructorInfo} />
</div> 