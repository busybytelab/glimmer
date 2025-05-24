<script lang="ts">
    import type { PracticeItem, ReviewStatus } from '$lib/types';
    import { QuestionViewType, QuestionType } from '$lib/types';
    import MultipleChoiceQuestion from './MultipleChoiceQuestion.svelte';
    import TrueFalseQuestion from './TrueFalseQuestion.svelte';
    import ShortAnswerQuestion from './ShortAnswerQuestion.svelte';
    import FillInBlankQuestion from './FillInBlankQuestion.svelte';
    import HintSystem from './HintSystem.svelte';
    import QuestionReviewControls from './QuestionReviewControls.svelte';

    export let item: PracticeItem;
    export let index: number;
    export let viewType: QuestionViewType = QuestionViewType.LEARNER;
    export let disabled = false;
    export let onAnswerChange: ((answer: string) => void) | undefined = undefined;
    export let printMode = false;
    export let isInstructor: boolean = false;
    export let showHints: boolean = false;
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
</script>

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

{#if showHints && viewType === QuestionViewType.LEARNER && item.hints && item.hints.length > 0 && !isInstructor}
    <div class="mt-4">
        <HintSystem 
            {item} 
            disabled={isDisabled}
            on:hintRequested={handleHintRequested} 
        />
    </div>
{/if}

{#if showReviewControls && onReviewStatusChange}
    <QuestionReviewControls
        {item}
        onReviewStatusChange={onReviewStatusChange}
    />
{/if}