<script lang="ts">
	import { onMount } from 'svelte';
	import pb from '$lib/pocketbase';
	import { goto } from '$app/navigation';
	import type { Learner } from '$lib/types';
	import LearnersList from '$components/learners/LearnersList.svelte';
	import FormButton from '$components/common/FormButton.svelte';
	import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
	import ErrorAlert from '$components/common/ErrorAlert.svelte';
	import { learnersService } from '$lib/services/learners';
	
	let learners: Learner[] = [];
	let loading = true;
	let error: string | null = null;

	async function loadLearners() {
		try {
			loading = true;
			error = null;

			// Check if user is authenticated
			const authData = pb.authStore.model;
			if (!authData) {
				error = 'You must be logged in to view learners';
				return;
			}

			
			// Get learners from the same account
			const result = await learnersService.getLearners();
			
			// Map expanded data to the learner objects
			learners = result.map(item => {
				console.log('learner Item:', item);
				const expandedData = item.expand;
				return {
					...item,
					account: expandedData?.account || null
				};
			}) as unknown as Learner[];
			
			console.log('Loaded learners:', learners);
		} catch (err) {
			console.error('Failed to load learners:', err);
			error = 'Failed to load learners';
		} finally {
			loading = false;
		}
	}

	onMount(async () => {
		console.log('Learners mounted, loading learners');
		await loadLearners();
	});

	function handleCreateNew() {
		goto('/account/learners/create');
	}

	function viewLearner(learner: Learner) {
		goto(`/account/learners/${learner.id}`);
	}

	function editLearner(learner: Learner) {
		goto(`/account/learners/${learner.id}/edit`);
	}
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Children Profiles</h1>
		<FormButton
			type="button"
			on:click={handleCreateNew}
		>
			Add Child
		</FormButton>
	</div>

	{#if loading}
		<div class="flex justify-center items-center h-64">
			<LoadingSpinner size="md" color="primary" />
		</div>
	{:else if error}
		<ErrorAlert message={error} />
	{:else}
	<LearnersList
	{learners}
	{loading}
	emptyMessage="No learners found. Create your first learner!"
	onClick={viewLearner}
	onEdit={editLearner}
	cardActions={[]}
/>
	{/if}
</div> 