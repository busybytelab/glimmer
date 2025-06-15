<script lang="ts">
    import type { PracticeItem, ReviewStatus, IconType } from '$lib/types';
    import { practiceItemService } from '$lib/services/practiceItem';
    import { goto } from '$app/navigation';
    import ReviewStatusButton from '../common/ReviewStatusButton.svelte';
    import { toast } from '$lib/stores/toast';
    /**
     * The practice item to be reviewed
     */
    export let item: PracticeItem;
    
    /**
     * The session ID if this item is being reviewed in a session context
     */
    export let sessionId: string | undefined = undefined;
    
    /**
     * Event handler for when review status changes
     */
    export let onReviewStatusChange: (itemId: string, status: ReviewStatus) => void;
    
    let isUpdating = false;
    
    async function handleReviewStatusChange(status: ReviewStatus) {
        if (isUpdating) return;
        
        isUpdating = true;
        try {
            const updatedItem = await practiceItemService.updateItem(item.id, {
                review_status: status
            });
            
            if (updatedItem.review_status) {
                onReviewStatusChange(item.id, updatedItem.review_status);
                // Update local item data
                item = { ...item, ...updatedItem };
            }
            
            if (status === 'NEED_EDIT') {
                const editUrl = new URL(`/account/practice-items/${item.id}/edit`, window.location.origin);
                if (sessionId) {
                    editUrl.searchParams.set('sessionId', sessionId);
                }
                goto(editUrl.pathname + editUrl.search);
            }
            
            toast.success(`Practice item marked as ${status?.toLowerCase().replace('_', ' ') || 'updated'}`);
        } catch (error) {
            console.error('Failed to update review status:', error);
            toast.error('Failed to update review status');
        } finally {
            isUpdating = false;
        }
    }
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
    <!-- Review Status Section -->
    <div class="space-y-3">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white">
            Review Status
        </h3>
        <p class="text-sm text-gray-500 dark:text-gray-400">
            Choose the appropriate status for this practice item
        </p>
        
        <!-- Review Controls Row -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <!-- Status Buttons (Left) -->
            <div class="flex flex-wrap gap-3">
                <ReviewStatusButton
                    status={item.review_status}
                    targetStatus="APPROVED"
                    label="Approve"
                    activeLabel="Approved"
                    reviewDate={item.review_date}
                    disabled={isUpdating}
                    onClick={() => handleReviewStatusChange('APPROVED')}
                    icon={'complete' as IconType}
                />
                
                <ReviewStatusButton
                    status={item.review_status}
                    targetStatus="NEED_EDIT"
                    label="Edit"
                    activeLabel="Edit"
                    reviewDate={item.review_date}
                    disabled={isUpdating}
                    onClick={() => handleReviewStatusChange('NEED_EDIT')}
                    icon={'edit' as IconType}
                />
                
                <ReviewStatusButton
                    status={item.review_status}
                    targetStatus="IGNORE"
                    label="Ignore"
                    activeLabel="Ignored"
                    reviewDate={item.review_date}
                    disabled={isUpdating}
                    onClick={() => handleReviewStatusChange('IGNORE')}
                    icon={'ignore' as IconType}
                />
            </div>

            <!-- Right Side Info & Actions -->
            <div class="flex flex-col sm:flex-row sm:items-center gap-3">
                <!-- Status Information -->
                {#if item.review_status && item.review_date}
                    <div class="flex items-center text-xs text-gray-500 dark:text-gray-400">
                        <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span>
                            {new Date(item.review_date).toLocaleString()}
                        </span>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</div> 