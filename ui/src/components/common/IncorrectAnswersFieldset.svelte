<script context="module" lang="ts">
    /**
     * Type definition for incorrect answer entries
     */
    export interface IncorrectAnswer {
        answer: string;
        explanation: string;
    }
</script>

<script lang="ts">
    import TextArea from './TextArea.svelte';
    
    /**
     * Array of incorrect answers with explanations
     */
    export let incorrectAnswers: IncorrectAnswer[] = [];
    
    /**
     * Whether the fieldset is disabled
     */
    export let disabled: boolean = false;
    
    /**
     * Add a new incorrect answer entry
     */
    function addIncorrectAnswer() {
        incorrectAnswers = [...incorrectAnswers, { answer: '', explanation: '' }];
    }
    
    /**
     * Remove an incorrect answer entry at the specified index
     */
    function removeIncorrectAnswer(index: number) {
        incorrectAnswers = incorrectAnswers.filter((_, i) => i !== index);
    }
</script>

<fieldset class="space-y-4" {disabled}>
    <legend class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
        Explanations for Incorrect Answers
        <span class="text-xs text-gray-500 dark:text-gray-400 ml-2">
            Optional - Add explanations for common incorrect answers
        </span>
    </legend>
    
    {#if incorrectAnswers.length === 0}
        <div class="text-center py-6 bg-gray-50 dark:bg-gray-800 rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-600">
            <p class="text-gray-500 dark:text-gray-400 text-sm mb-3">
                No incorrect answer explanations added yet
            </p>
            <button
                type="button"
                class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                on:click={addIncorrectAnswer}
                {disabled}
            >
                <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                Add First Incorrect Answer
            </button>
        </div>
    {:else}
        <div class="space-y-4">
            {#each incorrectAnswers as incorrect, index}
                <div class="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
                    <div class="flex items-start justify-between mb-3">
                        <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300">
                            Incorrect Answer #{index + 1}
                        </h4>
                        <button
                            type="button"
                            class="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 p-1 rounded-md hover:bg-red-50 dark:hover:bg-red-900/20 focus:outline-none focus:ring-2 focus:ring-red-500"
                            on:click={() => removeIncorrectAnswer(index)}
                            title="Remove this incorrect answer"
                            aria-label="Remove incorrect answer"
                            {disabled}
                        >
                            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                        </button>
                    </div>
                    
                    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
                        <div>
                            <label for={`incorrect-answer-${index}`} class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                Incorrect Answer
                            </label>
                            <input
                                type="text"
                                id={`incorrect-answer-${index}`}
                                bind:value={incorrect.answer}
                                class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white sm:text-sm"
                                placeholder="Enter the incorrect answer"
                                {disabled}
                            />
                        </div>
                        
                        <div>
                            <TextArea
                                id={`incorrect-explanation-${index}`}
                                label="Explanation"
                                bind:value={incorrect.explanation}
                                rows={3}
                                placeholder="Explain why this answer is incorrect"
                                {disabled}
                                cols="w-full"
                            />
                        </div>
                    </div>
                </div>
            {/each}
            
            <div class="flex justify-center pt-2">
                <button
                    type="button"
                    class="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                    on:click={addIncorrectAnswer}
                    {disabled}
                >
                    <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                    </svg>
                    Add Another Incorrect Answer
                </button>
            </div>
        </div>
    {/if}
</fieldset> 