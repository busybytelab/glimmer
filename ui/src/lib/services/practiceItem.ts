import pb from '$lib/pocketbase';
import type { PracticeItem } from '$lib/types';

class PracticeItemService {
    async getItem(id: string): Promise<PracticeItem> {
        try {
            const result = await pb.collection('practice_items').getOne(id, {
                expand: 'practice_topic,practice_session'
            });
            return result as PracticeItem;
        } catch (error: any) {
            throw new Error(error.message);
        }
    }

    async updateItem(id: string, data: Partial<PracticeItem>): Promise<PracticeItem> {
        try {
            const result = await pb.collection('practice_items').update(id, data);
            return result as PracticeItem;
        } catch (error: any) {
            throw new Error(error.message);
        }
    }
}

export const practiceItemService = new PracticeItemService(); 