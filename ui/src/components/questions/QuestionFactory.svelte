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
    export let onHintRequested: ((level: number) => void) | undefined = undefined;
    export let onReviewStatusChange: ((itemId: string, status: ReviewStatus) => void) | undefined = undefined;
    export let sessionId: string | undefined = undefined;

    // Derived values based on view type
    $: showAnswer = viewType === QuestionViewType.ANSWERED || viewType === QuestionViewType.PARENT || viewType === QuestionViewType.GENERATED;
    $: showInstructorInfo = viewType === QuestionViewType.PARENT || viewType === QuestionViewType.GENERATED;
    $: showReviewControls = viewType === QuestionViewType.GENERATED && isInstructor;
    
    // Controls whether the question can be interacted with
    // Disable in these cases:
    // 1. Explicitly disabled from parent
    // 2. Not in LEARNER view
    // 3. Parent using LEARNER view (to prevent accidental edits)
    $: isDisabled = disabled || 
                   viewType !== QuestionViewType.LEARNER || 
                   (isInstructor && viewType === QuestionViewType.LEARNER);

    // Show hint system if:
    // 1. Question has hints
    // 2. Not in print mode
    // 3. Not showing answer yet
    // 4. Not showing instructor info
    $: showHint = Boolean(item.hints?.length) && 
                  !printMode && 
                  !showAnswer && 
                  !showInstructorInfo;

    // Show explanation if:
    // 1. Not in print mode
    // 2. Question has been answered
    // 3. Question has an explanation
    $: showExplanation = !printMode && 
                        item.is_correct !== undefined && 
                        (item.explanation || Object.keys(item.explanation_for_incorrect || {}).length > 0);

    function handleHintRequested(event: CustomEvent<{ level: number }>) {
        if (onHintRequested) {
            onHintRequested(event.detail.level);
        }
    }
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg">
    
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
        {#if !showReviewControls}
            <hr class="my-4" />
        {/if}
    {/if}

    {#if showReviewControls && onReviewStatusChange}
        <div class="mt-4" transition:fade={{ duration: 300 }}>
        <QuestionReviewControls
            {item}
            {sessionId}
            onReviewStatusChange={onReviewStatusChange}
        />
        </div>
        <hr class="my-4" />
    {/if}
</div>