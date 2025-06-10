<script lang="ts">
  import { goto } from '$app/navigation';
  
  /**
   * Array of recommended topics
   */
  export let topics: Array<{
    id: string;
    title: string;
    description: string;
    difficulty: 'easy' | 'medium' | 'hard';
    estimatedMinutes: number;
  }>;

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case 'easy': return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300';
      case 'medium': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300';
      case 'hard': return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300';
      default: return '';
    }
  };

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
            <div class="flex items-center gap-2 mb-2">
              <h3 class="font-medium text-gray-900 dark:text-white">
                {topic.title}
              </h3>
              <span class="px-2 py-1 text-xs rounded-full {getDifficultyColor(topic.difficulty)}">
                {topic.difficulty}
              </span>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-300 mb-3">
              {topic.description}
            </p>
            <div class="flex items-center text-sm text-gray-500 dark:text-gray-400">
              <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {topic.estimatedMinutes} minutes
            </div>
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