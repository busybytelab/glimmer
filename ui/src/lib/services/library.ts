import pb from '$lib/pocketbase';
import type { PracticeTopicLibrary, PracticeSessionLibrary } from '$lib/types';

/**
 * Service for managing practice library collections
 * Handles community-shared topic and session templates
 */
class LibraryService {
    /**
     * Fetches practice topics from the library, sorted by usage
     * @param limit - Maximum number of topics to fetch (default: 50)
     * @returns Promise<PracticeTopicLibrary[]>
     */
    async getTopicsLibrary(limit: number = 50): Promise<PracticeTopicLibrary[]> {
        try {
            const result = await pb.collection('practice_topics_library').getList<PracticeTopicLibrary>(1, limit, {
                sort: '-total_usage,name,-last_used,-created',
                fields: 'id,name,description,category,country,target_age_range,target_grade_level,total_usage,last_used,created,updated'
            });
            return result.items;
        } catch (error: any) {
            console.error('Failed to fetch topics library:', error);
            throw new Error(error.message || 'Failed to fetch practice topics library');
        }
    }

    /**
     * Fetches the top N most used practice topics from the library
     * @param limit - Number of top topics to fetch (default: 3)
     * @returns Promise<PracticeTopicLibrary[]>
     */
    async getTopTopicsLibrary(limit: number = 3): Promise<PracticeTopicLibrary[]> {
        return this.getTopicsLibrary(limit);
    }

    /**
     * Fetches practice sessions from the library, optionally filtered by topic
     * @param topicLibraryId - Optional topic library ID to filter by
     * @param limit - Maximum number of sessions to fetch (default: 50)
     * @returns Promise<PracticeSessionLibrary[]>
     */
    async getSessionsLibrary(topicLibraryId?: string, limit: number = 50): Promise<PracticeSessionLibrary[]> {
        try {
            const options: any = {
                sort: '-total_usage,name,-last_used,-created',
                expand: 'practice_topic_library',
                fields: 'id,name,description,target_year,total_usage,last_used,practice_topic_library,created,updated,expand'
            };

            if (topicLibraryId) {
                options.filter = `practice_topic_library="${topicLibraryId}"`;
            }

            const result = await pb.collection('practice_sessions_library').getList<PracticeSessionLibrary>(1, limit, options);
            return result.items;
        } catch (error: any) {
            console.error('Failed to fetch sessions library:', error);
            throw new Error(error.message || 'Failed to fetch practice sessions library');
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
}

// Export singleton instance
export const libraryService = new LibraryService(); 