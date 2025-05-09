<script lang="ts">
	import { onMount } from 'svelte';
	import pb from '$lib/pocketbase';
	import type { PracticeTopic } from '$lib/types';
	import PracticeTopicForm from '../components/practice-topics/PracticeTopicForm.svelte';
	import PracticeTopicCard from '../components/practice-topics/PracticeTopicCard.svelte';
	import AppLayout from '../components/layout/AppLayout.svelte';
	import FormButton from '../components/common/FormButton.svelte';

	let topics: PracticeTopic[] = [];
	let loading = true;
	let error: string | null = null;
	let selectedTopic: PracticeTopic | null = null;
	let showCreateForm = false;
	let sidebarOpen = true;

	async function loadTopics() {
		try {
			loading = true;
			error = null;
			const result = await pb.collection('practice_topics').getFullList({
				sort: '-created',
				expand: 'account'
			});
			
			// Make sure tags are properly formatted as arrays
			topics = result.map((topic: any) => {
				if (topic.tags && !Array.isArray(topic.tags)) {
					// Handle case where PocketBase may return stringified JSON
					try {
						if (typeof topic.tags === 'string' && topic.tags.trim().startsWith('[')) {
							topic.tags = JSON.parse(topic.tags);
						} else if (typeof topic.tags === 'string') {
							topic.tags = topic.tags.split(',').map((tag: string) => tag.trim()).filter(Boolean);
						}
					} catch (e) {
						console.error('Error parsing tags:', e);
						topic.tags = [];
					}
				} else if (!topic.tags) {
					topic.tags = [];
				}
				
				return topic;
			}) as unknown as PracticeTopic[];
			
			console.log('Loaded topics:', topics);
		} catch (err) {
			console.error('Failed to load topics:', err);
			error = 'Failed to load practice topics';
		} finally {
			loading = false;
		}
	}

	onMount(async () => {
		console.log('PracticeTopics mounted, loading topics');
		
		// Check if edit parameter is in URL
		const urlParams = new URLSearchParams(window.location.search);
		const editTopicId = urlParams.get('edit');
		console.log('URL edit parameter:', editTopicId);
		
		// Load topics first
		await loadTopics();
		
		if (editTopicId) {
			console.log('Found edit parameter, will edit topic:', editTopicId);
			const topicToEdit = topics.find(t => t.id === editTopicId);
			if (topicToEdit) {
				console.log('Found topic to edit:', topicToEdit.name);
				handleTopicEdit(topicToEdit);
			} else {
				console.error('Could not find topic with id:', editTopicId);
			}
			
			// Clear the URL parameter to avoid reopening the edit form on refresh
			const url = new URL(window.location.href);
			url.searchParams.delete('edit');
			window.history.replaceState({}, '', url);
		}
	});

	function handleTopicSelect(topic: PracticeTopic) {
		// Navigate to view practice topic page using path-based routing
		(window as any).navigate(`/practice-topic/${topic.id}`);
	}

	function handleTopicEdit(topic: PracticeTopic) {
		selectedTopic = topic;
		showCreateForm = false;
	}

	function handleTopicUpdate(updatedTopic: PracticeTopic) {
		// Update the topic in the list
		topics = topics.map(t => t.id === updatedTopic.id ? updatedTopic : t);
		selectedTopic = null;
		showCreateForm = false;
		
		// Refresh the list to ensure it's up to date
		loadTopics();
	}

	function handleTopicDelete(topicId: string) {
		// Remove the topic from the list
		topics = topics.filter(t => t.id !== topicId);
		selectedTopic = null;
		showCreateForm = false;
	}

	function handleCreateNew() {
		selectedTopic = null;
		showCreateForm = true;
	}

	function handleCancel() {
		selectedTopic = null;
		showCreateForm = false;
	}
</script>

<AppLayout bind:sidebarOpen>
	<div class="container mx-auto px-4 py-8">
		{#if !selectedTopic && !showCreateForm}
			<div class="flex justify-between items-center mb-6">
				<h1 class="text-2xl font-bold text-gray-900">Practice Topics</h1>
				<FormButton
					type="button"
					on:click={handleCreateNew}
				>
					Create New Topic
				</FormButton>
			</div>
		{:else}
			<div class="flex justify-between items-center mb-6">
				<h1 class="text-2xl font-bold text-gray-900">
					{selectedTopic ? 'Edit Practice Topic' : 'Create New Topic'}
				</h1>
				<FormButton
					type="button"
					variant="secondary"
					on:click={handleCancel}
				>
					Back to Topics
				</FormButton>
			</div>
		{/if}

		{#if loading}
			<div class="flex justify-center items-center h-64">
				<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary"></div>
			</div>
		{:else if error}
			<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
				<strong class="font-bold">Error!</strong>
				<span class="block sm:inline"> {error}</span>
			</div>
		{:else if selectedTopic !== null}
			<div class="max-w-4xl mx-auto">
				<PracticeTopicForm
					topic={selectedTopic}
					on:update={({ detail }) => handleTopicUpdate(detail)}
					on:delete={({ detail }) => handleTopicDelete(detail)}
					on:cancel={handleCancel}
				/>
			</div>
		{:else if showCreateForm}
			<div class="max-w-4xl mx-auto">
				<PracticeTopicForm
					on:update={({ detail }) => handleTopicUpdate(detail)}
					on:cancel={handleCancel}
				/>
			</div>
		{:else if topics.length === 0}
			<div class="bg-yellow-50 border-l-4 border-yellow-400 p-4">
				<div class="flex">
					<div class="flex-shrink-0">
						<svg class="h-5 w-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
						</svg>
					</div>
					<div class="ml-3">
						<p class="text-sm text-yellow-700">
							No practice topics available. Create your first topic!
						</p>
					</div>
				</div>
			</div>
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each topics as topic}
					<PracticeTopicCard 
						{topic} 
						onClick={handleTopicSelect} 
						onEdit={handleTopicEdit} 
					/>
				{/each}
			</div>
		{/if}
	</div>
</AppLayout> 