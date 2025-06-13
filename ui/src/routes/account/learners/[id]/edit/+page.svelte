<script lang="ts">
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { goto } from "$app/navigation";
	import type { Learner } from "$lib/types";
	import LearnerForm from "$components/learners/LearnerForm.svelte";
	import FormButton from "$components/common/FormButton.svelte";
	import LoadingSpinner from "$components/common/LoadingSpinner.svelte";
	import ErrorAlert from "$components/common/ErrorAlert.svelte";
	import { learnersService } from "$lib/services/learners";
	let learner: Learner | null = null;
	let loading = true;
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

	function handleLearnerUpdate() {
		goto(`/account/learners/${learner?.id}`);
	}

	function handleLearnerDelete() {
		goto(`/account/learners`);
	}

	function handleCancel() {
		goto(`/account/learners/${learner?.id}`);
	}
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6 max-w-7xl">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white">
			{learner ? "Edit Profile" : "Loading..."}
		</h1>
		<FormButton type="button" variant="secondary" on:click={handleCancel}>
			Back to {learner?.nickname}'s profile
		</FormButton>
	</div>

	{#if loading}
		<div class="flex justify-center items-center h-64">
			<LoadingSpinner size="md" color="primary" />
		</div>
	{:else if error}
		<ErrorAlert message={error} />
	{:else if learner}
		<div class="w-full">
			<LearnerForm
				{learner}
				on:update={() => handleLearnerUpdate()}
				on:delete={() => handleLearnerDelete()}
				on:cancel={handleCancel}
			/>
		</div>
	{/if}
</div>
