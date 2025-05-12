<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import pb from '$lib/pocketbase';
	import type { Learner } from '$lib/types';
	import FormField from '../common/FormField.svelte';
	import FormButton from '../common/FormButton.svelte';
	import FormSection from '../common/FormSection.svelte';
	import ErrorAlert from '../common/ErrorAlert.svelte';
	import TextArea from '../common/TextArea.svelte';

	export let learner: Learner | null = null;

	const dispatch = createEventDispatcher<{
		update: Learner;
		delete: string;
		cancel: void;
	}>();

	type FormData = {
		nickname: string;
		age: number;
		grade_level: string;
		learning_preferences: string[];
		avatar: string;
	};

	let formData: FormData = learner ? {
		nickname: learner.nickname,
		age: learner.age,
		grade_level: learner.grade_level || '',
		learning_preferences: Array.isArray(learner.learning_preferences) ? learner.learning_preferences : [],
		avatar: learner.avatar || ''
	} : {
		nickname: '',
		age: 0,
		grade_level: '',
		learning_preferences: [],
		avatar: ''
	};

	let loading = false;
	let error: string | null = null;
	let learningPreferencesText = formData.learning_preferences.join('\n');

	$: {
		formData.learning_preferences = learningPreferencesText.split('\n').map(pref => pref.trim()).filter(Boolean);
	}

	// Validate form data
	function validateForm() {
		if (!formData.nickname.trim()) return "Nickname is required";
		if (formData.age < 0 || formData.age > 120) return "Age must be between 0 and 120";
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
			const dataToSend = { ...formData };
			
			// Make sure learning preferences are properly formatted as an array
			if (typeof dataToSend.learning_preferences === 'string') {
				dataToSend.learning_preferences = (dataToSend.learning_preferences as string).split('\n').map(pref => pref.trim()).filter(Boolean);
			}
						
			let result;
			if (learner) {
				// Update existing learner
				result = await pb.collection('learners').update(learner.id, dataToSend);
			} else {
				// Create new learner
				result = await pb.collection('learners').create(dataToSend);
			}
			
			// Dispatch the update event with the result
			dispatch('update', result as unknown as Learner);
		} catch (err) {
			console.error('Failed to save learner:', err);
			error = 'Failed to save learner';
		} finally {
			loading = false;
		}
	}

	async function handleDelete() {
		if (!learner) return;

		if (!confirm('Are you sure you want to delete this learner?')) {
			return;
		}

		try {
			loading = true;
			error = null;
			await pb.collection('learners').delete(learner.id);
			dispatch('delete', learner.id);
		} catch (err) {
			console.error('Failed to delete learner:', err);
			error = 'Failed to delete learner';
		} finally {
			loading = false;
		}
	}
</script>

<FormSection title={learner ? "Edit Learner" : "Create Learner"} description="Enter the details for this learner">
	<form on:submit|preventDefault={handleSubmit} class="w-full">
		{#if error}
			<div class="px-4 mb-6">
				<ErrorAlert message={error} />
			</div>
		{/if}

		<div class="px-4 py-4 bg-white dark:bg-gray-800 w-full">
			<div class="grid grid-cols-6 gap-6">
				<FormField
					label="Nickname"
					id="nickname"
					type="text"
					bind:value={formData.nickname}
					required
					cols="col-span-6"
					placeholder="Enter learner's nickname"
				/>

				<FormField
					label="Age"
					id="age"
					type="number"
					value={formData.age.toString()}
					on:input={(e) => {
						const target = e.target as HTMLInputElement;
						formData.age = parseInt(target.value) || 0;
					}}
					required
					cols="col-span-6 sm:col-span-3"
					placeholder="Enter age"
				/>

				<FormField
					label="Grade Level"
					id="gradeLevel"
					type="text"
					bind:value={formData.grade_level}
					cols="col-span-6 sm:col-span-3"
					placeholder="e.g. 3rd Grade"
				/>

				<TextArea
					id="learning_preferences"
					label="Learning Preferences"
					bind:value={learningPreferencesText}
					cols="col-span-6"
					placeholder="Enter learning preferences, one per line"
				/>

				<FormField
					label="Avatar URL"
					id="avatar"
					type="text"
					bind:value={formData.avatar}
					cols="col-span-6"
					placeholder="Enter avatar image URL"
				/>
			</div>
		</div>

		<div class="px-4 py-3 bg-gray-50 dark:bg-gray-700 border-t border-gray-200 dark:border-gray-600 flex justify-end space-x-4">
			<FormButton
				type="button"
				variant="secondary"
				on:click={() => dispatch('cancel')}
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