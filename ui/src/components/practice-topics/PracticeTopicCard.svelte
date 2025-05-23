<script lang="ts">
    import type { PracticeTopic } from '$lib/types';
    import { goto } from '$app/navigation';

    export let topic: PracticeTopic;
    export let href: string = '';
    export let isInstructor: boolean = false;

    function handleCardClick() {
        goto(href);
    }
    
    function handleEditClick(e: Event) {
        e.stopPropagation();
        goto(`/practice-topics/${topic.id}/edit`);
    }
</script>

<div
    class="relative bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow cursor-pointer"
    on:click={handleCardClick}
    on:keydown={(e) => e.key === 'Enter' && handleCardClick()}
    tabindex="0"
    role="button"
    aria-label={`View ${topic.name}`}
>
    <!-- Edit button positioned in the top-right corner -->
    {#if isInstructor}
        <button 
            class="absolute top-2 right-2 p-2 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full z-10" 
            on:click={handleEditClick}
            aria-label="Edit topic"
        >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
            </svg>
        </button>
    {/if}

    <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">{topic.name}</h2>
    
    <!-- Keep consistent spacing with min-height even when no description -->
    <div class="text-gray-600 dark:text-gray-300 mb-4 min-h-[1.5rem]">
        {#if topic.description}
            <p>{topic.description}</p>
        {/if}
    </div>
    
    <div class="flex flex-wrap gap-2">
        <span class="bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-300 text-xs font-medium px-2.5 py-0.5 rounded">
            {topic.subject}
        </span>
        
        {#if topic.target_age_range}
            <span class="bg-green-100 dark:bg-green-900/40 text-green-800 dark:text-green-300 text-xs font-medium px-2.5 py-0.5 rounded">
                Age: {topic.target_age_range}
            </span>
        {/if}
        
        {#if topic.target_grade_level}
            <span class="bg-purple-100 dark:bg-purple-900/40 text-purple-800 dark:text-purple-300 text-xs font-medium px-2.5 py-0.5 rounded">
                Grade: {topic.target_grade_level}
            </span>
        {/if}

        {#if topic.tags && Array.isArray(topic.tags) && topic.tags.length > 0}
            <div class="w-full mt-2 pt-2 border-t border-gray-50 dark:border-gray-700">
                <div class="flex flex-wrap gap-1">
                    {#each topic.tags as tag}
                        <span class="bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300 text-xs font-medium px-2 py-0.5 rounded">
                            {tag}
                        </span>
                    {/each}
                </div>
            </div>
        {/if}
    </div>
</div> 