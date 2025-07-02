<script lang="ts">
    import type { PracticeTopicLibrary } from '$lib/types';
    import PopularityBadge from '$components/common/PopularityBadge.svelte';

    /**
     * The topic data to display
     */
    export let topic: PracticeTopicLibrary;
    
    /**
     * Whether this topic is currently selected
     * @default false
     */
    export let isSelected: boolean = false;
    
    /**
     * Callback function when the topic card is clicked
     * @param topicId - The ID of the clicked topic
     * @param topic - The topic object that was clicked
     */
    export let onClick: (topicId: string, topic: PracticeTopicLibrary) => void;

    /**
     * Optional callback function when the Add button is clicked
     * @param topic - The topic object to add
     */
    export let onAdd: ((topic: PracticeTopicLibrary) => void) | undefined = undefined;

    function handleClick() {
        onClick(topic.id, topic);
    }

    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter' || event.key === ' ') {
            event.preventDefault();
            handleClick();
        }
    }

    function handleAddClick(event: MouseEvent) {
        event.stopPropagation(); // Prevent card click from firing
        if (onAdd) {
            onAdd(topic);
        }
    }
</script>

<div 
    class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow cursor-pointer {isSelected ? 'ring-2 ring-indigo-500 bg-indigo-50 dark:bg-indigo-900/20' : 'bg-white dark:bg-gray-800'}"
    on:click={handleClick}
    role="button"
    tabindex="0"
    on:keydown={handleKeyDown}
>
    <div class="flex items-start justify-between mb-2">
        <h3 class="font-medium text-gray-900 dark:text-white text-sm">{topic.name}</h3>
        {#if topic.total_usage}
            <PopularityBadge count={topic.total_usage} size="sm" />
        {/if}
    </div>
    
    {#if topic.description}
        <p class="text-xs text-gray-600 dark:text-gray-400 mb-3 line-clamp-2">{topic.description}</p>
    {/if}
    
    <div class="flex flex-wrap gap-1 mb-3">
        {#if topic.category}
            <span title="Subject" class="text-xs bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300 px-2 py-1 rounded">
                {topic.category}
            </span>
        {/if}
        {#if topic.target_grade_level}
            <span title="Grade" class="text-xs bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300 px-2 py-1 rounded">
                {topic.target_grade_level}
            </span>
        {/if}
    </div>

    {#if onAdd}
        <div class="flex justify-end">
            <button
                on:click={handleAddClick}
                title="Add to my topics"
                class="text-xs bg-indigo-600 dark:bg-indigo-500 text-white hover:bg-indigo-700 dark:hover:bg-indigo-600 px-3 py-1 rounded font-medium transition-colors"
            >
                Add
            </button>
        </div>
    {/if}
</div>

<style>
    .line-clamp-2 {
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        line-clamp: 2;
        overflow: hidden;
    }
</style> 