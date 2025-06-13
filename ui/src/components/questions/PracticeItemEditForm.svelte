<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    import { practiceItemService } from '$lib/services/practiceItem';
    import LoadingSpinner from '../common/LoadingSpinner.svelte';
    import ErrorAlert from '../common/ErrorAlert.svelte';
    import TextArea from '../common/TextArea.svelte';
    import FormField from '../common/FormField.svelte';

    /**
     * The practice item to be edited
     */
    export let item: PracticeItem;

    /**
     * Event handler for when the item is saved
     */
    export let onSave: (updatedItem: PracticeItem) => void;

    /**
     * Event handler for when editing is cancelled
     */
    export let onCancel: () => void;

    let questionText = item.question_text;
    let correctAnswer = typeof item.correct_answer === 'string' ? item.correct_answer : JSON.stringify(item.correct_answer);
    let explanation = item.explanation;
    let hints = Array.isArray(item.hints) ? item.hints.join('\n') : '';
    let options = Array.isArray(item.options) ? item.options.join('\n') : '';
    let incorrectAnswers: { answer: string; explanation: string }[] = [];
    let tags = Array.isArray(item.tags) ? item.tags.join(', ') : '';
    let loading = false;
    let error: string | null = null;

    // Initialize incorrect answers from the item
    if (item.explanation_for_incorrect) {
        try {
            incorrectAnswers = Object.entries(item.explanation_for_incorrect).map(([answer, explanation]) => ({
                answer,
                explanation
            }));
        } catch (err) {
            console.error('Failed to parse explanation_for_incorrect:', err);
            incorrectAnswers = [];
        }
    }

    function addIncorrectAnswer() {
        incorrectAnswers = [...incorrectAnswers, { answer: '', explanation: '' }];
    }

    function removeIncorrectAnswer(index: number) {
        incorrectAnswers = incorrectAnswers.filter((_, i) => i !== index);
    }

    async function handleSubmit() {
        if (!item.id) return;

        loading = true;
        error = null;

        try {
            // Parse hints back into array
            const hintsArray = hints.split('\n').filter(hint => hint.trim());
            
            // Parse options back into array
            const optionsArray = options.split('\n').filter(option => option.trim());
            
            // Convert incorrect answers to object and stringify
            const explanationForIncorrect = incorrectAnswers.reduce((acc, { answer, explanation }) => {
                if (answer.trim() && explanation.trim()) {
                    acc[answer.trim()] = explanation.trim();
                }
                return acc;
            }, {} as Record<string, string>);
            
            // Parse tags back into array
            const tagsArray = tags.split(',').map(tag => tag.trim()).filter(tag => tag);

            const updatedItem = await practiceItemService.updateItem(item.id, {
                question_text: questionText,
                correct_answer: correctAnswer,
                explanation,
                hints: hintsArray,
                options: optionsArray,
                explanation_for_incorrect: explanationForIncorrect,
                tags: tagsArray,
                review_status: 'APPROVED',
                review_date: new Date().toISOString()
            });

            onSave(updatedItem);
        } catch (err) {
            console.error('Failed to update practice item:', err);
            error = err instanceof Error ? err.message : 'Failed to update practice item';
        } finally {
            loading = false;
        }
    }
</script>

<div class="mt-4 p-4 bg-white dark:bg-gray-800 rounded-lg shadow">
    <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Edit Practice Item</h3>

    {#if error}
        <ErrorAlert message={error} />
    {/if}

    <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        <TextArea
            id="questionText"
            label="Question Text"
            bind:value={questionText}
            rows={3}
            required
        />

        <TextArea
            id="correctAnswer"
            label="Correct Answer"
            bind:value={correctAnswer}
            rows={2}
            required
        />

        <TextArea
            id="explanation"
            label="Explanation"
            bind:value={explanation}
            rows={3}
            required
        />

        <TextArea
            id="hints"
            label="Hints (one per line)"
            bind:value={hints}
            rows={3}
            placeholder="Enter each hint on a new line"
        />

        <TextArea
            id="options"
            label="Options (one per line)"
            bind:value={options}
            rows={3}
            placeholder="Enter each option on a new line"
        />

        <div>
            <fieldset class="space-y-4">
                <legend class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Explanations for Incorrect Answers
                </legend>
                {#each incorrectAnswers as incorrect, index}
                    <div class="flex gap-4 items-start">
                        <div class="flex-1">
                            <label for={`incorrect-answer-${index}`} class="sr-only">Incorrect answer</label>
                            <input
                                type="text"
                                id={`incorrect-answer-${index}`}
                                bind:value={incorrect.answer}
                                class="mt-1 block w-full rounded-md border-gray-300 dark:border-gray-600 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:bg-gray-700 dark:text-white sm:text-sm"
                                placeholder="Incorrect answer"
                            />
                        </div>
                        <div class="flex-1">
                            <TextArea
                                id={`incorrect-explanation-${index}`}
                                bind:value={incorrect.explanation}
                                rows={2}
                                placeholder="Explanation for this incorrect answer"
                                label=""
                            />
                        </div>
                        <button
                            type="button"
                            class="mt-1 text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                            on:click={() => removeIncorrectAnswer(index)}
                        >
                            Remove
                        </button>
                    </div>
                {/each}
                <button
                    type="button"
                    class="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                    on:click={addIncorrectAnswer}
                >
                    Add Incorrect Answer
                </button>
            </fieldset>
        </div>

        <div>
            <FormField
                id="tags"
                label="Tags (comma-separated)"
                type="text"
                bind:value={tags}
                placeholder="Enter tags separated by commas"
                cols="col-span-6"
            />
            {#if tags.split(',').map(tag => tag.trim()).filter(Boolean).length > 0}
                <div class="mt-2 flex flex-wrap gap-2">
                    {#each tags.split(',').map(tag => tag.trim()).filter(Boolean) as tag}
                        <span class="bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-300 text-xs font-medium px-2.5 py-0.5 rounded">
                            {tag}
                        </span>
                    {/each}
                </div>
            {/if}
        </div>

        <div class="flex justify-end space-x-3">
            <button
                type="button"
                class="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                on:click={onCancel}
                disabled={loading}
            >
                Cancel
            </button>
            <button
                type="submit"
                class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                disabled={loading}
            >
                {#if loading}
                    <LoadingSpinner size="sm" color="white" />
                {/if}
                Save Changes
            </button>
        </div>
    </form>
</div> 