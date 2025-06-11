<script lang="ts">
  import { goto } from '$app/navigation';
  
  /**
   * Array of recommended topics with their progress
   */
  export let topics: Array<{
    id: string;
    title: string;
    description: string;
  }>;

  const startPractice = (topicId: string) => {
    // TODO: fix this, learner should not create session, they start session and url must be /learnerss/[id]/practice-topics/[topicId]
    goto(`/practice-topics/${topicId}/create-session`);
  };
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6 mb-6">
  <div class="flex items-center justify-between mb-4">
    <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
      Recommended Practice Topics
    </h2>
    <a 
      href="/practice-topics" 
      class="text-blue-600 dark:text-blue-400 hover:underline text-sm"
    >
      View All Topics
    </a>
  </div>

  <div class="space-y-4">
    {#each topics as topic (topic.id)}
      <div class="border dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <h3 class="font-medium text-gray-900 dark:text-white mb-2">
              {topic.title}
            </h3>
            <p class="text-sm text-gray-600 dark:text-gray-300">
              {topic.description}
            </p>
          </div>
          <button
            on:click={() => startPractice(topic.id)}
            class="ml-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm font-medium"
          >
            Start Practice
          </button>
        </div>
      </div>
    {/each}
  </div>
</div> 