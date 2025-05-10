<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import type { PracticeTopic } from '$lib/types';
    import pb from '$lib/pocketbase';
    import PracticeTopicForm from '../../../../components/practice-topics/PracticeTopicForm.svelte';
    import Breadcrumbs from '../../../../components/common/Breadcrumbs.svelte';
    import ActionToolbar from '../../../../components/common/ActionToolbar.svelte';

    // Define the breadcrumb item type
    type BreadcrumbItem = {
        label: string;
        href?: string;
        icon?: string;
    };

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
                    icon: 'topic'
                },
                {
                    label: 'Edit Topic',
                    icon: 'edit'
                }
            ];
            return;
        }
        
        breadcrumbItems = [
            {
                label: 'Topics',
                href: '/practice-topics',
                icon: 'topic'
            },
            {
                label: topic.name,
                href: `/practice-topics/${topic.id}`,
                icon: 'topic'
            },
            {
                label: `Edit ${topic.name}`,
                icon: 'edit'
            }
        ];
    }

    function handleTopicUpdate(_updatedTopic: PracticeTopic) {
        goto('/practice-topics');
    }

    function handleTopicDelete(_topicId: string) {
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

<div class="container mx-auto px-4 py-8 max-w-7xl">
    <div class="flex justify-between items-center mb-6">
        <div>
            <Breadcrumbs items={breadcrumbItems} />
        </div>
        <ActionToolbar actions={topicActions} />
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary"></div>
        </div>
    {:else if error}
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
            <strong class="font-bold">Error!</strong>
            <span class="block sm:inline"> {error}</span>
        </div>
    {:else if topic}
        <div class="w-full">
            <PracticeTopicForm
                {topic}
                on:update={({ detail }) => handleTopicUpdate(detail)}
                on:delete={({ detail }) => handleTopicDelete(detail)}
                on:cancel={handleCancel}
            />
        </div>
    {/if}
</div> 