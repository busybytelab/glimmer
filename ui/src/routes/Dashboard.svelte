<script lang="ts">
	import { user, error } from '$lib/stores';
	import { onMount } from 'svelte';
	import type { Instructor, Learner } from '$lib/types';
	import AppLayout from '../components/layout/AppLayout.svelte';
	import DashboardContent from '../components/dashboard/DashboardContent.svelte';

	let sidebarOpen = true;
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
				//const account = await pb.collection('accounts').getOne(instructor.account);
				
				// Fetch stats for instructor's account
				// const [learners, topics, sessions] = await Promise.all([
				// 	pb.collection('learners').getList(1, 1, { filter: `account="${account.id}"` }),
				// 	pb.collection('practice_topics').getList(1, 1, { filter: `account="${account.id}"` }),
				// 	pb.collection('practice_sessions').getList(1, 1, { filter: `account="${account.id}"` })
				// ]);

				stats = {
					activeLearners: 0, //learners.totalItems,
					practiceTopics: 0, //topics.totalItems,
					completedSessions: 0 //sessions.totalItems
				};
			} else {
				// Learner dashboard
				const learner = currentUser as Learner;
				console.log(learner);

				// TODO: fix collection relationships and rules
				//const account = await pb.collection('accounts').getOne(learner.account);
				
				// Fetch stats for learner
				// const [topics, sessions] = await Promise.all([
				// 	pb.collection('practice_topics').getList(1, 1, { filter: `account="${account.id}"` }),
				// 	pb.collection('practice_sessions').getList(1, 1, { filter: `learner="${learner.id}"` })
				// ]);

				stats = {
					activeLearners: 1, // The learner themselves
					practiceTopics: 0, //topics.totalItems,
					completedSessions: 0 //sessions.totalItems
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

<AppLayout bind:sidebarOpen>
	<DashboardContent
		{activeLearners}
		{practiceTopics}
		{completedSessions}
		{isLoading}
	/>
</AppLayout>

<style>
	:global(html) {
		--primary: #2c3e50;
		--secondary: #3498db;
	}

	.text-primary {
		color: var(--primary);
	}

	.focus\:ring-secondary:focus {
		--tw-ring-color: var(--secondary);
	}
</style> 