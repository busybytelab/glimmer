import PocketBase, { RecordService } from 'pocketbase';

// Base types for system fields that PocketBase adds to all records
interface BaseSystemFields {
    id: string;
    created: string;
    updated: string;
}

// Your collection types
export interface PracticeItem {
    id: string;
    question_text: string;
    question_type: string;
    options?: string;
    correct_answer: string;
    explanation: string;
    explanation_for_incorrect?: string;
    hints?: string;
    difficulty_level?: string;
    status: string;
    tags?: string[];
    practice_topic: string;
    account: string;
    created: string;
    updated: string;
    user_answer?: string;
    hints_used?: string;
}

export interface PracticeTopic extends BaseSystemFields {
    name: string;
    subject: string;
    description?: string;
    target_age_range?: string;
    target_grade_level?: string;
    learning_goals?: string[];
    base_prompt: string;
    system_prompt?: string;
    tags?: string[];
    llm_model?: string;
    difficulty_level?: string;
}

export interface PracticeSession {
    id: string;
    name?: string;
    status: string;
    assigned_at: string;
    completed_at?: string;
    generation_prompt?: string;
    learner: string;
    practice_topic: string;
    account: string;
    practice_items: string;
    score?: number;
    created: string;
    updated: string;
    expand?: {
        learner?: {
            id: string;
            name: string;
        };
        practice_topic?: {
            id: string;
            name: string;
        };
        practice_items?: PracticeItem[];
    };
}

export interface Account extends BaseSystemFields {
    owner: User;
    llm_api_key?: string;
    ollama_server_url?: string;
    default_llm_model?: string;
    default_language?: string;
}

export interface User extends BaseSystemFields {
    email: string;
    name: string;
}

export interface Instructor extends BaseSystemFields {
    nickname: string;
    account: string;
    user: User;
}

export interface Learner extends BaseSystemFields {
    nickname: string;
    age: number;
    grade_level?: string;
    learning_preferences?: string[];
    avatar?: string;
    account: string;
    user: User;
    expand?: {
        user?: User;
        account?: Account;
    };
}

// Type for the PocketBase client with your collections
export interface TypedPocketBase extends PocketBase {
    collection(idOrName: string): RecordService; // fallback for any other collection
    collection(idOrName: 'practice_items'): RecordService<PracticeItem>;
    collection(idOrName: 'practice_topics'): RecordService<PracticeTopic>;
    collection(idOrName: 'practice_sessions'): RecordService<PracticeSession>;
    collection(idOrName: 'accounts'): RecordService<Account>;
    collection(idOrName: 'users'): RecordService<User>;
    collection(idOrName: 'instructors'): RecordService<Instructor>;
    collection(idOrName: 'learners'): RecordService<Learner>;
} 