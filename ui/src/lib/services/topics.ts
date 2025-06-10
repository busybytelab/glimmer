import pb from '$lib/pocketbase';
import type { PracticeTopic, PracticeSession, TopicFormData } from '$lib/types';
import { accountService } from './accounts';

class TopicsService {
    // TODO: remove these auth methods, replace usage with auth
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

    async getTopics(): Promise<PracticeTopic[]> {
        return await pb.collection('practice_topics').getFullList({
            sort: '-created',
            expand: 'account'
        });
    }

    async getTopic(topicId: string): Promise<PracticeTopic> {
        await this.ensureAuth();

        try {
            const result = await pb.collection('practice_topics').getOne(topicId);
            result.tags = this.formatTags(result.tags);
            return result as unknown as PracticeTopic;
        } catch (error: any) {
            if (error.status === 401) {
                await this.refreshAuth();
                try {
                    const result = await pb.collection('practice_topics').getOne(topicId);
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

    async createTopic(formData: TopicFormData): Promise<PracticeTopic> {
        await this.ensureAuth();

        try {
            const account = await accountService.getAccount();
            const dataToSend = {
                ...formData,
                account: account.id
            };

            const result = await pb.collection('practice_topics').create(dataToSend);
            result.tags = this.formatTags(result.tags);
            return result as unknown as PracticeTopic;
        } catch (error: any) {
            if (error.status === 401) {
                await this.refreshAuth();
                try {
                    const account = await accountService.getAccount();
                    const dataToSend = {
                        ...formData,
                        account: account.id
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

    private isAgeInRange(learnerAge: number | undefined, targetAgeRange: string | undefined): boolean {
        if ( !targetAgeRange || !learnerAge) {
            return true;
        }
        const [minAge, maxAge] = targetAgeRange.split('-').map(age => parseInt(age.trim()));
        return learnerAge >= minAge && learnerAge <= maxAge;
    }

    private isGradeInRange(learnerGrade: string | undefined, targetGradeLevel: string | undefined): boolean {
        if (!targetGradeLevel || !learnerGrade) {
            return true;
        }
        
        // Extract numeric part from learner's grade (e.g., "8th" -> 8)
        const learnerGradeNum = parseInt(learnerGrade.replace(/[^0-9]/g, ''));
        
        // Handle both single grade and range formats
        if (targetGradeLevel.includes('-')) {
            const [minGrade, maxGrade] = targetGradeLevel.split('-').map(grade => parseInt(grade.trim()));
            return learnerGradeNum >= minGrade && learnerGradeNum <= maxGrade;
        } else {
            // Single grade level
            const targetGrade = parseInt(targetGradeLevel);
            return learnerGradeNum === targetGrade;
        }
    }

    async getTopicsForLearner(learnerAge: number | undefined, learnerGrade: string | undefined): Promise<PracticeTopic[]> {
        await this.ensureAuth();

        try {
            // First get all topics for the account
            const allTopics = await this.getTopics();
            
            // Filter topics based on learner's age and grade
            const filteredTopics = allTopics.filter(topic => {
                const ageMatch = topic.target_age_range ? this.isAgeInRange(learnerAge, topic.target_age_range) : true;
                const gradeMatch = topic.target_grade_level ? this.isGradeInRange(learnerGrade, topic.target_grade_level) : true;
                return ageMatch && gradeMatch;
            });

            // If no matching topics found, return all topics
            if (filteredTopics.length === 0) {
                console.log('No matching topics found for learner, returning all topics');
                return allTopics;
            }

            return filteredTopics;
        } catch (error: any) {
            if (error.status === 401) {
                await this.refreshAuth();
                try {
                    const allTopics = await this.getTopics();
                    const filteredTopics = allTopics.filter(topic => {
                        const ageMatch = topic.target_age_range ? this.isAgeInRange(learnerAge, topic.target_age_range) : true;
                        const gradeMatch = topic.target_grade_level ? this.isGradeInRange(learnerGrade, topic.target_grade_level) : true;
                        console.log(topic.name, 'ageMatch:', ageMatch, 'gradeMatch:', gradeMatch);
                        return ageMatch && gradeMatch;
                    });

                    // If no matching topics found, return all topics
                    if (filteredTopics.length === 0) {
                        console.log('No matching topics found for learner, returning all topics');
                        return allTopics;
                    }

                    return filteredTopics;
                } catch (retryError: any) {
                    throw new Error(retryError.message);
                }
            }
            throw new Error(error.message);
        }
    }
}

export const topicsService = new TopicsService(); 