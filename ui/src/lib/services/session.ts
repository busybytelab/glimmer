import pb from '$lib/pocketbase';
import type { PracticeSession, PracticeItem, Learner, PracticeSessionStats, PracticeTopic, LearnerProgress } from '$lib/types';
import { PracticeSessionStatusChecker } from '$lib/types';

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
                expand: 'learner,practice_topic,practice_items',
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
            
            // Split items into active and completed using the status checker
            return {
                active: items.filter(item => 
                    PracticeSessionStatusChecker.isInProgress(item) || PracticeSessionStatusChecker.needsAttention(item)
                ),
                completed: items.filter(item => 
                    PracticeSessionStatusChecker.isRecentlyCompleted(item) && !PracticeSessionStatusChecker.needsAttention(item)
                )
            };
        } catch (error: any) {
            throw new Error(error.message);
        }
    }

    /**
     * Processes a list of session stats to generate a learner progress report.
     * This is a private helper to be reused by both single-learner and all-learner progress fetches.
     * @param items - A list of practice session stats.
     * @returns A LearnerProgress object.
     */
    private _processStatsForProgress(items: PracticeSessionStats[]): LearnerProgress {
        // Group sessions by topic to analyze performance
        const topicPerformance = new Map<string, {
            totalSessions: number;
            wrongAnswers: number;
            totalItems: number;
        }>();

        items.forEach(item => {
            const existing = topicPerformance.get(item.topic_name) || {
                totalSessions: 0,
                wrongAnswers: 0,
                totalItems: 0
            };
            
            topicPerformance.set(item.topic_name, {
                totalSessions: existing.totalSessions + 1,
                wrongAnswers: existing.wrongAnswers + item.wrong_answers_count,
                totalItems: existing.totalItems + item.total_items
            });
        });

        // Identify topics needing help and doing well
        const needsHelpWith: string[] = [];
        const doingWellIn: string[] = [];
        
        topicPerformance.forEach((perf, topic) => {
            const errorRate = perf.totalItems > 0 ? perf.wrongAnswers / perf.totalItems : 0;
            if (errorRate > 0.3 && perf.totalSessions >= 2) {
                needsHelpWith.push(topic);
            } else if (errorRate < 0.1 && perf.totalSessions >= 2) {
                doingWellIn.push(topic);
            }
        });

        // Calculate overall statistics
        const totalSessions = items.length;
        const completedSessions = items.filter(item => PracticeSessionStatusChecker.isRecentlyCompleted(item)).length;
        
        const totalScore = items.reduce((sum, item) => sum + item.total_score, 0);
        const averageScore = totalSessions > 0 ? Math.round(totalScore / totalSessions) : 0;

        return {
            needsAttention: items.filter(item => PracticeSessionStatusChecker.needsAttention(item)).slice(0, 3),
            inProgress: items.filter(item => PracticeSessionStatusChecker.isInProgress(item)).slice(0, 3),
            recentlyCompleted: items.filter(item => PracticeSessionStatusChecker.isRecentlyCompleted(item)).slice(0, 3),
            overallProgress: {
                totalSessions,
                completedSessions,
                averageScore,
                needsHelpWith,
                doingWellIn
            }
        };
    }

    /**
     * Get learner progress in a parent-friendly format for a single learner.
     * @param learnerId The ID of the learner.
     * @returns Parent-friendly progress information.
     */
    async getLearnerProgressForParent(learnerId: string): Promise<LearnerProgress> {
        try {
            const items = await pb.collection('pbc_practice_session_stats').getFullList<PracticeSessionStats>({
                filter: `learner_id="${learnerId}"`,
                sort: '-last_answer_time'
            });

            return this._processStatsForProgress(items);
        } catch (error: any) {
            console.error(`Failed to get progress for learner ${learnerId}:`, error);
            throw new Error(error.message);
        }
    }

    /**
     * Get progress for all learners
     * @returns A record mapping each learner's ID to their progress report.
     */
    async getLearnersProgressForAccount(): Promise<Record<string, LearnerProgress>> {
        if (!pb.authStore.isValid || !pb.authStore.model) {
            throw new Error('User is not authenticated.');
        }

        try {            
            const allStats = await pb.collection('pbc_practice_session_stats').getFullList<PracticeSessionStats>({
                sort: '-last_answer_time'
            });

            // Group stats by learner ID
            const statsByLearner = allStats.reduce((acc, stat) => {
                const learnerId = stat.learner_id;
                if (!acc[learnerId]) {
                    acc[learnerId] = [];
                }
                acc[learnerId].push(stat);
                return acc;
            }, {} as Record<string, PracticeSessionStats[]>);
            
            // Process stats for each learner
            const allLearnerProgress: Record<string, LearnerProgress> = {};
            for (const learnerId in statsByLearner) {
                allLearnerProgress[learnerId] = this._processStatsForProgress(statsByLearner[learnerId]);
            }

            return allLearnerProgress;
        } catch (error: any) {
            console.error('Failed to get progress for all learners in account:', error);
            throw new Error(error.message);
        }
    }

    /**
     * Get detailed stats for a specific session
     * @param sessionId The ID of the session
     * @returns Session statistics including completion status and scores
     */
    async getSessionStats(sessionId: string): Promise<PracticeSessionStats | null> {
        try {
            const result = await pb.collection('pbc_practice_session_stats').getFirstListItem(
                `id="${sessionId}"`,
                { sort: '-last_answer_time' }
            ) as PracticeSessionStats;
            return result || null;
        } catch (error: any) {
            console.error('Failed to get session stats:', error);
            throw new Error('Failed to get session statistics');
        }
    }
}

export const sessionService = new SessionService(); 