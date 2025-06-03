<script lang="ts">
    import type { PracticeItem, ReviewStatus } from '$lib/types';
    import { QuestionViewType, QuestionType } from '$lib/types';
    import MultipleChoiceQuestion from './MultipleChoiceQuestion.svelte';
    import TrueFalseQuestion from './TrueFalseQuestion.svelte';
    import ShortAnswerQuestion from './ShortAnswerQuestion.svelte';
    import FillInBlankQuestion from './FillInBlankQuestion.svelte';
    import HintSystem from './HintSystem.svelte';
    import QuestionReviewControls from './QuestionReviewControls.svelte';
    import QuestionExplanation from './QuestionExplanation.svelte';
    import { fade } from 'svelte/transition';

    export let item: PracticeItem;
    export let index: number;
    export let viewType: QuestionViewType = QuestionViewType.LEARNER;
    export let disabled = false;
    export let onAnswerChange: ((answer: string) => void) | undefined = undefined;
    export let printMode = false;
    export let isInstructor: boolean = false;
    export const showHints: boolean = false;
    export let onHintRequested: ((level: number) => void) | undefined = undefined;
    export let onReviewStatusChange: ((itemId: string, status: ReviewStatus) => void) | undefined = undefined;

    // Derived values based on view type
    $: showAnswer = viewType === QuestionViewType.ANSWERED || viewType === QuestionViewType.INSTRUCTOR || viewType === QuestionViewType.GENERATED;
    $: showInstructorInfo = viewType === QuestionViewType.INSTRUCTOR || viewType === QuestionViewType.GENERATED;
    $: showReviewControls = viewType === QuestionViewType.GENERATED && isInstructor;
    
    // Controls whether the question can be interacted with
    // Disable in these cases:
    // 1. Explicitly disabled from parent
    // 2. Not in LEARNER view
    // 3. Instructor using LEARNER view (to prevent accidental edits)
    $: isDisabled = disabled || 
                   viewType !== QuestionViewType.LEARNER || 
                   (isInstructor && viewType === QuestionViewType.LEARNER);

    // Handle hint requests from HintSystem
    function handleHintRequested(event: CustomEvent<{level: number}>) {
        const { level } = event.detail;
        if (onHintRequested) {
            onHintRequested(level);
        }
    }

    // State for showing/hiding hints and explanations
    let showHint = false;
    let showExplanation = false;

    $: if (viewType === QuestionViewType.LEARNER && item.is_correct) {
        showExplanation = true;
    }
</script>

<div class="relative">
    <div class="absolute top-4 right-4 flex items-center space-x-2">
        {#if viewType === QuestionViewType.LEARNER && item.hints && item.hints.length > 0 && !isInstructor && !item.is_correct}
            <button
                class="group p-1 rounded-full hover:bg-indigo-100 dark:hover:bg-indigo-900 transition-colors duration-200 focus:outline-none"
                style="box-shadow: none; border: none;"
                on:click={() => showHint = !showHint}
                title="Show hint"
                aria-label="Show hint">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-500 group-hover:text-indigo-700 dark:text-indigo-400 dark:group-hover:text-indigo-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
            </button>
        {/if}
        {#if (viewType === QuestionViewType.INSTRUCTOR && item.user_answer) || (viewType === QuestionViewType.LEARNER && item.is_correct)}
            <button
                class="group p-1 rounded-full hover:bg-indigo-100 dark:hover:bg-indigo-900 transition-colors duration-200 focus:outline-none"
                style="box-shadow: none; border: none;"
                on:click={() => showExplanation = !showExplanation}
                title="Explanation"
                aria-label="Show explanation">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-500 group-hover:text-indigo-700 dark:text-indigo-400 dark:group-hover:text-indigo-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
            </button>
        {/if}
    </div>

    {#if item.question_type === QuestionType.MULTIPLE_CHOICE}
        <MultipleChoiceQuestion
            {item}
            {index}
            disabled={isDisabled}
            showAnswer={showAnswer}
            showInstructorInfo={showInstructorInfo}
            {onAnswerChange}
            {printMode}
        />
    {:else if item.question_type === QuestionType.TRUE_FALSE}
        <TrueFalseQuestion
            {item}
            {index}
            disabled={isDisabled}
            showAnswer={showAnswer}
            showInstructorInfo={showInstructorInfo}
            {onAnswerChange}
            {printMode}
        />
    {:else if item.question_type === QuestionType.SHORT_ANSWER}
        <ShortAnswerQuestion
            {item}
            {index}
            disabled={isDisabled}
            showAnswer={showAnswer}
            showInstructorInfo={showInstructorInfo}
            {onAnswerChange}
            {printMode}
        />
    {:else if item.question_type === QuestionType.FILL_IN_BLANK}
        <FillInBlankQuestion
            {item}
            {index}
            disabled={isDisabled}
            showAnswer={showAnswer}
            showInstructorInfo={showInstructorInfo}
            {onAnswerChange}
            {printMode}
        />
    {:else}
        <div class="border border-red-200 dark:border-red-800 rounded-lg p-4 bg-red-50 dark:bg-red-900/30">
            <p class="text-red-600 dark:text-red-400">Unsupported question type: {item.question_type}, id: {item.id}</p>
        </div>
    {/if}

    {#if showHint}
        <div class="mt-4" transition:fade={{ duration: 300 }}>
            <HintSystem 
                {item} 
                disabled={isDisabled}
                on:hintRequested={handleHintRequested} 
                minimal={true}
            />
        </div>
    {/if}

    {#if showExplanation}
        <div class="mt-4" transition:fade={{ duration: 300 }}>
            <QuestionExplanation {item} />
        </div>
    {/if}

    {#if showReviewControls && onReviewStatusChange}
        <QuestionReviewControls
            {item}
            onReviewStatusChange={onReviewStatusChange}
        />
    {/if}
</div>