<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import pb from '$lib/pocketbase';
	import type { PracticeTopic } from '$lib/types';
	import FormField from '../common/FormField.svelte';
	import TextArea from '../common/TextArea.svelte';
	import FormButton from '../common/FormButton.svelte';
	import FormSection from '../common/FormSection.svelte';
	import ErrorAlert from '../common/ErrorAlert.svelte';
	import SelectField from '../common/SelectField.svelte';
	import { llmService } from '$lib/services/llm';

	export let topic: PracticeTopic | null = null;

	const dispatch = createEventDispatcher<{
		update: PracticeTopic;
		delete: string;
		cancel: void;
	}>();

	type FormData = {
		name: string;
		subject: string;
		description: string;
		target_age_range: string;
		target_grade_level: string;
		learning_goals: string[];
		base_prompt: string;
		system_prompt: string;
		tags: string[];
		instructor?: string;
		account?: string;
		llm_model?: string;
	};

	let formData: FormData = topic ? {
		name: topic.name,
		subject: topic.subject,
		description: topic.description || '',
		target_age_range: topic.target_age_range || '',
		target_grade_level: topic.target_grade_level || '',
		learning_goals: Array.isArray(topic.learning_goals) ? topic.learning_goals : [],
		base_prompt: topic.base_prompt || '',
		system_prompt: topic.system_prompt || '',
		tags: Array.isArray(topic.tags) ? topic.tags : [],
		llm_model: topic.llm_model || ''
	} : {
		name: '',
		subject: '',
		description: '',
		target_age_range: '',
		target_grade_level: '',
		learning_goals: [],
		base_prompt: '',
		system_prompt: '',
		tags: [],
		llm_model: ''
	};

	let loading = false;
	let error: string | null = null;
	let learningGoalsText = formData.learning_goals.join('\n');
	let tagsText = formData.tags.join(', ');
	let formSubmitted = false;
	let isLoadingModels = false;
	let modelError: string | null = null;
	let availableModels: { id: string; name: string; isDefault?: boolean }[] = [];

	// Fetch available models from backend
	async function fetchModels() {
		isLoadingModels = true;
		modelError = null;
		
		try {
			const data = await llmService.getInfo();
			
			// Transform the models data into the format we need
			availableModels = data.platforms.flatMap(platform => 
				platform.models.map(model => ({
					id: model.name,
					name: `${model.name}${model.isDefault ? ' (Default)' : ''}`,
					isDefault: model.isDefault
				}))
			);
			
			// If no models were found, add a default option
			if (availableModels.length === 0) {
				availableModels = [{ id: "", name: "Default" }];
			}
			
			// Select the default model if available and no model is selected
			if (!formData.llm_model) {
				const defaultModel = availableModels.find(m => m.isDefault);
				if (defaultModel) {
					formData.llm_model = defaultModel.id;
				}
			}
		} catch (err) {
			console.error('Error fetching models:', err);
			modelError = err instanceof Error ? err.message : 'Failed to fetch models';
			// Fallback to default model
			availableModels = [{ id: "", name: "Default" }];
		} finally {
			isLoadingModels = false;
		}
	}

	// Fetch models when component mounts
	fetchModels();

	$: {
		formData.learning_goals = learningGoalsText.split('\n').map(goal => goal.trim()).filter(Boolean);
		formData.tags = tagsText.split(',').map(tag => tag.trim()).filter(Boolean);
	}

	// Validate form data
	function validateForm() {
		if (!formData.name) return "Name is required";
		if (!formData.subject) return "Subject is required";
		if (!formData.base_prompt) return "Base prompt is required";
		return null;
	}

	async function handleSubmit() {
		formSubmitted = true;
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
			if (!currentUser && !topic) {
				error = 'You must be logged in to create a topic';
				loading = false;
				return;
			}

			// Prepare data to send
			const dataToSend = { ...formData };
			
			// Make sure tags are properly formatted as an array
			if (typeof dataToSend.tags === 'string') {
				dataToSend.tags = (dataToSend.tags as string).split(',').map(tag => tag.trim()).filter(Boolean);
			}
			
			console.log('Data being sent to server:', dataToSend);
			
			if (!topic) {
				// For new topics, get the user's account and instructor info
				try {
					// Make sure currentUser exists
					if (!currentUser) {
						error = 'You must be logged in to create a topic';
						loading = false;
						return;
					}
					
					// First try to get instructor record for current user
					const instructors = await pb.collection('instructors').getList(1, 1, {
						filter: `user.id = "${currentUser.id}"`
					});
					
					if (instructors && instructors.items.length > 0) {
						dataToSend.instructor = instructors.items[0].id;
						dataToSend.account = instructors.items[0].account;
					} else {
						// If not an instructor, try to get account directly
						const accounts = await pb.collection('accounts').getList(1, 1, {
							filter: `owner.id = "${currentUser.id}"`
						});
						
						if (accounts && accounts.items.length > 0) {
							dataToSend.account = accounts.items[0].id;
						} else {
							error = 'Could not determine account for user';
							loading = false;
							return;
						}
					}
				} catch (err) {
					console.error('Failed to get user account info:', err);
					error = 'Failed to get account information';
					loading = false;
					return;
				}
			}

			let result;
			if (topic) {
				// Update existing topic
				result = await pb.collection('practice_topics').update(topic.id, dataToSend);
			} else {
				// Create new topic
				result = await pb.collection('practice_topics').create(dataToSend);
			}
			
			// Log the response from the server to see what tags look like
			console.log('Response from server:', result);
			
			// Dispatch the update event with the result
			dispatch('update', result as unknown as PracticeTopic);
		} catch (err) {
			console.error('Failed to save topic:', err);
			error = 'Failed to save practice topic';
		} finally {
			loading = false;
		}
	}

	function handleTagsInput(event: Event) {
		const inputElement = event.target as HTMLInputElement;
		tagsText = inputElement.value;
		formData.tags = tagsText.split(',').map(tag => tag.trim()).filter(Boolean);
		console.log('Tags updated through input:', formData.tags);
	}

	async function handleDelete() {
		if (!topic) return;

		if (!confirm('Are you sure you want to delete this practice topic?')) {
			return;
		}

		try {
			loading = true;
			error = null;
			await pb.collection('practice_topics').delete(topic.id);
			dispatch('delete', topic.id);
		} catch (err) {
			console.error('Failed to delete topic:', err);
			error = 'Failed to delete practice topic';
		} finally {
			loading = false;
		}
	}
</script>

<div class="mx-auto px-4 sm:px-6 md:px-8">
<FormSection title={topic ? "Edit Practice Topic" : "Create Practice Topic"} description="Enter the details for this practice topic">
	<form on:submit|preventDefault={handleSubmit} class="p-4">
		{#if error}
			<div class="mb-6">
				<ErrorAlert message={error} />
			</div>
		{/if}

		<div class="py-4 bg-white">
			<div class="grid grid-cols-6 gap-6">
				<FormField 
					id="name"
					label="Name"
					bind:value={formData.name}
					disabled={loading}
					required={true}
					cols="col-span-6"
					placeholder="Enter topic name"
				/>

				<FormField 
					id="subject"
					label="Subject"
					bind:value={formData.subject}
					disabled={loading}
					required={true}
					cols="col-span-6 sm:col-span-3"
					placeholder="e.g. Mathematics, Science, History"
				/>

				<FormField 
					id="target_age_range"
					label="Target Age Range"
					bind:value={formData.target_age_range}
					disabled={loading}
					cols="col-span-6 sm:col-span-3"
					placeholder="e.g. 7-10"
				/>

				<FormField 
					id="target_grade_level"
					label="Target Grade Level"
					bind:value={formData.target_grade_level}
					disabled={loading}
					cols="col-span-6 sm:col-span-3"
					placeholder="e.g. 3-5"
				/>

				<TextArea
					id="description"
					label="Description"
					bind:value={formData.description}
					disabled={loading}
					cols="col-span-6"
					placeholder="Describe the topic and its purpose"
				/>

				<TextArea
					id="learning_goals"
					label="Learning Goals (one per line)"
					bind:value={learningGoalsText}
					disabled={loading}
					cols="col-span-6"
					placeholder="Enter learning goals, one per line"
				/>

				<TextArea
					id="base_prompt"
					label="Base Prompt"
					bind:value={formData.base_prompt}
					disabled={loading}
					required={true}
					cols="col-span-6"
					placeholder="Enter a base prompt for generating practice items"
				/>

				<TextArea
					id="system_prompt"
					label="System Prompt"
					bind:value={formData.system_prompt}
					disabled={loading}
					cols="col-span-6"
					placeholder="Enter a system prompt for the AI assistant"
				/>

				<SelectField 
					id="llm_model"
					label="LLM Model"
					bind:value={formData.llm_model}
					disabled={loading || isLoadingModels}
					cols="col-span-6 sm:col-span-3"
				>
					{#if isLoadingModels}
						<option value="">Loading models...</option>
					{:else if modelError}
						<option value="">Error loading models</option>
					{:else}
						{#each availableModels as model}
							<option value={model.id}>{model.name}</option>
						{/each}
					{/if}
				</SelectField>

				<FormField 
					id="tags"
					label="Tags (comma-separated)"
					type="text"
					bind:value={tagsText}
					on:input={handleTagsInput}
					disabled={loading}
					cols="col-span-6"
					placeholder="e.g. math, arithmetic, addition"
				/>
				{#if formData.tags.length > 0}
					<div class="mt-2 flex flex-wrap gap-2">
						{#each formData.tags as tag}
							<span class="bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded">
								{tag}
							</span>
						{/each}
					</div>
				{/if}
			</div>
		</div>

		<div class="px-4 py-3 bg-gray-50 text-right sm:px-6 mt-4 -mx-4 -mb-4 border-t border-gray-200 flex justify-end space-x-4">
			<FormButton 
				type="button" 
				variant="secondary"
				disabled={loading}
				on:click={() => dispatch('cancel')}
			>
				Cancel
			</FormButton>

			{#if topic}
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
				isLoading={loading}
				loadingText="Saving..."
			>
				Save
			</FormButton>
		</div>
	</form>
</FormSection>
</div> 