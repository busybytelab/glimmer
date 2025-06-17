<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import pb from '$lib/pocketbase';
	import type { Learner } from '$lib/types';
	import FormField from '../common/FormField.svelte';
	import FormButton from '../common/FormButton.svelte';
	import FormSection from '../common/FormSection.svelte';
	import ErrorAlert from '../common/ErrorAlert.svelte';
	import TextArea from '../common/TextArea.svelte';
	import AvatarSelector from './AvatarSelector.svelte';

	export let learner: Learner | null = null;

	const dispatch = createEventDispatcher<{
		update: { detail: Learner };
		delete: { detail: string };
		cancel: { detail: void };
	}>();

	type FormData = {
		nickname: string;
		age: number;
		grade_level: string;
		learning_preferences: string[];
		avatar: string;
		account?: string; // Optional for updates, required for create
	};

	let formData: FormData = learner ? {
		nickname: learner.nickname,
		age: learner.age,
		grade_level: learner.grade_level || '',
		learning_preferences: Array.isArray(learner.learning_preferences) ? learner.learning_preferences : [],
		avatar: learner.avatar || ''
	} : {
		nickname: '',
		age: 4,
		grade_level: '',
		learning_preferences: [],
		avatar: ''
	};

	let ageInput = formData.age.toString();

	let loading = false;
	let error: string | null = null;
	let learningPreferencesText = formData.learning_preferences.join('\n');

	$: {
		formData.learning_preferences = learningPreferencesText.split('\n').map(pref => pref.trim()).filter(Boolean);
		const parsed = parseInt(ageInput);
		formData.age = isNaN(parsed) ? 1 : Math.min(30, Math.max(1, parsed));
	}

	// Validate form data
	function validateForm() {
		console.log('Validating form data:', formData);
		if (!formData.nickname?.trim()) return "Nickname is required";
		const age = parseInt(formData.age.toString());
		console.log('Parsed age:', age);
		if (isNaN(age) || age < 1 || age > 30) return "Age must be between 1 and 30";
		return null;
	}

	async function handleSubmit() {
		const validationError = validateForm();
		if (validationError) {
			error = validationError;
			return;
		}
	
		try {
			loading = true;
			error = null;

			// Get the current user
			const currentUser = pb.authStore.model;
			if (!currentUser && !learner) {
				error = 'You must be logged in to create a learner';
				loading = false;
				return;
			}

			// Prepare data to send
			const dataToSend = {
				...formData,
				nickname: formData.nickname.trim(),
				age: parseInt(formData.age.toString())
			};
			
			// Make sure learning preferences are properly formatted as an array
			if (typeof dataToSend.learning_preferences === 'string') {
				dataToSend.learning_preferences = (dataToSend.learning_preferences as string).split('\n').map(pref => pref.trim()).filter(Boolean);
			}

			// Add required account field for new learners
			if (!learner) {
				// Get the user's account
				try {
					const account = await pb.collection('accounts').getFirstListItem(`owner="${currentUser?.id}"`);
					dataToSend.account = account.id;
				} catch (err) {
					console.error('Failed to get user account:', err);
					error = 'Failed to get user account';
					loading = false;
					return;
				}
			}

			console.log('Sending data:', dataToSend);
						
			let result;
			if (learner) {
				// Update existing learner
				result = await pb.collection('learners').update(learner.id, dataToSend);
			} else {
				// Create new learner
				result = await pb.collection('learners').create(dataToSend);
			}
			
			// Dispatch the update event with the result
			dispatch('update', { detail: result as unknown as Learner });
		} catch (err: any) {
			console.error('Failed to save learner:', err);
			// Show more detailed error message from PocketBase if available
			if (err?.response?.data?.message) {
				error = err.response.data.message;
			} else if (err?.response?.data) {
				// Handle validation errors
				const validationErrors = Object.entries(err.response.data)
					.map(([field, message]) => `${field}: ${message}`)
					.join(', ');
				error = `Validation failed: ${validationErrors}`;
			} else if (err?.message) {
				error = err.message;
			} else {
				error = 'Failed to save learner';
			}
		} finally {
			loading = false;
		}
	}

	async function handleDelete() {
		if (!learner) return;

		if (!confirm('Are you sure you want to delete this profile?')) {
			return;
		}

		try {
			loading = true;
			error = null;
			await pb.collection('learners').delete(learner.id);
			dispatch('delete', { detail: learner.id });
		} catch (err) {
			console.error('Failed to delete learner:', err);
			error = 'Failed to delete learner';
		} finally {
			loading = false;
		}
	}
</script>

<FormSection title={learner ? "Edit Child Profile" : "Add Child"} description="Enter the details for child profile">
	<form on:submit|preventDefault={handleSubmit} class="w-full">
		{#if error}
			<div class="px-4 mb-6">
				<ErrorAlert message={error} />
			</div>
		{/if}

		<div class="px-4 py-4 bg-white dark:bg-gray-800 w-full">
			<div class="grid grid-cols-6 gap-6">
				<FormField
					label="Name"
					id="name"
					type="text"
					bind:value={formData.nickname}
					required
					cols="col-span-6 sm:col-span-3"
					placeholder="Enter child's name"
				/>

				<FormField
					label="Age"
					id="age"
					type="number"
					bind:value={ageInput}
					required
					cols="col-span-6 sm:col-span-3"
					placeholder="Enter age (1-30)"
				/>

				<FormField
					label="Grade Level"
					id="gradeLevel"
					type="text"
					bind:value={formData.grade_level}
					cols="col-span-6 sm:col-span-3"
					placeholder="e.g. 3rd, 4th, etc"
				/>

				<TextArea
					id="learning_preferences"
					label="Learning Preferences"
					bind:value={learningPreferencesText}
					cols="col-span-6"
					placeholder="Enter learning preferences such as visual, auditory, kinesthetic, one per line"
				/>

				<div class="col-span-6">
					<label for="avatar" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
						Avatar
					</label>
					<AvatarSelector
						selectedAvatar={formData.avatar}
						on:select={(e) => formData.avatar = e.detail}
					/>
				</div>
			</div>
		</div>

		<div class="px-4 py-3 bg-gray-50 dark:bg-gray-700 border-t border-gray-200 dark:border-gray-600 flex justify-end space-x-4">
			<FormButton
				type="button"
				variant="secondary"
				on:click={() => dispatch('cancel', { detail: undefined })}
			>
				Cancel
			</FormButton>
			{#if learner}
				<FormButton
					type="button"
					variant="danger"
					on:click={handleDelete}
					disabled={loading}
				>
					Delete
				</FormButton>
			{/if}
			<FormButton
				type="submit"
				isLoading={loading}
				loadingText="Saving..."
			>
				Save
			</FormButton>
		</div>
	</form>
</FormSection> 