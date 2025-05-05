<script lang="ts">
	import { user, error } from '$lib/stores';
	import { onMount } from 'svelte';
	import AppLayout from '../components/layout/AppLayout.svelte';
	import EditProfile from '../components/settings/EditProfile.svelte';
	import pb from '$lib/pocketbase';

	let sidebarOpen = true;
	let isLoading = false;
	let userProfile: { name: string; email: string } | null = null;
	let savedSuccessfully = false;

	onMount(async () => {
		await fetchUserProfile();
	});

	async function fetchUserProfile() {
		isLoading = true;
		error.set(null);

		try {
			const currentUser = $user;
			if (!currentUser) {
				throw new Error('User not found');
			}

			// Here we would fetch additional profile data if needed
			userProfile = {
				name: currentUser.user?.name || '',
				email: currentUser.user?.email || '',
			};
		} catch (err) {
			console.error('Error fetching user profile:', err);
			error.set(err instanceof Error ? err.message : 'Failed to fetch user profile');
		} finally {
			isLoading = false;
		}
	}

	async function handleSaveProfile(event: CustomEvent<{ name: string }>) {
		isLoading = true;
		error.set(null);
		savedSuccessfully = false;

		try {
			const { name } = event.detail;
			const currentUser = $user;
			
			if (!currentUser || !currentUser.user) {
				throw new Error('User not found');
			}
			
			// Update user in PocketBase
			const userId = currentUser.user.id;
			await pb.collection('users').update(userId, { name });
			
			// Update local user data
			if (currentUser.user) {
				currentUser.user.name = name;
				user.set(currentUser);
			}
			
			// Update local profile
			if (userProfile) {
				userProfile = {
					...userProfile,
					name
				};
			}
			
			// Show success message briefly
			savedSuccessfully = true;
			setTimeout(() => {
				savedSuccessfully = false;
			}, 3000);
			
		} catch (err) {
			console.error('Error updating profile:', err);
			error.set(err instanceof Error ? err.message : 'Failed to update profile');
		} finally {
			isLoading = false;
		}
	}
</script>

<AppLayout bind:sidebarOpen>
	{#if savedSuccessfully}
		<div class="max-w-3xl mx-auto px-4 sm:px-6 md:px-8 my-4">
			<div class="bg-green-50 border-l-4 border-green-500 text-green-700 p-4 rounded-md shadow-sm" role="alert">
				<div class="flex">
					<div class="flex-shrink-0">
						<svg class="h-5 w-5 text-green-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
						</svg>
					</div>
					<div class="ml-3">
						<p class="text-sm font-medium">Your profile has been updated successfully!</p>
					</div>
				</div>
			</div>
		</div>
	{/if}
	
	<EditProfile 
		profile={userProfile} 
		{isLoading}
		on:save={handleSaveProfile}
	/>
</AppLayout> 