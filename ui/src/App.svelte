<script lang="ts">
	import { onMount } from 'svelte';
	import { isAuthenticated, user, isAuthLoading, error } from '$lib/stores';
	import type { Instructor, Learner } from '$lib/types';
	import pb from '$lib/pocketbase';
	import { getAuthToken, clearAuthToken } from '$lib/auth';
	import Login from './routes/Login.svelte';
	import Dashboard from './routes/Dashboard.svelte';
	import Profile from './routes/Profile.svelte';
	import Settings from './routes/Settings.svelte';
	import PracticeTopics from './routes/PracticeTopics.svelte';
	import ViewPracticeTopic from './routes/ViewPracticeTopic.svelte';
	import CreatePracticeSession from './routes/CreatePracticeSession.svelte';
	import Learners from './routes/Learners.svelte';
	import Chat from './routes/Chat.svelte';

	// Create a promise to track when auth is initialized
	let authInitialized = false;
	
	// Current route - default to dashboard
	let currentRoute = 'dashboard';
	
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
			authInitialized = true;
		}
	}

	onMount(() => {
		// Initialize auth
		initializeAuth();
		
		// Handle routing based on URL path
		function handleRouting() {
			const path = window.location.pathname || '/';
			console.log('Routing path:', path);
			
			// Check if this is a dynamic route with parameters
			if (path.includes('/practice-topic/')) {
				console.log('Setting route to practice-topic view');
				currentRoute = 'practice-topic';
			} else if (path.includes('/create-practice/')) {
				console.log('Setting route to create-practice');
				currentRoute = 'create-practice';
			} else if (path.includes('/practice-session/')) {
				console.log('Setting route to practice-session');
				currentRoute = 'practice-session';
			} else if (path === '/') {
				console.log('Setting route to dashboard');
				currentRoute = 'dashboard';
			} else {
				// Remove leading slash and use the path as route
				const routeName = path.startsWith('/') ? path.slice(1) : path;
				console.log('Setting route to:', routeName);
				currentRoute = routeName;
			}
		}
		
		// Initial route
		handleRouting();
		
		// Add navigation event listener for SPA navigation
		window.addEventListener('popstate', () => {
			console.log('Navigation event, path:', window.location.pathname);
			handleRouting();
		});
		
		// Intercept link clicks for SPA navigation
		document.addEventListener('click', (e) => {
			const target = e.target as HTMLElement;
			const anchor = target.closest('a');
			
			if (anchor && anchor.getAttribute('href')?.startsWith('/') && !anchor.hasAttribute('target')) {
				const href = anchor.getAttribute('href') || '/';
				
				// Skip interception for asset files and other static resources
				if (href.startsWith('/assets/') || 
					href.endsWith('.svg') || 
					href.endsWith('.png') || 
					href.endsWith('.jpg') || 
					href.endsWith('.jpeg') || 
					href.endsWith('.gif') || 
					href.endsWith('.ico')) {
					return;
				}
				
				e.preventDefault();
				history.pushState(null, '', href);
				handleRouting();
			}
		});
	});

	// Navigate function to use throughout the app
	function navigate(path: string) {
		history.pushState(null, '', path);
		
		const newPath = window.location.pathname || '/';
		console.log('Navigated to path:', newPath);
		
		// Update current route based on new path
		if (newPath.includes('/practice-topic/')) {
			currentRoute = 'practice-topic';
		} else if (newPath.includes('/create-practice/')) {
			currentRoute = 'create-practice';
		} else if (newPath.includes('/practice-session/')) {
			currentRoute = 'practice-session';
		} else if (newPath === '/') {
			currentRoute = 'dashboard';
		} else {
			const routeName = newPath.startsWith('/') ? newPath.slice(1) : newPath;
			currentRoute = routeName;
		}
	}

	// Add navigate to the window so components can use it
	if (typeof window !== 'undefined') {
		(window as any).navigate = navigate;
	}

	// Handle browser navigation
	window.addEventListener('popstate', (event) => {
		console.log('Popstate event triggered', event);
		
		// Let the main routing handle this
		const path = window.location.pathname || '/';
		console.log('Current path:', path);
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
		{#if currentRoute === 'dashboard'}
			<Dashboard />
		{:else if currentRoute === 'profile'}
			<Profile />
		{:else if currentRoute === 'settings'}
			<Settings />
		{:else if currentRoute === 'practice-topics'}
			<PracticeTopics />
		{:else if currentRoute === 'practice-topic'}
			<ViewPracticeTopic />
		{:else if currentRoute === 'create-practice'}
			<CreatePracticeSession />
		{:else if currentRoute === 'learners'}
			<Learners />
		{:else if currentRoute === 'chat'}
			<Chat />
		{:else}
			<Dashboard />
		{/if}
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
