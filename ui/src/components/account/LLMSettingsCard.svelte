<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { error } from '$lib/stores'; 
    import FormField from '../common/FormField.svelte';
    import FormButton from '../common/FormButton.svelte';
    import ErrorAlert from '../common/ErrorAlert.svelte';
    
    // This would typically come from an API or store
    let apiSettings = {
        openaiApiKey: '',
        ollamaEndpoint: 'http://localhost:11434',
        defaultModel: 'gpt-3.5-turbo'
    };
    
    export let isLoading = false;
    let savedSuccessfully = false;
    
    const dispatch = createEventDispatcher();
    
    const modelOptions = [
        { value: 'gpt-3.5-turbo', label: 'GPT-3.5 Turbo (OpenAI)' },
        { value: 'gpt-4', label: 'GPT-4 (OpenAI)' },
        { value: 'llama2', label: 'Llama 2 (Ollama)' },
        { value: 'mistral', label: 'Mistral (Ollama)' }
    ];
    
    function handleSubmit() {
        isLoading = true;
        error.set(null);
        
        // Simulate API call to save settings
        setTimeout(() => {
            // Here you would actually save the settings to your backend
            console.log('Saving LLM settings:', apiSettings);
            
            savedSuccessfully = true;
            isLoading = false;
            
            // Reset success message after 3 seconds
            setTimeout(() => {
                savedSuccessfully = false;
            }, 3000);
            
            // Dispatch event for parent components
            dispatch('save', apiSettings);
        }, 500);
    }
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden mb-6">
    <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
        <h2 class="text-lg font-medium text-gray-900 dark:text-white">LLM Settings</h2>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Configure your language model providers and preferences
        </p>
    </div>
    
    {#if savedSuccessfully}
        <div class="px-6 py-4 bg-green-50 dark:bg-green-900/20 border-b border-green-100 dark:border-green-900/30">
            <div class="flex items-center">
                <svg class="h-5 w-5 text-green-500 dark:text-green-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <p class="ml-3 text-sm font-medium text-green-700 dark:text-green-400">
                    Your LLM settings have been updated successfully!
                </p>
            </div>
        </div>
    {/if}
    
    {#if $error}
        <div class="px-6 py-4">
            <ErrorAlert message={$error} />
        </div>
    {/if}
    
    <div class="px-6 pb-6">
        <form on:submit|preventDefault={handleSubmit}>
            <div class="py-4">
                <div class="grid grid-cols-6 gap-6">
                    <!-- OpenAI API Key -->
                    <FormField 
                        id="openaiApiKey"
                        label="OpenAI API Key"
                        bind:value={apiSettings.openaiApiKey}
                        type="password"
                        placeholder="sk-..."
                        disabled={isLoading}
                        cols="col-span-6"
                    />
                    <p class="col-span-6 mt-1 text-sm text-gray-500 dark:text-gray-400">
                        Your OpenAI API key for GPT models. Leave empty to use system defaults.
                    </p>
                    
                    <!-- Ollama Endpoint -->
                    <FormField 
                        id="ollamaEndpoint"
                        label="Ollama Endpoint"
                        bind:value={apiSettings.ollamaEndpoint}
                        placeholder="http://localhost:11434"
                        disabled={isLoading}
                        cols="col-span-6 sm:col-span-4"
                    />
                    <p class="col-span-6 sm:col-span-4 mt-1 text-sm text-gray-500 dark:text-gray-400">
                        The URL of your Ollama instance (leave default for local installation)
                    </p>
                    
                    <!-- Default Model Dropdown -->
                    <div class="col-span-6 sm:col-span-3">
                        <label for="defaultModel" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                            Default Model
                        </label>
                        <select
                            id="defaultModel"
                            bind:value={apiSettings.defaultModel}
                            disabled={isLoading}
                            class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 dark:bg-gray-700 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
                        >
                            {#each modelOptions as option}
                                <option value={option.value}>{option.label}</option>
                            {/each}
                        </select>
                        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
                            The model that will be used by default for generating content
                        </p>
                    </div>
                </div>
            </div>
            
            <div class="flex justify-end mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
                <FormButton 
                    type="submit" 
                    {isLoading}
                    loadingText="Saving..."
                >
                    Save
                </FormButton>
            </div>
        </form>
    </div>
</div> 