import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { topicsService } from '$lib/services/topics';

export const load: PageServerLoad = async ({ params }) => {
    try {
        // Verify that the learner ID is provided
        const learnerId = params.id;
        if (!learnerId) {
            throw error(400, 'Invalid learner ID');
        }

        // Load practice topics
        const topics = await topicsService.getTopic(learnerId);

        return {
            topics
        };
    } catch (err) {
        console.error('Error in learner home page load:', err);
        throw error(500, 'Failed to load practice topics');
    }
}; 