<script lang="ts">
	import { onMount } from 'svelte';
	import { isAuthenticated, user, isAuthLoading, error } from '$lib/stores';
	import type { Instructor, Learner } from '$lib/types';
	import pb from '$lib/pocketbase';
	import Login from './routes/Login.svelte';
	import Dashboard from './routes/Dashboard.svelte';

	// Create a promise to track when auth is initialized
	let authInitialized = false;
	let authStatePromise: Promise<void>;
	
	// Function to handle the auth flow
	async function initializeAuth() {
		isAuthLoading.set(true);
		error.set(null);
		
		try {
			// Check for saved token
			if (localStorage.getItem('authToken')) {
				pb.authStore.loadFromCookie(`pb_auth=${localStorage.getItem('authToken')}`);
			}
			
			// Only proceed with auth verification if we have a token
			if (pb.authStore.token) {
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
							const learner = await pb.collection('learners').getFirstListItem(`user="${pb.authStore.record?.id}"`, { requestKey: null });
							if (learner) {
								learner.user = userData;
								user.set(learner as unknown as Learner);
								isAuthenticated.set(true);
								return;
							}
						} catch (err) {
							// No learner found
						}

						// If we get here, user exists in main users collection but not in instructors/learners
						// Clear auth state and show login
						pb.authStore.clear();
						localStorage.removeItem('authToken');
						isAuthenticated.set(false);
					}
				} catch (err) {
					// Token refresh failed, clear auth state
					console.error('Auth refresh failed:', err);
					pb.authStore.clear();
					localStorage.removeItem('authToken');
					isAuthenticated.set(false);
				}
			} else {
				// No token, user is not authenticated
				isAuthenticated.set(false);
			}
		} catch (err) {
			// Catch any other errors
			console.error('Authentication initialization failed:', err);
			pb.authStore.clear();
			localStorage.removeItem('authToken');
			isAuthenticated.set(false);
		} finally {
			isAuthLoading.set(false);
			authInitialized = true;
		}
	}

	onMount(() => {
		authStatePromise = initializeAuth();
	});

	// Handle browser navigation
	window.addEventListener('popstate', () => {
		if (authInitialized) {
			// Only reload if auth is already initialized to prevent loops
			window.location.reload();
		}
	});
</script>

<main>
	{#if $isAuthLoading}
		<div class="flex justify-center items-center h-screen">
			<div class="animate-spin rounded-full h-32 w-32 border-t-2 border-b-2 border-gray-900 dark:border-white"></div>
		</div>
	{:else if $error}
		<div class="flex flex-col items-center justify-center h-screen">
			<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
				<strong class="font-bold">Error!</strong>
				<span class="block sm:inline">{$error}</span>
			</div>
			<button 
				on:click={() => window.location.reload()} 
				class="mt-4 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
			>
				Try Again
			</button>
		</div>
	{:else if $isAuthenticated && $user !== null}
		<Dashboard />
	{:else}
		<Login />
	{/if}
</main>

<style>
	main {
		max-width: 1280px;
		margin: 0 auto;
		padding: 2rem;
		text-align: center;
	}
</style>
