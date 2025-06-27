<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import type { Learner, PracticeTopic, LearnerProgress, AccountStats } from '$lib/types';
    import { learnersService } from '$lib/services/learners';
    import { topicsService } from '$lib/services/topics';
    import { sessionService } from '$lib/services/session';
    import { accountService } from '$lib/services/accounts';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    import LearnerSummary from '$components/learners/LearnerSummary.svelte';
    import StatCard from '$components/dashboard/StatCard.svelte';
    import Icon from '$components/common/Icon.svelte';

    let learners: Learner[] = [];
    let topics: PracticeTopic[] = [];
    let learnerProgress: Record<string, LearnerProgress> = {};
    let accountStats: AccountStats | null = null;
    let loading = false;
    let error: string | null = null;

    /**
     * Load all necessary data for the dashboard.
     * Fetches learners, topics, progress for all learners, and account stats in parallel.
     */
    async function loadDashboardData(): Promise<void> {
        loading = true;
        error = null;
        
        try {
            const [learnersData, topicsData, progressData, stats] = await Promise.all([
                learnersService.getLearners(),
                topicsService.getTopics(),
                sessionService.getLearnersProgressForAccount(),
                accountService.getAccountStats()
            ]);
            
            learners = learnersData;
            topics = topicsData;
            learnerProgress = progressData;
            accountStats = stats;

        } catch (err) {
            console.error('Dashboard initialization failed:', err);
            error = 'Failed to load dashboard data. Please try again.';
        } finally {
            loading = false;
        }
    }

    onMount(loadDashboardData);

    // Navigation handlers
    function handleAddChild() {
        goto('/account/learners/create');
    }

    function handleAddTopic() {
        goto('/account/practice-topics/create');
    }

    function handleViewProgress(learnerId: string) {
        goto(`/account/learners/${learnerId}`);
    }
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
    <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Family Dashboard</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-2">Monitor your children's progress and manage their learning materials</p>
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <LoadingSpinner size="lg" color="primary" />
        </div>
    {:else if error}
        <ErrorAlert message={error} />
    {:else}
        <!-- Quick Stats -->
        <div class="grid grid-cols-1 gap-5 sm:grid-cols-3 mb-8">
            <StatCard 
                title="Learning Topics" 
                value={accountStats?.total_practice_topics ?? 0} 
                isLoading={loading} 
                icon="topic"
            />
            
            <StatCard 
                title="Study Materials" 
                value={accountStats?.total_practice_items ?? 0} 
                isLoading={loading} 
                icon="practice"
            />

            <StatCard 
                title="Total Exercises Done" 
                value={accountStats?.total_practice_results ?? 0} 
                isLoading={loading} 
                icon="answers"
            />
        </div>

        <!-- No Topics Guide -->
        {#if topics.length === 0}
            <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-8 mb-8 text-center">
                <div class="w-16 h-16 mx-auto mb-4 text-blue-500 dark:text-blue-400">
                    <Icon type="practice" class_name="w-full h-full" />
                </div>
                <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                    Create Your First Learning Topic
                </h3>
                <p class="text-gray-600 dark:text-gray-400 mb-6 max-w-md mx-auto">
                    Get started by creating topics for your children to study. Topics help organize study materials by subject or grade level.
                </p>
                <button
                    on:click={handleAddTopic}
                    class="inline-flex items-center px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors duration-200"
                >
                    <Icon type="add" class_name="w-5 h-5 mr-2" />
                    Add First Topic
                </button>
            </div>
        {/if}

        <!-- No Children Guide -->
        {#if learners.length === 0}
            <div class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-8 mb-8 text-center">
                <div class="w-16 h-16 mx-auto mb-4 text-green-500 dark:text-green-400">
                    <Icon type="learner" class_name="w-full h-full" />
                </div>
                <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                    Add Your Child
                </h3>
                <p class="text-gray-600 dark:text-gray-400 mb-6 max-w-md mx-auto">
                    Create profiles for your children to track their progress and personalize their learning experience.
                </p>
                <button
                    on:click={handleAddChild}
                    class="inline-flex items-center px-6 py-3 bg-green-600 hover:bg-green-700 text-white font-medium rounded-lg transition-colors duration-200"
                >
                    <Icon type="add" class_name="w-5 h-5 mr-2" />
                    Add Child
                </button>
            </div>
        {/if}

        <!-- Children's Progress Section -->
        {#if learners.length > 0}
            <div class="space-y-8">
                <div class="flex items-center justify-between">
                    <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
                        Children's Progress
                    </h2>
                
                </div>

                {#each learners as learner}
                    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700">
                        <!-- Learner Header -->
                        <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                            <div class="flex items-center justify-between">
                                <div class="flex items-center space-x-4">
                                    <div class="w-12 h-12 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
                                        {#if learner.avatar}
                                            <img
                                                src={learner.avatar}
                                                alt={`${learner.nickname}'s avatar`}
                                                class="w-full h-full rounded-full object-cover"
                                            />
                                        {:else}
                                            <span class="text-xl text-gray-500 dark:text-gray-400">
                                                {learner.nickname.charAt(0).toUpperCase()}
                                            </span>
                                        {/if}
                                    </div>
                                    <div>
                                        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                                            {learner.nickname}'s Progress
                                        </h3>
                                        {#if learner.grade_level}
                                            <p class="text-sm text-gray-600 dark:text-gray-400">
                                                Grade {learner.grade_level}
                                            </p>
                                        {/if}
                                    </div>
                                </div>
                                <button
                                    on:click={() => handleViewProgress(learner.id)}
                                    class="inline-flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg transition-colors duration-200"
                                >
                                    View Full Progress
                                    <Icon type="next" class_name="w-4 h-4 ml-2" />
                                </button>
                            </div>
                        </div>

                        <!-- Learning Summary Component -->
                        <div class="p-6">
                            {#if learnerProgress[learner.id]}
                                <LearnerSummary progress={learnerProgress[learner.id]} />
                            {:else}
                                <div class="text-center text-gray-500 dark:text-gray-400 p-4">
                                    No progress data available for this learner yet.
                                </div>
                            {/if}
                        </div>
                    </div>
                {/each}
            </div>
        {/if}

        <!-- Quick Actions Section -->
        {#if learners.length > 0 && topics.length > 0}
            <div class="mt-12 bg-gray-50 dark:bg-gray-800/50 rounded-lg p-6">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Quick Actions</h3>
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                    <button
                        on:click={handleAddTopic}
                        class="flex items-center p-4 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 hover:shadow-md transition-shadow duration-200"
                    >
                        <div class="w-10 h-10 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center mr-3">
                            <Icon type="practice" class_name="w-5 h-5 text-blue-600 dark:text-blue-400" />
                        </div>
                        <div class="text-left">
                            <p class="font-medium text-gray-900 dark:text-white">New Topic</p>
                            <p class="text-sm text-gray-600 dark:text-gray-400">Add learning materials</p>
                        </div>
                    </button>

                    <button
                        on:click={handleAddChild}
                        class="flex items-center p-4 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 hover:shadow-md transition-shadow duration-200"
                    >
                        <div class="w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center mr-3">
                            <Icon type="learner" class_name="w-5 h-5 text-green-600 dark:text-green-400" />
                        </div>
                        <div class="text-left">
                            <p class="font-medium text-gray-900 dark:text-white">Add Child</p>
                            <p class="text-sm text-gray-600 dark:text-gray-400">Add child profile</p>
                        </div>
                    </button>

                    <a
                        href="/account/practice-topics"
                        class="flex items-center p-4 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 hover:shadow-md transition-shadow duration-200"
                    >
                        <div class="w-10 h-10 bg-purple-100 dark:bg-purple-900/30 rounded-lg flex items-center justify-center mr-3">
                            <Icon type="topic" class_name="w-5 h-5 text-purple-600 dark:text-purple-400" />
                        </div>
                        <div class="text-left">
                            <p class="font-medium text-gray-900 dark:text-white">All Topics</p>
                            <p class="text-sm text-gray-600 dark:text-gray-400">View learning materials</p>
                        </div>
                    </a>
                </div>
            </div>
        {/if}
    {/if}
</div> 