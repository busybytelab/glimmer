<script lang="ts">
    export let items: Array<{
        label: string;
        href?: string;
        icon?: string;
    }> = [];

    export let divider: string = '/';
    export let showHomeIcon: boolean = false;

    type IconMap = {
        [key: string]: string;
    }

    function getIconPath(icon: string): string {
        const icons: IconMap = {
            home: 'M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z',
            topic: 'M4 14h4v-4H4v4zm0 5h4v-4H4v4zM4 9h4V5H4v4zm5 5h12v-4H9v4zm0 5h12v-4H9v4zM9 5v4h12V5H9z',
            session: 'M14 10H2v2h12v-2zm0-4H2v2h12V6zm4 8v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zM2 16h8v-2H2v2z',
            learner: 'M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z',
            edit: 'M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34a.9959.9959 0 0 0-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z',
            create: 'M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z'
        };
        
        return icons[icon] || '';
    }
</script>

<nav class="flex" aria-label="Breadcrumb">
    <ol class="inline-flex items-center space-x-1 md:space-x-2">
        {#if showHomeIcon && items.length > 0}
            <li class="inline-flex items-center">
                <a href="/" class="inline-flex items-center text-sm text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
                    <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
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
                            class="inline-flex items-center text-sm font-medium text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
                        >
                            {#if item.icon}
                                <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                    <path d={getIconPath(item.icon)}></path>
                                </svg>
                            {/if}
                            {item.label}
                        </a>
                    {:else}
                        <span 
                            class="inline-flex items-center text-sm font-medium text-gray-500 dark:text-gray-400"
                            aria-current={index === items.length - 1 ? 'page' : undefined}
                        >
                            {#if item.icon}
                                <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
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