<script lang="ts">
    import type { PracticeItem, ReviewStatus } from '$lib/types';
    import { practiceItemService } from '$lib/services/practiceItem';
    import PracticeItemEditForm from './PracticeItemEditForm.svelte';
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
    
    async function handleReviewStatusChange(status: PracticeItem['review_status']) {
        try {
            const updatedItem = await practiceItemService.updatePracticeItem(item.id, {
                review_status: status
            });
            if (updatedItem.review_status) {
                onReviewStatusChange(item.id, updatedItem.review_status);
            }
            if (status === 'NEED_EDIT') {
                isEditing = true;
            }
            toast.success(`Practice item marked as ${status?.toLowerCase() || 'updated'}`);
        } catch (error) {
            console.error('Failed to update review status:', error);
            toast.error('Failed to update review status');
        }
    }

    function handleEditSave(updatedItem: PracticeItem) {
        if (updatedItem.review_status) {
            onReviewStatusChange(item.id, updatedItem.review_status);
        }
        isEditing = false;
        toast.success('Practice item updated successfully');
    }

    function handleEditCancel() {
        isEditing = false;
    }
</script>

<div class="flex flex-col gap-4">
    {#if isEditing}
        <PracticeItemEditForm
            {item}
            onSave={handleEditSave}
            onCancel={handleEditCancel}
        />
    {:else}
        <div class="flex flex-wrap gap-2">
            <button
                class="px-3 py-1.5 text-sm rounded-md transition-colors {item.review_status === 'APPROVED'
                    ? 'bg-green-100 text-green-800'
                    : 'bg-gray-100 text-gray-800 hover:bg-gray-200'}"
                on:click={() => handleReviewStatusChange('APPROVED')}
                title={item.review_status === 'APPROVED' && item.review_date 
                    ? `Approved by ${item.expand?.reviewer?.expand?.user?.name || 'Unknown'} on ${new Date(item.review_date).toLocaleString()}`
                    : 'Approve this item'}
            >
                {item.review_status === 'APPROVED' ? 'Approved' : 'Approve'}
            </button>
            <button
                class="px-3 py-1.5 text-sm rounded-md transition-colors {item.review_status === 'NEED_EDIT'
                    ? 'bg-yellow-100 text-yellow-800'
                    : 'bg-gray-100 text-gray-800 hover:bg-gray-200'}"
                on:click={() => handleReviewStatusChange('NEED_EDIT')}
                title={item.review_status === 'NEED_EDIT' && item.review_date 
                    ? `Marked for editing by ${item.expand?.reviewer?.expand?.user?.name || 'Unknown'} on ${new Date(item.review_date).toLocaleString()}`
                    : 'Mark this item for editing'}
            >
                Needs Edit
            </button>
            <button
                class="px-3 py-1.5 text-sm rounded-md transition-colors {item.review_status === 'IGNORE'
                    ? 'bg-red-100 text-red-800 dark:bg-red-900/40 dark:text-red-300'
                    : 'bg-gray-100 text-gray-800 hover:bg-gray-200'}"
                on:click={() => handleReviewStatusChange('IGNORE')}
                title={item.review_status === 'IGNORE' && item.review_date 
                    ? `Ignored by ${item.expand?.reviewer?.expand?.user?.name || 'Unknown'} on ${new Date(item.review_date).toLocaleString()}`
                    : 'Ignore this item'}
            >
                {item.review_status === 'IGNORE' ? 'Ignored' : 'Ignore'}
            </button>
        </div>
    {/if}
</div> 