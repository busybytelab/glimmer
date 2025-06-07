import pb from '$lib/pocketbase';
import type { PracticeTopic, PracticeSession, TopicFormData } from '$lib/types';

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

    private async getUserAccountInfo(): Promise<{ instructor?: string; account?: string }> {
        const currentUser = pb.authStore.model;
        if (!currentUser) {
            throw new Error('You must be logged in to create a topic');
        }

        try {
            // First try to get instructor record for current user
            const instructors = await pb.collection('instructors').getList(1, 1, {
                filter: `user.id = "${currentUser.id}"`
            });
            
            if (instructors && instructors.items.length > 0) {
                return {
                    instructor: instructors.items[0].id,
                    account: instructors.items[0].account
                };
            }

            // If not an instructor, try to get account directly
            const accounts = await pb.collection('accounts').getList(1, 1, {
                filter: `owner.id = "${currentUser.id}"`
            });
            
            if (accounts && accounts.items.length > 0) {
                return { account: accounts.items[0].id };
            }

            throw new Error('Could not determine account for user');
        } catch (err) {
            console.error('Failed to get user account info:', err);
            throw new Error('Failed to get account information');
        }
    }

    async createTopic(formData: TopicFormData): Promise<PracticeTopic> {
        await this.ensureAuth();

        try {
            const accountInfo = await this.getUserAccountInfo();
            const dataToSend = {
                ...formData,
                ...accountInfo
            };

            const result = await pb.collection('practice_topics').create(dataToSend);
            result.tags = this.formatTags(result.tags);
            return result as unknown as PracticeTopic;
        } catch (error: any) {
            if (error.status === 401) {
                await this.refreshAuth();
                try {
                    const accountInfo = await this.getUserAccountInfo();
                    const dataToSend = {
                        ...formData,
                        ...accountInfo
                    };
                    const result = await pb.collection('practice_topics').create(dataToSend);
                    result.tags = this.formatTags(result.tags);
                    return result as unknown as PracticeTopic;
                } catch (retryError: any) {
                    throw new Error(retryError.message);
                }
            }
            throw new Error(error.message);
        }
    }

    async updateTopic(id: string, formData: TopicFormData): Promise<PracticeTopic> {
        await this.ensureAuth();

        try {
            const result = await pb.collection('practice_topics').update(id, formData);
            result.tags = this.formatTags(result.tags);
            return result as unknown as PracticeTopic;
        } catch (error: any) {
            if (error.status === 401) {
                await this.refreshAuth();
                try {
                    const result = await pb.collection('practice_topics').update(id, formData);
                    result.tags = this.formatTags(result.tags);
                    return result as unknown as PracticeTopic;
                } catch (retryError: any) {
                    throw new Error(retryError.message);
                }
            }
            throw new Error(error.message);
        }
    }

    async deleteTopic(id: string): Promise<void> {
        await this.ensureAuth();

        try {
            await pb.collection('practice_topics').delete(id);
        } catch (error: any) {
            if (error.status === 401) {
                await this.refreshAuth();
                try {
                    await pb.collection('practice_topics').delete(id);
                } catch (retryError: any) {
                    throw new Error(retryError.message);
                }
            }
            throw new Error(error.message);
        }
    }
}

export const topicsService = new TopicsService(); 