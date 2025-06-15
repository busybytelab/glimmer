<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    import { QuestionViewType } from '$lib/types';
    import QuestionFactory from '../questions/QuestionFactory.svelte';
    import VerticalProgressStepper from '../common/VerticalProgressStepper.svelte';

    export let practiceItems: PracticeItem[] = [];
    export let currentStep: number = 0;
    export let stepResults: ('correct' | 'incorrect' | 'pending')[] = [];
    export let selectedViewType: QuestionViewType;
    export let sessionStatus: string;
    export let savingItems: Set<number>;

    export let onStepClick: (index: number) => void;
    export let onAnswerChange: (index: number, answer: string) => void;
    export let onHintRequest: (index: number, level: number) => void;

    $: currentItem = practiceItems[currentStep];
</script>

<div class="flex space-x-6">
    <div class="w-12">
        <VerticalProgressStepper
            steps={practiceItems.length}
            {currentStep}
            results={stepResults}
            {onStepClick}
        />
    </div>
    <div class="flex-1">
        {#if currentItem}
            <div class="question-container">
                <QuestionFactory
                    item={currentItem}
                    index={currentStep}
                    viewType={selectedViewType}
                    disabled={selectedViewType !== QuestionViewType.LEARNER || sessionStatus === 'Completed' || savingItems.has(currentStep)}
                    onAnswerChange={(answer) => onAnswerChange(currentStep, answer)}
                    isInstructor={false}
                    onHintRequested={(level) => onHintRequest(currentStep, level)}
                />
            </div>
        {/if}
    </div>
</div> 