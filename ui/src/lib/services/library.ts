import pb from '$lib/pocketbase';
import type { PracticeTopicLibrary, PracticeSessionLibrary } from '$lib/types';
import { sessionImportExportService } from '$lib/services/sessionImportExport';
import type { PracticeSession } from '$lib/types';
import type { ExportedSession } from '$lib/types';

/**
 * Service for managing practice library collections
 * Handles community-shared topic and session templates
 */
class LibraryService {
    /**
     * Extracts the numeric value from a grade level string
     * @param gradeLevel - Grade level string (e.g., "5th", "Year 5", "5")
     * @returns number | undefined
     */
    private parseGradeLevel(gradeLevel?: string): number | undefined {
        if (!gradeLevel) return undefined;
        
        // Extract the first number from the string
        const match = gradeLevel.match(/\d+/);
        return match ? parseInt(match[0]) : undefined;
    }

    /**
     * Checks if a grade level matches a target grade level string
     * @param grade - The grade level to check
     * @param targetGrade - The target grade level string (can be a range like "2-7" or a single grade)
     * @returns boolean
     */
    private matchesGradeLevel(grade?: number, targetGrade?: string | number): boolean {
        if (!grade || !targetGrade) {
            return true;
        }
        
        const targetGradeStr = targetGrade.toString();
        
        // Handle range format (e.g., "2-7")
        if (targetGradeStr.includes('-')) {
            const [min, max] = targetGradeStr.split('-').map(g => parseInt(g.trim()));
            return grade >= min && grade <= max;
        }
        
        // Handle single grade
        return grade === parseInt(targetGradeStr.trim());
    }

    /**
     * Get top topics from the library
     * @param limit - Number of topics to return
     * @param gradeLevel - Optional grade level to filter by (e.g., "5th", "Year 5", "5")
     * @returns Promise<PracticeTopicLibrary[]>
     */
    async getTopTopicsLibrary(limit: number = 3, gradeLevel?: string): Promise<PracticeTopicLibrary[]> {
        try {
            const numericGrade = this.parseGradeLevel(gradeLevel);
            
            // Get all topics and filter on the client side
            const topics = await pb.collection('practice_topics_library').getList(1, limit * 3, {
                sort: '-total_usage'
            });

            // Filter topics by grade level if specified
            const filteredTopics = topics.items.filter(topic => 
                this.matchesGradeLevel(numericGrade, (topic as PracticeTopicLibrary).target_grade_level)
            );

            // Return only the requested number of topics
            return filteredTopics.slice(0, limit) as PracticeTopicLibrary[];
        } catch (err) {
            console.error('Failed to get top topics from library:', err);
            throw err;
        }
    }

    /**
     * Get all topics from the library
     * @param gradeLevel - Optional grade level to filter by (e.g., "5th", "Year 5", "5")
     * @returns Promise<PracticeTopicLibrary[]>
     */
    async getTopicsLibrary(gradeLevel?: string): Promise<PracticeTopicLibrary[]> {
        try {
            const numericGrade = this.parseGradeLevel(gradeLevel);
            
            // Get all topics
            const topics = await pb.collection('practice_topics_library').getFullList({
                sort: '-total_usage'
            });

            // Filter topics by grade level if specified
            return topics.filter(topic => 
                this.matchesGradeLevel(numericGrade, (topic as PracticeTopicLibrary).target_grade_level)
            ) as PracticeTopicLibrary[];
        } catch (err) {
            console.error('Failed to get topics from library:', err);
            throw err;
        }
    }

    /**
     * Get sessions from the library
     * @param topicId - Optional topic ID to filter by
     * @param gradeLevel - Optional grade level to filter by (e.g., "5th", "Year 5", "5")
     * @returns Promise<PracticeSessionLibrary[]>
     */
    async getSessionsLibrary(topicId?: string, gradeLevel?: string): Promise<PracticeSessionLibrary[]> {
        try {
            const numericGrade = this.parseGradeLevel(gradeLevel);
            let filter = '';
            
            // Only filter by topic ID in the database query
            if (topicId) {
                filter = `practice_topic_library = "${topicId}"`;
            }

            const sessions = await pb.collection('practice_sessions_library').getFullList({
                sort: '-total_usage',
                expand: 'practice_topic_library',
                filter
            });

            // Filter sessions by grade level on the client side
            return sessions.filter(session => 
                this.matchesGradeLevel(numericGrade, (session as PracticeSessionLibrary).target_year)
            ) as PracticeSessionLibrary[];
        } catch (err) {
            console.error('Failed to get sessions from library:', err);
            throw err;
        }
    }

    /**
     * Fetches a single practice topic from the library by ID
     * @param id - The topic library ID
     * @returns Promise<PracticeTopicLibrary>
     */
    async getTopicLibrary(id: string): Promise<PracticeTopicLibrary> {
        try {
            const result = await pb.collection('practice_topics_library').getOne<PracticeTopicLibrary>(id);
            return result;
        } catch (error: any) {
            console.error('Failed to fetch topic library:', error);
            throw new Error(error.message || 'Topic library not found');
        }
    }

    /**
     * Fetches a single practice session from the library by ID
     * @param id - The session library ID
     * @returns Promise<PracticeSessionLibrary>
     */
    async getSessionLibrary(id: string): Promise<PracticeSessionLibrary> {
        try {
            const result = await pb.collection('practice_sessions_library').getOne<PracticeSessionLibrary>(id, {
                expand: 'practice_topic_library,practice_items'
            });
            return result;
        } catch (error: any) {
            console.error('Failed to fetch session library:', error);
            throw new Error(error.message || 'Session library not found');
        }
    }

    /**
     * Searches topics library by name or description
     * @param query - Search query string
     * @param limit - Maximum number of results (default: 20)
     * @returns Promise<PracticeTopicLibrary[]>
     */
    async searchTopicsLibrary(query: string, limit: number = 20): Promise<PracticeTopicLibrary[]> {
        try {
            const result = await pb.collection('practice_topics_library').getList<PracticeTopicLibrary>(1, limit, {
                filter: `name ~ "${query}" || description ~ "${query}" || category ~ "${query}"`,
                sort: '-total_usage,name,-last_used,-created'
            });
            return result.items;
        } catch (error: any) {
            console.error('Failed to search topics library:', error);
            throw new Error(error.message || 'Failed to search practice topics library');
        }
    }

    /**
     * Searches sessions library by name or description
     * @param query - Search query string
     * @param topicLibraryId - Optional topic library ID to filter by
     * @param limit - Maximum number of results (default: 20)
     * @returns Promise<PracticeSessionLibrary[]>
     */
    async searchSessionsLibrary(query: string, topicLibraryId?: string, limit: number = 20): Promise<PracticeSessionLibrary[]> {
        try {
            let filter = `name ~ "${query}" || description ~ "${query}"`;
            if (topicLibraryId) {
                filter += ` && practice_topic_library="${topicLibraryId}"`;
            }

            const result = await pb.collection('practice_sessions_library').getList<PracticeSessionLibrary>(1, limit, {
                filter,
                sort: '-total_usage,name,-last_used,-created',
                expand: 'practice_topic_library'
            });
            return result.items;
        } catch (error: any) {
            console.error('Failed to search sessions library:', error);
            throw new Error(error.message || 'Failed to search practice sessions library');
        }
    }

    /**
     * Imports a practice session from the library to the user's account
     * @param librarySession - The library session to import
     * @param learnerId - The ID of the learner to assign the session to
     * @returns Promise<PracticeSession> - The newly created session
     */
    async importSessionFromLibrary(librarySession: PracticeSessionLibrary, learnerId: string): Promise<PracticeSession> {
        try {
            // First, get the complete session data with items
            const completeSession = await this.getSessionLibrary(librarySession.id);
            
            if (!completeSession.expand?.practice_topic_library || !completeSession.expand?.practice_items) {
                throw new Error('Failed to load complete session data');
            }

            // Create the export data structure for import
            const exportData: ExportedSession = {
                session: {
                    id: completeSession.id,
                    created: completeSession.created,
                    updated: completeSession.updated,
                    collectionId: completeSession.collectionId,
                    collectionName: completeSession.collectionName,
                    name: completeSession.name,
                    status: 'Imported',
                    assigned_at: new Date().toISOString(),
                    completed_at: undefined,
                    generation_prompt: completeSession.generation_prompt || '',
                    practice_topic: completeSession.practice_topic_library,
                    practice_items: JSON.stringify(completeSession.practice_items || []),
                    score: 0
                },
                topic: {
                    id: completeSession.expand.practice_topic_library.id,
                    name: completeSession.expand.practice_topic_library.name,
                    subject: completeSession.expand.practice_topic_library.category || 'General',
                    description: completeSession.expand.practice_topic_library.description,
                    target_age_range: completeSession.expand.practice_topic_library.target_age_range,
                    target_grade_level: completeSession.expand.practice_topic_library.target_grade_level,
                    learning_goals: [],
                    base_prompt: completeSession.expand.practice_topic_library.base_prompt,
                    system_prompt: completeSession.expand.practice_topic_library.system_prompt,
                    tags: [],
                    llm_model: undefined
                },
                practice_items: completeSession.expand.practice_items.map(item => ({
                    id: item.id,
                    created: item.created,
                    updated: item.updated,
                    collectionId: item.collectionId,
                    collectionName: item.collectionName,
                    question_text: item.question_text,
                    question_type: item.question_type,
                    options: item.options,
                    correct_answer: item.correct_answer,
                    explanation: item.explanation,
                    explanation_for_incorrect: item.explanation_for_incorrect,
                    hints: item.hints,
                    difficulty_level: item.difficulty_level,
                    status: 'Imported',
                    tags: item.tags,
                    practice_topic: item.practice_topic,
                    review_status: item.review_status,
                    review_date: item.review_date
                }))
            };

            // Use the sessionImportExportService to handle the actual import
            const importedSession = await sessionImportExportService.importPracticeSession(
                exportData,
                learnerId
            );

            return importedSession;
        } catch (error: any) {
            console.error('Failed to import session from library:', error);
            throw new Error(error.message || 'Failed to import session from library');
        }
    }
}

// Export singleton instance
export const libraryService = new LibraryService(); 