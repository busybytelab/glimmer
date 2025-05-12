<script lang="ts">
    import { onMount } from 'svelte';

    // Type definition for toolbar actions
    export let actions: Array<{
        id: string;
        label: string;
        icon: string;
        onClick: () => void;
        variant?: 'primary' | 'secondary' | 'danger' | 'success' | 'warning';
        disabled?: boolean;
    }> = [];

    export let showDropdownBreakpoint = 640; // sm breakpoint (in px)
    export let hideTextBreakpoint = 768; // md breakpoint (in px)
    
    let windowWidth: number;
    let showDropdown = false;
    let toolbarRef: HTMLDivElement;
    
    // Function to convert icon name to SVG path
    function getIconPath(icon: string): string {
        type IconType = 'print' | 'edit' | 'delete' | 'view' | 'download' | 'share' | 'duplicate' | 'add' | 'start' | 'complete' | 'reset' | 'back' | 'next' | 'more';
        
        const icons: Record<IconType, string> = {
            print: 'M5 4h14a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2zm3 15v-5h8v5H8zm10 0v-5a2 2 0 0 0-2-2H8a2 2 0 0 0-2 2v5H5a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h14a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1h-1zm-9-9a1 1 0 1 0 0-2 1 1 0 0 0 0 2z',
            edit: 'M16.77 8l1.94-2a1 1 0 0 0 0-1.41l-3.34-3.3a1 1 0 0 0-1.41 0l-1.97 2-8.5 8.5v4.21h4.21l8.5-8.5 1.57-1.5zm-9.8 6.1-1.39 1.37h-1.11v-1.11l1.38-1.37 1.12 1.11zm7.63-5.28-3.88 3.88-1.12-1.12 3.88-3.88 1.12 1.12z',
            delete: 'M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z',
            view: 'M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z',
            download: 'M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z',
            share: 'M18 16.08c-.76 0-1.44.3-1.96.77L8.91 12.7c.05-.23.09-.46.09-.7s-.04-.47-.09-.7l7.05-4.11c.54.5 1.25.81 2.04.81 1.66 0 3-1.34 3-3s-1.34-3-3-3-3 1.34-3 3c0 .24.04.47.09.7L8.04 9.81C7.5 9.31 6.79 9 6 9c-1.66 0-3 1.34-3 3s1.34 3 3 3c.79 0 1.5-.31 2.04-.81l7.12 4.16c-.05.21-.08.43-.08.65 0 1.61 1.31 2.92 2.92 2.92 1.61 0 2.92-1.31 2.92-2.92s-1.31-2.92-2.92-2.92z',
            duplicate: 'M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z',
            add: 'M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z',
            start: 'M8 5v14l11-7z',
            complete: 'M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z',
            reset: 'M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z',
            back: 'M20 11H7.83l5.59-5.59L12 4l-8 8 8 8 1.41-1.41L7.83 13H20v-2z',
            next: 'M12 4l-1.41 1.41L16.17 11H4v2h12.17l-5.58 5.59L12 20l8-8z',
            more: 'M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z'
        };
        
        return icons[icon as IconType] || icons.more;
    }
    
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
                                <svg class="mr-3 h-5 w-5 text-gray-400 dark:text-gray-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                                    <path d={getIconPath(action.icon)} />
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
                    <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                        <path d={getIconPath(action.icon)} />
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