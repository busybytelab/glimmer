import pb from '$lib/pocketbase';

export interface CreatePracticeSessionRequest {
    learnerId: string;
    practiceTopicId: string;
    systemPrompt?: string;
    basePrompt?: string;
}

export interface PracticeSession {
    id: string;
    name: string;
    status: string;
    practice_topic: string;
    learner: string;
    practice_items: string;
    assigned_at: string;
    created: string;
    updated: string;
}

class PracticeService {
    private async ensureAuth(): Promise<void> {
        if (!pb.authStore.isValid) {
            throw new Error('You must be logged in to create a practice session');
        }
    }

    private async refreshAuth(): Promise<void> {
        try {
            await pb.collection('users').authRefresh();
        } catch (err) {
            // If refresh fails, clear the auth store and throw
            pb.authStore.clear();
            throw new Error('Your session has expired. Please log in again.');
        }
    }

    private async makeRequest(request: CreatePracticeSessionRequest): Promise<PracticeSession> {
        const response = await fetch('/api/glimmer/v1/practice/session', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': pb.authStore.token
            },
            body: JSON.stringify(request)
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw {
                status: response.status,
                message: errorData.message || `Server error: ${response.status}`
            };
        }

        return response.json();
    }

    async createSession(request: CreatePracticeSessionRequest): Promise<PracticeSession> {
        await this.ensureAuth();

        try {
            return await this.makeRequest(request);
        } catch (error: any) {
            // If unauthorized, try to refresh the token once
            if (error.status === 401) {
                await this.refreshAuth();
                try {
                    return await this.makeRequest(request);
                } catch (retryError: any) {
                    throw new Error(retryError.message);
                }
            }
            throw new Error(error.message);
        }
    }
}

export const practiceService = new PracticeService(); 