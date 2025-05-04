declare module 'pocketbase' {
    export default class PocketBase {
        constructor(baseUrl: string);
        collection(collection: string): any;
        authStore: {
            model: any;
            token: string | null;
            isValid: boolean;
            clear(): void;
            record?: {
                id: string;
                email: string;
                name: string;
                created: string;
                updated: string;
            };
        };
        authWithPassword(email: string, password: string): Promise<any>;
        authRefresh(): Promise<any>;
    }

    export interface PracticeItem {
        id: string;
        question_text: string;
        question_type: string;
        options?: string[];
        correct_answer: string;
        explanation: string;
        explanation_for_incorrect?: Record<string, string>;
        hints?: string[];
        difficulty_level?: 'Easy' | 'Medium' | 'Hard';
        status: 'Generated' | 'NeedsReview' | 'Approved' | 'Rejected';
        tags?: string[];
        created: string;
        updated: string;
    }

    export interface PracticeTopic {
        id: string;
        name: string;
        subject: string;
        description?: string;
        target_age_range?: string;
        target_grade_level?: string;
        learning_goals?: string[];
        base_prompt: string;
        system_prompt?: string;
        tags?: string[];
        created: string;
        updated: string;
    }

    export interface Account {
        id: string;
        llm_api_key?: string;
        ollama_server_url?: string;
        default_llm_model?: string;
        default_language?: string;
        created: string;
        updated: string;
    }

    export interface User {
        id: string;
        email: string;
        name: string;
        created: string;
        updated: string;
    }

    export interface Instructor {
        id: string;
        nickname: string;
        account: string;
        user: User;
        created: string;
        updated: string;
    }

    export interface Learner {
        id: string;
        nickname: string;
        age: number;
        grade_level?: string;
        learning_preferences?: string[];
        avatar?: string;
        account: string;
        user: User;
        created: string;
        updated: string;
    }
} 