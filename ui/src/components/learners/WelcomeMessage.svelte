<script lang="ts">
  import type { LatestAchievement } from '$lib/types';
  import { AchievementIconMap } from '$lib/types';

  /**
   * The learner's name
   */
  export let name: string;
  
  /**
   * The learner's age
   */
  export let age: number;
  
  /**
   * Latest achievement badge details
   * @type {LatestAchievement | null}
   */
  export let latestAchievement: LatestAchievement | null = null;

  // Get appropriate greeting based on time of day
  const getTimeBasedGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return 'Good morning';
    if (hour < 17) return 'Good afternoon';
    return 'Good evening';
  };

  // Get age-appropriate message
  const getAgeAppropriateMessage = (age: number) => {
    if (age < 8) return 'Ready for some fun learning?';
    if (age < 12) return 'Ready to take on today\'s challenges?';
    return 'Ready to level up your skills?';
  };
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6 mb-6">
  <div class="flex items-start justify-between">
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">
        {getTimeBasedGreeting()}, {name}! 👋
      </h1>
      <p class="text-gray-600 dark:text-gray-300 text-lg">
        {getAgeAppropriateMessage(age)}
      </p>
    </div>
    
    {#if latestAchievement}
      <div class="flex items-center bg-blue-50 dark:bg-blue-900/30 rounded-lg p-3">
        <svg 
          class="w-12 h-12 mr-3 text-blue-600 dark:text-blue-400" 
          viewBox="0 0 24 24" 
          fill="currentColor"
        >
          <path d={AchievementIconMap[latestAchievement.icon]} />
        </svg>
        <div>
          <p class="font-medium text-blue-900 dark:text-blue-100">
            Latest Achievement
          </p>
          <p class="text-sm text-blue-700 dark:text-blue-300">
            {latestAchievement.title}
          </p>
        </div>
      </div>
    {/if}
  </div>
</div> 