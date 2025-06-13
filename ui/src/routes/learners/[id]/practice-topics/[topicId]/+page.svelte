<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import type { PracticeTopic, BreadcrumbItem, BreadcrumbIcon } from '$lib/types';
	import { topicsService } from '$lib/services/topics';
	import { sessionService } from '$lib/services/session';
	import ActionToolbar from '$components/common/ActionToolbar.svelte';
	import Breadcrumbs from '$components/common/Breadcrumbs.svelte';
	import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
	import ErrorAlert from '$components/common/ErrorAlert.svelte';

	let topic: PracticeTopic | null = null;
	let pastPractices: any[] = [];
	let loading = true;
	let error: string | null = null;
	let topicId: string | null = null;
	let learnerId: string | null = null;
	let breadcrumbItems: BreadcrumbItem[] = [];

	onMount(async () => {
		try {
			// Get the topic ID and learner ID from the page store
			topicId = $page.params.topicId;
			learnerId = $page.params.id;
			
			// Load data if we have an ID
			if (topicId && learnerId) {
				await loadTopic(topicId);
				await loadPastPractices(topicId, learnerId);
				updateBreadcrumbs();
			} else {
				error = 'Invalid topic ID or learner ID';
				loading = false;
			}
		} catch (err) {
			console.error('Error in onMount:', err);
			error = err instanceof Error ? err.message : 'An unexpected error occurred';
			loading = false;
		}
	});

	async function loadTopic(id: string) {
		try {
			loading = true;
			error = null;
			
			if (!id) {
				throw new Error('Topic ID is required');
			}

			const result = await topicsService.getTopic(id);
			if (!result) {
				throw new Error('Topic not found');
			}

			// Parse tags if they're stored as a string
			if (result.tags) {
				const tagsValue = result.tags as string | string[];
				if (typeof tagsValue === 'string') {
					try {
						if (tagsValue.trim().startsWith('[')) {
							result.tags = JSON.parse(tagsValue);
						} else {
							result.tags = tagsValue.split(',').map((tag: string) => tag.trim()).filter(Boolean);
						}
					} catch (err) {
						console.error('Error parsing tags:', err);
						result.tags = [];
					}
				}
			} else {
				result.tags = [];
			}

			topic = result;
		} catch (err) {
			console.error('Failed to load topic:', err);
			error = 'Failed to load practice topic';
		} finally {
			loading = false;
		}
	}

	async function loadPastPractices(topicId: string, learnerId: string) {
		try {
			const result = await sessionService.getSessionsForTopicAndLearner(topicId, learnerId);
			pastPractices = result;
		} catch (err) {
			console.error('Failed to load past practices:', err);
			error = 'Failed to load past practice sessions';
		}
	}

	function goBack() {
		window.history.back();
	}

	function updateBreadcrumbs() {
		if (!topic) return;
		
		breadcrumbItems = [
			{
				label: 'Practice Topics',
				href: `/learners/${learnerId}/practice-topics`,
				icon: 'topic' as BreadcrumbIcon
			},
			{
				label: topic.name,
				icon: 'topic' as BreadcrumbIcon
			}
		];
	}

	// Actions for the toolbar
	$: topicActions = [
		{
			id: 'back',
			label: 'Back',
			icon: 'back',
			variant: 'secondary' as const,
			onClick: goBack
		}
	];
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
	<div class="flex justify-between items-center mb-6">
		<div>
			<Breadcrumbs items={breadcrumbItems} />
		</div>
		<ActionToolbar actions={topicActions} />
	</div>

	{#if loading}
		<div class="flex justify-center items-center h-64">
			<LoadingSpinner size="md" color="primary" />
		</div>
	{:else if error}
		<ErrorAlert message={error} />
	{:else if topic}
		<div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6 mb-6">
			<h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">{topic.name}</h2>
			
			{#if topic.description}
				<p class="text-gray-600 dark:text-gray-300 mb-4">{topic.description}</p>
			{/if}
			
			<div class="flex flex-wrap gap-2 mb-4">
				<span class="bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-300 text-xs font-medium px-2.5 py-0.5 rounded">
					{topic.subject}
				</span>
				
				{#if topic.target_age_range}
					<span class="bg-green-100 dark:bg-green-900/40 text-green-800 dark:text-green-300 text-xs font-medium px-2.5 py-0.5 rounded">
						Age: {topic.target_age_range}
					</span>
				{/if}
				
				{#if topic.target_grade_level}
					<span class="bg-purple-100 dark:bg-purple-900/40 text-purple-800 dark:text-purple-300 text-xs font-medium px-2.5 py-0.5 rounded">
						Grade: {topic.target_grade_level}
					</span>
				{/if}
			</div>
		</div>

		{#if pastPractices.length > 0}
			<div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
				<h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Past Practice Sessions</h3>
				
				<div class="overflow-x-auto">
					<table class="min-w-full bg-white dark:bg-gray-800">
						<thead>
							<tr class="border-b border-gray-200 dark:border-gray-700">
								<th class="px-4 py-2 text-left text-sm font-semibold text-gray-700 dark:text-gray-300">Session Name</th>
								<th class="px-4 py-2 text-left text-sm font-semibold text-gray-700 dark:text-gray-300">Date</th>
								<th class="px-4 py-2 text-left text-sm font-semibold text-gray-700 dark:text-gray-300">Status</th>
								<th class="px-4 py-2 text-left text-sm font-semibold text-gray-700 dark:text-gray-300">Actions</th>
							</tr>
						</thead>
						<tbody>
							{#each pastPractices as session}
								<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700">
									<td class="px-4 py-2 text-sm text-gray-700 dark:text-gray-300">{session.name || 'Practice Session'}</td>
									<td class="px-4 py-2 text-sm text-gray-700 dark:text-gray-300">
										{new Date(session.created).toLocaleDateString()}
									</td>
									<td class="px-4 py-2 text-sm">
										<span class={`px-2 py-1 rounded-full text-xs font-medium 
											${session.status === 'Completed' ? 'bg-green-100 dark:bg-green-900/40 text-green-800 dark:text-green-300' : 
											  session.status === 'InProgress' ? 'bg-yellow-100 dark:bg-yellow-900/40 text-yellow-800 dark:text-yellow-300' : 
											  'bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-300'}`}>
											{session.status}
										</span>
									</td>
									<td class="px-4 py-2 text-sm text-gray-700 dark:text-gray-300">
										<a 
											href={`/learners/${learnerId}/practice-sessions/${session.id}`}
											class="text-primary hover:text-primary-dark"
										>
											View Session
										</a>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		{:else}
			<div class="bg-yellow-50 dark:bg-yellow-900/20 border-l-4 border-yellow-400 dark:border-yellow-600 p-4">
				<div class="flex">
					<div class="flex-shrink-0">
						<svg class="h-5 w-5 text-yellow-400 dark:text-yellow-300" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
						</svg>
					</div>
					<div class="ml-3">
						<p class="text-sm text-yellow-700 dark:text-yellow-200">
							No practice sessions available for this topic.
						</p>
					</div>
				</div>
			</div>
		{/if}
	{/if}
</div> 