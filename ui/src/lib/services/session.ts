import pb from '$lib/pocketbase';
import type { PracticeSession, PracticeItem, Learner, PracticeSessionStats, PracticeTopic } from '$lib/types';

export interface SessionWithExpandedData extends PracticeSession {
    expand?: {
        learner?: Learner;
        practice_topic?: PracticeTopic;
        practice_items?: PracticeItem[];
    };
}

class SessionService {
    async loadSession(id: string): Promise<SessionWithExpandedData> {
        try {
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

    async updateSession(id: string, data: Partial<PracticeSession>): Promise<PracticeSession> {
        return await pb.collection('practice_sessions').update(id, data) as PracticeSession;
    }

    async deleteSession(id: string): Promise<void> {
        await pb.collection('practice_sessions').delete(id);
    }

    async getSessions(page: number = 1, perPage: number = 10, filter?: string): Promise<PracticeSession[]> {
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

    async getSessionsForTopicAndLearner(topicId: string, learnerId: string): Promise<PracticeSession[]> {
        try {
            const result = await pb.collection('practice_sessions').getList(1, 50, {
                filter: `practice_topic="${topicId}" && learner="${learnerId}"`,
                sort: '-created',
                expand: 'learner,practice_topic'
            });
            return result.items as PracticeSession[];
        } catch (error: any) {
            throw new Error(error.message);
        }
    }

    async getSessionStatsForLearner(learnerId: string): Promise<{
        active: PracticeSessionStats[];
        completed: PracticeSessionStats[];
    }> {
        try {
            const result = await pb.collection('pbc_practice_session_stats').getList(1, 100, {
                filter: `learner_id="${learnerId}"`,
                sort: '-last_answer_time'
            });

            const items = result.items as PracticeSessionStats[];
            
            // Split items into active and completed
            return {
                active: items.filter(item => 
                    item.answered_items < item.total_items || item.wrong_answers_count > 0
                ),
                completed: items.filter(item => 
                    item.answered_items === item.total_items && item.wrong_answers_count === 0
                )
            };
        } catch (error: any) {
            throw new Error(error.message);
        }
    }
}

export const sessionService = new SessionService(); 