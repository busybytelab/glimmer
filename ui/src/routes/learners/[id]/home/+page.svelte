<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import type { Learner, AchievementIcon, Achievement, PracticeSessionStats } from '$lib/types';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    import { learnersService } from '$lib/services/learners';
    import { sessionService } from '$lib/services/session';
    import WelcomeMessage from '$components/learners/WelcomeMessage.svelte';
    import AvailableSessions from '$components/learners/AvailableSessions.svelte';
    import NextAchievement from '$components/learners/NextAchievement.svelte';

    let loading = false;
    let error: string | null = null;
    let learnerId: string = '';
    let learner: Learner | null = null;
    let activeSessionStats: PracticeSessionStats[] = [];
    let completedSessionStats: PracticeSessionStats[] = [];

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
            const [learnerData, sessionStats] = await Promise.all([
                learnersService.getLearner(learnerId),
                sessionService.getSessionStatsForLearner(learnerId)
            ]);
            
            learner = learnerData;
            activeSessionStats = sessionStats.active;
            completedSessionStats = sessionStats.completed;
        } catch (e) {
            error = 'Failed to load learner data';
            console.error(e);
        } finally {
            loading = false;
        }
    });

    // Transform session stats into the format expected by AvailableSessions
    $: availableSessions = activeSessionStats;
    $: completedSessions = completedSessionStats;
</script>

<div class="container mx-auto px-4 py-8">

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

            <AvailableSessions sessions={availableSessions} learnerId={learnerId} />

            {#if completedSessions.length > 0}
                <AvailableSessions 
                    sessions={completedSessions} 
                    learnerId={learnerId} 
                    title="Completed Sessions"
                    mode="completed"
                />
            {/if}

            <NextAchievement achievement={mockData.nextAchievement} />
        </div>
    {/if}
</div> 