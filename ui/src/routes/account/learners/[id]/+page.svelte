<script lang="ts">
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { goto } from "$app/navigation";
	import type { Learner } from "$lib/types";
	import FormButton from "$components/common/FormButton.svelte";
	import LoadingSpinner from "$components/common/LoadingSpinner.svelte";
	import ErrorAlert from "$components/common/ErrorAlert.svelte";
	import { learnersService } from "$lib/services/learners";
	import LearnerProgress from "$components/learners/LearnerProgress.svelte";

	let learner: Learner | null = null;
	let loading = false;
	let error: string | null = null;

	onMount(async () => {
		try {
			loading = true;

			// Get learner ID from the URL parameter
			const learnerId = $page.params.id;

			if (!learnerId) {
				error = "No learner ID provided";
				loading = false;
				return;
			}

			learner = await learnersService.getLearner(learnerId);
		} catch (err) {
			console.error("Error loading learner:", err);
			error = "Failed to load learner";
		} finally {
			loading = false;
		}
	});

	function editLearner() {
		if (!learner) return;
		goto(`/account/learners/${learner.id}/edit`);
	}
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white">
			{learner?.nickname}'s Learning Journey
		</h1>
		<FormButton type="button" variant="secondary" on:click={editLearner}>
			Edit
		</FormButton>
	</div>

	{#if loading}
		<div class="flex justify-center items-center h-64">
			<LoadingSpinner size="md" color="primary" />
		</div>
	{:else if error}
		<ErrorAlert message={error} />
	{:else if learner}
		<div class="space-y-6">
			<!-- Basic Profile Information -->
			<div class="bg-white dark:bg-gray-800 shadow-sm rounded-lg p-6">
				<div class="flex items-center space-x-4">
					{#if learner.avatar}
						<img
							src={learner.avatar}
							alt={learner.nickname || "Learner"}
							class="w-20 h-20 rounded-full"
						/>
					{:else}
						<div
							class="w-20 h-20 rounded-full bg-gradient-to-br from-blue-500 to-purple-500 flex items-center justify-center"
						>
							<span class="text-3xl font-bold text-white"
								>{learner.nickname?.[0]?.toUpperCase() || "L"}</span
							>
						</div>
					{/if}
					<div class="flex-1">
						<h2
							class="text-xl font-semibold text-gray-900 dark:text-white mb-2"
						>
							{learner.nickname || "Unknown learner"}
						</h2>
						<div class="flex flex-wrap gap-2">
							{#if learner.age}
								<span class="bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300 text-sm font-medium px-2.5 py-0.5 rounded">
									Age: {learner.age}
								</span>
							{/if}
							
							{#if learner.grade_level}
								<span class="bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-300 text-sm font-medium px-2.5 py-0.5 rounded">
									Grade: {learner.grade_level}
								</span>
							{/if}
						</div>

						{#if learner.learning_preferences && learner.learning_preferences.length > 0}
							<div class="mt-3">
								<p class="text-sm text-gray-600 dark:text-gray-400 mb-1">Learning Style:</p>
								<div class="flex flex-wrap gap-2">
									{#each learner.learning_preferences as preference}
										<span
											class="bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300 text-xs font-medium px-2.5 py-0.5 rounded"
										>
											{preference}
										</span>
									{/each}
								</div>
							</div>
						{/if}
					</div>
				</div>
			</div>

			<!-- Learning Progress -->
			{#if learner.id}
				<LearnerProgress learnerId={learner.id} />
			{/if}
		</div>
	{/if}
</div>
