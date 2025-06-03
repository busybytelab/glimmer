import pb from '$lib/pocketbase';
import type { PracticeItem } from '$lib/types';

class PracticeItemService {
    private async ensureAuth(): Promise<void> {
        if (!pb.authStore.isValid) {
            throw new Error('You must be logged in to perform this action');
        }
    }

    /**
     * Updates a practice item with new data
     * @param id The ID of the practice item to update
     * @param data The data to update the practice item with
     * @returns The updated practice item
     */
    async updatePracticeItem(id: string, data: Partial<PracticeItem>): Promise<PracticeItem> {
        try {
            await this.ensureAuth();

            // First, get the instructor record for the current user
            const instructor = await pb.collection('instructors').getFirstListItem(`user = "${pb.authStore.model?.id}"`);
            
            if (!instructor) {
                throw new Error('Instructor record not found');
            }

            // Update the practice item
            const updatedItem = await pb.collection('practice_items').update(id, {
                ...data,
                reviewer: instructor.id,
                review_date: new Date().toISOString()
            });

            return updatedItem as PracticeItem;
        } catch (err) {
            console.error('Failed to update practice item:', err);
            if (err instanceof Error) {
                throw err;
            }
            throw new Error('Failed to update practice item');
        }
    }
}

export const practiceItemService = new PracticeItemService(); 