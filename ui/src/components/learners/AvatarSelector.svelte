<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import robot1 from '../../assets/robot1.svg';
    import bunny1 from '../../assets/bunny1.svg';
    import octopus1 from '../../assets/octopus1.svg';

    export let selectedAvatar: string = '';

    const dispatch = createEventDispatcher<{
        select: string;
    }>();

    const avatars = [
        { src: '/robot1.svg', preview: robot1, alt: 'Robot Avatar' },
        { src: '/bunny1.svg', preview: bunny1, alt: 'Bunny Avatar' },
        { src: '/octopus1.svg', preview: octopus1, alt: 'Octopus Avatar' }
    ];

    function handleSelect(src: string) {
        selectedAvatar = src;
        dispatch('select', src);
    }
</script>

<div class="grid grid-cols-4 sm:grid-cols-6 lg:grid-cols-8 gap-3">
    {#each avatars as avatar}
        <button
            type="button"
            class="relative w-20 h-20 rounded-xl border-2 overflow-hidden transition-all duration-200 hover:scale-105 {selectedAvatar === avatar.src ? 'border-primary ring-2 ring-primary ring-opacity-50' : 'border-gray-200 dark:border-gray-700'}"
            on:click={() => handleSelect(avatar.src)}
        >
            <img
                src={avatar.preview}
                alt={avatar.alt}
                class="w-full h-full object-contain p-2"
            />
        </button>
    {/each}
</div> 