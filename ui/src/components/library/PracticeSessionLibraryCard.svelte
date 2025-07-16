<script lang="ts">
    import type { PracticeSessionLibrary } from '$lib/types';
    import PopularityBadge from '$components/common/PopularityBadge.svelte';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';

    /**
     * The session data to display
     */
    export let session: PracticeSessionLibrary;

    /**
     * Whether the session is currently being imported
     * @default false
     */
    export let isImporting: boolean = false;

    /**
     * Whether the add button should be disabled
     * @default false
     */
    export let disabled: boolean = false;

    /**
     * The selected learner for contextual button labels
     */
    export let selectedLearner: { nickname: string } | undefined = undefined;

    /**
     * Optional callback function when the Add button is clicked
     * @param session - The session object to add
     */
    export let onAdd: ((session: PracticeSessionLibrary) => void) | undefined = undefined;

    // Contextual button label
    $: buttonLabel = 'Add';
</script>

<div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow bg-white dark:bg-gray-800">
    <div class="flex items-start justify-between mb-2">
        <h3 class="font-medium text-gray-900 dark:text-white text-sm">{session.name}</h3>
        {#if session.total_usage}
            <PopularityBadge count={session.total_usage} size="sm" />
        {/if}
    </div>
    
    {#if session.description}
        <p class="text-xs text-gray-600 dark:text-gray-400 mb-3 line-clamp-2">{session.description}</p>
    {/if}
    
    <div class="flex flex-wrap gap-1 mb-3">
        {#if session.expand?.practice_topic_library?.category}
            <span class="text-xs bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-300 px-2 py-1 rounded">
                {session.expand.practice_topic_library.category}
            </span>
        {/if}
        {#if session.target_year}
            <span class="text-xs bg-orange-100 dark:bg-orange-900/30 text-orange-800 dark:text-orange-300 px-2 py-1 rounded">
                Year {session.target_year}
            </span>
        {/if}
    </div>

    <div class="flex justify-end">
        <button 
            class="text-xs bg-indigo-600 dark:bg-indigo-500 text-white hover:bg-indigo-700 dark:hover:bg-indigo-600 px-3 py-1 rounded font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            on:click={() => onAdd?.(session)}
            {disabled}
            title={selectedLearner ? `Add this activity` : 'Select a child first'}
        >
            {#if isImporting}
                <LoadingSpinner size="sm" color="primary" />
            {:else}
                {buttonLabel}
            {/if}
        </button>
    </div>
</div> 