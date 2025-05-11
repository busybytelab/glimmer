<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    import { createEventDispatcher } from 'svelte';

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
    <div class="mt-2">
        <button
            class="inline-flex items-center px-2 py-1 text-sm text-indigo-600 hover:text-indigo-800 disabled:opacity-50 disabled:cursor-not-allowed"
            on:click={requestHint}
            disabled={disabled || !hasMoreHints}
            title={hasMoreHints ? 'Get a hint' : 'No more hints available'}
        >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Hint {currentHintLevel + 1}/{hints.length}
        </button>

        {#if currentHintLevel > 0}
            <div class="mt-2 p-3 bg-indigo-50 rounded-md">
                <h5 class="text-sm font-medium text-indigo-900 mb-1">Current Hint:</h5>
                <p class="text-sm text-indigo-700">{hints[currentHintLevel - 1]}</p>
            </div>
        {/if}
    </div>
{/if} 