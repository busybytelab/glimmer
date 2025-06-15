<script lang="ts">
    import { QuestionViewType, IconTypeMap, type IconType } from '$lib/types';
    
    /**
     * Current view type selected
     */
    export let viewType: QuestionViewType;
    
    /**
     * Event handler for when view type changes
     */
    export let onViewChange: (newViewType: QuestionViewType) => void;
    
    /**
     * Whether the user is an instructor (determines available views)
     */
    export let isInstructor: boolean = false;
    
    /**
     * Whether to show descriptions (helpful for first-time users)
     * @default true
     */
    export let showDescriptions: boolean = true;
    
    interface ViewOption {
        value: QuestionViewType;
        label: string;
        description: string;
        icon: IconType;
        color: string;
    }
    
    const viewOptions: ViewOption[] = [
        { 
            value: QuestionViewType.LEARNER, 
            label: 'Practice Mode', 
            description: 'See questions exactly as your child sees them during practice',
            icon: 'practice',
            color: 'blue'
        },
        { 
            value: QuestionViewType.PARENT, 
            label: 'Progress Overview', 
            description: 'See your child\'s learning progress, strengths, and areas to improve',
            icon: 'progress',
            color: 'green'
        },
        { 
            value: QuestionViewType.ANSWERED, 
            label: 'Questions & Answers', 
            description: 'See your child\'s responses alongside the correct answers',
            icon: 'answers',
            color: 'orange'
        },
        { 
            value: QuestionViewType.GENERATED, 
            label: 'Question Review', 
            description: 'Check, approve, or edit questions to ensure they\'re appropriate',
            icon: 'review',
            color: 'purple'
        }
    ];
    
    // Filter options based on user role
    $: availableOptions = isInstructor 
        ? viewOptions 
        : viewOptions.filter(opt => opt.value !== QuestionViewType.PARENT);
    
    /**
     * Get Tailwind color classes for each view option
     */
    function getColorClasses(color: string, isActive: boolean): string {
        const colorMap = {
            blue: isActive 
                ? 'bg-blue-600 text-white border-blue-600' 
                : 'bg-blue-50 text-blue-700 border-blue-200 hover:bg-blue-100 dark:bg-blue-900/20 dark:text-blue-300 dark:border-blue-800 dark:hover:bg-blue-900/30',
            green: isActive 
                ? 'bg-green-600 text-white border-green-600' 
                : 'bg-green-50 text-green-700 border-green-200 hover:bg-green-100 dark:bg-green-900/20 dark:text-green-300 dark:border-green-800 dark:hover:bg-green-900/30',
            orange: isActive 
                ? 'bg-orange-600 text-white border-orange-600' 
                : 'bg-orange-50 text-orange-700 border-orange-200 hover:bg-orange-100 dark:bg-orange-900/20 dark:text-orange-300 dark:border-orange-800 dark:hover:bg-orange-900/30',
            purple: isActive 
                ? 'bg-purple-600 text-white border-purple-600' 
                : 'bg-purple-50 text-purple-700 border-purple-200 hover:bg-purple-100 dark:bg-purple-900/20 dark:text-purple-300 dark:border-purple-800 dark:hover:bg-purple-900/30'
        };
        return colorMap[color as keyof typeof colorMap] || colorMap.blue;
    }
</script>

<div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-4 mb-6">
    <div class="mb-3">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-1">
            Choose Your View
        </h3>
        <p class="text-sm text-gray-600 dark:text-gray-400">
            Select how you'd like to see the questions and progress
        </p>
    </div>
    
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-{availableOptions.length} gap-3">
        {#each availableOptions as option}
            <button 
                class="relative p-4 border-2 rounded-lg transition-all duration-200 text-left focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 {getColorClasses(option.color, viewType === option.value)}"
                on:click={() => onViewChange(option.value)}
                type="button"
                aria-pressed={viewType === option.value ? 'true' : 'false'}
                title={option.description}
            >
                <!-- Selection indicator -->
                {#if viewType === option.value}
                    <div class="absolute top-2 right-2">
                        <div class="w-2 h-2 bg-white rounded-full"></div>
                    </div>
                {/if}
                
                <!-- Icon and label -->
                <div class="flex items-center mb-2">
                    <svg class="w-6 h-6 mr-3 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                        <path stroke-linecap="round" stroke-linejoin="round" d={IconTypeMap[option.icon]} />
                    </svg>
                    <span class="font-semibold text-base">
                        {option.label}
                    </span>
                </div>
                
                <!-- Description -->
                {#if showDescriptions}
                    <p class="text-sm opacity-90 leading-relaxed">
                        {option.description}
                    </p>
                {/if}
            </button>
        {/each}
    </div>
    
    <!-- Help text for first-time users -->
    <div class="mt-4 p-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
        <div class="flex items-start">
            <svg class="w-5 h-5 text-blue-500 mr-2 mt-0.5 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d={IconTypeMap.hint} />
            </svg>
            <div class="text-sm text-gray-700 dark:text-gray-300">
                <span class="font-medium">Tip:</span> 
                You can switch between views anytime to see different perspectives of your child's learning progress.
            </div>
        </div>
    </div>
</div> 