<script lang="ts">
    import { isAuthLoading } from '$lib/stores';
    import { building } from '$app/environment';
    import LoadingSpinner from '../common/LoadingSpinner.svelte';

    // Skip auth check during build
    if (building) {
        isAuthLoading.set(false);
    }
</script>

{#if $isAuthLoading}
    <div class="flex justify-center items-center h-screen">
        <LoadingSpinner size="lg" color="gray" />
    </div>
{:else}
    <slot />
{/if} 