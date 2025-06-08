<script lang="ts">
    import { goto } from '$app/navigation';
    import type { PracticeTopic, BreadcrumbItem, BreadcrumbIcon } from '$lib/types';
    import PracticeTopicForm from '../../../components/practice-topics/PracticeTopicForm.svelte';
    import Breadcrumbs from '../../../components/common/Breadcrumbs.svelte';
    import ActionToolbar from '../../../components/common/ActionToolbar.svelte';

    let breadcrumbItems: BreadcrumbItem[] = [
        {
            label: 'Topics',
            href: '/practice-topics',
            icon: 'topic' as BreadcrumbIcon
        },
        {
            label: 'New Topic',
            icon: 'create' as BreadcrumbIcon
        }
    ];

    function handleTopicUpdate(_updatedTopic: PracticeTopic) {
        // After creating a topic, navigate back to the topics list
        goto('/practice-topics');
    }

    function handleCancel() {
        goto('/practice-topics');
    }

    // Actions for the toolbar
    const topicActions = [
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

    <div class="w-full">
        <PracticeTopicForm
            on:update={({ detail }) => handleTopicUpdate(detail)}
            on:cancel={handleCancel}
        />
    </div>
</div> 