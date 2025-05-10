<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import type { Learner } from '$lib/types';
	import pb from '$lib/pocketbase';
	import LearnerForm from '../../../../components/learners/LearnerForm.svelte';
	import FormButton from '../../../../components/common/FormButton.svelte';

	let learner: Learner | null = null;
	let loading = true;
	let error: string | null = null;

	onMount(async () => {
		// Get learner ID from the URL parameter
		const learnerId = $page.params.id;
		
		if (!learnerId) {
			error = 'No learner ID provided';
			loading = false;
			return;
		}
		
		await loadLearner(learnerId);
	});

	async function loadLearner(id: string) {
		try {
			loading = true;
			error = null;
			
			const result = await pb.collection('learners').getOne<Learner>(id, {
				expand: 'user'
			});
			
			// Map expanded user data to the learner object
			const expandedData = result.expand;
			learner = {
				...result,
				user: expandedData?.user || { name: 'Unknown user' }
			} as unknown as Learner;
		} catch (err) {
			console.error('Failed to load learner:', err);
			error = 'Failed to load learner';
		} finally {
			loading = false;
		}
	}

	function handleLearnerUpdate(updatedLearner: Learner) {
		goto('/learners');
	}

	function handleLearnerDelete(learnerId: string) {
		goto('/learners');
	}

	function handleCancel() {
		goto('/learners');
	}
</script>

<div class="container mx-auto px-4 py-8 max-w-7xl">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-2xl font-bold text-gray-900">
			{learner ? 'Edit Learner' : 'Loading...'}
		</h1>
		<FormButton
			type="button"
			variant="secondary"
			on:click={handleCancel}
		>
			Back to Learners
		</FormButton>
	</div>

	{#if loading}
		<div class="flex justify-center items-center h-64">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary"></div>
		</div>
	{:else if error}
		<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
			<strong class="font-bold">Error!</strong>
			<span class="block sm:inline"> {error}</span>
		</div>
	{:else if learner}
		<div class="w-full">
			<LearnerForm
				{learner}
				on:update={({ detail }) => handleLearnerUpdate(detail)}
				on:delete={({ detail }) => handleLearnerDelete(detail)}
				on:cancel={handleCancel}
			/>
		</div>
	{/if}
</div> 