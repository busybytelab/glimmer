import pb from '$lib/pocketbase';
import type { PracticeTopic, PracticeSession } from '$lib/types';

class TopicsService {
    private async ensureAuth(): Promise<void> {
        if (!pb.authStore.isValid) {
            throw new Error('You must be logged in to access practice topics');
        }
    }

    private async refreshAuth(): Promise<void> {
        try {
            await pb.collection('users').authRefresh();
        } catch (err) {
            pb.authStore.clear();
            throw new Error('Your session has expired. Please log in again.');
        }
    }

    private formatTags(tags: any): string[] {
        if (!tags) return [];
        
        if (Array.isArray(tags)) return tags;
        
        try {
            if (typeof tags === 'string' && tags.trim().startsWith('[')) {
                return JSON.parse(tags);
            } else if (typeof tags === 'string') {
                return tags.split(',').map((tag: string) => tag.trim()).filter(Boolean);
            }
        } catch (e) {
            console.error('Error parsing tags:', e);
        }
        
        return [];
    }

    async getTopic(id: string): Promise<PracticeTopic> {
        await this.ensureAuth();

        try {
            const result = await pb.collection('practice_topics').getOne(id);
            result.tags = this.formatTags(result.tags);
            return result as unknown as PracticeTopic;
        } catch (error: any) {
            if (error.status === 401) {
                await this.refreshAuth();
                try {
                    const result = await pb.collection('practice_topics').getOne(id);
                    result.tags = this.formatTags(result.tags);
                    return result as unknown as PracticeTopic;
                } catch (retryError: any) {
                    throw new Error(retryError.message);
                }
            }
            throw new Error(error.message);
        }
    }

    async getPastPractices(topicId: string): Promise<PracticeSession[]> {
        await this.ensureAuth();

        try {
            const result = await pb.collection('practice_sessions').getList(1, 10, {
                filter: `practice_topic="${topicId}"`,
                sort: '-created',
                expand: 'learner,practice_topic'
            });
            
            return result.items as unknown as PracticeSession[];
        } catch (error: any) {
            if (error.status === 401) {
                await this.refreshAuth();
                try {
                    const result = await pb.collection('practice_sessions').getList(1, 10, {
                        filter: `practice_topic="${topicId}"`,
                        sort: '-created',
                        expand: 'learner,practice_topic'
                    });
                    
                    return result.items as unknown as PracticeSession[];
                } catch (retryError: any) {
                    throw new Error(retryError.message);
                }
            }
            throw new Error(error.message);
        }
    }
}

export const topicsService = new TopicsService(); 