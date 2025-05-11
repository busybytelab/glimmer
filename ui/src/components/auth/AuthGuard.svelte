<script lang="ts">
    import { isAuthLoading } from '$lib/stores';
    import { building } from '$app/environment';

    // Skip auth check during build
    if (building) {
        isAuthLoading.set(false);
    }
</script>

{#if $isAuthLoading}
    <div class="flex justify-center items-center h-screen">
        <div class="animate-spin rounded-full h-32 w-32 border-t-2 border-b-2 border-gray-900 dark:border-white"></div>
    </div>
{:else}
    <slot />
{/if} 