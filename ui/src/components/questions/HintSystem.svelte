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
    function requestPreviousHint() {
        if (disabled || currentHintLevel <= 1) return;
        dispatch('hintRequested', { level: currentHintLevel - 1 });
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
                        <div class="flex items-center space-x-2">
                            {#if currentHintLevel > 1}
                                <button
                                    class="p-1 rounded-full text-indigo-600 hover:bg-indigo-100"
                                    on:click={requestPreviousHint}
                                    disabled={disabled}
                                    aria-label="Previous hint"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                                    </svg>
                                </button>
                            {:else}
                                <div class="w-6 h-6"></div>
                            {/if}
                            <h5 class="text-sm font-medium text-indigo-900 dark:text-indigo-200">Hint {currentHintLevel}/{hints.length}</h5>
                            {#if hasMoreHints}
                                <button
                                    class="p-1 rounded-full text-indigo-600 hover:bg-indigo-100"
                                    on:click={requestHint}
                                    disabled={disabled}
                                    aria-label="Next hint"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                                    </svg>
                                </button>
                            {:else}
                                <div class="w-6 h-6"></div>
                            {/if}
                        </div>
                    </div>
                    <p class="text-sm text-indigo-700 dark:text-indigo-300 ml-7">{hints[currentHintLevel - 1]}</p>
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
                        <div class="flex items-center space-x-2 mb-1.5">
                            {#if currentHintLevel > 1}
                                <button
                                    class="p-1 rounded-full text-indigo-600 hover:bg-indigo-100"
                                    on:click={requestPreviousHint}
                                    disabled={disabled}
                                    aria-label="Previous hint"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                                    </svg>
                                </button>
                            {:else}
                                <div class="w-7 h-7"></div>
                            {/if}
                            <h5 class="text-sm font-medium text-indigo-900 dark:text-indigo-200">Hint {currentHintLevel}/{hints.length}</h5>
                            {#if hasMoreHints}
                                <button
                                    class="p-1 rounded-full text-indigo-600 hover:bg-indigo-100"
                                    on:click={requestHint}
                                    disabled={disabled}
                                    aria-label="Next hint"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                                    </svg>
                                </button>
                            {:else}
                                <div class="w-7 h-7"></div>
                            {/if}
                        </div>
                        <p class="text-sm text-indigo-700 dark:text-indigo-300">{hints[currentHintLevel - 1]}</p>
                    </div>
                </div>
            </div>
        </div>
    {/if}
{/if} 