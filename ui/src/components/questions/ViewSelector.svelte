<script lang="ts">
    import { QuestionViewType } from '$lib/types';
    
    /**
     * Current view type selected
     */
    export let viewType: QuestionViewType;
    
    /**
     * Event handler for when view type changes
     */
    export let onViewChange: (newViewType: QuestionViewType) => void;
    
    /**
     * Whether the user is an instructor (determines available views)
     */
    export let isInstructor: boolean = false;
    
    const viewOptions = [
        { value: QuestionViewType.LEARNER, label: 'Learner View', icon: 'person' },
        { value: QuestionViewType.PARENT, label: 'Parent View', icon: 'school' },
        { value: QuestionViewType.ANSWERED, label: 'Answered View', icon: 'check_circle' },
        { value: QuestionViewType.GENERATED, label: 'Review View', icon: 'rate_review' }
    ];
    
    // Filter options based on user role
    $: availableOptions = isInstructor 
        ? viewOptions 
        : viewOptions.filter(opt => opt.value !== QuestionViewType.PARENT);
</script>

<div class="flex items-center bg-gray-100 dark:bg-gray-700 p-2 rounded-lg mb-4">
    <span class="text-sm text-gray-600 dark:text-gray-300 mr-3">View as:</span>
    <div class="flex space-x-2">
        {#each availableOptions as option}
            <button 
                class="px-3 py-1 text-sm rounded-md transition-colors {viewType === option.value 
                    ? 'bg-indigo-600 text-white' 
                    : 'bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'}"
                on:click={() => onViewChange(option.value)}
                type="button"
            >
                {option.label}
            </button>
        {/each}
    </div>
</div> 