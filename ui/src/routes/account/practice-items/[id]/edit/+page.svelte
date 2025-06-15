<script lang="ts">
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
    import { onMount } from 'svelte';
    import type { PracticeItem, BreadcrumbItem, IconType } from '$lib/types';
    import { practiceItemService } from '$lib/services/practiceItem';
    import PracticeItemEditForm from '$components/questions/PracticeItemEditForm.svelte';
    import Breadcrumbs from '$components/common/Breadcrumbs.svelte';
    import ActionToolbar from '$components/common/ActionToolbar.svelte';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    
    let item: PracticeItem | null = null;
    let loading = true;
    let error: string | null = null;
    let breadcrumbItems: BreadcrumbItem[] = [];
    
    // Get session ID from URL if available
    $: sessionId = $page.url.searchParams.get('sessionId');
    
    onMount(async () => {
        try {
            const itemId = $page.params.id;
            if (itemId) {
                await loadItem(itemId);
                updateBreadcrumbs();
            } else {
                error = 'Item ID is required';
                loading = false;
            }
        } catch (err) {
            console.error('Error in onMount:', err);
            error = err instanceof Error ? err.message : 'An unexpected error occurred';
            loading = false;
        }
    });
    
    async function loadItem(id: string) {
        try {
            loading = true;
            error = null;
            
            if (!id) {
                throw new Error('Item ID is required');
            }
            
            item = await practiceItemService.getItem(id);
            
        } catch (err) {
            console.error('Failed to load item:', err);
            error = 'Failed to load practice item';
        } finally {
            loading = false;
        }
    }
    
    function handleItemUpdate() {
        // After successful update, navigate back to the session overview if we have a session ID
        if (sessionId) {
            goto(`/account/practice-sessions/${sessionId}/overview`);
        } else {
            // If no session context, go back to topics
            goto('/account/practice-topics');
        }
    }
    
    function handleCancel() {
        // Navigate back based on context
        if (sessionId) {
            goto(`/account/practice-sessions/${sessionId}/overview`);
        } else {
            goto('/account/practice-topics');
        }
    }
    
    function updateBreadcrumbs() {
        if (!item) return;
        
        const items: BreadcrumbItem[] = [
            {
                label: 'Topics',
                href: '/account/practice-topics',
                icon: 'topic' as IconType
            }
        ];
        
        if (item.expand?.practice_topic) {
            items.push({
                label: item.expand.practice_topic.name,
                href: `/account/practice-topics/${item.practice_topic}`,
                icon: 'topic' as IconType
            });
        }
        
        if (sessionId) {
            items.push({
                label: 'Practice Session',
                href: `/account/practice-sessions/${sessionId}/overview`,
                icon: 'session' as IconType
            });
        }
        
        items.push({
            label: 'Edit Practice Item',
            icon: 'edit' as IconType
        });
        
        breadcrumbItems = items;
    }
    
    // Actions for the toolbar
    $: itemActions = [
        {
            id: 'back',
            label: 'Back',
            icon: 'back' as IconType,
            variant: 'secondary' as const,
            onClick: handleCancel
        }
    ];
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6 max-w-7xl">
    <div class="flex justify-between items-center mb-6">
        <div>
            <Breadcrumbs items={breadcrumbItems} />
        </div>
        <ActionToolbar actions={itemActions} />
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <LoadingSpinner size="md" color="primary" />
        </div>
    {:else if error}
        <ErrorAlert message={error} />
    {:else if item}
        <div class="w-full">
            <PracticeItemEditForm
                {item}
                onSave={handleItemUpdate}
                onCancel={handleCancel}
            />
        </div>
    {/if}
</div> 