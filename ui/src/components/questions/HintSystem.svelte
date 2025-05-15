<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    import { createEventDispatcher } from 'svelte';
    import { fade } from 'svelte/transition';

    export let item: PracticeItem;
    export let disabled = false;

    const dispatch = createEventDispatcher<{
        hintRequested: { level: number };
    }>();

    // Parse hints if they exist
    $: hints = item.hints ? (typeof item.hints === 'string' ? JSON.parse(item.hints) : item.hints) : [];
    $: currentHintLevel = item.hint_level_reached || 0;
    $: hasMoreHints = currentHintLevel < hints.length;

    function requestHint() {
        if (disabled || !hasMoreHints) return;
        dispatch('hintRequested', { level: currentHintLevel + 1 });
    }
</script>

{#if hints.length > 0}
    <div class="mt-2" transition:fade={{ duration: 300 }}>
        <button
            class="inline-flex items-center px-3 py-1.5 text-sm text-indigo-600 bg-indigo-50 hover:bg-indigo-100 dark:text-indigo-300 dark:bg-indigo-900/30 dark:hover:bg-indigo-800/50 rounded-md disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
            on:click={requestHint}
            disabled={disabled || !hasMoreHints}
            title={hasMoreHints ? 'Get a hint' : 'No more hints available'}
        >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {#if currentHintLevel === 0}
                Need a hint?
            {:else}
                Hint {currentHintLevel}/{hints.length}
            {/if}
        </button>

        {#if currentHintLevel > 0}
            <div class="mt-3 p-4 bg-indigo-50 dark:bg-indigo-900/30 rounded-md border border-indigo-100 dark:border-indigo-800" transition:fade={{ duration: 300 }}>
                <h5 class="text-sm font-medium text-indigo-900 dark:text-indigo-200 mb-1">Hint {currentHintLevel}:</h5>
                <p class="text-sm text-indigo-700 dark:text-indigo-300">{hints[currentHintLevel - 1]}</p>
            </div>
        {/if}
    </div>
{/if} 