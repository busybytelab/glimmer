<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { goto } from '$app/navigation';
	import pb from '$lib/pocketbase';
	import type { PracticeSession } from '$lib/types';
	import FormField from '../common/FormField.svelte';
	import FormButton from '../common/FormButton.svelte';
	import FormSection from '../common/FormSection.svelte';
	import ErrorAlert from '../common/ErrorAlert.svelte';
	import SelectField from '../common/SelectField.svelte';
	import LearnersList from '../learners/LearnersList.svelte';

	export let session: PracticeSession | null = null;

	const dispatch = createEventDispatcher<{
		update: PracticeSession;
		delete: string;
		cancel: void;
	}>();

	type FormData = {
		name: string;
		status: string;
		learner: string;
	};

	let formData: FormData = session ? {
		name: session.name || '',
		status: session.status,
		learner: session.learner
	} : {
		name: '',
		status: 'Generated',
		learner: ''
	};

	let loading = false;
	let error: string | null = null;
	let learners: any[] = [];
	let loadingLearners = false;
	let selectedLearner: any = null;
	let showLearnerSelection = false;

	// Load learners when component mounts
	fetchLearners();

	async function fetchLearners() {
		try {
			loadingLearners = true;
			const result = await pb.collection('learners').getList(1, 50, {
				sort: 'nickname',
				expand: 'user'
			});
			learners = result.items;

			// If we have the current learner, find them in the list
			if (session?.learner) {
				selectedLearner = learners.find(l => l.id === session.learner);
			}
		} catch (err) {
			console.error('Error fetching learners:', err);
			error = 'Failed to load learners';
		} finally {
			loadingLearners = false;
		}
	}

	// Validate form data
	function validateForm() {
		if (!formData.name) return "Name is required";
		if (!formData.status) return "Status is required";
		if (!formData.learner) return "Learner is required";
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

			let result;
			if (session) {
				// Update existing session
				result = await pb.collection('practice_sessions').update(session.id, formData);
			} else {
				// Create new session (should not happen in this form)
				error = "Cannot create new sessions from this form";
				loading = false;
				return;
			}
						
			// Dispatch the update event with the result
			dispatch('update', result as unknown as PracticeSession);
			
			// Navigate to the instructor view
			goto(`/practice-sessions/${session.id}/instructor`);
		} catch (err) {
			console.error('Failed to save session:', err);
			error = 'Failed to save practice session';
		} finally {
			loading = false;
		}
	}

	async function handleDelete() {
		if (!session) return;

		if (!confirm('Are you sure you want to delete this practice session?')) {
			return;
		}

		try {
			loading = true;
			error = null;
			await pb.collection('practice_sessions').delete(session.id);
			dispatch('delete', session.id);
			
			// Navigate back to practice topics or dashboard
			if (session.expand?.practice_topic) {
				goto(`/practice-topics/${session.expand.practice_topic.id}`);
			} else {
				goto('/dashboard');
			}
		} catch (err) {
			console.error('Failed to delete session:', err);
			error = 'Failed to delete practice session';
		} finally {
			loading = false;
		}
	}

	function handleCancel() {
		if (!session) return;
		dispatch('cancel');
		goto(`/practice-sessions/${session.id}/instructor`);
	}

	function handleLearnerSelect(learner: any) {
		selectedLearner = learner;
		formData.learner = learner.id;
		showLearnerSelection = false;
	}
</script>

<FormSection title="Edit Practice Session" description="Modify the details for this practice session">
	<form on:submit|preventDefault={handleSubmit} class="w-full">
		{#if error}
			<div class="px-4 mb-6">
				<ErrorAlert message={error} />
			</div>
		{/if}

		<div class="px-4 py-4 bg-white dark:bg-gray-800 w-full">
			<div class="grid grid-cols-6 gap-6">
				<FormField 
					id="name"
					label="Name"
					bind:value={formData.name}
					disabled={loading}
					required={true}
					cols="col-span-6"
					placeholder="Enter session name"
				/>

				<SelectField 
					id="status"
					label="Status"
					bind:value={formData.status}
					disabled={loading}
					required={true}
					cols="col-span-6 sm:col-span-3"
				>
					<option value="Generated">Generated</option>
					<option value="InProgress">In Progress</option>
					<option value="Completed">Completed</option>
				</SelectField>

				<div class="col-span-6">
					<div class="flex flex-col">
						<label for="learner-select" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							Learner <span class="text-red-500">*</span>
						</label>
						
						{#if selectedLearner}
							<div class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-md mb-2">
								<div class="flex items-center">
									<div class="bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300 p-2 rounded-full mr-3">
										<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
										</svg>
									</div>
									<div>
										<p class="font-medium">{selectedLearner.nickname}</p>
										{#if selectedLearner.grade_level}
											<p class="text-xs text-gray-500 dark:text-gray-400">Grade: {selectedLearner.grade_level}</p>
										{/if}
									</div>
								</div>
								<button 
									type="button"
									class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 text-sm"
									on:click={() => showLearnerSelection = true}
								>
									Change
								</button>
							</div>
						{:else}
							<button 
								id="learner-select"
								type="button"
								class="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-offset-gray-800"
								on:click={() => showLearnerSelection = true}
							>
								Select Learner
							</button>
						{/if}
					</div>

					{#if showLearnerSelection}
						<div class="mt-4 border border-gray-200 dark:border-gray-700 rounded-lg p-4 bg-gray-50 dark:bg-gray-700">
							<div class="flex justify-between items-center mb-4">
								<h3 class="text-lg font-medium text-gray-900 dark:text-white">Select a Learner</h3>
								<button 
									type="button"
									aria-label="Close learner selection"
									class="text-gray-500 hover:text-gray-700 dark:text-gray-300 dark:hover:text-gray-100"
									on:click={() => showLearnerSelection = false}
								>
									<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
										<path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
									</svg>
								</button>
							</div>
							
							<LearnersList
								{learners}
								loading={loadingLearners}
								emptyMessage="No learners found. Please add learners to your account first."
								gridCols="grid-cols-1 md:grid-cols-2"
								showPreferences={false}
								onClick={handleLearnerSelect}
							/>
						</div>
					{/if}
				</div>
			</div>
		</div>

		<div class="px-4 py-3 bg-gray-50 dark:bg-gray-700 border-t border-gray-200 dark:border-gray-600 flex justify-end space-x-4">
			<FormButton 
				type="button" 
				variant="secondary"
				disabled={loading}
				on:click={handleCancel}
			>
				Cancel
			</FormButton>

			{#if session}
				<FormButton 
					type="button" 
					variant="danger"
					disabled={loading}
					on:click={handleDelete}
				>
					Delete
				</FormButton>
			{/if}

			<FormButton 
				type="submit" 
				variant="primary"
				disabled={loading}
			>
				Save Changes
			</FormButton>
		</div>
	</form>
</FormSection> 