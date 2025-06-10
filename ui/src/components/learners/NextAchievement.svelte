<script lang="ts">
  import type { Achievement } from '$lib/types';
  import { AchievementIconMap } from '$lib/types';

  /**
   * Next achievement details including progress and requirements
   * @type {Achievement}
   */
  export let achievement: Achievement;

  $: progressPercentage = Math.round((achievement.progress / achievement.requiredProgress) * 100);
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
  <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
    Next Achievement
  </h2>

  <div class="flex items-start gap-4">
    <div class="flex-shrink-0">
      <svg 
        class="w-16 h-16 text-gray-400 dark:text-gray-500" 
        viewBox="0 0 24 24" 
        fill="currentColor"
      >
        <path d={AchievementIconMap[achievement.icon]} />
      </svg>
    </div>

    <div class="flex-1">
      <h3 class="font-medium text-gray-900 dark:text-white">
        {achievement.title}
      </h3>
      
      <p class="text-sm text-gray-600 dark:text-gray-300 mt-1">
        {achievement.description}
      </p>

      <div class="mt-3">
        <div class="flex justify-between text-sm mb-1">
          <span class="text-gray-600 dark:text-gray-300">Progress</span>
          <span class="text-gray-900 dark:text-white font-medium">{progressPercentage}%</span>
        </div>
        <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
          <div 
            class="bg-blue-600 h-2 rounded-full transition-all duration-300" 
            style="width: {progressPercentage}%"
          ></div>
        </div>
      </div>

      {#if achievement.actions.length > 0}
        <div class="mt-4">
          <h4 class="text-sm font-medium text-gray-900 dark:text-white mb-2">
            To unlock:
          </h4>
          <ul class="space-y-1">
            {#each achievement.actions as action}
              <li class="text-sm text-gray-600 dark:text-gray-300 flex items-center">
                <svg class="w-4 h-4 mr-2" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4" />
                </svg>
                {action}
              </li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
  </div>
</div> 