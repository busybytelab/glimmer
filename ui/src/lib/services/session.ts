import pb from '$lib/pocketbase';
import type { PracticeSession, PracticeItem, Learner } from '$lib/types';
import { authService } from './auth';

export interface SessionWithExpandedData extends PracticeSession {
    expand?: {
        learner?: Learner;
        practice_topic?: { id: string; name: string };
        practice_items?: PracticeItem[];
    };
}

class SessionService {
    private async ensureAuth(): Promise<void> {
        if (!pb.authStore.isValid) {
            throw new Error('You must be logged in to access practice sessions');
        }
    }

    async loadSession(id: string): Promise<SessionWithExpandedData> {
        try {
            await this.ensureAuth();

            const result = await pb.collection('practice_sessions').getOne(id, {
                expand: 'learner,practice_topic,practice_items',
                fields: 'id,name,status,assigned_at,completed_at,generation_prompt,learner,practice_topic,practice_items,expand'
            });

            if (!result) {
                throw new Error('Session not found');
            }

            return result as unknown as SessionWithExpandedData;
        } catch (err) {
            console.error('Failed to load session:', err);
            if (err instanceof Error) {
                throw err;
            }
            throw new Error('Failed to load practice session');
        }
    }

    async loadSessionForLearner(id: string): Promise<SessionWithExpandedData> {
        try {
            await this.ensureAuth();

            const result = await pb.collection('practice_sessions').getOne(id, {
                expand: 'learner,learner.user,practice_topic,practice_items',
                fields: 'id,name,status,assigned_at,completed_at,generation_prompt,learner,practice_topic,practice_items,expand'
            });

            if (!result) {
                throw new Error('Session not found');
            }

            // Filter practice items to only show approved or unreviewed items
            if (result.expand?.practice_items) {
                result.expand.practice_items = result.expand.practice_items.filter(
                    // TODO: we probably need to have a flag in topic or session to only allow approved items
                    item => !item.review_status || item.review_status === 'APPROVED'
                );
            }

            return result as unknown as SessionWithExpandedData;
        } catch (err) {
            console.error('Failed to load session for learner:', err);
            if (err instanceof Error) {
                throw err;
            }
            throw new Error('Failed to load practice session');
        }
    }

    parsePracticeItems(session: SessionWithExpandedData): PracticeItem[] {
        if (session.expand?.practice_items) {
            return session.expand.practice_items;
        }
        // If practice_items is not expanded, it's an array of IDs
        // We should never reach here because we always expand practice_items in loadSession
        throw new Error('Practice items not expanded. This is a data integrity error.');
    }

    /**
     * Updates a practice session
     * @param id Session ID
     * @param data Updated session data
     * @returns Updated session
     */
    async updateSession(id: string, data: Partial<PracticeSession>): Promise<PracticeSession> {
        await this.ensureAuth();
        return await pb.collection('practice_sessions').update(id, data) as PracticeSession;
    }

    /**
     * Deletes a practice session
     * @param id Session ID
     */
    async deleteSession(id: string): Promise<void> {
        await this.ensureAuth();
        await pb.collection('practice_sessions').delete(id);
    }

    /**
     * Gets a list of practice sessions
     * @param page Page number (1-based)
     * @param perPage Number of items per page
     * @param filter Optional filter string
     * @returns List of practice sessions
     */
    async getSessions(page: number = 1, perPage: number = 10, filter?: string): Promise<PracticeSession[]> {
        await this.ensureAuth();
        const options: any = {
            sort: '-created',
            expand: 'learner,practice_topic'
        };
        if (filter) {
            options.filter = filter;
        }
        const result = await pb.collection('practice_sessions').getList(page, perPage, options);
        return result.items as PracticeSession[];
    }

    /**
     * Gets practice sessions for a specific topic and learner
     * @param topicId Topic ID
     * @param learnerId Learner ID
     * @returns List of practice sessions
     */
    async getSessionsForTopicAndLearner(topicId: string, learnerId: string): Promise<PracticeSession[]> {
        await this.ensureAuth();
        try {
            const result = await pb.collection('practice_sessions').getList(1, 50, {
                filter: `practice_topic="${topicId}" && learner="${learnerId}"`,
                sort: '-created',
                expand: 'learner,practice_topic'
            });
            return result.items as PracticeSession[];
        } catch (error: any) {
            if (error.status === 401) {
                await authService.refreshAuthToken();
                try {
                    const result = await pb.collection('practice_sessions').getList(1, 50, {
                        filter: `practice_topic="${topicId}" && learner="${learnerId}"`,
                        sort: '-created',
                        expand: 'learner,practice_topic'
                    });
                    return result.items as PracticeSession[];
                } catch (retryError: any) {
                    throw new Error(retryError.message);
                }
            }
            throw new Error(error.message);
        }
    }
}

export const sessionService = new SessionService(); 