import pb from '$lib/pocketbase';
import type { PracticeTopic, PracticeSession, TopicFormData, PracticeTopicLibrary } from '$lib/types';
import { accountService } from './accounts';

export class TopicsService {
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
        const response = await pb.collection('practice_topics').getList<PracticeTopic>(1, 50, {
            sort: '-created',
        });
        return response.items;
    }

    async getTopic(topicId: string): Promise<PracticeTopic> {
        try {
            const result = await pb.collection('practice_topics').getOne(topicId);
            result.tags = this.formatTags(result.tags);
            return result as unknown as PracticeTopic;
        } catch (error: any) {
            throw new Error(error.message);
        }
    }

    async getPastPractices(topicId: string): Promise<PracticeSession[]> {
        try {
            const result = await pb.collection('practice_sessions').getList(1, 10, {
                filter: `practice_topic="${topicId}"`,
                sort: '-created',
                expand: 'learner,practice_topic'
            });
            
            return result.items as unknown as PracticeSession[];
        } catch (error: any) {
            throw new Error(error.message);
        }
    }

    async createTopic(formData: TopicFormData): Promise<PracticeTopic> {
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
            throw new Error(error.message);
        }
    }

    async updateTopic(id: string, formData: TopicFormData): Promise<PracticeTopic> {
        try {
            const result = await pb.collection('practice_topics').update(id, formData);
            result.tags = this.formatTags(result.tags);
            return result as unknown as PracticeTopic;
        } catch (error: any) {
            throw new Error(error.message);
        }
    }

    async deleteTopic(id: string): Promise<void> {
        try {
            await pb.collection('practice_topics').delete(id);
        } catch (error: any) {
            throw new Error(error.message);
        }
    }

    /**
     * Imports a practice topic from the library to the user's account
     * Creates a new topic based on library template data
     * @param libraryTopic - The library topic to import
     * @param customizations - Optional customizations to apply during import
     * @returns Promise<PracticeTopic> - The newly created topic
     */
    async importFromLibrary(
        libraryTopic: PracticeTopicLibrary, 
        customizations: Partial<TopicFormData> = {}
    ): Promise<PracticeTopic> {
        try {
            const account = await accountService.getAccount();
            
            // Check if similar topic already exists
            const existingSimilarTopic = await this.findSimilarTopic(libraryTopic.name);
            if (existingSimilarTopic) {
                throw new Error(`A similar topic "${existingSimilarTopic.name}" already exists in your account. Please choose a different name or modify the existing topic.`);
            }

            // Validate that required fields are present
            if (!libraryTopic.base_prompt || libraryTopic.base_prompt.trim() === '') {
                throw new Error(`The library topic "${libraryTopic.name}" is missing a required field and cannot be imported. Please contact the library maintainer.`);
            }

            // Map library topic data to account topic format
            const topicData: TopicFormData = {
                name: customizations.name || libraryTopic.name,
                subject: customizations.subject || libraryTopic.category || 'General',
                description: customizations.description || libraryTopic.description || '',
                target_age_range: customizations.target_age_range || libraryTopic.target_age_range || '',
                target_grade_level: customizations.target_grade_level || libraryTopic.target_grade_level || '',
                learning_goals: customizations.learning_goals || [],
                base_prompt: customizations.base_prompt || libraryTopic.base_prompt,
                system_prompt: customizations.system_prompt || libraryTopic.system_prompt || '',
                tags: customizations.tags || [], // Library topics don't have tags field
                llm_model: customizations.llm_model || '',
                account: account.id
            };

            // Additional validation before creating
            if (!topicData.base_prompt || topicData.base_prompt.trim() === '') {
                throw new Error('Base prompt is required but was not provided in the library topic or customizations.');
            }

            // Create the new topic
            const result = await pb.collection('practice_topics').create(topicData);
            result.tags = this.formatTags(result.tags);
            
            console.log('Successfully imported topic from library:', result.name);
            return result as unknown as PracticeTopic;
        } catch (error: any) {
            console.error('Failed to import topic from library:', error);
            throw new Error(error.message || 'Failed to import topic from library');
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
            throw new Error(error.message);
        }
    }

    /**
     * Finds topics with similar names to the given topic name
     * Returns the best matching topic if similarity is above threshold
     */
    async findSimilarTopic(topicName: string): Promise<PracticeTopic | null> {
        try {
            // Get all topics
            const topics = await this.getTopics();
            
            // Find the best match using string similarity
            let bestMatch: { topic: PracticeTopic; similarity: number } | null = null;
            
            for (const topic of topics) {
                // Calculate similarity score (basic implementation)
                const similarity = this.calculateSimilarity(
                    this.normalizeString(topic.name),
                    this.normalizeString(topicName)
                );
                
                if (similarity > 0.8 && (!bestMatch || similarity > bestMatch.similarity)) {
                    bestMatch = { topic, similarity };
                }
            }
            
            return bestMatch?.topic || null;
        } catch (error) {
            console.error('Error finding similar topics:', error);
            return null;
        }
    }

    /**
     * Calculates similarity between two strings
     * Returns a score between 0 and 1, where 1 is exact match
     */
    private calculateSimilarity(str1: string, str2: string): number {
        const len1 = str1.length;
        const len2 = str2.length;
        
        // If either string is empty, return 0
        if (len1 === 0 || len2 === 0) return 0;
        
        // If strings are identical, return 1
        if (str1 === str2) return 1;
        
        // Calculate Levenshtein distance
        const matrix: number[][] = Array(len1 + 1).fill(null).map(() => Array(len2 + 1).fill(0));
        
        for (let i = 0; i <= len1; i++) matrix[i][0] = i;
        for (let j = 0; j <= len2; j++) matrix[0][j] = j;
        
        for (let i = 1; i <= len1; i++) {
            for (let j = 1; j <= len2; j++) {
                const cost = str1[i - 1] === str2[j - 1] ? 0 : 1;
                matrix[i][j] = Math.min(
                    matrix[i - 1][j] + 1,      // deletion
                    matrix[i][j - 1] + 1,      // insertion
                    matrix[i - 1][j - 1] + cost // substitution
                );
            }
        }
        
        // Convert distance to similarity score
        const maxLen = Math.max(len1, len2);
        const distance = matrix[len1][len2];
        return 1 - (distance / maxLen);
    }

    /**
     * Normalizes a string for comparison by:
     * - Converting to lowercase
     * - Removing special characters
     * - Removing extra whitespace
     */
    private normalizeString(str: string): string {
        return str
            .toLowerCase()
            .replace(/[^a-z0-9\s]/g, '')
            .replace(/\s+/g, ' ')
            .trim();
    }
}

export const topicsService = new TopicsService(); 