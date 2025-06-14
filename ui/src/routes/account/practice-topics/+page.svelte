<script lang="ts">
	import { onMount } from 'svelte';
	import { topicsService } from '$lib/services/topics';
	import { goto } from '$app/navigation';
	import type { PracticeTopic, IconType } from '$lib/types';
	import PracticeTopicCard from '$components/practice-topics/PracticeTopicCard.svelte';
	import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
	import ErrorAlert from '$components/common/ErrorAlert.svelte';
	import ActionToolbar from '$components/common/ActionToolbar.svelte';

	let topics: PracticeTopic[] = [];
	let loading = true;
	let error: string | null = null;

	async function loadTopics() {
		try {
			loading = true;
			error = null;
			const result = await topicsService.getTopics();
			
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
		await loadTopics();
	});

	function handleCreateNew() {
		goto('/account/practice-topics/create');
	}

	function handleImport() {
		goto('/account/practice-topics/import');
	}

	// Actions for the toolbar
	const topicActions = [
		{
			id: 'import',
			label: 'Import Session',
			icon: 'download' as IconType,
			variant: 'secondary' as const,
			onClick: handleImport
		},
		{
			id: 'create',
			label: 'Create Topic',
			icon: 'add' as IconType,
			variant: 'primary' as const,
			onClick: handleCreateNew
		}
	];
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Practice Topics</h1>
		<ActionToolbar actions={topicActions} />
	</div>

	{#if loading}
		<div class="flex justify-center items-center h-64">
			<LoadingSpinner size="md" color="primary" />
		</div>
	{:else if error}
		<ErrorAlert message={error} />
	{:else if topics.length === 0}
		<div class="bg-yellow-50 dark:bg-yellow-900/20 border-l-4 border-yellow-400 dark:border-yellow-600 p-4">
			<div class="flex">
				<div class="flex-shrink-0">
					<svg class="h-5 w-5 text-yellow-400 dark:text-yellow-300" viewBox="0 0 20 20" fill="currentColor">
						<path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
					</svg>
				</div>
				<div class="ml-3">
					<p class="text-sm text-yellow-700 dark:text-yellow-200">
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
					href={`/account/practice-topics/${topic.id}`}
				/>
			{/each}
		</div>
	{/if}
</div> 