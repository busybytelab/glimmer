<script lang="ts">
  import { onMount } from 'svelte';
  import { isAuthenticated, user, isAuthLoading, error, theme } from '$lib/stores';
  import pb from '$lib/pocketbase';
  import { getAuthToken, clearAuthToken } from '$lib/auth';
  import type { Instructor, Learner } from '$lib/types';
  import SideNav from '../components/layout/SideNav.svelte';
  import LoadingSpinner from '../components/common/LoadingSpinner.svelte';
  import ErrorAlert from '../components/common/ErrorAlert.svelte';
  import Toast from '../components/common/Toast.svelte';
  import { isPublicRoute } from '$lib/auth';
  import '../app.css';

  // Sidebar state for layout
  let sidebarOpen = false; // Start with closed sidebar on mobile

  // Toggle sidebar function
  function toggleSidebar() {
    sidebarOpen = !sidebarOpen;
  }

  // Check if current route is public
  $: isPublic = isPublicRoute(window.location.pathname);

  // Function to handle the auth flow
  async function initializeAuth() {
    isAuthLoading.set(true);
    error.set(null);
    try {
      // Get token using our utility function
      const token = getAuthToken();
      // Only proceed with auth verification if we have a token
      if (token) {
        // Set token in PocketBase if it's not already set
        if (pb.authStore.token !== token) {
          // Since token is read-only, we need to clear and recreate the auth store
          pb.authStore.clear();
          // Then manually save the auth data with the token
          localStorage.setItem('pocketbase_auth', JSON.stringify({
            token: token,
            model: pb.authStore.model
          }));
        }

        
        try {
          // Refresh auth state, which validates the token and gets fresh user data
          await pb.collection('users').authRefresh();
          // Token is valid, get user data
          if (pb.authStore.isValid) {
            const userData = await pb.collection('users').getOne(pb.authStore.record?.id ?? '');
            // First check if user is an instructor
            try {
              const instructor = await pb.collection('instructors').getFirstListItem(`user="${pb.authStore.record?.id}"`);
              if (instructor) {
                instructor.user = userData;
                user.set(instructor as unknown as Instructor);
                isAuthenticated.set(true);
                return;
              }
            } catch (err) {
              // No instructor found, continue to check for learner
            }

            // Then check if user is a learner
            try {
              const learner = await pb.collection('learners').getFirstListItem(`user.id="${pb.authStore.record?.id}"`, { requestKey: null });
              if (learner) {
                learner.user = userData;
                user.set(learner as unknown as Learner);
                isAuthenticated.set(true);
                return;
              }
            } catch (err) {
              // No learner found
            }

            console.log('WARNING: no instructor or learner found');

            // If we get here, user exists in main users collection but not in instructors/learners
            // Clear auth state and show login
            clearAuthToken();
            isAuthenticated.set(false);
          }
        } catch (err) {
          // Token refresh failed, clear auth state
          console.error('Auth refresh failed:', err);
          clearAuthToken();
          isAuthenticated.set(false);
        }
      } else {
        // No token, user is not authenticated
        isAuthenticated.set(false);
      }
    } catch (err) {
      // Catch any other errors
      console.error('Authentication initialization failed:', err);
      clearAuthToken();
      isAuthenticated.set(false);
    } finally {
      isAuthLoading.set(false);
    }
  }

  onMount(() => {
    // Initialize auth
    initializeAuth();
    
    // Initialize theme
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme) {
      theme.set(savedTheme as 'light' | 'dark');
    }
  });
</script>

<style lang="postcss">
  :global(html) {
    --primary: #2c3e50;
    --secondary: #3498db;
  }
  @media print {
    :global(body) {
      height: auto !important;
      overflow: visible !important;
    }
  }
</style>

{#if $isAuthLoading}
  <div class="flex justify-center items-center h-screen bg-content">
    <LoadingSpinner size="lg" color="gray" />
  </div>
{:else if $error}
  <div class="flex flex-col items-center justify-center h-screen bg-content">
    <ErrorAlert message={$error} />
    <button 
      on:click={() => window.location.reload()} 
      class="mt-4 bg-blue-500 hover:bg-blue-700 dark:bg-blue-600 dark:hover:bg-blue-800 text-white font-bold py-2 px-4 rounded"
    >
      Try Again
    </button>
  </div>
{:else}
  {#if $isAuthenticated && $user && !isPublic}
    <!-- Authenticated layout -->
    <div class="h-screen flex overflow-hidden bg-gray-100 dark:bg-gray-900 print:h-auto print:overflow-visible">
      <!-- Mobile sidebar - overlay when opened -->
      {#if sidebarOpen}
        <div class="md:hidden fixed inset-0 flex z-40 print:hidden">
          <!-- Overlay backdrop - behind the sidebar -->
          <div 
            class="fixed inset-0 bg-gray-600 bg-opacity-75 dark:bg-black dark:bg-opacity-80"
            on:click={toggleSidebar}
            on:keydown={(e) => e.key === 'Escape' && toggleSidebar()}
            aria-label="Close sidebar overlay"
            role="button"
            tabindex="0"
          ></div>
          
          <!-- Mobile SideNav - positioned on top of overlay -->
          <div class="relative flex-shrink-0 w-64 max-w-sm z-50">
            <SideNav toggleSidebar={toggleSidebar} />
          </div>
        </div>
      {/if}
      
      <!-- Desktop sidebar - always visible on desktop -->
      <div class="hidden md:flex md:flex-shrink-0 print:hidden">
        <SideNav toggleSidebar={toggleSidebar} />
      </div>
      
      <!-- Main content -->
      <div class="flex flex-col w-0 flex-1 overflow-hidden h-screen print:h-auto print:overflow-visible print:relative print:z-0">
        <!-- Mobile header with hamburger menu -->
        <div class="md:hidden pl-1 pt-1 sm:pl-3 sm:pt-3 print:hidden">
          <button
            on:click={toggleSidebar}
            class="-ml-0.5 -mt-0.5 h-12 w-12 inline-flex items-center justify-center rounded-md text-gray-500 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white focus:outline-none focus:ring-2 focus:ring-inset focus:ring-secondary"
          >
            <span class="sr-only">Open sidebar</span>
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
        </div>
        <!-- Main content area -->
        <main class="flex-1 relative z-0 overflow-y-auto focus:outline-none h-full bg-content print:h-auto print:overflow-visible print:relative print:z-0">
          <slot />
        </main>
      </div>
    </div>
  {:else}
    <!-- Unauthenticated layout -->
    <div class="bg-content min-h-screen">
      <slot />
    </div>
  {/if}
{/if}

<Toast /> 