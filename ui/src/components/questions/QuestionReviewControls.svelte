<script lang="ts">
    import type { PracticeItem, ReviewStatus } from '$lib/types';
    import { practiceItemService } from '$lib/services/practiceItem';
    import PracticeItemEditForm from './PracticeItemEditForm.svelte';
    import ReviewStatusButton from '../common/ReviewStatusButton.svelte';
    import { toast } from '$lib/stores/toast';
    /**
     * The practice item to be reviewed
     */
    export let item: PracticeItem;
    
    /**
     * Event handler for when review status changes
     */
    export let onReviewStatusChange: (itemId: string, status: ReviewStatus) => void;
    
    let isEditing = false;
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
                isEditing = true;
            }
            
            toast.success(`Practice item marked as ${status?.toLowerCase().replace('_', ' ') || 'updated'}`);
        } catch (error) {
            console.error('Failed to update review status:', error);
            toast.error('Failed to update review status');
        } finally {
            isUpdating = false;
        }
    }

    function handleEditSave(updatedItem: PracticeItem) {
        if (updatedItem.review_status) {
            onReviewStatusChange(item.id, updatedItem.review_status);
        }
        // Update local item data
        item = { ...item, ...updatedItem };
        isEditing = false;
        toast.success('Practice item updated successfully');
    }

    function handleEditCancel() {
        isEditing = false;
    }

    function handleEditClick() {
        isEditing = true;
    }
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
    {#if isEditing}
        <div class="border-t border-gray-200 dark:border-gray-700 pt-4">
            <PracticeItemEditForm
                {item}
                onSave={handleEditSave}
                onCancel={handleEditCancel}
            />
        </div>
    {:else}
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
                    />
                    
                    <ReviewStatusButton
                        status={item.review_status}
                        targetStatus="NEED_EDIT"
                        label="Needs Edit"
                        activeLabel="Needs Edit"
                        reviewDate={item.review_date}
                        disabled={isUpdating}
                        onClick={() => handleReviewStatusChange('NEED_EDIT')}
                    />
                    
                    <ReviewStatusButton
                        status={item.review_status}
                        targetStatus="IGNORE"
                        label="Ignore"
                        activeLabel="Ignored"
                        reviewDate={item.review_date}
                        disabled={isUpdating}
                        onClick={() => handleReviewStatusChange('IGNORE')}
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
                    
                    <!-- Edit Button -->
                    <button
                        type="button"
                        class="inline-flex items-center px-3 py-1.5 border border-gray-300 dark:border-gray-600 shadow-sm text-xs font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                        on:click={handleEditClick}
                        disabled={isUpdating}
                        title="Edit this practice item"
                    >
                        <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                        Edit
                    </button>
                </div>
            </div>
        </div>
    {/if}
</div> 