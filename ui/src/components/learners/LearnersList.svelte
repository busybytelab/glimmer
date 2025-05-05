<script lang="ts">
    import type { Learner } from '$lib/types';
    import LearnerCard from './LearnerCard.svelte';
    
    export let learners: Learner[] = [];
    export let loading: boolean = false;
    export let emptyMessage: string = 'No learners found.';
    export let gridCols: string = 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3';
    export let showPreferences: boolean = true;
    
    export let onClick: (learner: Learner) => void = () => {};
    export let cardActions: Array<{
        label: string;
        color?: 'primary' | 'secondary' | 'success' | 'danger' | 'warning';
        onClick: (learner: Learner) => void;
        condition?: (learner: Learner) => boolean;
    }> = [];
    
    // Filter actions based on conditions if provided
    function getActionsForLearner(learner: Learner) {
        return cardActions.filter(action => !action.condition || action.condition(learner));
    }
</script>

{#if loading}
    <div class="flex justify-center items-center h-32">
        <div class="animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-primary"></div>
    </div>
{:else if learners.length === 0}
    <div class="bg-yellow-50 border-l-4 border-yellow-400 p-4">
        <div class="flex">
            <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                </svg>
            </div>
            <div class="ml-3">
                <p class="text-sm text-yellow-700">
                    {emptyMessage}
                </p>
            </div>
        </div>
    </div>
{:else}
    <div class={`grid ${gridCols} gap-6`}>
        {#each learners as learner}
            <LearnerCard 
                {learner} 
                clickable={onClick !== undefined}
                {onClick}
                {showPreferences}
                actions={getActionsForLearner(learner)}
            />
        {/each}
    </div>
{/if} 