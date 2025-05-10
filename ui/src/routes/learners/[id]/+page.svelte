<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import type { Learner } from '$lib/types';
	import pb from '$lib/pocketbase';
	import FormButton from '../../../components/common/FormButton.svelte';

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

	function editLearner() {
		if (!learner) return;
		goto(`/learners/edit/${learner.id}`);
	}
</script>

<div class="container mx-auto px-4 py-8">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-2xl font-bold text-gray-900">Learner Profile</h1>
		<FormButton
			type="button"
			variant="primary"
			on:click={editLearner}
		>
			Edit Learner
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
		<div class="bg-white shadow-md rounded-lg p-6">
			<div class="flex items-center space-x-4 mb-6">
				{#if learner.avatar}
					<img src={learner.avatar} alt={learner.nickname} class="w-16 h-16 rounded-full" />
				{:else}
					<div class="w-16 h-16 rounded-full bg-gray-200 flex items-center justify-center">
						<span class="text-2xl text-gray-500">{learner.nickname[0].toUpperCase()}</span>
					</div>
				{/if}
				<div>
					<h2 class="text-xl font-semibold text-gray-900">{learner.nickname}</h2>
					<p class="text-gray-600">{learner.user.name}</p>
				</div>
			</div>

			<div class="grid grid-cols-2 gap-4 mb-6">
				<div class="bg-gray-50 p-3 rounded">
					<div class="text-sm font-semibold text-gray-700 mb-1">Age:</div>
					<div>{learner.age}</div>
				</div>
				
				<div class="bg-gray-50 p-3 rounded">
					<div class="text-sm font-semibold text-gray-700 mb-1">Grade Level:</div>
					<div>{learner.grade_level || 'Not specified'}</div>
				</div>
			</div>

			{#if learner.learning_preferences && learner.learning_preferences.length > 0}
				<div class="mb-6">
					<div class="text-sm font-semibold text-gray-700 mb-2">Learning Preferences:</div>
					<div class="flex flex-wrap gap-2">
						{#each learner.learning_preferences as preference}
							<span class="bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded">
								{preference}
							</span>
						{/each}
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div> 