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
			Edit Profile
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
			<!-- Learning Progress -->
			{#if learner.id}
				<LearnerProgress learnerId={learner.id} learnerName={learner.nickname} />
			{/if}
		</div>
	{/if}
</div>
