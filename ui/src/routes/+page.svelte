<script lang="ts">
  import { onMount } from 'svelte';
  import { isAuthenticated } from '$lib/stores';
  import LoadingSpinner from '../components/common/LoadingSpinner.svelte';

  // Redirect to dashboard if authenticated, otherwise to login
  onMount(() => {
    // Use standard browser navigation
    const path = $isAuthenticated ? '/dashboard' : '/login';
    window.location.href = path;
  });
</script>

<div class="flex justify-center items-center h-screen">
  <LoadingSpinner size="lg" color="gray" />
</div> 