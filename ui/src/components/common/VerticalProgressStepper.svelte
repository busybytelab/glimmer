<script lang="ts">
    export let steps: number;
    export let currentStep: number = 0;
    export let results: ('correct' | 'incorrect' | 'pending')[] = [];
    export let onStepClick: (index: number) => void = () => {};
</script>

<div class="flex flex-col space-y-2">
    {#each Array(steps) as _, index}
        <div class="flex items-center">
            <!-- Step circle -->
            <button
                class="flex items-center justify-center w-6 h-6 rounded-full border-2 transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 text-sm"
                class:border-primary-500={index === currentStep}
                class:border-gray-300={index !== currentStep}
                class:bg-white={true}
                class:text-primary-500={index === currentStep}
                class:text-gray-500={index !== currentStep}
                on:click={() => onStepClick(index)}
            >
                {index + 1}
            </button>

            <!-- Connecting line -->
            <div 
                class="w-0.5 h-8 ml-3"
                class:bg-primary-500={index < currentStep}
                class:bg-gray-300={index >= currentStep}
            ></div>
        </div>
    {/each}
</div> 