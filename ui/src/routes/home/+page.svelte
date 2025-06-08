<script lang="ts">
	import { error } from '$lib/stores';
	import { onMount } from 'svelte';
	import HomeContent from '../../components/dashboard/HomeContent.svelte';

	let activeLearners = 0;
	let practiceTopics = 0;
	let completedSessions = 0;
	let isLoading = true;

	interface HomeStats {
		activeLearners: number;
		practiceTopics: number;
		completedSessions: number;
	}

	onMount(async () => {
		await fetchStats();
	});


	// FIXME: this needs a complete rewrite
	async function fetchStats() {
		isLoading = true;
		error.set(null);

		try {
			let stats: HomeStats = {
				activeLearners: 0,
				practiceTopics: 0,
				completedSessions: 0
			};

			activeLearners = stats.activeLearners;
			practiceTopics = stats.practiceTopics;
			completedSessions = stats.completedSessions;
		} catch (err) {
			console.error('Error fetching home stats:', err);
			error.set(err instanceof Error ? err.message : 'Failed to fetch home stats');
		} finally {
			isLoading = false;
		}
	}
</script>

<HomeContent
	{activeLearners}
	{practiceTopics}
	{completedSessions}
	{isLoading}
/>

<style lang="postcss">
	:global(html) {
		--primary: #2c3e50;
		--secondary: #3498db;
	}
</style> 