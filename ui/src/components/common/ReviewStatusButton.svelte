<script lang="ts">
    import type { ReviewStatus, IconType } from '$lib/types';
    import Icon from './Icon.svelte';
    
    /**
     * The current review status to display
     */
    export let status: ReviewStatus | undefined;
    
    /**
     * The target status this button represents
     */
    export let targetStatus: ReviewStatus;
    
    /**
     * The display label for the button
     */
    export let label: string;
    
    /**
     * The label when the status is active
     */
    export let activeLabel: string;
    
    /**
     * Whether the button is disabled
     */
    export let disabled: boolean = false;
    
    /**
     * The review date for tooltip
     */
    export let reviewDate: string | undefined = undefined;
    
    /**
     * Optional icon to display
     */
    export let icon: IconType | undefined = undefined;
    
    /**
     * Click handler
     */
    export let onClick: () => void;
    
    $: isActive = status === targetStatus;
    $: displayLabel = isActive ? activeLabel : label;
    
    // Get button styles based on target status and active state
    function getButtonStyles(targetStatus: ReviewStatus, isActive: boolean): string {
        const baseStyles = 'inline-flex items-center gap-2 px-3 py-1.5 text-sm font-medium rounded-md transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed';
        
        if (isActive) {
            switch (targetStatus) {
                case 'APPROVED':
                    return `${baseStyles} bg-green-100 text-green-800 border border-green-200 dark:bg-green-900/40 dark:text-green-300 dark:border-green-700 focus:ring-green-500`;
                case 'NEED_EDIT':
                    return `${baseStyles} bg-yellow-100 text-yellow-800 border border-yellow-200 dark:bg-yellow-900/40 dark:text-yellow-300 dark:border-yellow-700 focus:ring-yellow-500`;
                case 'IGNORE':
                    return `${baseStyles} bg-red-100 text-red-800 border border-red-200 dark:bg-red-900/40 dark:text-red-300 dark:border-red-700 focus:ring-red-500`;
                default:
                    return `${baseStyles} bg-gray-100 text-gray-800 border border-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-600 focus:ring-gray-500`;
            }
        } else {
            return `${baseStyles} bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 hover:border-gray-400 dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-700 focus:ring-indigo-500`;
        }
    }
    
    // Get tooltip text
    function getTooltipText(): string {
        if (isActive && reviewDate) {
            const formattedDate = new Date(reviewDate).toLocaleString();
            switch (targetStatus) {
                case 'APPROVED':
                    return `Approved on ${formattedDate}`;
                case 'NEED_EDIT':
                    return `Marked for editing on ${formattedDate}`;
                case 'IGNORE':
                    return `Ignored on ${formattedDate}`;
                default:
                    return formattedDate;
            }
        } else {
            switch (targetStatus) {
                case 'APPROVED':
                    return 'Approve this item';
                case 'NEED_EDIT':
                    return 'Mark this item for editing';
                case 'IGNORE':
                    return 'Ignore this item';
                default:
                    return '';
            }
        }
    }
</script>

<button
    class={getButtonStyles(targetStatus, isActive)}
    on:click={onClick}
    {disabled}
    title={getTooltipText()}
    aria-pressed={isActive}
>
    {#if icon}
        <Icon type={icon} class_name="h-4 w-4" />
    {/if}
    {displayLabel}
</button> 