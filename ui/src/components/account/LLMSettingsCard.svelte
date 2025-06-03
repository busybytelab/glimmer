<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { error } from '$lib/stores'; 
    import FormField from '../common/FormField.svelte';
    import FormButton from '../common/FormButton.svelte';
    import ErrorAlert from '../common/ErrorAlert.svelte';
    import { toast } from '$lib/stores/toast';
    
    // This would typically come from an API or store
    let apiSettings = {
        openaiApiKey: '',
        ollamaEndpoint: 'http://localhost:11434',
        defaultModel: 'gpt-3.5-turbo'
    };
    
    export let isLoading = false;
    
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
            
            toast.success('Your LLM settings have been updated successfully!');
            isLoading = false;
            
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
    
    {#if $error}
        <div class="px-6 py-4">
            <ErrorAlert message={$error} />
        </div>
    {/if}
    
    <div class="px-6 py-4">
        <form on:submit|preventDefault={handleSubmit} class="space-y-4">
            <FormField
                id="openaiApiKey"
                label="OpenAI API Key"
                type="password"
                bind:value={apiSettings.openaiApiKey}
                placeholder="Enter your OpenAI API key"
                disabled={isLoading}
            />
            
            <FormField
                id="ollamaEndpoint"
                label="Ollama Endpoint"
                type="text"
                bind:value={apiSettings.ollamaEndpoint}
                placeholder="Enter Ollama server URL"
                disabled={isLoading}
            />
            
            <div class="mb-4">
                <label for="defaultModel" class="block text-sm font-medium text-gray-700 mb-1 text-left">
                    Default Model
                </label>
                <select
                    id="defaultModel"
                    bind:value={apiSettings.defaultModel}
                    disabled={isLoading}
                    class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md text-sm shadow-sm
                    focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500
                    dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600
                    disabled:bg-gray-50 disabled:text-gray-500 disabled:border-gray-200 disabled:shadow-none"
                >
                    {#each modelOptions as option}
                        <option value={option.value}>{option.label}</option>
                    {/each}
                </select>
            </div>
            
            <div class="flex justify-end">
                <FormButton
                    type="submit"
                    isLoading={isLoading}
                    loadingText="Saving..."
                >
                    Save Settings
                </FormButton>
            </div>
        </form>
    </div>
</div> 