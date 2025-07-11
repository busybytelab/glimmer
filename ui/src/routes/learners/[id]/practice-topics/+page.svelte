<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import type { PracticeTopic, Learner, BreadcrumbItem, IconType } from '$lib/types';
	import PracticeTopicCard from '$components/practice-topics/PracticeTopicCard.svelte';
	import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
	import ErrorAlert from '$components/common/ErrorAlert.svelte';
	import Breadcrumbs from '$components/common/Breadcrumbs.svelte';
	import { topicsService } from '$lib/services/topics';
	import { learnersService } from '$lib/services/learners';

	let topics: PracticeTopic[] = [];
	let loading = true;
	let error: string | null = null;
	let learnerId: string = '';
	let breadcrumbItems: BreadcrumbItem[] = [];
	let learner: any = null;

	onMount(async () => {
		try {
			// Get the learner ID from the page store
			learnerId = $page.params.id;
			
			if (learnerId) {
				// First load the learner
				learner = await learnersService.getLearner(learnerId);
				if (!learner) {
					throw new Error('Learner not found');
				}
				topics = await loadTopics(learner);
				updateBreadcrumbs();
			} else {
				error = 'Invalid learner ID';
				loading = false;
			}
		} catch (err) {
			console.error('Error in onMount:', err);
			error = err instanceof Error ? err.message : 'An unexpected error occurred';
			loading = false;
		}
	});

	async function loadTopics(learner: Learner): Promise<PracticeTopic[]> {
		try {
			loading = true;
			error = null;
			if (!learner) {
				throw new Error('Learner information not available');
			}
			const result = await topicsService.getTopicsForLearner(learner.age, learner.grade_level);
			console.log('Loaded topics for learner:', result);
			return result;
	
		} catch (err) {
			console.error('Failed to load topics:', err);
			error = 'Failed to load practice topics';
			return [];
		} finally {
			loading = false;
		}
	}

	function updateBreadcrumbs() {
		breadcrumbItems = [
			{
				label: `${learner.nickname}'s Profile`,
				href: '/',
				icon: 'home' as IconType
			},
			{
				label: 'Practice Topics',
				icon: 'topic' as IconType
			}
		];
	}
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
	<div class="flex justify-between items-center mb-6">
		<div>
			<Breadcrumbs items={breadcrumbItems} />
		</div>
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
						No practice topics available.
					</p>
				</div>
			</div>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each topics as topic}
				<PracticeTopicCard 
					{topic} 
					href={`/learners/${learnerId}/practice-topics/${topic.id}`}
					showEditButton={false}
				/>
			{/each}
		</div>
	{/if}
</div> 