<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import type { PracticeTopic, BreadcrumbItem, BreadcrumbIcon } from '$lib/types';
    import pb from '$lib/pocketbase';
    import PracticeTopicForm from '../../../../components/practice-topics/PracticeTopicForm.svelte';
    import Breadcrumbs from '../../../../components/common/Breadcrumbs.svelte';
    import ActionToolbar from '../../../../components/common/ActionToolbar.svelte';
    import LoadingSpinner from '../../../../components/common/LoadingSpinner.svelte';
    import ErrorAlert from '../../../../components/common/ErrorAlert.svelte';

    let topic: PracticeTopic | null = null;
    let loading = true;
    let error: string | null = null;
    let breadcrumbItems: BreadcrumbItem[] = [];

    onMount(async () => {
        // Get topic ID from the URL parameter
        const topicId = $page.params.id;
        
        if (!topicId) {
            error = 'No topic ID provided';
            loading = false;
            return;
        }
        
        await loadTopic(topicId);
        updateBreadcrumbs();
    });

    async function loadTopic(id: string) {
        try {
            loading = true;
            error = null;
            
            const result = await pb.collection('practice_topics').getOne<PracticeTopic>(id);
            
            // Make sure tags are properly formatted as arrays
            if (result.tags) {
                const tagsValue = result.tags as unknown;
                if (typeof tagsValue === 'string') {
                    try {
                        if (tagsValue.trim().startsWith('[')) {
                            result.tags = JSON.parse(tagsValue) as string[];
                        } else {
                            result.tags = tagsValue.split(',').map((tag: string) => tag.trim()).filter(Boolean);
                        }
                    } catch (e) {
                        console.error('Error parsing tags:', e);
                        result.tags = [];
                    }
                }
            } else {
                result.tags = [];
            }
            
            topic = result;
        } catch (err) {
            console.error('Failed to load topic:', err);
            error = 'Failed to load practice topic';
        } finally {
            loading = false;
        }
    }

    function updateBreadcrumbs() {
        if (!topic) {
            // Basic breadcrumbs while topic is still loading
            breadcrumbItems = [
                {
                    label: 'Topics',
                    href: '/practice-topics',
                    icon: 'topic' as BreadcrumbIcon
                },
                {
                    label: 'Edit Topic',
                    icon: 'edit' as BreadcrumbIcon
                }
            ];
            return;
        }
        
        breadcrumbItems = [
            {
                label: 'Topics',
                href: '/practice-topics',
                icon: 'topic' as BreadcrumbIcon
            },
            {
                label: topic.name,
                href: `/practice-topics/${topic.id}`,
                icon: 'topic' as BreadcrumbIcon
            },
            {
                label: `Edit ${topic.name}`,
                icon: 'edit' as BreadcrumbIcon
            }
        ];
    }

    function handleTopicUpdate() {
        goto('/practice-topics');
    }

    function handleTopicDelete() {
        goto('/practice-topics');
    }

    function handleCancel() {
        goto('/practice-topics');
    }

    // Actions for the toolbar
    $: topicActions = [
        {
            id: 'back',
            label: 'Back',
            icon: 'back',
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
        <ActionToolbar actions={topicActions} />
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <LoadingSpinner size="md" color="primary" />
        </div>
    {:else if error}
        <ErrorAlert message={error} />
    {:else if topic}
        <div class="w-full">
            <PracticeTopicForm
                {topic}
                on:update={() => handleTopicUpdate()}
                on:delete={() => handleTopicDelete()}
                on:cancel={handleCancel}
            />
        </div>
    {/if}
</div> 