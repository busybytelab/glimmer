import pb from '$lib/pocketbase';
import type { PracticeSession, PracticeItem, User } from '$lib/types';

export interface SessionWithExpandedData extends PracticeSession {
    expand?: {
        learner?: { 
            id: string; 
            expand?: {
                user?: User;
            }
        };
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
                expand: 'learner,learner.user,practice_topic,practice_items',
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

    async checkUserRole(): Promise<boolean> {
        await this.ensureAuth();

        try {
            const authData = pb.authStore.model;
            if (!authData) {
                throw new Error('User not authenticated');
            }

            try {
                const instructorRecord = await pb.collection('instructors').getFirstListItem(`user="${authData.id}"`);
                return !!instructorRecord;
            } catch (err) {
                return false;
            }
        } catch (err) {
            console.error('Failed to check user role:', err);
            return false;
        }
    }

    parsePracticeItems(session: SessionWithExpandedData): PracticeItem[] {
        if (session.expand?.practice_items) {
            return session.expand.practice_items;
        } else if (session.practice_items) {
            try {
                return JSON.parse(session.practice_items as string) as PracticeItem[];
            } catch (e) {
                console.error('Error parsing practice items:', e);
                return [];
            }
        }
        return [];
    }
}

export const sessionService = new SessionService(); 