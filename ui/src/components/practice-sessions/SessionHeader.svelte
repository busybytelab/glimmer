<script lang="ts">
    import type { SessionWithExpandedData } from '$lib/services/session';

    /**
     * The session data to display in the header
     */
    export let session: SessionWithExpandedData;
</script>

<div class="mb-6">
    <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-4">
        {session.name || 'Practice Session'}
    </h2>

    {#if session.expand?.practice_topic}
        <p class="text-gray-600 dark:text-gray-300 mb-4">
            Topic: {session.expand.practice_topic.name}
        </p>
    {/if}

    {#if session.expand?.learner}
        <p class="text-gray-600 dark:text-gray-300 mb-4">
            Learner: {session.expand?.learner?.nickname || 'Unknown Learner'}
        </p>
    {/if}
</div> 