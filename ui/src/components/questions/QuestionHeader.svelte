<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    import { page } from '$app/stores';
    import { fade } from 'svelte/transition';
    import { writable } from 'svelte/store';
    import { onMount } from 'svelte';

    export let item: PracticeItem;
    export let index: number;

    const showCopyTooltip = writable(false);
    let tooltipTimeout: NodeJS.Timeout;

    function getQuestionUrl() {
        const baseUrl = window.location.origin;
        const path = $page.url.pathname;
        return `${baseUrl}${path}#${item.id}`;
    }

    async function copyToClipboard() {
        const url = getQuestionUrl();
        await navigator.clipboard.writeText(url);
        showCopyTooltip.set(true);
        
        if (tooltipTimeout) clearTimeout(tooltipTimeout);
        tooltipTimeout = setTimeout(() => {
            showCopyTooltip.set(false);
        }, 2000);
    }

    onMount(() => {
        // Check if the current hash matches this question's ID
        if (window.location.hash === `#${item.id}`) {
            // Use setTimeout to ensure the DOM is fully rendered
            setTimeout(() => {
                const element = document.getElementById(item.id);
                if (element) {
                    element.scrollIntoView({ behavior: 'smooth', block: 'start' });
                }
            }, 100);
        }
    });
</script>

<div id={item.id} class="flex items-center gap-2 group scroll-mt-20">
    <h4 class="text-md font-medium text-gray-900 dark:text-white mb-2">
        Question {index + 1}
    </h4>
    
    <button
        class="opacity-0 group-hover:opacity-100 transition-opacity duration-200 p-1 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md"
        on:click={copyToClipboard}
        aria-label="Copy question link"
    >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
        </svg>
    </button>

    {#if $showCopyTooltip}
        <div
            class="absolute bg-gray-900 text-white text-xs rounded py-1 px-2 -mt-8"
            transition:fade
        >
            Link copied!
        </div>
    {/if}
</div> 