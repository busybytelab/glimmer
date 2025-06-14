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
    import LearnerCard from '$components/learners/LearnerCard.svelte';

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

    // Get the selected learner
    $: selectedLearner = learners.find(l => l.id === selectedLearnerId);
    // Get the selected topic
    $: selectedTopic = topics.find(t => t.id === selectedTopicId);

    onMount(async () => {
        try {
            learners = await learnersService.getLearners();
            // Automatically select the child if there's only one
            if (learners.length === 1) {
                selectedLearnerId = learners[0].id;
            }
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
    <div class="flex items-center gap-3 mb-8">
        <svg class="w-8 h-8 text-blue-600 dark:text-blue-400" viewBox="0 0 24 24">
            <path fill="currentColor" d="M14 10H2v2h12v-2zm0-4H2v2h12V6zm4 8v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zM2 16h8v-2H2v2z" />
        </svg>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Import Practice Session</h1>
    </div>

    <div class="bg-blue-50 dark:bg-blue-900/40 p-4 rounded-lg mb-6">
        <h2 class="text-lg font-semibold text-blue-700 dark:text-blue-300 mb-2">About Importing Practice Sessions</h2>
        <div class="text-blue-600 dark:text-blue-200 space-y-2">
            <p>Importing a practice session allows you to:</p>
            <ul class="list-disc list-inside ml-4 space-y-1">
                <li>Reuse practice materials shared by other parents or teachers</li>
                <li>Add pre-made practice sessions to your child's study plan</li>
                <li>Save time by using carefully crafted questions and topics</li>
            </ul>
            <p class="mt-4">The imported session will include:</p>
            <ul class="list-disc list-inside ml-4 space-y-1">
                <li>A complete set of practice questions</li>
                <li>The topic information and learning goals</li>
                <li>Correct answers and helpful explanations</li>
            </ul>
        </div>
    </div>

    <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
        <div class="space-y-6">
            <div>
                <label for="file" class="block text-lg font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Step 1: Choose Your Session File
                </label>
                <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
                    Select the practice session file (JSON format) that was shared with you.
                </p>
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
                    <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
                        Step 2: Choose How to Handle the Practice Topic
                    </h3>
                    <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
                        You can either create a new topic from the imported content or use an existing topic from your collection.
                    </p>
                    
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
                                Create a new topic (Recommended for new content)
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
                                Use an existing topic (Good for adding to existing materials)
                                {#if checkingSimilarTopics}
                                    <span class="text-gray-500 dark:text-gray-400">(Finding similar topics...)</span>
                                {:else if similarTopic}
                                    <span class="text-green-600 dark:text-green-400">(We found a matching topic!)</span>
                                {/if}
                            </label>
                        </div>

                        {#if !createNewTopic}
                            <div class="pl-6">
                                {#if loadingTopics}
                                    <LoadingSpinner size="md" color="primary" />
                                {:else if topics.length === 0}
                                    <p class="text-sm text-gray-600 dark:text-gray-400">You don't have any topics yet. We recommend creating a new topic instead.</p>
                                {:else}
                                    <div class="relative">
                                        <select
                                            id="topic"
                                            bind:value={selectedTopicId}
                                            class="
                                                block w-full rounded-md
                                                px-4 py-2.5
                                                text-base
                                                bg-white dark:bg-gray-800
                                                border border-gray-300 dark:border-gray-600
                                                focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
                                                appearance-none
                                                relative
                                                transition-colors duration-200
                                            "
                                            aria-label="Select a practice topic"
                                        >
                                            <option value="" disabled selected>Select a topic</option>
                                            {#each topics as topic}
                                                <option value={topic.id}>
                                                    {topic.name}
                                                    {#if similarTopic && topic.id === similarTopic.id}
                                                        (Best Match)
                                                    {/if}
                                                </option>
                                            {/each}
                                        </select>
                                        <div class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2">
                                            <svg class="h-5 w-5 text-gray-400 dark:text-gray-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                                                <path fill-rule="evenodd" d="M10 3a.75.75 0 01.55.24l3.25 3.5a.75.75 0 11-1.1 1.02L10 4.852 7.3 7.76a.75.75 0 01-1.1-1.02l3.25-3.5A.75.75 0 0110 3zm-3.76 9.2a.75.75 0 011.06.04l2.7 2.908 2.7-2.908a.75.75 0 111.1 1.02l-3.25 3.5a.75.75 0 01-1.1 0l-3.25-3.5a.75.75 0 01.04-1.06z" clip-rule="evenodd" />
                                            </svg>
                                        </div>
                                    </div>

                                    {#if selectedTopic}
                                        <div class="mt-4">
                                            <PracticeTopicCard 
                                                topic={selectedTopic}
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
                    <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
                        Step 3: Select Your Child
                    </h3>
                    {#if learners.length === 0}
                        <div class="text-sm text-red-600 dark:text-red-400 mb-3">
                            No children found in your account. Please add a child first to continue.
                        </div>
                    {:else if learners.length === 1}
                        <div class="flex items-center gap-2 bg-blue-50 dark:bg-blue-900/40 p-3 rounded-md mb-3">
                            <svg class="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0" viewBox="0 0 24 24">
                                <path fill="currentColor" d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z" />
                            </svg>
                            <p class="text-sm text-blue-700 dark:text-blue-300">
                                This practice session will be set up for {learners[0].nickname}
                            </p>
                        </div>
                    {:else}
                        <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
                            Choose which child will practice with these materials.
                        </p>
                    {/if}
                    <select
                        id="learner"
                        bind:value={selectedLearnerId}
                        class="
                            block w-full rounded-md
                            px-4 py-2.5
                            text-base
                            bg-white dark:bg-gray-800
                            border border-gray-300 dark:border-gray-600
                            focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
                            disabled:bg-gray-50 disabled:dark:bg-gray-800/50
                            disabled:cursor-not-allowed
                            disabled:border-gray-200 dark:disabled:border-gray-700
                            appearance-none
                            relative
                            transition-colors duration-200
                            {learners.length === 1 ? 'opacity-75' : ''}
                        "
                        disabled={learners.length === 1}
                        aria-label={learners.length === 1 ? `Practice session will be set up for ${learners[0].nickname}` : 'Select a child'}
                    >
                        <option value="" disabled selected>Select a child</option>
                        {#each learners as learner}
                            <option value={learner.id}>{learner.nickname}</option>
                        {/each}
                    </select>
                    <div class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2">
                        <svg class="h-5 w-5 text-gray-400 dark:text-gray-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                            <path fill-rule="evenodd" d="M10 3a.75.75 0 01.55.24l3.25 3.5a.75.75 0 11-1.1 1.02L10 4.852 7.3 7.76a.75.75 0 01-1.1-1.02l3.25-3.5A.75.75 0 0110 3zm-3.76 9.2a.75.75 0 011.06.04l2.7 2.908 2.7-2.908a.75.75 0 111.1 1.02l-3.25 3.5a.75.75 0 01-1.1 0l-3.25-3.5a.75.75 0 01.04-1.06z" clip-rule="evenodd" />
                        </svg>
                    </div>

                    {#if selectedLearner}
                        <div class="mt-4">
                            <LearnerCard 
                                learner={selectedLearner}
                                shadow={false}
                                showPreferences={true}
                            />
                        </div>
                    {/if}
                </div>

                <FormButton 
                    type="button"
                    on:click={handleImport} 
                    disabled={!selectedFile || !selectedLearnerId || (!createNewTopic && !selectedTopicId) || isLoading}
                    isLoading={isLoading}
                    loadingText="Setting up your practice session..."
                    variant="primary"
                >
                    Import and Start Practice Session
                </FormButton>
            {/if}
        </div>
    </div>
</div> 