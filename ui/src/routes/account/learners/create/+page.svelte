<script lang="ts">
	import { goto } from '$app/navigation';
	import LearnerForm from '$components/learners/LearnerForm.svelte';
	import FormButton from '$components/common/FormButton.svelte';
	import ErrorAlert from '$components/common/ErrorAlert.svelte';

	let error: string | null = null;

	function handleLearnerCreate() {
		goto('/home');
	}

	function handleCancel() {
		goto('/home');
	}
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6 max-w-7xl">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Add Child</h1>
		<FormButton
			type="button"
			variant="secondary"
			on:click={handleCancel}
		>
			Back to Home
		</FormButton>
	</div>

	{#if error}
		<ErrorAlert message={error} />
	{:else}
		<div class="w-full">
			<LearnerForm
				learner={null}
				on:update={() => handleLearnerCreate()}
				on:cancel={handleCancel}
			/>
		</div>
	{/if}
</div> 