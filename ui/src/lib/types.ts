/**
 * Application Data Models and UI Types
 * 
 * This file contains all data model interfaces that represent the application's
 * domain entities and UI-specific type definitions.
 * 
 * USAGE GUIDELINES:
 * - Add all data model interfaces that represent database collections here
 * - Add UI-specific type definitions that are used across components
 * - All data models should extend PocketBaseRecord for consistency
 * - Keep types focused on data structure, not service behavior
 * - Document complex types with JSDoc comments
 * 
 * DO NOT ADD:
 * - Service interfaces (add to pocketbase-types.ts instead)
 * - PocketBase client types (add to pocketbase-types.ts instead)
 * - Component-specific types (keep those in the component files)
 */

// Base record type that matches PocketBase's structure
interface PocketBaseRecord {
    id: string;
    created: string;
    updated: string;
    collectionId: string;
    collectionName: string;
    expand?: Record<string, any>;
}

// -------------------------------------------------------------------------
// Enums and Constants
// -------------------------------------------------------------------------

/**
 * Question view types for controlling display and interaction modes
 * 
 * @description Determines how questions are displayed and interacted with based on user role and context
 */
export type QuestionViewType = 'learner' | 'answered' | 'parent' | 'hint' | 'correction' | 'generated';

// Constants for QuestionViewType values
export const QuestionViewType = {
    /** Learner actively answering questions */
    LEARNER: 'learner' as QuestionViewType,
    
    /** Learner viewing their answered questions (read-only) */
    ANSWERED: 'answered' as QuestionViewType,
    
    /** Parent viewing full details including correct answers and explanations */
    PARENT: 'parent' as QuestionViewType,
    
    /** Learner viewing question with hint */
    HINT: 'hint' as QuestionViewType,
    
    /** Learner correcting previously answered question */
    CORRECTION: 'correction' as QuestionViewType,
    
    /** Parent viewing newly generated questions (for review/approval) */
    GENERATED: 'generated' as QuestionViewType
};

/**
 * Question types supported by the application
 * 
 * @description These values must match the backend's question type identifiers
 */
export type QuestionType = 'multiple_choice' | 'true_false' | 'short_answer' | 'fill_in_blank';

// Constants for QuestionType values
export const QuestionType = {
    /** Multiple choice questions with radio button options */
    MULTIPLE_CHOICE: 'multiple_choice' as QuestionType,
    
    /** True/False questions */
    TRUE_FALSE: 'true_false' as QuestionType,
    
    /** Short answer questions with text input */
    SHORT_ANSWER: 'short_answer' as QuestionType,
    
    /** Fill-in-the-blank questions */
    FILL_IN_BLANK: 'fill_in_blank' as QuestionType
};

/**
 * Review status values for practice items
 * @values 'APPROVED' | 'IGNORE' | 'NEED_EDIT'
 */
export type ReviewStatus = 'APPROVED' | 'IGNORE' | 'NEED_EDIT';

// -------------------------------------------------------------------------
// Collection Models
// -------------------------------------------------------------------------

// Your collection types
export interface PracticeItem extends PocketBaseRecord {
    /** The question text to be presented to the learner */
    question_text: string;
    
    /** The type of question (e.g., multiple choice, true/false) */
    question_type: string;
    
    /** Available options for multiple choice questions */
    options?: Record<string, any>;
    
    /** The correct answer(s) for the question */
    correct_answer: string;
    
    /** Explanation of why the answer is correct */
    explanation: string;
    
    /** Specific explanations for incorrect answer choices */
    explanation_for_incorrect?: Record<string, string>;
    
    /** Optional hints to help learners */
    hints?: string[];
    
    /** The difficulty level of the question */
    difficulty_level?: string;
    
    /** Current status of the practice item */
    status: string;
    
    /** Associated tags for categorization */
    tags?: Record<string, any>;
    
    /** Reference to the practice topic this item belongs to */
    practice_topic: string;
    
    /** Reference to the account that owns this item */
    account: string;

    /** Reference to the instructor who reviewed this item */
    reviewer?: string;

    /** Date when the item was reviewed */
    review_date?: string;

    /** Current review status */
    review_status?: ReviewStatus;

    // Fields added at runtime (not in database)
    user_answer?: string;
    is_correct?: boolean;
    score?: number;
    feedback?: string;
    hint_level_reached?: number;
    attempt_number?: number;

    // Expand types for relations
    expand?: {
        practice_topic?: PracticeTopic;
        account?: Account;
    };
}

// Practice result interface for tracking user responses
export interface PracticeResult extends PocketBaseRecord {
    practice_item: string;
    practice_session: string;
    learner: string;
    answer: string;
    is_correct: boolean;
    score: number;
    feedback: string;
    hint_level_reached: number;
    attempt_number: number;
    started_at: string;
    submitted_at: string;
}

// -------------------------------------------------------------------------
// Other Collection Models
// -------------------------------------------------------------------------

export interface PracticeTopic extends PocketBaseRecord {
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

export interface PracticeSession extends PocketBaseRecord {
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
    expand?: {
        learner?: {
            id: string;
            expand?: {
                user?: User;
            };
        };
        practice_topic?: {
            id: string;
            name: string;
        };
        practice_items?: PracticeItem[];
    };
}

export interface Account extends PocketBaseRecord {
    owner: User;
    llm_api_key?: string;
    ollama_server_url?: string;
    default_llm_model?: string;
    default_language?: string;
}

export interface User extends PocketBaseRecord {
    email: string;
    name: string;
}

export interface Learner extends PocketBaseRecord {
    age: number;
    grade_level?: string;
    learning_preferences?: string[];
    avatar?: string;
    account: string;
    nickname: string;
    expand?: {
        account?: Account;
    };
}

// -------------------------------------------------------------------------
// UI-specific Types
// -------------------------------------------------------------------------

export interface BreadcrumbItem {
    label: string;
    href?: string;
    icon?: string;
}

export interface RegistrationForm {
	email: string;
	password: string;
	passwordConfirm: string;
}

export type TopicFormData = {
    name: string;
    subject: string;
    description: string;
    target_age_range: string;
    target_grade_level: string;
    learning_goals: string[];
    base_prompt: string;
    system_prompt: string;
    tags: string[];
    instructor?: string;
    account?: string;
    llm_model?: string;
}; 