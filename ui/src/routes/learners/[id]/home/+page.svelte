<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import type { BreadcrumbItem, Learner, AchievementIcon, Achievement, PracticeSessionStats } from '$lib/types';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    import Breadcrumbs from '$components/common/Breadcrumbs.svelte';
    import { learnersService } from '$lib/services/learners';
    import { sessionService } from '$lib/services/session';
    import WelcomeMessage from '$components/learners/WelcomeMessage.svelte';
    import RecommendedTopics from '$components/learners/RecommendedTopics.svelte';
    import NextAchievement from '$components/learners/NextAchievement.svelte';

    let loading = false;
    let error: string | null = null;
    let learnerId: string = '';
    let learner: Learner | null = null;
    let sessionStats: PracticeSessionStats[] = [];
    let breadcrumbs: BreadcrumbItem[] = [];

    // Mock data for achievements - will be replaced with actual data later
    const mockData = {
        latestAchievement: {
            title: 'math-apprentice',
            icon: 'math-apprentice' as AchievementIcon,
            description: 'Completed 5 math practice sessions'
        },
        nextAchievement: {
            title: 'Math Whiz',
            description: 'Complete 10 math practice sessions with a score of 80% or higher',
            icon: 'math-whiz' as AchievementIcon,
            progress: 7,
            requiredProgress: 10,
            actions: [
                'Complete 3 more math practice sessions',
                'Maintain a score of 80% or higher',
                'Try different topics'
            ]
        } satisfies Achievement
    };

    onMount(async () => {
        try {
            loading = true;
            learnerId = $page.params.id;
            
            // Load learner and session stats in parallel
            const [learnerData, stats] = await Promise.all([
                learnersService.getLearner(learnerId),
                sessionService.getSessionStatsForLearner(learnerId)
            ]);
            
            learner = learnerData;
            sessionStats = stats;
            
            breadcrumbs = [
                { label: 'Home', href: '/', icon: 'home' },
                { label: learner.nickname, href: `/learners/${learnerId}/home`, icon: 'learner' }
            ];
        } catch (e) {
            error = 'Failed to load learner data';
            console.error(e);
        } finally {
            loading = false;
        }
    });

    // Transform session stats into the format expected by RecommendedTopics
    $: recommendedTopics = sessionStats.map(stat => ({
        id: stat.id,
        title: stat.topic_name,
        description: `Progress: ${stat.answered_items}/${stat.total_items} items completed. Score: ${stat.total_score}%`
    }));
</script>

<div class="container mx-auto px-4 py-8">
    <Breadcrumbs items={breadcrumbs} />

    {#if loading}
        <div class="flex justify-center items-center min-h-[400px]">
            <LoadingSpinner />
        </div>
    {:else if error}
        <ErrorAlert message={error} />
    {:else if learner}
        <div class="space-y-8">
            <WelcomeMessage 
                name={learner.nickname}
                age={learner.age}
                latestAchievement={mockData.latestAchievement}
            />

            <RecommendedTopics topics={recommendedTopics} />

            <NextAchievement achievement={mockData.nextAchievement} />
        </div>
    {/if}
</div> 