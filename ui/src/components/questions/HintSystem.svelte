<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    import { createEventDispatcher } from 'svelte';
    import { fade } from 'svelte/transition';

    export let item: PracticeItem;
    export let disabled = false;
    export let minimal = false;

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
    {#if minimal}
        <div class="pt-1" transition:fade={{ duration: 300 }}>
            <div class="flex items-start">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-500 dark:text-indigo-400 mt-0.5 mr-2 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <div class="flex-1">
                    <div class="flex justify-between items-center mb-1.5">
                        <h5 class="text-sm font-medium text-indigo-900 dark:text-indigo-200">Hint {currentHintLevel}/{hints.length}</h5>
                        {#if hasMoreHints}
                            <button
                                class="text-sm text-indigo-600 hover:text-indigo-700 dark:text-indigo-400 dark:hover:text-indigo-300 font-medium"
                                on:click={requestHint}
                                disabled={disabled}
                            >
                                Next Hint
                            </button>
                        {/if}
                    </div>
                    <p class="text-sm text-indigo-700 dark:text-indigo-300">{hints[currentHintLevel - 1]}</p>
                </div>
            </div>
        </div>
    {:else}
        <div class="bg-indigo-50 dark:bg-indigo-900/30 rounded-md border border-indigo-100 dark:border-indigo-800 shadow-sm" transition:fade={{ duration: 300 }}>
            <div class="p-3">
                <div class="flex items-start">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-500 dark:text-indigo-400 mt-0.5 mr-2 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <div class="flex-1">
                        <div class="flex justify-between items-center mb-1.5">
                            <h5 class="text-sm font-medium text-indigo-900 dark:text-indigo-200">Hint {currentHintLevel}/{hints.length}</h5>
                            {#if hasMoreHints}
                                <button
                                    class="text-sm text-indigo-600 hover:text-indigo-700 dark:text-indigo-400 dark:hover:text-indigo-300 font-medium"
                                    on:click={requestHint}
                                    disabled={disabled}
                                >
                                    Next Hint
                                </button>
                            {/if}
                        </div>
                        <p class="text-sm text-indigo-700 dark:text-indigo-300">{hints[currentHintLevel - 1]}</p>
                    </div>
                </div>
            </div>
        </div>
    {/if}
{/if} 