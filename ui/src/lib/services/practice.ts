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
            throw new Error(errorData.message || `Server error: ${response.status}`);
        }

        return response.json();
    }

    async createSession(request: CreatePracticeSessionRequest): Promise<PracticeSession> {
        try {
            return await this.makeRequest(request);
        } catch (error: any) {
            throw new Error(error.message);
        }
    }
}

export const practiceService = new PracticeService(); 