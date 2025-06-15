<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    import { practiceItemService } from '$lib/services/practiceItem';
    import { 
        practiceItemToFormData, 
        formDataToPracticeItemUpdate, 
        validatePracticeItemForm,
        type PracticeItemFormData 
    } from '$lib/utils/practiceItemForm';
    import LoadingSpinner from '../common/LoadingSpinner.svelte';
    import ErrorAlert from '../common/ErrorAlert.svelte';
    import TextArea from '../common/TextArea.svelte';
    import FormSection from '../common/FormSection.svelte';
    import FormButton from '../common/FormButton.svelte';
    import IncorrectAnswersFieldset from '../common/IncorrectAnswersFieldset.svelte';
    import TagsInput from '../common/TagsInput.svelte';

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

    // Form state
    let formData: PracticeItemFormData = practiceItemToFormData(item);
    let loading = false;
    let validationErrors: string[] = [];
    let serverError: string | null = null;

    /**
     * Handle form submission
     */
    async function handleSubmit() {
        if (!item.id) return;

        // Clear previous errors
        validationErrors = [];
        serverError = null;

        // Validate form data
        validationErrors = validatePracticeItemForm(formData);
        if (validationErrors.length > 0) {
            return;
        }

        loading = true;

        try {
            const updateData = formDataToPracticeItemUpdate(formData);
            const updatedItem = await practiceItemService.updateItem(item.id, updateData);
            onSave(updatedItem);
        } catch (err) {
            console.error('Failed to update practice item:', err);
            serverError = err instanceof Error ? err.message : 'Failed to update practice item';
        } finally {
            loading = false;
        }
    }

    /**
     * Handle cancel action
     */
    function handleCancel() {
        onCancel();
    }
</script>

<FormSection title="Edit Practice Item" description="Update the practice item details and content">
    <div class="px-4 py-5 sm:p-6 space-y-6">
        <!-- Error Messages -->
        {#if serverError}
            <ErrorAlert message={serverError} />
        {/if}

        {#if validationErrors.length > 0}
            <ErrorAlert 
                title="Please fix the following errors:"
                message={validationErrors.join(', ')}
            />
        {/if}

        <form on:submit|preventDefault={handleSubmit} class="space-y-6">
            <!-- Basic Information Section -->
            <div class="grid grid-cols-1 gap-6">
                <TextArea
                    id="questionText"
                    label="Question Text"
                    bind:value={formData.questionText}
                    rows={3}
                    required
                    disabled={loading}
                    placeholder="Enter the question that will be presented to learners"
                    cols="col-span-1"
                />

                <TextArea
                    id="correctAnswer"
                    label="Correct Answer"
                    bind:value={formData.correctAnswer}
                    rows={2}
                    required
                    disabled={loading}
                    placeholder="Enter the correct answer"
                    cols="col-span-1"
                />

                <TextArea
                    id="explanation"
                    label="Explanation"
                    bind:value={formData.explanation}
                    rows={3}
                    required
                    disabled={loading}
                    placeholder="Explain why this answer is correct"
                    cols="col-span-1"
                />
            </div>

            <!-- Additional Content Section -->
            <div class="border-t border-gray-200 dark:border-gray-700 pt-6">
                <h4 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
                    Additional Content
                </h4>
                
                <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    <TextArea
                        id="hints"
                        label="Hints"
                        bind:value={formData.hints}
                        rows={4}
                        disabled={loading}
                        placeholder="Enter each hint on a new line"
                        cols="col-span-1"
                    />

                    <TextArea
                        id="options"
                        label="Answer Options"
                        bind:value={formData.options}
                        rows={4}
                        disabled={loading}
                        placeholder="Enter each option on a new line"
                        cols="col-span-1"
                    />
                </div>
            </div>

            <!-- Incorrect Answers Section -->
            <div class="border-t border-gray-200 dark:border-gray-700 pt-6">
                <IncorrectAnswersFieldset
                    bind:incorrectAnswers={formData.incorrectAnswers}
                    disabled={loading}
                />
            </div>

            <!-- Tags Section -->
            <div class="border-t border-gray-200 dark:border-gray-700 pt-6">
                <TagsInput
                    id="tags"
                    label="Tags"
                    bind:value={formData.tags}
                    disabled={loading}
                    placeholder="Enter tags separated by commas (e.g., math, algebra, equations)"
                    cols="col-span-1"
                />
            </div>

            <!-- Form Actions -->
            <div class="border-t border-gray-200 dark:border-gray-700 pt-6">
                <div class="flex justify-end space-x-3">
                    <button
                        type="button"
                        class="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                        on:click={handleCancel}
                        disabled={loading}
                    >
                        Cancel
                    </button>
                    
                    <FormButton
                        type="submit"
                        variant="primary"
                        disabled={loading}
                        isLoading={loading}
                        loadingText="Saving Changes..."
                    >
                        {#if loading}
                            <LoadingSpinner size="sm" color="white" />
                            <span class="ml-2">Saving Changes...</span>
                        {:else}
                            Save Changes
                        {/if}
                    </FormButton>
                </div>
            </div>
        </form>
    </div>
</FormSection> 