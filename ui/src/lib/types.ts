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
        learner?: Learner;
        practice_topic?: {
            id: string;
            name: string;
        };
        practice_items?: PracticeItem[];
    };
}

/**
 * Represents a practice session for export/import functionality
 * Excludes user-specific data (learner, account) to allow session reuse
 */
export interface ExportedSession {
    /** The practice session data without user-specific fields */
    session: Omit<PracticeSession, 'learner' | 'account' | 'expand'>;
    
    /** The associated practice topic data */
    topic: {
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
        llm_model?: string;
    };
    
    /** The practice items associated with this session */
    practice_items: Array<Omit<PracticeItem, 'account' | 'reviewer' | 'expand'>>;
}

/**
 * Statistics for a practice session
 * Generated from the practice_session_stats view
 * Matches all fields from the migration 1746341017_created_practice_session_stats.go
 */
export interface PracticeSessionStats extends PocketBaseRecord {
    /** session_id */
    id: string;
    /** Name of the practice session */
    session_name: string;
    /** Name of the associated topic */
    topic_name: string;
    /** Total number of items in the session */
    total_items: number;
    /** Number of items answered by the learner */
    answered_items: number;
    /** Number of wrong answers */
    wrong_answers_count: number;
    /** Total score for the session */
    total_score: number;
    /** Status of the session (e.g., completed, in-progress) */
    session_status: string;
    /** Timestamp of the last answer */
    last_answer_time: string;
    /** Learner ID associated with the session */
    learner_id: string;
    /** Account ID associated with the session */
    account: string;
    /** Number of approved items */
    approved_items: number;
    /** Number of edited items */
    edited_items: number;
    /** Number of not reviewed items */
    not_reviewed_items: number;
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

/**
 * Available icons for breadcrumb items
 */
export type BreadcrumbIcon = 'home' | 'topic' | 'session' | 'learner' | 'edit' | 'create';

/**
 * Map of breadcrumb icons to their SVG path data
 */
export const BreadcrumbIconMap: Record<BreadcrumbIcon, string> = {
    home: 'M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z',
    topic: 'M4 14h4v-4H4v4zm0 5h4v-4H4v4zM4 9h4V5H4v4zm5 5h12v-4H9v4zm0 5h12v-4H9v4zM9 5v4h12V5H9z',
    session: 'M14 10H2v2h12v-2zm0-4H2v2h12V6zm4 8v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zM2 16h8v-2H2v2z',
    learner: 'M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z',
    edit: 'M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34a.9959.9959 0 0 0-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z',
    create: 'M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z'
};

export interface BreadcrumbItem {
    label: string;
    href?: string;
    icon?: BreadcrumbIcon;
}

/**
 * Available icons for achievements
 */
export type AchievementIcon = 'math-whiz' | 'math-apprentice' | 'science-explorer' | 'reading-master' | 'writing-pro' | 'problem-solver' | 'quick-thinker' | 'perfect-score' | 'practice-streak' | 'helper' | 'explorer';

/**
 * Map of achievement icons to their SVG path data
 */
export const AchievementIconMap: Record<AchievementIcon, string> = {
    'math-whiz': 'M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-5 14h-4v-4H6v-2h4V7h4v4h4v2h-4v4z',
    'math-apprentice': 'M19 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2zM9 17H7v-7h2v7zm4 0h-2V7h2v10zm4 0h-2v-4h2v4z',
    'science-explorer': 'M9.5 3A6.5 6.5 0 003 9.5c0 1.61.59 3.09 1.56 4.23l.27.27A2 2 0 006.5 15H8a1 1 0 011 1v1a2 2 0 002 2h2a2 2 0 002-2v-1a1 1 0 011-1h1.5a2 2 0 001.67-1.23l.27-.27A6.5 6.5 0 1016 9.5a6.5 6.5 0 00-6.5-6.5z',
    'reading-master': 'M18 2H6c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zM6 4h5v8l-2.5-1.5L6 12V4z',
    'writing-pro': 'M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z',
    'problem-solver': 'M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z',
    'quick-thinker': 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z',
    'perfect-score': 'M12 17.27L18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z',
    'practice-streak': 'M12.65 10C11.83 7.67 9.61 6 7 6c-3.31 0-6 2.69-6 6 0 3.31 2.69 6 6 6 2.61 0 4.83-1.67 5.65-4H17v4h4v-4h2v-2H12.65zM7 16c-2.21 0-4-1.79-4-4s1.79-4 4-4 4 1.79 4 4-1.79 4-4 4z',
    'helper': 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.94-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z',
    'explorer': 'M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z'
};

/**
 * Interface for achievement data with progress tracking
 */
export interface Achievement {
    /** Title of the achievement */
    title: string;
    /** Description explaining how to earn the achievement */
    description: string;
    /** Icon identifier for the achievement */
    icon: AchievementIcon;
    /** Current progress towards completion */
    progress: number;
    /** Required progress to complete the achievement */
    requiredProgress: number;
    /** List of actions needed to complete the achievement */
    actions: string[];
}

/**
 * Interface for displaying latest achievement badge
 */
export interface LatestAchievement {
    /** Title of the achievement */
    title: string;
    /** Icon identifier for the achievement */
    icon: AchievementIcon;
    /** Description of what was achieved */
    description: string;
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