<script lang="ts">
    import type { Learner } from '$lib/types';
    import { IconTypeMap } from '$lib/types';
    
    export let learner: Learner;
    export let actions: Array<{
        label: string;
        color?: 'primary' | 'secondary' | 'success' | 'danger' | 'warning';
        onClick: (learner: Learner) => void;
    }> = [];
    export let clickable: boolean = false;
    export let onClick: (learner: Learner) => void = () => {};
    export let shadow: boolean = true;
    export let showPreferences: boolean = true;
    export let onEdit: ((learner: Learner) => void) | undefined = undefined;
    export let isSelected: boolean = false;
    
    // Helper function to get button color classes
    function getColorClasses(color: string = 'primary') {
        switch(color) {
            case 'secondary':
                return 'bg-gray-100 hover:bg-gray-200 text-gray-800 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-gray-200';
            case 'success':
                return 'bg-green-100 hover:bg-green-200 text-green-800 dark:bg-green-900/30 dark:hover:bg-green-900/50 dark:text-green-300';
            case 'danger':
                return 'bg-red-100 hover:bg-red-200 text-red-800 dark:bg-red-900/30 dark:hover:bg-red-900/50 dark:text-red-300';
            case 'warning':
                return 'bg-yellow-100 hover:bg-yellow-200 text-yellow-800 dark:bg-yellow-900/30 dark:hover:bg-yellow-900/50 dark:text-yellow-300';
            case 'primary':
            default:
                return 'bg-blue-100 hover:bg-blue-200 text-blue-800 dark:bg-blue-900/30 dark:hover:bg-blue-900/50 dark:text-blue-300';
        }
    }

    function handleEditClick(e: Event) {
        e.stopPropagation();
        if (onEdit) onEdit(learner);
    }

    // Get the user's name safely
    $: userName = learner?.nickname || 'Unknown learner';
</script>

{#if clickable}
<div 
    class={`
        relative bg-white dark:bg-gray-800 rounded-lg w-full text-left
        ${shadow ? 'shadow-md' : 'border border-gray-100 dark:border-gray-700'} 
        p-6 
        hover:shadow-lg transition-shadow cursor-pointer
        ${isSelected ? 'ring-2 ring-indigo-500 bg-indigo-50 dark:bg-indigo-900/20' : ''}
    `}
    on:click={() => onClick(learner)}
    on:keydown={(e) => e.key === 'Enter' && onClick(learner)}
    tabindex="0"
    role="button"
    aria-label={`Select ${userName}`}
>
    {#if onEdit}
        <button 
            class="absolute top-2 right-2 p-2 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md transition-colors duration-200 ease-in-out" 
            on:click={handleEditClick}
            aria-label="Edit learner"
        >
            <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d={IconTypeMap.edit} />
            </svg>
        </button>
    {/if}

    <div class="flex items-center mb-4">
        <div class="bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300 p-3 rounded-full mr-4">
            {#if learner.avatar}
                <img 
                    src={learner.avatar} 
                    alt={userName} 
                    class="h-8 w-8 rounded-full object-cover"
                />
            {:else}
                <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d={IconTypeMap.user} />
                </svg>
            {/if}
        </div>
        <div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{userName}</h3>
        </div>
    </div>
    
    <div class="flex flex-wrap gap-2 mb-3">
        {#if learner.age}
            <span class="bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300 text-xs font-medium px-2.5 py-0.5 rounded">
                Age: {learner.age}
            </span>
        {/if}
        
        {#if learner.grade_level}
            <span class="bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-300 text-xs font-medium px-2.5 py-0.5 rounded">
                Grade: {learner.grade_level}
            </span>
        {/if}
    </div>
    
    {#if showPreferences && learner.learning_preferences && learner.learning_preferences.length > 0}
        <div class="mt-3 mb-4">
            <p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Learning preferences:</p>
            <div class="flex flex-wrap gap-1">
                <span class="bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300 text-xs font-medium px-2 py-0.5 rounded">
                    {learner.learning_preferences}
                </span>
            </div>
        </div>
    {/if}
    
    {#if actions.length > 0}
        <div class="flex flex-wrap gap-2 mt-3">
            {#each actions as action}
                <div
                    role="button"
                    tabindex="0"
                    class={`px-3 py-1.5 text-xs font-medium rounded-md cursor-pointer ${getColorClasses(action.color)}`}
                    on:click|stopPropagation={(e) => {
                        e.preventDefault();
                        action.onClick(learner);
                    }}
                    on:keypress|stopPropagation={(e) => {
                        if (e.key === 'Enter' || e.key === ' ') {
                            e.preventDefault();
                            action.onClick(learner);
                        }
                    }}
                >
                    {action.label}
                </div>
            {/each}
        </div>
    {/if}
</div>
{:else}
<div 
    class={`
        relative bg-white dark:bg-gray-800 rounded-lg
        ${shadow ? 'shadow-md' : 'border border-gray-100 dark:border-gray-700'} 
        p-6 
        ${isSelected ? 'ring-2 ring-indigo-500 bg-indigo-50 dark:bg-indigo-900/20' : ''}
    `}
>
    {#if onEdit}
        <button 
            class="absolute top-2 right-2 p-2 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md transition-colors duration-200 ease-in-out" 
            on:click={handleEditClick}
            aria-label="Edit learner"
        >
            <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d={IconTypeMap.edit} />
            </svg>
        </button>
    {/if}

    <div class="flex items-center mb-4">
        <div class="bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300 p-3 rounded-full mr-4">
            {#if learner.avatar}
                <img 
                    src={learner.avatar} 
                    alt={userName} 
                    class="h-8 w-8 rounded-full object-cover"
                />
            {:else}
                <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d={IconTypeMap.user} />
                </svg>
            {/if}
        </div>
        <div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{userName}</h3>
        </div>
    </div>
    
    <div class="flex flex-wrap gap-2 mb-3">
        {#if learner.age}
            <span class="bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300 text-xs font-medium px-2.5 py-0.5 rounded">
                Age: {learner.age}
            </span>
        {/if}
        
        {#if learner.grade_level}
            <span class="bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-300 text-xs font-medium px-2.5 py-0.5 rounded">
                Grade: {learner.grade_level}
            </span>
        {/if}
    </div>
    
    {#if showPreferences && learner.learning_preferences && learner.learning_preferences.length > 0}
        <div class="mt-3 mb-4">
            <p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Learning preferences:</p>
            <div class="flex flex-wrap gap-1">
                {#each learner.learning_preferences as pref}
                    <span class="bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300 text-xs font-medium px-2 py-0.5 rounded">
                        {pref}
                    </span>
                {/each}
            </div>
        </div>
    {/if}
    
    {#if actions.length > 0}
        <div class="flex flex-wrap gap-2 mt-3">
            {#each actions as action}
                <div
                    role="button"
                    tabindex="0"
                    class={`px-3 py-1.5 text-xs font-medium rounded-md cursor-pointer ${getColorClasses(action.color)}`}
                    on:click|stopPropagation={(e) => {
                        e.preventDefault();
                        action.onClick(learner);
                    }}
                    on:keypress|stopPropagation={(e) => {
                        if (e.key === 'Enter' || e.key === ' ') {
                            e.preventDefault();
                            action.onClick(learner);
                        }
                    }}
                >
                    {action.label}
                </div>
            {/each}
        </div>
    {/if}
</div>
{/if} 