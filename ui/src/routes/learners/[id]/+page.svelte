<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import type { Learner } from '$lib/types';
	import FormButton from '$components/common/FormButton.svelte';
	import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
	import ErrorAlert from '$components/common/ErrorAlert.svelte';
	import { learnersService } from '$lib/services/learners';
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
		
		learner = await learnersService.getLearner(learnerId);
	});

	

	function editLearner() {
		if (!learner) return;
		goto(`/learners/${learner.id}/edit`);
	}
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Learner Profile</h1>
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
			<LoadingSpinner size="md" color="primary" />
		</div>
	{:else if error}
		<ErrorAlert message={error} />
	{:else if learner}
		<div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
			<div class="flex items-center space-x-4 mb-6">
				{#if learner.avatar}
					<img src={learner.avatar} alt={learner.nickname || 'Learner'} class="w-16 h-16 rounded-full" />
				{:else}
					<div class="w-16 h-16 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
						<span class="text-2xl text-gray-500 dark:text-gray-300">{learner.nickname?.[0]?.toUpperCase() || 'L'}</span>
					</div>
				{/if}
				<div>
					<h2 class="text-xl font-semibold text-gray-900 dark:text-white">{learner.nickname || 'Unknown learner'}</h2>
				</div>
			</div>

			<div class="grid grid-cols-2 gap-4 mb-6">
				<div class="bg-gray-50 dark:bg-gray-700 p-3 rounded">
					<div class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">Age:</div>
					<div class="dark:text-gray-200">{learner.age}</div>
				</div>
				
				<div class="bg-gray-50 dark:bg-gray-700 p-3 rounded">
					<div class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">Grade Level:</div>
					<div class="dark:text-gray-200">{learner.grade_level || 'Not specified'}</div>
				</div>
			</div>

			{#if learner.learning_preferences && learner.learning_preferences.length > 0}
				<div class="mb-6">
					<div class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">Learning Preferences:</div>
					<div class="flex flex-wrap gap-2">
						{#each learner.learning_preferences as preference}
							<span class="bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300 text-xs font-medium px-2.5 py-0.5 rounded">
								{preference}
							</span>
						{/each}
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div> 