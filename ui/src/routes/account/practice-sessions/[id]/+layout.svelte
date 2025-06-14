<script lang="ts">
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';

    export let loading = false;
    export let error: string | null = null;
</script>

{#if loading}
    <div class="flex justify-center items-center h-screen">
        <LoadingSpinner size="lg" color="primary" />
    </div>
{:else if error}
    <div class="container mx-auto px-4 py-8">
        <ErrorAlert message={error} />
    </div>
{/if}

<slot /> 