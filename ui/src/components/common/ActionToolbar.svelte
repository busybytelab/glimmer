<script lang="ts">
    import { onMount } from 'svelte';
    import { IconTypeMap, type IconType } from '$lib/types';

    // Type definition for toolbar actions
    export let actions: Array<{
        id: string;
        label: string;
        icon: IconType;
        onClick: () => void;
        variant?: 'primary' | 'secondary' | 'danger' | 'success' | 'warning';
        disabled?: boolean;
    }> = [];

    export let showDropdownBreakpoint = 640; // sm breakpoint (in px)
    export let hideTextBreakpoint = 768; // md breakpoint (in px)
    
    let windowWidth: number;
    let showDropdown = false;
    let toolbarRef: HTMLDivElement;
    
    // Get button classes based on variant
    function getButtonClasses(variant: string = 'primary'): string {
        const baseClass = 'inline-flex items-center justify-center px-3 py-2 border shadow-sm rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 transition-all duration-200 ease-in-out disabled:opacity-50 disabled:cursor-not-allowed';
        
        type VariantType = 'primary' | 'secondary' | 'danger' | 'success' | 'warning';
        
        const variantClasses: Record<VariantType, string> = {
            primary: 'bg-indigo-600 hover:bg-indigo-700 focus:ring-indigo-500 border-indigo-600 text-white dark:bg-indigo-700 dark:hover:bg-indigo-800',
            secondary: 'bg-gray-600 hover:bg-gray-700 focus:ring-gray-500 border-gray-600 text-white dark:bg-gray-700 dark:hover:bg-gray-800',
            danger: 'bg-red-600 hover:bg-red-700 focus:ring-red-500 border-red-600 text-white dark:bg-red-700 dark:hover:bg-red-800',
            success: 'bg-green-600 hover:bg-green-700 focus:ring-green-500 border-green-600 text-white dark:bg-green-700 dark:hover:bg-green-800',
            warning: 'bg-yellow-600 hover:bg-yellow-700 focus:ring-yellow-500 border-yellow-600 text-white dark:bg-yellow-700 dark:hover:bg-yellow-800'
        };
        
        return `${baseClass} ${variantClasses[variant as VariantType]}`;
    }
    
    // Function to handle click outside of dropdown
    function handleClickOutside(event: MouseEvent) {
        if (showDropdown && toolbarRef && !toolbarRef.contains(event.target as Node)) {
            showDropdown = false;
        }
    }
    
    onMount(() => {
        // Add click outside listener
        document.addEventListener('click', handleClickOutside);
        
        // Check window width for initial UI state
        windowWidth = window.innerWidth;
        
        // Add resize listener
        const handleResize = () => {
            windowWidth = window.innerWidth;
        };
        
        window.addEventListener('resize', handleResize);
        
        // Cleanup
        return () => {
            document.removeEventListener('click', handleClickOutside);
            window.removeEventListener('resize', handleResize);
        };
    });
    
    $: showAsDropdown = windowWidth < showDropdownBreakpoint;
    $: hideText = windowWidth < hideTextBreakpoint && !showAsDropdown;
</script>

<div class="action-toolbar" bind:this={toolbarRef}>
    {#if showAsDropdown}
        <!-- Dropdown for very small screens -->
        <div class="relative">
            <button 
                class="inline-flex items-center justify-center px-3 py-2 border border-gray-300 dark:border-gray-600 shadow-sm text-sm font-medium rounded-md text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-offset-gray-800"
                on:click={() => showDropdown = !showDropdown}
            >
                <span>Actions</span>
                <svg class="ml-2 -mr-0.5 h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                    <path d={showDropdown ? 'M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z' : 'M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z'} />
                </svg>
            </button>
            
            {#if showDropdown}
                <div class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white dark:bg-gray-800 ring-1 ring-black ring-opacity-5 dark:ring-gray-700 focus:outline-none z-10">
                    <div class="py-1">
                        {#each actions as action}
                            <button
                                class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center"
                                on:click={() => {
                                    showDropdown = false;
                                    action.onClick();
                                }}
                                disabled={action.disabled}
                            >
                                <svg class="mr-3 h-5 w-5 text-gray-400 dark:text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                                    <path stroke-linecap="round" stroke-linejoin="round" d={IconTypeMap[action.icon]} />
                                </svg>
                                {action.label}
                            </button>
                        {/each}
                    </div>
                </div>
            {/if}
        </div>
    {:else}
        <!-- Regular buttons for larger screens -->
        <div class="flex space-x-2">
            {#each actions as action}
                <button
                    class={getButtonClasses(action.variant)}
                    on:click={action.onClick}
                    disabled={action.disabled}
                    title={action.label}
                >
                    <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                        <path stroke-linecap="round" stroke-linejoin="round" d={IconTypeMap[action.icon]} />
                    </svg>
                    {#if !hideText}
                        <span class="ml-2">{action.label}</span>
                    {/if}
                </button>
            {/each}
        </div>
    {/if}
</div>

<style>
    .action-toolbar {
        display: flex;
        justify-content: flex-end;
    }
</style> 