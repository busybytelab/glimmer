<script lang="ts">
    import { sessionImportExportService } from '$lib/services/sessionImportExport';
    import { goto } from '$app/navigation';
    import { toast } from '$lib/stores/toast';
    import FormButton from '$components/common/FormButton.svelte';
    import type { ExportedSession, Learner, PracticeTopic } from '$lib/types';
    import { learnersService } from '$lib/services/learners';
    import { topicsService } from '$lib/services/topics';
    import { onMount } from 'svelte';
    import PracticeTopicCard from '$components/practice-topics/PracticeTopicCard.svelte';
    import LoadingSpinner from '$components/common/LoadingSpinner.svelte';

    let selectedFile: File | null = null;
    let selectedLearnerId: string = '';
    let learners: Learner[] = [];
    let topics: PracticeTopic[] = [];
    let isLoading = false;
    let loadingTopics = false;
    let importData: ExportedSession | null = null;
    let createNewTopic = true;
    let selectedTopicId: string = '';
    let similarTopic: PracticeTopic | null = null;
    let checkingSimilarTopics = false;

    onMount(async () => {
        try {
            learners = await learnersService.getLearners();
            loadingTopics = true;
            topics = await topicsService.getTopics();
        } catch (error) {
            console.error('Failed to load initial data:', error);
            toast.error('Failed to load initial data');
        } finally {
            loadingTopics = false;
        }
    });

    async function checkForSimilarTopics(topicName: string) {
        checkingSimilarTopics = true;
        try {
            similarTopic = await topicsService.findSimilarTopic(topicName);
            if (similarTopic) {
                // Automatically select the similar topic
                createNewTopic = false;
                selectedTopicId = similarTopic.id;
            }
        } catch (error) {
            console.error('Failed to find similar topics:', error);
        } finally {
            checkingSimilarTopics = false;
        }
    }

    async function handleFileSelect(event: Event) {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files[0]) {
            selectedFile = input.files[0];
            try {
                const fileContent = await input.files[0].text();
                importData = JSON.parse(fileContent) as ExportedSession;
                
                // Reset topic selection
                createNewTopic = true;
                selectedTopicId = '';
                similarTopic = null;
                
                // Check for similar topics
                if (importData.topic.name) {
                    await checkForSimilarTopics(importData.topic.name);
                }
            } catch (error) {
                console.error('Failed to parse file:', error);
                toast.error('Failed to parse file. Please make sure it is a valid exported session.');
                selectedFile = null;
                importData = null;
            }
        }
    }

    async function handleImport() {
        if (!selectedFile || !selectedLearnerId || (!createNewTopic && !selectedTopicId)) {
            toast.error('Please fill in all required fields');
            return;
        }

        if (!importData) {
            toast.error('Invalid import data');
            return;
        }

        isLoading = true;
        try {
            const session = await sessionImportExportService.importPracticeSession(
                importData, 
                selectedLearnerId,
                createNewTopic ? null : selectedTopicId
            );
            toast.success('Session imported successfully');
            
            // Navigate to the imported session
            goto('/practice-sessions/' + session.id + '/instructor');
        } catch (error) {
            console.error('Failed to import session:', error);
            toast.error('Failed to import session');
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
    <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Import Practice Session</h1>

    <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
        <div class="space-y-6">
            <div>
                <label for="file" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Session File (JSON)
                </label>
                <input
                    id="file"
                    type="file"
                    accept=".json"
                    on:change={handleFileSelect}
                    class="block w-full text-sm text-gray-900 dark:text-gray-300 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 dark:file:bg-blue-900/40 dark:file:text-blue-300"
                />
            </div>

            {#if importData}
                <div class="border-t border-gray-200 dark:border-gray-700 pt-6">
                    <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Practice Topic</h3>
                    
                    <div class="space-y-4">
                        <div class="flex items-center space-x-2">
                            <input
                                type="radio"
                                id="createNew"
                                bind:group={createNewTopic}
                                value={true}
                                class="text-blue-600 dark:text-blue-400"
                            />
                            <label for="createNew" class="text-sm text-gray-700 dark:text-gray-300">
                                Create new topic from import data
                            </label>
                        </div>

                        {#if createNewTopic && importData.topic}
                            <div class="pl-6">
                                <PracticeTopicCard topic={{
                                    ...importData.topic,
                                    created: new Date().toISOString(),
                                    updated: new Date().toISOString(),
                                    collectionId: '',
                                    collectionName: 'practice_topics'
                                }} showEditButton={false} />
                            </div>
                        {/if}

                        <div class="flex items-center space-x-2">
                            <input
                                type="radio"
                                id="useExisting"
                                bind:group={createNewTopic}
                                value={false}
                                class="text-blue-600 dark:text-blue-400"
                            />
                            <label for="useExisting" class="text-sm text-gray-700 dark:text-gray-300">
                                Use existing topic
                                {#if checkingSimilarTopics}
                                    <span class="text-gray-500 dark:text-gray-400">(checking for similar topics...)</span>
                                {:else if similarTopic}
                                    <span class="text-green-600 dark:text-green-400">(best match found)</span>
                                {/if}
                            </label>
                        </div>

                        {#if !createNewTopic}
                            <div class="pl-6">
                                {#if loadingTopics}
                                    <LoadingSpinner size="md" color="primary" />
                                {:else if topics.length === 0}
                                    <p class="text-sm text-gray-600 dark:text-gray-400">No topics available. Please create a topic first.</p>
                                {:else}
                                    <select
                                        id="topic"
                                        bind:value={selectedTopicId}
                                        class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md dark:bg-gray-700 dark:text-white"
                                    >
                                        <option value="">Select a topic</option>
                                        {#each topics as topic}
                                            <option value={topic.id}>
                                                {topic.name}
                                                {#if similarTopic && topic.id === similarTopic.id}
                                                    (Best Match)
                                                {/if}
                                            </option>
                                        {/each}
                                    </select>

                                    {#if selectedTopicId}
                                        <div class="mt-4">
                                            <PracticeTopicCard 
                                                topic={topics.find(t => t.id === selectedTopicId) || topics[0]} 
                                                showEditButton={false}
                                            />
                                        </div>
                                    {/if}
                                {/if}
                            </div>
                        {/if}
                    </div>
                </div>

                <div class="border-t border-gray-200 dark:border-gray-700 pt-6">
                    <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Select Learner</h3>
                    <select
                        id="learner"
                        bind:value={selectedLearnerId}
                        class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md dark:bg-gray-700 dark:text-white"
                    >
                        <option value="">Select a learner</option>
                        {#each learners as learner}
                            <option value={learner.id}>{learner.nickname}</option>
                        {/each}
                    </select>
                </div>

                <FormButton 
                    type="button"
                    on:click={handleImport} 
                    disabled={!selectedFile || !selectedLearnerId || (!createNewTopic && !selectedTopicId) || isLoading}
                    isLoading={isLoading}
                    loadingText="Importing..."
                    variant="primary"
                >
                    Import Session
                </FormButton>
            {/if}
        </div>
    </div>
</div> 