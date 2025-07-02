<script lang="ts">
    /**
     * The popularity count to display
     */
    export let count: number;
    
    /**
     * The size variant of the badge
     * - 'sm': Small badge for compact displays
     * - 'md': Medium badge for standard use
     * @default 'sm'
     */
    export let size: 'sm' | 'md' = 'sm';
    
    const sizeClasses = {
        sm: 'text-xs px-2 py-1',
        md: 'text-sm px-3 py-1.5'
    };
    
    const iconSizeClasses = {
        sm: 'h-3 w-3',
        md: 'h-4 w-4'
    };
    
    /**
     * Format large numbers into compact display format
     * Examples: 1234 → 1.2K, 1234567 → 1.2M
     */
    function formatCount(count: number): string {
        if (count < 1000) {
            return count.toString();
        } else if (count < 1000000) {
            const formatted = (count / 1000).toFixed(1);
            // Remove trailing .0
            return formatted.endsWith('.0') ? formatted.slice(0, -2) + 'K' : formatted + 'K';
        } else if (count < 1000000000) {
            const formatted = (count / 1000000).toFixed(1);
            // Remove trailing .0
            return formatted.endsWith('.0') ? formatted.slice(0, -2) + 'M' : formatted + 'M';
        } else {
            const formatted = (count / 1000000000).toFixed(1);
            // Remove trailing .0
            return formatted.endsWith('.0') ? formatted.slice(0, -2) + 'B' : formatted + 'B';
        }
    }
</script>

<div class="flex items-center space-x-1 text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 rounded {sizeClasses[size]}">
    <svg xmlns="http://www.w3.org/2000/svg" class="{iconSizeClasses[size]} text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
        <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
    </svg>
    <span>{formatCount(count)}</span>
</div> 