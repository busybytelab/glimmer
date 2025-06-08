<script lang="ts">
    import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
    import { page } from '$app/stores';
	import { authService } from '$lib/services/auth';
	import LoadingSpinner from '../../../components/common/LoadingSpinner.svelte';
	import ErrorAlert from '../../../components/common/ErrorAlert.svelte';

    let loading = true;
    let error: string | null = null;

    onMount(async () => {
        try {
            const sessionId = $page.params.id;
			if (!sessionId) {
				error = 'Session ID is required';
				return;
			}

			// Check if user is authenticated
			const authenticated = await authService.isAuthenticated();
			if (!authenticated) {
				goto('/login');
            return;
        }
        
		goto(`/practice-sessions/${sessionId}/instructor`);
		// TODO: for learner view, we need to move this under leaner route
		//goto(`/practice-sessions/${sessionId}/learner`);
			
        } catch (err) {
			console.error('Error in practice session route:', err);
			error = 'Failed to load practice session';
        } finally {
			loading = false;
		}
	});
</script>

    {#if loading}
	<div class="flex items-center justify-center min-h-[400px]">
		<LoadingSpinner />
        </div>
    {:else if error}
	<div class="p-4">
        <ErrorAlert message={error} />
</div>
{/if}