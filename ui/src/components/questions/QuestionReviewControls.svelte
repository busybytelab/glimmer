<script lang="ts">
    import type { PracticeItem, ReviewStatus } from '$lib/types';
    import pb from '$lib/pocketbase';
    
    /**
     * The practice item to be reviewed
     */
    export let item: PracticeItem;
    
    /**
     * Event handler for when review status changes
     */
    export let onReviewStatusChange: (itemId: string, status: ReviewStatus) => void;
    
    let loading = false;
    let error: string | null = null;
    
    async function handleReviewStatusChange(status: ReviewStatus) {
        if (!item.id) return;
        
        loading = true;
        error = null;
        
        try {
            // First, get the instructor record for the current user
            const instructor = await pb.collection('instructors').getFirstListItem(`user = "${pb.authStore.model?.id}"`);
            
            if (!instructor) {
                throw new Error('Instructor record not found');
            }
            
            // Update the practice item with the new review status
            await pb.collection('practice_items').update(item.id, {
                review_status: status,
                review_date: new Date().toISOString(),
                reviewer: instructor.id // Use the instructor record ID
            });
            
            // Notify parent component
            onReviewStatusChange(item.id, status);
        } catch (err) {
            console.error('Failed to update review status:', err);
            error = err instanceof Error ? err.message : 'Failed to update review status';
        } finally {
            loading = false;
        }
    }
</script>

<div class="mt-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-md">
    <h5 class="text-sm font-medium text-gray-900 dark:text-white mb-3">Review Status:</h5>
    
    {#if error}
        <p class="text-red-600 dark:text-red-400 text-sm mb-3">{error}</p>
    {/if}
    
    <div class="flex space-x-2">
        <button
            class="px-3 py-1 text-sm rounded-md transition-colors {item.review_status === 'APPROVED' 
                ? 'bg-green-600 text-white' 
                : 'bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'}"
            on:click={() => handleReviewStatusChange('APPROVED')}
            disabled={loading}
            type="button"
        >
            Approve
        </button>
        
        <button
            class="px-3 py-1 text-sm rounded-md transition-colors {item.review_status === 'NEED_EDIT' 
                ? 'bg-yellow-600 text-white' 
                : 'bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'}"
            on:click={() => handleReviewStatusChange('NEED_EDIT')}
            disabled={loading}
            type="button"
        >
            Needs Edit
        </button>
        
        <button
            class="px-3 py-1 text-sm rounded-md transition-colors {item.review_status === 'IGNORE' 
                ? 'bg-red-600 text-white' 
                : 'bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'}"
            on:click={() => handleReviewStatusChange('IGNORE')}
            disabled={loading}
            type="button"
        >
            Ignore
        </button>
    </div>
    
    {#if item.review_date}
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-2">
            Reviewed on: {new Date(item.review_date).toLocaleDateString()}
        </p>
    {/if}
</div> 