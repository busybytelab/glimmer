<script lang="ts">
	import { user, error } from '$lib/stores';
	import { onMount } from 'svelte';
	import type { Instructor, Learner } from '$lib/types';
	import DashboardContent from '../../components/dashboard/DashboardContent.svelte';

	let activeLearners = 0;
	let practiceTopics = 0;
	let completedSessions = 0;
	let isLoading = true;

	interface DashboardStats {
		activeLearners: number;
		practiceTopics: number;
		completedSessions: number;
	}

	onMount(async () => {
		await fetchStats();
	});

	async function fetchStats() {
		isLoading = true;
		error.set(null);

		try {
			const currentUser = $user;
			if (!currentUser) {
				throw new Error('User not found');
			}

			let stats: DashboardStats = {
				activeLearners: 0,
				practiceTopics: 0,
				completedSessions: 0
			};

			if ('account' in currentUser) {
				// Instructor dashboard
				const instructor = currentUser as Instructor;
				console.log(instructor);
				// TODO: fix collection relationships and rules

				stats = {
					activeLearners: 0,
					practiceTopics: 0,
					completedSessions: 0
				};
			} else {
				// Learner dashboard
				const learner = currentUser as Learner;
				console.log(learner);

				stats = {
					activeLearners: 1, // The learner themselves
					practiceTopics: 0,
					completedSessions: 0
				};
			}

			activeLearners = stats.activeLearners;
			practiceTopics = stats.practiceTopics;
			completedSessions = stats.completedSessions;
		} catch (err) {
			console.error('Error fetching dashboard stats:', err);
			error.set(err instanceof Error ? err.message : 'Failed to fetch dashboard stats');
		} finally {
			isLoading = false;
		}
	}
</script>

<DashboardContent
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