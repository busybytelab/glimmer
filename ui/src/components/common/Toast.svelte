<script lang="ts">
    import { toast } from '$lib/stores/toast';
    import { fade, fly } from 'svelte/transition';
    import { onMount } from 'svelte';

    let container: HTMLDivElement;

    onMount(() => {
        // Ensure the container is positioned correctly
        if (container) {
            container.style.position = 'fixed';
            container.style.top = '1rem';
            container.style.right = '1rem';
            container.style.zIndex = '50';
        }
    });
</script>

<div
    bind:this={container}
    class="fixed top-4 right-4 z-50 flex flex-col gap-2"
>
    {#each $toast as notification (notification.id)}
        <div
            class="min-w-[300px] max-w-md rounded-lg shadow-lg p-4 text-white transform transition-all duration-300 ease-in-out"
            class:bg-green-500={notification.type === 'success'}
            class:bg-red-500={notification.type === 'error'}
            class:bg-blue-500={notification.type === 'info'}
            class:bg-yellow-500={notification.type === 'warning'}
            in:fly={{ y: -20, duration: 300 }}
            out:fade={{ duration: 200 }}
        >
            <div class="flex items-start justify-between">
                <p class="text-sm font-medium">{notification.message}</p>
                <button
                    class="ml-4 text-white hover:text-gray-200 focus:outline-none"
                    on:click={() => toast.remove(notification.id)}
                >
                    <span class="sr-only">Close</span>
                    <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                        <path
                            fill-rule="evenodd"
                            d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                            clip-rule="evenodd"
                        />
                    </svg>
                </button>
            </div>
        </div>
    {/each}
</div> 