<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import type { BreadcrumbItem, Learner, AchievementIcon, Achievement } from '$lib/types';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    import Breadcrumbs from '$components/common/Breadcrumbs.svelte';
    import { learnersService } from '$lib/services/learners';
    import WelcomeMessage from '$components/learners/WelcomeMessage.svelte';
    import RecommendedTopics from '$components/learners/RecommendedTopics.svelte';
    import NextAchievement from '$components/learners/NextAchievement.svelte';

    type Difficulty = 'easy' | 'medium' | 'hard';

    let loading = false;
    let error: string | null = null;
    let learnerId: string = '';
    let learner: Learner | null = null;
    let breadcrumbs: BreadcrumbItem[] = [];

    // Mock data for now - will be replaced with actual data from services
    const mockData = {
        recommendedTopics: [
            {
                id: '1',
                title: 'Multiplication Tables',
                description: 'Practice multiplication tables from 1 to 12 with interactive exercises.',
                difficulty: 'medium' as Difficulty,
                estimatedMinutes: 15
            },
            {
                id: '2',
                title: 'Fractions Basics',
                description: 'Learn about fractions and how to work with them.',
                difficulty: 'hard' as Difficulty,
                estimatedMinutes: 20
            },
            {
                id: '3',
                title: 'Basic Geometry',
                description: 'Explore shapes, angles, and basic geometric concepts.',
                difficulty: 'easy' as Difficulty,
                estimatedMinutes: 10
            }
        ],
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
            learner = await learnersService.getLearner(learnerId);
            
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

            <RecommendedTopics topics={mockData.recommendedTopics} />

            <NextAchievement achievement={mockData.nextAchievement} />
        </div>
    {/if}
</div> 