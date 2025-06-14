<script lang="ts">
    import { IconTypeMap } from '$lib/types';
    import type { BreadcrumbItem } from '$lib/types';

    export let items: BreadcrumbItem[] = [];

    export let divider: string = '/';
    export let showHomeIcon: boolean = false;

    function getIconPath(icon: string | undefined): string {
        if (!icon) return '';
        return IconTypeMap[icon] || '';
    }
</script>

<nav class="flex py-1" aria-label="Breadcrumb">
    <ol class="inline-flex items-center space-x-1 md:space-x-2">
        {#if showHomeIcon && items.length > 0}
            <li class="inline-flex items-center">
                <a href="/" class="inline-flex items-center text-base font-medium text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
                    <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path d={getIconPath('home')}></path>
                    </svg>
                    Home
                </a>
            </li>
        {/if}
        
        {#each items as item, index}
            <li>
                <div class="flex items-center">
                    {#if index > 0 || showHomeIcon}
                        <span class="mx-1 md:mx-2 text-gray-400 dark:text-gray-500">{divider}</span>
                    {/if}
                    
                    {#if item.href && index !== items.length - 1}
                        <a 
                            href={item.href} 
                            class="inline-flex items-center text-base font-medium text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
                        >
                            {#if item.icon}
                                <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                    <path d={getIconPath(item.icon)}></path>
                                </svg>
                            {/if}
                            {item.label}
                        </a>
                    {:else}
                        <span 
                            class="inline-flex items-center text-base font-medium text-gray-500 dark:text-gray-400"
                            aria-current={index === items.length - 1 ? 'page' : undefined}
                        >
                            {#if item.icon}
                                <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                    <path d={getIconPath(item.icon)}></path>
                                </svg>
                            {/if}
                            {item.label}
                        </span>
                    {/if}
                </div>
            </li>
        {/each}
    </ol>
</nav> 