import pb from '$lib/pocketbase';
import type { PracticeResult } from '$lib/types';

class ResultsService {
    /**
     * Fetches existing practice results for a session and learner
     * @param sessionId - The ID of the practice session
     * @param learnerId - The ID of the learner
     * @returns List of practice results
     */
    async getResults(sessionId: string, learnerId: string): Promise<PracticeResult[]> {
        try {
            const results = await pb.collection('practice_results').getList(1, 100, {
                filter: `practice_session = "${sessionId}" && learner = "${learnerId}"`,
                expand: 'practice_item,learner',
                sort: '-created'
            });

            return results.items;
        } catch (err) {
            console.error('Failed to fetch practice results:', err);
            if (err instanceof Error) {
                throw err;
            }
            throw new Error('Failed to fetch practice results');
        }
    }

    /**
     * Updates an existing practice result
     * @param resultId - The ID of the result to update
     * @param data - The data to update
     */
    async updateResult(resultId: string, data: Partial<PracticeResult>): Promise<void> {
        try {
            await pb.collection('practice_results').update(resultId, data);
        } catch (err) {
            console.error('Failed to update practice result:', err);
            if (err instanceof Error) {
                throw err;
            }
            throw new Error('Failed to update practice result');
        }
    }

    /**
     * Creates a new practice result
     * @param data - The result data to create
     */
    async createResult(data: {
        practice_item: string;
        practice_session: string;
        learner: string;
        answer: string;
        is_correct: boolean;
        started_at: string;
        submitted_at: string;
        attempt_number: number;
        score?: number;
        feedback?: string;
        hint_level_reached?: number;
    }): Promise<PracticeResult> {
        try {
            // return the created result
            return await pb.collection('practice_results').create(data);
        } catch (err) {
            console.error('Failed to create practice result:', err);
            if (err instanceof Error) {
                throw err;
            }
            throw new Error('Failed to create practice result');
        }
    }

    /**
     * Gets the latest result for a practice item in a session
     * @param practiceItemId - The ID of the practice item
     * @param sessionId - The ID of the practice session
     * @returns The latest result or null if none exists
     */
    async getLatestResult(practiceItemId: string, sessionId: string): Promise<PracticeResult | null> {
        try {
            const results = await pb.collection('practice_results').getList(1, 1, {
                filter: `practice_item = "${practiceItemId}" && practice_session = "${sessionId}"`,
                sort: '-created'
            });

            return results.items[0] || null;
        } catch (err) {
            console.error('Failed to fetch latest practice result:', err);
            if (err instanceof Error) {
                throw err;
            }
            throw new Error('Failed to fetch latest practice result');
        }
    }
}

export const resultsService = new ResultsService(); 