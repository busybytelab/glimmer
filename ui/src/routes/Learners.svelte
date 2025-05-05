<script lang="ts">
	import { onMount } from 'svelte';
	import pb from '$lib/pocketbase';
	import type { Learner } from '$lib/types';
	import AppLayout from '../components/layout/AppLayout.svelte';
	import FormButton from '../components/common/FormButton.svelte';
	import LearnersList from '../components/learners/LearnersList.svelte';

	let learners: Learner[] = [];
	let loading = true;
	let error: string | null = null;
	let sidebarOpen = true;
	
	onMount(() => {
		loadLearners();
	});
	
	async function loadLearners() {
		try {
			loading = true;
			error = null;
			
			// Get current user info
			const authData = pb.authStore.model;
			if (!authData) {
				console.error('User not authenticated');
				error = 'You must be logged in to view learners';
				return;
			}
			
			// Get instructor record
			const instructorRecord = await pb.collection('instructors').getFirstListItem(`user="${authData.id}"`);
			
			// Get learners from the same account
			const result = await pb.collection('learners').getList(1, 100, {
				filter: `account="${instructorRecord.account}"`,
				sort: 'nickname',
				expand: 'user'
			});
			
			learners = result.items as unknown as Learner[];
			console.log('Loaded learners:', learners);
		} catch (err) {
			console.error('Failed to load learners:', err);
			error = 'Failed to load learners';
		} finally {
			loading = false;
		}
	}
	
	function viewLearner(learner: Learner) {
		console.log('View learner profile:', learner);
		// Use path-based routing
		(window as any).navigate(`/learner/${learner.id}`);
	}
	
	function editLearner(learner: Learner) {
		console.log('Edit learner:', learner);
		// Use path-based routing
		(window as any).navigate(`/edit-learner/${learner.id}`);
	}
	
	function deleteLearner(learner: Learner) {
		if (confirm(`Are you sure you want to delete ${learner.nickname}?`)) {
			console.log('Delete learner:', learner);
			// TODO: Implement learner deletion
		}
	}
	
	function createNewLearner() {
		console.log('Create new learner');
		// Use path-based routing
		(window as any).navigate('/create-learner');
	}
</script>

<AppLayout bind:sidebarOpen>
	<div class="container mx-auto px-4 py-8">
		<div class="flex justify-between items-center mb-6">
			<h1 class="text-2xl font-bold text-gray-900">Learners</h1>
			<FormButton
				type="button"
				on:click={createNewLearner}
			>
				Add New Learner
			</FormButton>
		</div>
		
		{#if error}
			<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-6" role="alert">
				<strong class="font-bold">Error!</strong>
				<span class="block sm:inline"> {error}</span>
			</div>
		{/if}
		
		<LearnersList
			{learners}
			{loading}
			emptyMessage="No learners found. Add your first learner to get started!"
			onClick={viewLearner}
			cardActions={[
				{
					label: 'View',
					color: 'primary',
					onClick: viewLearner
				},
				{
					label: 'Edit',
					color: 'secondary', 
					onClick: editLearner
				},
				{
					label: 'Delete',
					color: 'danger',
					onClick: deleteLearner
				}
			]}
		/>
	</div>
</AppLayout> 