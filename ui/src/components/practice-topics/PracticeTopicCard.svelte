<script lang="ts">
    import type { PracticeTopic } from '$lib/types';
    import { goto } from '$app/navigation';

    export let topic: PracticeTopic;
    export let href: string = '';

    function handleCardClick() {
        goto(href);
    }
    
    function handleEditClick(e: Event) {
        e.stopPropagation();
        goto(`/practice-topics/edit/${topic.id}`);
    }
</script>

<div
    class="relative bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow cursor-pointer"
    on:click={handleCardClick}
    on:keydown={(e) => e.key === 'Enter' && handleCardClick()}
    tabindex="0"
    role="button"
    aria-label={`View ${topic.name}`}
>
    <!-- Edit button positioned in the top-right corner -->
    <button 
        class="absolute top-2 right-2 p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-full z-10" 
        on:click={handleEditClick}
        aria-label="Edit topic"
    >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
        </svg>
    </button>

    <h2 class="text-xl font-semibold text-gray-900 mb-2">{topic.name}</h2>
    
    <!-- Keep consistent spacing with min-height even when no description -->
    <div class="text-gray-600 mb-4 min-h-[1.5rem]">
        {#if topic.description}
            <p>{topic.description}</p>
        {/if}
    </div>
    
    <div class="flex flex-wrap gap-2">
        <span class="bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded">
            {topic.subject}
        </span>
        
        {#if topic.target_age_range}
            <span class="bg-green-100 text-green-800 text-xs font-medium px-2.5 py-0.5 rounded">
                Age: {topic.target_age_range}
            </span>
        {/if}
        
        {#if topic.target_grade_level}
            <span class="bg-purple-100 text-purple-800 text-xs font-medium px-2.5 py-0.5 rounded">
                Grade: {topic.target_grade_level}
            </span>
        {/if}

        {#if topic.tags && Array.isArray(topic.tags) && topic.tags.length > 0}
            <div class="w-full mt-2 pt-2 border-t border-gray-50">
                <div class="flex flex-wrap gap-1">
                    {#each topic.tags as tag}
                        <span class="bg-gray-100 text-gray-800 text-xs font-medium px-2 py-0.5 rounded">
                            {tag}
                        </span>
                    {/each}
                </div>
            </div>
        {/if}
    </div>
</div> 