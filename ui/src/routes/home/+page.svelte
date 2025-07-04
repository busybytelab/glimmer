<script lang="ts">
    import { goto } from '$app/navigation';
    import ErrorAlert from '$components/common/ErrorAlert.svelte';
    import type { Learner } from '$lib/types';
    import { setCurrentLearner } from '$lib/stores/learnerStore';

    export let data: {
        learners?: Learner[];
        error?: string | null;
    };

    $: ({ learners = [], error = null } = data || {});

    function handleLearnerSelect(learner: Learner) {
        // Set the current learner in the store
        setCurrentLearner(learner);
        // Navigate to the learner's home page
        goto(`/learners/${learner.id}/home`);
    }

    function handleAccountDashboard() {
        // Clear the current learner when going to account settings
        setCurrentLearner(null);
        goto('/account/dashboard');
    }

    function handleAddChild() {
        // Navigate to create learner page
        goto('/account/learners/create');
    }
</script>

<div class="min-h-screen bg-gray-100 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-7xl mx-auto">
        <div class="text-center mb-12">
            <p class="text-lg text-gray-600 dark:text-gray-300">Choose your profile to start learning</p>
        </div>

        {#if error}
            <div class="max-w-md mx-auto">
                <ErrorAlert message={error} />
            </div>
        {:else}
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <!-- Learner Cards -->
                {#each learners as learner}
                    <button
                        on:click={() => handleLearnerSelect(learner)}
                        class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow duration-200 flex flex-col items-center text-left"
                    >
                        <div class="w-24 h-24 rounded-full bg-gray-200 dark:bg-gray-700 mb-4 flex items-center justify-center">
                            {#if learner.avatar}
                                <img
                                    src={learner.avatar}
                                    alt={`${learner.nickname}'s avatar`}
                                    class="w-full h-full rounded-full object-cover"
                                />
                            {:else}
                                <span class="text-4xl text-gray-500 dark:text-gray-400">
                                    {learner.nickname.charAt(0).toUpperCase()}
                                </span>
                            {/if}
                        </div>
                        <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                            {learner.nickname}
                        </h3>
                        <p class="text-gray-600 dark:text-gray-300 text-sm">
                            Continue as {learner.nickname}
                        </p>
                    </button>
                {/each}

                <!-- Add Child Card -->
                <button
                    on:click={handleAddChild}
                    class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow duration-200 flex flex-col items-center text-left"
                >
                    <div class="w-24 h-24 rounded-full bg-gray-200 dark:bg-gray-700 mb-4 flex items-center justify-center">
                        <svg
                            class="w-12 h-12 text-gray-500 dark:text-gray-400"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M12 4v16m8-8H4"
                            />
                        </svg>
                    </div>
                    <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                        Add Child
                    </h3>
                    <p class="text-gray-600 dark:text-gray-300 text-sm">
                        Set up a profile for your child
                    </p>
                </button>

                <!-- Account Dashboard Card -->
                <button
                    on:click={handleAccountDashboard}
                    class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow duration-200 flex flex-col items-center text-left"
                >
                    <div class="w-24 h-24 rounded-full bg-gray-200 dark:bg-gray-700 mb-4 flex items-center justify-center">
                        <svg
                            class="w-12 h-12 text-gray-500 dark:text-gray-400"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
                            />
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                            />
                        </svg>
                    </div>
                    <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                        Account
                    </h3>
                    <p class="text-gray-600 dark:text-gray-300 text-sm">
                        Manage your account and children profiles
                    </p>
                </button>
            </div>
        {/if}
    </div>
</div> 