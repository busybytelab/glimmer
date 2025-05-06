import { getAuthToken } from '$lib/auth';

export type ModelInfo = {
    name: string;
    sizeHuman?: string;
    isDefault: boolean;
};

export type PlatformInfo = {
    name: string;
    isDefault: boolean;
    models: ModelInfo[];
};

export type LLMInfo = {
    platforms: PlatformInfo[];
};

export type Usage = {
    LlmModelName?: string;
    CacheHit?: boolean;
    Cost?: number;
    PromptTokens: number;
    CompletionTokens: number;
    TotalTokens: number;
};

export type ChatResponse = {
    response: string;
    usage?: Usage;
};

class LLMService {
    private static instance: LLMService;
    private baseUrl = '/api/llm';

    private constructor() {}

    public static getInstance(): LLMService {
        if (!LLMService.instance) {
            LLMService.instance = new LLMService();
        }
        return LLMService.instance;
    }

    /**
     * Fetches available LLM models and platform information
     */
    public async getInfo(): Promise<LLMInfo> {
        const authToken = getAuthToken();
        if (!authToken) {
            throw new Error('Please log in again.');
        }

        const response = await fetch(`${this.baseUrl}/info`, {
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });
        if (!response.ok) {
            throw new Error(`Error: ${response.status} ${response.statusText}`);
        }
        return response.json();
    }

    /**
     * Sends a chat message to the LLM service
     */
    public async chat(prompt: string, systemPrompt: string, model?: string): Promise<ChatResponse> {
        const authToken = getAuthToken();
        if (!authToken) {
            throw new Error('Please log in again.');
        }

        const requestBody: Record<string, string> = {
            prompt,
            systemPrompt
        };

        if (model) {
            requestBody.model = model;
        }

        const response = await fetch(`${this.baseUrl}/chat`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            },
            body: JSON.stringify(requestBody)
        });

        if (!response.ok) {
            throw new Error(`Error: ${response.status} ${response.statusText}`);
        }

        return response.json();
    }
}

// Export a singleton instance
export const llmService = LLMService.getInstance(); 