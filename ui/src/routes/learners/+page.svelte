<script lang="ts">
	import { onMount } from 'svelte';
	import pb from '$lib/pocketbase';
	import { goto } from '$app/navigation';
	import type { Learner } from '$lib/types';
	import LearnersList from '../../components/learners/LearnersList.svelte';
	import FormButton from '../../components/common/FormButton.svelte';
	import LoadingSpinner from '../../components/common/LoadingSpinner.svelte';
	import ErrorAlert from '../../components/common/ErrorAlert.svelte';

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

			// Get instructor record to get the account
			const instructorRecord = await pb.collection('instructors').getFirstListItem(`user="${authData.id}"`);
			
			// Get learners from the same account
			const result = await pb.collection('learners').getFullList({
				filter: `account="${instructorRecord.account}"`,
				sort: '-created',
				expand: 'user,account'
			});
			
			// Map expanded data to the learner objects
			learners = result.map(item => {
				console.log('learner Item:', item);
				const expandedData = item.expand;
				return {
					...item,
					user: expandedData?.user || { name: 'Unknown user' },
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
		goto('/learners/create');
	}

	function viewLearner(learner: Learner) {
		goto(`/learners/${learner.id}`);
	}
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Learners</h1>
		<FormButton
			type="button"
			on:click={handleCreateNew}
		>
			Create New Learner
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
	cardActions={[]}
/>
	{/if}
</div> 