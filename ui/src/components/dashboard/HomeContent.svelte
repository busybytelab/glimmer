<script lang="ts">
	import { error } from '$lib/stores';
	import ErrorAlert from '../common/ErrorAlert.svelte';
	import StatCard from './StatCard.svelte';
	import type { IconType } from '$lib/types';

	export let activeLearners: number;
	export let practiceTopics: number;
	export let completedSessions: number;
	export let isLoading: boolean = false;

	// Use proper IconType values
	const icons: Record<string, IconType> = {
		learners: 'learner',
		topics: 'topic',
		sessions: 'complete'
	};
</script>

<div class="py-6">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
		<h1 class="text-2xl font-semibold text-gray-900 dark:text-white">Home</h1>
	</div>
	{#if $error}
		<div class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8 mt-4">
			<ErrorAlert message={$error} />
		</div>
	{/if}
	<div class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
		<div class="py-4">
			<div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
				<!-- Quick Stats -->
				<StatCard 
					title="Active Learners" 
					value={activeLearners} 
					isLoading={isLoading} 
					icon={icons.learners} 
				/>
				
				<StatCard 
					title="Practice Topics" 
					value={practiceTopics} 
					isLoading={isLoading} 
					icon={icons.topics} 
				/>
				
				<StatCard 
					title="Completed Sessions" 
					value={completedSessions} 
					isLoading={isLoading} 
					icon={icons.sessions} 
				/>
			</div>
		</div>
	</div>
</div> 