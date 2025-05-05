<script lang="ts">
    import type { Learner } from '$lib/types';
    
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
    
    // Helper function to get button color classes
    function getColorClasses(color: string = 'primary') {
        switch(color) {
            case 'secondary':
                return 'bg-gray-100 hover:bg-gray-200 text-gray-800';
            case 'success':
                return 'bg-green-100 hover:bg-green-200 text-green-800';
            case 'danger':
                return 'bg-red-100 hover:bg-red-200 text-red-800';
            case 'warning':
                return 'bg-yellow-100 hover:bg-yellow-200 text-yellow-800';
            case 'primary':
            default:
                return 'bg-blue-100 hover:bg-blue-200 text-blue-800';
        }
    }
</script>

{#if clickable}
<button 
    class={`
        bg-white rounded-lg w-full text-left
        ${shadow ? 'shadow-md' : 'border border-gray-100'} 
        p-6 
        hover:shadow-lg transition-shadow cursor-pointer
    `}
    on:click={() => onClick(learner)}
    aria-label={`Select ${learner.nickname}`}
>
    <div class="flex items-center mb-4">
        <div class="bg-blue-100 text-blue-800 p-3 rounded-full mr-4">
            {#if learner.avatar}
                <img 
                    src={learner.avatar} 
                    alt={learner.nickname} 
                    class="h-8 w-8 rounded-full object-cover"
                />
            {:else}
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
            {/if}
        </div>
        <div>
            <h3 class="text-lg font-semibold text-gray-900">{learner.nickname}</h3>
            <p class="text-sm text-gray-600">{learner.user?.name || 'Unknown user'}</p>
        </div>
    </div>
    
    <div class="flex flex-wrap gap-2 mb-3">
        {#if learner.age}
            <span class="bg-green-100 text-green-800 text-xs font-medium px-2.5 py-0.5 rounded">
                Age: {learner.age}
            </span>
        {/if}
        
        {#if learner.grade_level}
            <span class="bg-purple-100 text-purple-800 text-xs font-medium px-2.5 py-0.5 rounded">
                Grade: {learner.grade_level}
            </span>
        {/if}
    </div>
    
    {#if showPreferences && learner.learning_preferences && learner.learning_preferences.length > 0}
        <div class="mt-3 mb-4">
            <p class="text-xs text-gray-500 mb-1">Learning preferences:</p>
            <div class="flex flex-wrap gap-1">
                <span class="bg-gray-100 text-gray-800 text-xs font-medium px-2 py-0.5 rounded">
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
</button>
{:else}
<div 
    class={`
        bg-white rounded-lg
        ${shadow ? 'shadow-md' : 'border border-gray-100'} 
        p-6 
    `}
>
    <div class="flex items-center mb-4">
        <div class="bg-blue-100 text-blue-800 p-3 rounded-full mr-4">
            {#if learner.avatar}
                <img 
                    src={learner.avatar} 
                    alt={learner.nickname} 
                    class="h-8 w-8 rounded-full object-cover"
                />
            {:else}
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
            {/if}
        </div>
        <div>
            <h3 class="text-lg font-semibold text-gray-900">{learner.nickname}</h3>
            <p class="text-sm text-gray-600">{learner.user?.name || 'Unknown user'}</p>
        </div>
    </div>
    
    <div class="flex flex-wrap gap-2 mb-3">
        {#if learner.age}
            <span class="bg-green-100 text-green-800 text-xs font-medium px-2.5 py-0.5 rounded">
                Age: {learner.age}
            </span>
        {/if}
        
        {#if learner.grade_level}
            <span class="bg-purple-100 text-purple-800 text-xs font-medium px-2.5 py-0.5 rounded">
                Grade: {learner.grade_level}
            </span>
        {/if}
    </div>
    
    {#if showPreferences && learner.learning_preferences && learner.learning_preferences.length > 0}
        <div class="mt-3 mb-4">
            <p class="text-xs text-gray-500 mb-1">Learning preferences:</p>
            <div class="flex flex-wrap gap-1">
                {#each learner.learning_preferences as pref}
                    <span class="bg-gray-100 text-gray-800 text-xs font-medium px-2 py-0.5 rounded">
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