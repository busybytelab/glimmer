<script lang="ts">
    import { onMount } from 'svelte';
    import EditProfile from './EditProfile.svelte';
    import { toast } from '$lib/stores/toast';
    import type { User } from '$lib/types';
    import { userService } from '$lib/services/user';

    let isLoading = false;
    let userProfile: User | null = null;

    onMount(async () => {
        userProfile = await userService.getCurrentUser();
    });

    async function handleSaveProfile(event: CustomEvent) {
        try {
            await saveProfile(event.detail);
            toast.success('Your profile has been updated successfully!');
        } catch (error) {
            console.error('Failed to save profile:', error);
            toast.error('Failed to update profile');
        }
    }

    async function saveProfile(data: { name: string }) {
        isLoading = true;

        try {
            const currentUser = await userService.getCurrentUser();
            
            if (!currentUser) {
                throw new Error('No user found');
            }

            currentUser.name = data.name;

            // Update the user's profile
            const updatedUser = await userService.updateUser(currentUser);

            // Update the local user store
            if (currentUser) {
                userProfile = updatedUser;
            }
        } catch (err) {
            console.error('Error updating profile:', err);
            toast.error('Failed to update profile');
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden mb-6">
    <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
        <h2 class="text-lg font-medium text-gray-900 dark:text-white">Profile</h2>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Manage your personal information
        </p>
    </div>

    <div class="p-0">
        <EditProfile 
            profile={userProfile} 
            {isLoading}
            on:save={handleSaveProfile}
        />
    </div>
</div> 