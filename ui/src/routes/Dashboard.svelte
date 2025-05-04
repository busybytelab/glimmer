<script lang="ts">
	import { user, error } from '$lib/stores';
	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';
	import type { Instructor, Learner } from '$lib/types';

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

	function handleLogout() {
		pb.authStore.clear();
		localStorage.removeItem('authToken');
		window.location.href = '/';
	}
</script>

<div class="h-screen flex overflow-hidden bg-gray-100">
	<!-- Mobile sidebar -->
	{#if sidebarOpen}
		<div class="md:hidden fixed inset-0 flex z-40">
			<button
				class="fixed inset-0 bg-gray-600 bg-opacity-75"
				on:click={() => sidebarOpen = false}
				on:keydown={(e) => e.key === 'Escape' && (sidebarOpen = false)}
				aria-label="Close sidebar overlay"
			></button>
			<div class="relative flex-1 flex flex-col max-w-xs w-full bg-white">
				<div class="absolute top-0 right-0 -mr-12 pt-2">
					<button
						on:click={() => sidebarOpen = false}
						class="ml-1 flex items-center justify-center h-10 w-10 rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
					>
						<span class="sr-only">Close sidebar</span>
						<svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
						</svg>
					</button>
				</div>
				<div class="flex-1 h-0 pt-5 pb-4 overflow-y-auto">
					<div class="flex-shrink-0 flex items-center px-4">
						<div class="flex items-center">
							<img src="/glimmer.svg" alt="Glimmer Logo" class="h-8 w-8 mr-2" />
							<h1 class="text-2xl font-bold text-primary">Glimmer</h1>
						</div>
					</div>
					<nav class="mt-5 px-2 space-y-1">
						<a href="/dashboard" class="group flex items-center px-2 py-2 text-base font-medium rounded-md text-gray-900 bg-gray-100">
							<svg class="mr-4 h-6 w-6 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
							</svg>
							Dashboard
						</a>
						<a href="/topics" class="group flex items-center px-2 py-2 text-base font-medium rounded-md text-gray-600 hover:bg-gray-50 hover:text-gray-900">
							<svg class="mr-4 h-6 w-6 text-gray-400 group-hover:text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
							</svg>
							Practice Topics
						</a>
						<a href="/learners" class="group flex items-center px-2 py-2 text-base font-medium rounded-md text-gray-600 hover:bg-gray-50 hover:text-gray-900">
							<svg class="mr-4 h-6 w-6 text-gray-400 group-hover:text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
							</svg>
							Learners
						</a>
						<a href="/settings" class="group flex items-center px-2 py-2 text-base font-medium rounded-md text-gray-600 hover:bg-gray-50 hover:text-gray-900">
							<svg class="mr-4 h-6 w-6 text-gray-400 group-hover:text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
							</svg>
							Settings
						</a>
					</nav>
				</div>
				<div class="flex-shrink-0 flex border-t border-gray-200 p-4">
					<div class="flex-shrink-0 w-full group block">
						<div class="flex items-center">
							<div class="ml-3">
								<p class="text-base font-medium text-gray-700 group-hover:text-gray-900">
									{$user?.user?.name || 'User'}
								</p>
								<button
									on:click={handleLogout}
									class="text-sm font-medium text-gray-500 group-hover:text-gray-700"
								>
									Sign out
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Desktop sidebar -->
	<div class="hidden md:flex md:flex-shrink-0">
		<div class="flex flex-col w-64">
			<div class="flex flex-col h-0 flex-1 border-r border-gray-200 bg-white">
				<div class="flex-1 flex flex-col pt-5 pb-4 overflow-y-auto">
					<div class="flex items-center flex-shrink-0 px-4">
						<div class="flex items-center">
							<img src="/glimmer.svg" alt="Glimmer Logo" class="h-8 w-8 mr-2" />
							<h1 class="text-2xl font-bold text-primary">Glimmer</h1>
						</div>
					</div>
					<nav class="mt-5 flex-1 px-2 bg-white space-y-1">
						<a href="/dashboard" class="group flex items-center px-2 py-2 text-sm font-medium rounded-md text-gray-900 bg-gray-100">
							<svg class="mr-3 h-6 w-6 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
							</svg>
							Dashboard
						</a>
						<a href="/topics" class="group flex items-center px-2 py-2 text-sm font-medium rounded-md text-gray-600 hover:bg-gray-50 hover:text-gray-900">
							<svg class="mr-3 h-6 w-6 text-gray-400 group-hover:text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
							</svg>
							Practice Topics
						</a>
						<a href="/learners" class="group flex items-center px-2 py-2 text-sm font-medium rounded-md text-gray-600 hover:bg-gray-50 hover:text-gray-900">
							<svg class="mr-3 h-6 w-6 text-gray-400 group-hover:text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
							</svg>
							Learners
						</a>
						<a href="/settings" class="group flex items-center px-2 py-2 text-sm font-medium rounded-md text-gray-600 hover:bg-gray-50 hover:text-gray-900">
							<svg class="mr-3 h-6 w-6 text-gray-400 group-hover:text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
							</svg>
							Settings
						</a>
					</nav>
				</div>
				<div class="flex-shrink-0 flex border-t border-gray-200 p-4">
					<div class="flex-shrink-0 w-full group block">
						<div class="flex items-center">
							<div class="ml-3">
								<p class="text-sm font-medium text-gray-700 group-hover:text-gray-900">
									{$user?.user?.name || 'User'}
								</p>
								<button
									on:click={handleLogout}
									class="text-xs font-medium text-gray-500 group-hover:text-gray-700"
								>
									Sign out
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Main content -->
	<div class="flex flex-col w-0 flex-1 overflow-hidden h-screen">
		<div class="md:hidden pl-1 pt-1 sm:pl-3 sm:pt-3">
			<button
				on:click={() => sidebarOpen = !sidebarOpen}
				class="-ml-0.5 -mt-0.5 h-12 w-12 inline-flex items-center justify-center rounded-md text-gray-500 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-secondary"
			>
				<span class="sr-only">Open sidebar</span>
				<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
				</svg>
			</button>
		</div>
		<main class="flex-1 relative z-0 overflow-y-auto focus:outline-none h-full">
			<div class="py-6">
				<div class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
					<h1 class="text-2xl font-semibold text-gray-900">Dashboard</h1>
				</div>
				{#if $error}
					<div class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8 mt-4">
						<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
							<strong class="font-bold">Error!</strong>
							<span class="block sm:inline">{$error}</span>
						</div>
					</div>
				{/if}
				<div class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
					<div class="py-4">
						<div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
							<!-- Quick Stats -->
							<div class="bg-white overflow-hidden shadow rounded-lg">
								<div class="p-5">
									<div class="flex items-center">
										<div class="flex-shrink-0">
											<svg class="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
											</svg>
										</div>
										<div class="ml-5 w-0 flex-1">
											<dl>
												<dt class="text-sm font-medium text-gray-500 truncate">
													Active Learners
												</dt>
												<dd class="flex items-baseline">
													<div class="text-2xl font-semibold text-gray-900">
														{#if isLoading}
															<svg class="animate-spin h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
																<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
																<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
															</svg>
														{:else}
															{activeLearners}
														{/if}
													</div>
												</dd>
											</dl>
										</div>
									</div>
								</div>
							</div>

							<div class="bg-white overflow-hidden shadow rounded-lg">
								<div class="p-5">
									<div class="flex items-center">
										<div class="flex-shrink-0">
											<svg class="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
											</svg>
										</div>
										<div class="ml-5 w-0 flex-1">
											<dl>
												<dt class="text-sm font-medium text-gray-500 truncate">
													Practice Topics
												</dt>
												<dd class="flex items-baseline">
													<div class="text-2xl font-semibold text-gray-900">
														{#if isLoading}
															<svg class="animate-spin h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
																<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
																<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
															</svg>
														{:else}
															{practiceTopics}
														{/if}
													</div>
												</dd>
											</dl>
										</div>
									</div>
								</div>
							</div>

							<div class="bg-white overflow-hidden shadow rounded-lg">
								<div class="p-5">
									<div class="flex items-center">
										<div class="flex-shrink-0">
											<svg class="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
											</svg>
										</div>
										<div class="ml-5 w-0 flex-1">
											<dl>
												<dt class="text-sm font-medium text-gray-500 truncate">
													Completed Sessions
												</dt>
												<dd class="flex items-baseline">
													<div class="text-2xl font-semibold text-gray-900">
														{#if isLoading}
															<svg class="animate-spin h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
																<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
																<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
															</svg>
														{:else}
															{completedSessions}
														{/if}
													</div>
												</dd>
											</dl>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</main>
	</div>
</div>

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