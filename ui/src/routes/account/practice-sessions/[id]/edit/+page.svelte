<script lang="ts">
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
    import { onMount } from 'svelte';
    import type { PracticeSession, BreadcrumbItem, IconType } from '$lib/types';
    import PracticeSessionBasicForm from '$components/practice-sessions/PracticeSessionBasicForm.svelte';
    import Breadcrumbs from '$components/common/Breadcrumbs.svelte';
    import ActionToolbar from '$components/common/ActionToolbar.svelte';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    import { sessionService } from '$lib/services/session';
    
    let session: PracticeSession | null = null;
    let loading = true;
    let error: string | null = null;
    let breadcrumbItems: BreadcrumbItem[] = [];
    
    onMount(async () => {
        try {
            const sessionId = $page.params.id;
            if (sessionId) {
                await loadSession(sessionId);
                updateBreadcrumbs();
            } else {
                error = 'Session ID is required';
                loading = false;
            }
        } catch (err) {
            console.error('Error in onMount:', err);
            error = err instanceof Error ? err.message : 'An unexpected error occurred';
            loading = false;
        }
    });
    
    async function loadSession(id: string) {
        try {
            loading = true;
            error = null;
            
            if (!id) {
                throw new Error('Session ID is required');
            }
            
            session = await sessionService.loadSession(id);
            
        } catch (err) {
            console.error('Failed to load session:', err);
            error = 'Failed to load practice session';
        } finally {
            loading = false;
        }
    }
    
    function handleSessionUpdate() {
        if (!session) {
            return;
        }
        goto(`/account/practice-sessions/${session.id}/edit`);
    }
    
    function handleSessionDelete() {
        if (!session) {
            return;
        }
        // If we know the practice topic, go back to it, otherwise go to home
        if (session.expand?.practice_topic) {
            goto(`/account/practice-topics/${session.practice_topic}`);
        } else {
            goto('/account/practice-topics');
        }
    }
    
    function handleCancel() {
        if (!session) {
            return;
        }
        goto(`/account/practice-topics/${session.practice_topic}`);
    }
    
    function updateBreadcrumbs() {
        if (!session) {
            return;
        }
        
        const items: BreadcrumbItem[] = [
            {
                label: 'Topics',
                href: '/account/practice-topics',
                icon: 'topic' as IconType
            }
        ];
        
        if (session.expand?.practice_topic) {
            items.push({
                label: session.expand.practice_topic.name,
                href: `/account/practice-topics/${session.practice_topic}`,
                icon: 'topic' as IconType
            });
        }
        
        // Use the session name or a fallback
        const sessionName = session.name || 'Practice Session';
        
        items.push({
            label: sessionName,
            href: `/account/practice-sessions/${session.id}/overview`,
            icon: 'session' as IconType
        });
        
        items.push({
            label: 'Edit',
            icon: 'edit' as IconType
        });
        
        breadcrumbItems = items;
    }
    
    // Back action for the toolbar
    $: sessionActions = [
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
        <ActionToolbar actions={sessionActions} />
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <LoadingSpinner size="md" color="primary" />
        </div>
    {:else if error}
        <ErrorAlert message={error} />
    {:else if session}
        <div class="w-full">
            <PracticeSessionBasicForm
                {session}
                on:update={() => handleSessionUpdate()}
                on:delete={() => handleSessionDelete()}
                on:cancel={handleCancel}
            />
        </div>
    {/if}
</div> 