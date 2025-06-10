import pb from '$lib/pocketbase';
import type { Learner } from '$lib/types';

class LearnersService {
    /**
     * Gets a list of learners
     * @param page Page number (1-based)
     * @param perPage Number of items per page
     * @returns List of learners
     */
    async getLearners(page: number = 1, perPage: number = 50): Promise<Learner[]> {
        const result = await pb.collection('learners').getList(page, perPage);
        return result.items as Learner[];
    }

    /**
     * Gets a single learner by ID
     * @param id Learner ID
     * @returns Learner object
     */
    async getLearner(id: string): Promise<Learner> {
        return await pb.collection('learners').getOne(id) as Learner;
    }

    /**
     * Creates a new learner
     * @param data Learner data
     * @returns Created learner
     */
    async createLearner(data: Partial<Learner>): Promise<Learner> {
        return await pb.collection('learners').create(data) as Learner;
    }

    /**
     * Updates an existing learner
     * @param id Learner ID
     * @param data Updated learner data
     * @returns Updated learner
     */
    async updateLearner(id: string, data: Partial<Learner>): Promise<Learner> {
        return await pb.collection('learners').update(id, data) as Learner;
    }

    /**
     * Deletes a learner
     * @param id Learner ID
     */
    async deleteLearner(id: string): Promise<void> {
        await pb.collection('learners').delete(id);
    }

    /**
     * Gets the first learner matching the filter
     * @param filter Filter string
     * @returns First matching learner or null
     */
    async getFirstLearner(filter: string): Promise<Learner | null> {
        try {
            return await pb.collection('learners').getFirstListItem(filter) as Learner;
        } catch (error) {
            return null;
        }
    }
}

export const learnersService = new LearnersService(); 