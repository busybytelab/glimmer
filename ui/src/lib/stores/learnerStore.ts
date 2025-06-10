import { writable } from 'svelte/store';
import type { Learner } from '$lib/types';

// Store for current learner
export const currentLearnerStore = writable<{
    learner: Learner | null;
    loading: boolean;
    error: string | null;
}>({
    learner: null,
    loading: false,
    error: null
});

// Action to set the current learner
export function setCurrentLearner(learner: Learner | null) {
    currentLearnerStore.update(state => ({
        ...state,
        learner,
        error: null
    }));
}

// Action to clear the current learner
export function clearCurrentLearner() {
    currentLearnerStore.update(state => ({
        ...state,
        learner: null,
        error: null
    }));
} 