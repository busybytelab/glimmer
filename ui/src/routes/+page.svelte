<script lang="ts">
  import { onMount } from 'svelte';
  import { isAuthenticated } from '$lib/stores';

  // Redirect to dashboard if authenticated, otherwise to login
  onMount(() => {
    // Use standard browser navigation
    const path = $isAuthenticated ? '/dashboard' : '/login';
    window.location.href = path;
  });
</script>

<div class="flex justify-center items-center h-screen">
  <div class="animate-spin rounded-full h-32 w-32 border-t-2 border-b-2 border-gray-900 dark:border-white"></div>
</div> 