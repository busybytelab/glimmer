<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    import MultipleChoiceQuestion from './MultipleChoiceQuestion.svelte';
    import TrueFalseQuestion from './TrueFalseQuestion.svelte';
    import ShortAnswerQuestion from './ShortAnswerQuestion.svelte';
    import FillInBlankQuestion from './FillInBlankQuestion.svelte';

    export let item: PracticeItem;
    export let index: number;
    export let viewType: 'learner' | 'answered' | 'instructor' = 'learner';
    export let disabled = false;
    export let onAnswerChange: ((answer: string) => void) | undefined = undefined;
    export let printMode = false;

    $: showAnswer = viewType === 'answered' || viewType === 'instructor';
    $: showInstructorInfo = viewType === 'instructor';
    $: isDisabled = disabled || viewType !== 'learner';
</script>
{#if item.question_type === 'multiple_choice'}
    <MultipleChoiceQuestion
        {item}
        {index}
        disabled={isDisabled}
        showAnswer={showAnswer}
        showInstructorInfo={showInstructorInfo}
        {onAnswerChange}
        {printMode}
    />
{:else if item.question_type === 'true_false'}
    <TrueFalseQuestion
        {item}
        {index}
        disabled={isDisabled}
        showAnswer={showAnswer}
        showInstructorInfo={showInstructorInfo}
        {onAnswerChange}
        {printMode}
    />
{:else if item.question_type === 'short_answer'}
    <ShortAnswerQuestion
        {item}
        {index}
        disabled={isDisabled}
        showAnswer={showAnswer}
        showInstructorInfo={showInstructorInfo}
        {onAnswerChange}
        {printMode}
    />
{:else if item.question_type === 'fill_in_blank'}
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
    <div class="border border-red-200 rounded-lg p-4 bg-red-50">
        <p class="text-red-600">Unsupported question type: {item.question_type}, id: {item.id}</p>
    </div>
{/if}