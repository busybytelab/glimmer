<script lang="ts">
	import { error } from '$lib/stores';
	import ErrorAlert from '../common/ErrorAlert.svelte';
	import StatCard from './StatCard.svelte';

	export let activeLearners: number;
	export let practiceTopics: number;
	export let completedSessions: number;
	export let isLoading: boolean = false;

	// SVG path data for icons
	const icons = {
		learners: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />',
		topics: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />',
		sessions: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />'
	};
</script>

<div class="py-6">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
		<h1 class="text-2xl font-semibold text-gray-900">Dashboard</h1>
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