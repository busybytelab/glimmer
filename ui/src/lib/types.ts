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
 * Icon types supported by the application
 * Used for both breadcrumbs and action buttons
 */
export type IconType = 'home' | 'topic' | 'session' | 'learner' | 'edit' | 'create' | 
    'print' | 'delete' | 'view' | 'download' | 'share' | 'duplicate' | 'add' | 
    'start' | 'complete' | 'reset' | 'back' | 'next' | 'more' | 'user' |
    'practice' | 'progress' | 'answers' | 'review' | 'hint' | 'ignore';

/**
 * Map of icon types to their SVG path data
 */
export const IconTypeMap: Record<IconType, string> = {
    // Breadcrumb icons
    home: 'M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z',
    topic: 'M4 14h4v-4H4v4zm0 5h4v-4H4v4zM4 9h4V5H4v4zm5 5h12v-4H9v4zm0 5h12v-4H9v4zM9 5v4h12V5H9z',
    session: 'M14 10H2v2h12v-2zm0-4H2v2h12V6zm4 8v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zM2 16h8v-2H2v2z',
    learner: 'M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z',
    edit: 'M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10',
    create: 'M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z',
    
    // Action toolbar icons
    print: 'M6.72 13.829c-.24.03-.48.062-.72.096m.72-.096a42.415 42.415 0 0110.56 0m-10.56 0L6.34 18m10.94-4.171c.24.03.48.062.72.096m-.72-.096L17.66 18m0 0l.229 2.523a1.125 1.125 0 01-1.12 1.227H7.231c-.662 0-1.18-.568-1.12-1.227L6.34 18m11.318 0h1.091A2.25 2.25 0 0021 15.75V9.456c0-1.081-.768-2.015-1.837-2.175a48.055 48.055 0 00-1.913-.247M6.34 18H5.25A2.25 2.25 0 013 15.75V9.456c0-1.081.768-2.015 1.837-2.175a48.041 48.041 0 011.913-.247m10.5 0a48.536 48.536 0 00-10.5 0m10.5 0V3.375c0-.621-.504-1.125-1.125-1.125h-8.25c-.621 0-1.125.504-1.125 1.125v3.659M18 10.5h.008v.008H18V10.5zm-3 0h.008v.008H15V10.5z',
    delete: 'M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0',
    view: 'M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z M15 12a3 3 0 11-6 0 3 3 0 016 0z',
    download: 'M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3',
    share: 'M7.217 10.907a2.25 2.25 0 100 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186l9.566-5.314m-9.566 7.5l9.566 5.314m0 0a2.25 2.25 0 103.935 2.186 2.25 2.25 0 00-3.935-2.186zm0-12.814a2.25 2.25 0 103.933-2.185 2.25 2.25 0 00-3.933 2.185z',
    duplicate: 'M15.75 17.25v3.375c0 .621-.504 1.125-1.125 1.125h-9.75a1.125 1.125 0 01-1.125-1.125V7.875c0-.621.504-1.125 1.125-1.125H6.75a9.06 9.06 0 011.5.124m7.5 10.376h3.375c.621 0 1.125-.504 1.125-1.125V11.25c0-4.46-3.243-8.161-7.5-8.876a9.06 9.06 0 00-1.5-.124H9.375c-.621 0-1.125.504-1.125 1.125v3.5m7.5 10.375H9.375a1.125 1.125 0 01-1.125-1.125v-9.25m12 6.625v-1.875a3.375 3.375 0 00-3.375-3.375h-1.5a1.125 1.125 0 01-1.125-1.125v-1.5a3.375 3.375 0 00-3.375-3.375H9.75',
    add: 'M12 4.5v15m7.5-7.5h-15',
    start: 'M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.348a1.125 1.125 0 010 1.971l-11.54 6.347a1.125 1.125 0 01-1.667-.985V5.653z',
    complete: 'M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    reset: 'M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99',
    back: 'M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18',
    next: 'M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3',
    more: 'M12 6.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 12.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 18.75a.75.75 0 110-1.5.75.75 0 010 1.5z',
    user: 'M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z',
    practice: 'M4.26 10.147a60.436 60.436 0 00-.491 6.347A48.627 48.627 0 0112 20.904a48.627 48.627 0 018.232-4.41 60.46 60.46 0 00-.491-6.347m-15.482 0a50.57 50.57 0 00-2.658-.813A59.905 59.905 0 0112 3.493a59.902 59.902 0 0110.399 5.84c-.896.248-1.783.52-2.658.814m-15.482 0A50.697 50.697 0 0112 13.489a50.702 50.702 0 017.74-3.342M6.75 15a.75.75 0 100-1.5.75.75 0 000 1.5zm0 0v-3.675A55.378 55.378 0 0112 8.443m-7.007 11.55A5.981 5.981 0 006.75 15.75v-1.5',
    progress: 'M7.5 14.25v2.25m3-4.5v4.5m3-6.75v6.75m3-9v9M6 20.25h12A2.25 2.25 0 0020.25 18V6A2.25 2.25 0 0018 3.75H6A2.25 2.25 0 003.75 6v12A2.25 2.25 0 006 20.25z',
    answers: 'M9 12.75L11.25 15 15 9.75M21 12c0 1.268-.63 2.39-1.593 3.068a3.745 3.745 0 01-1.043 3.296 3.745 3.745 0 01-3.296 1.043A3.745 3.745 0 0112 21c-1.268 0-2.39-.63-3.068-1.593a3.746 3.746 0 01-3.296-1.043 3.745 3.745 0 01-1.043-3.296A3.745 3.745 0 013 12c0-1.268.63-2.39 1.593-3.068a3.745 3.745 0 011.043-3.296 3.746 3.746 0 013.296-1.043A3.746 3.746 0 0112 3c1.268 0 2.39.63 3.068 1.593a3.746 3.746 0 013.296 1.043 3.746 3.746 0 011.043 3.296A3.745 3.745 0 0121 12z',
    review: 'M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L6.832 19.82a4.5 4.5 0 01-1.897 1.13l-2.685.8.8-2.685a4.5 4.5 0 011.13-1.897L16.863 4.487zm0 0L19.5 7.125',
    hint: 'M9 9a3 3 0 0 1 3-3m-2 15h4m0-3c0-4.1 4-4.9 4-9A6 6 0 1 0 6 9c0 4 4 5 4 9h4Z',
    // Eye with slash icon for ignore
    ignore: 'M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88'
};

// Remove old BreadcrumbIcon type and map
// export type BreadcrumbIcon = 'home' | 'topic' | 'session' | 'learner' | 'edit' | 'create';
// export const BreadcrumbIconMap: Record<BreadcrumbIcon, string> = { ... };

// Update BreadcrumbItem to use IconType
export interface BreadcrumbItem {
    label: string;
    href?: string;
    icon?: IconType;
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

/**
 * Interface representing overall progress metrics for a learner
 */
export interface LearnerOverallProgress {
    /** Total number of practice sessions */
    totalSessions: number;
    /** Number of successfully completed sessions */
    completedSessions: number;
    /** Average score across all sessions */
    averageScore: number;
    /** Topics where the learner needs additional help */
    needsHelpWith: string[];
    /** Topics where the learner is performing well */
    doingWellIn: string[];
}

/**
 * Interface representing a learner's progress across all practice sessions
 * Used by getLearnerProgressForParent to provide a comprehensive view
 */
export interface LearnerProgress {
    /** Sessions that have wrong answers and need review */
    needsAttention: PracticeSessionStats[];
    /** Active sessions that are not completed */
    inProgress: PracticeSessionStats[];
    /** Recently completed sessions with no wrong answers */
    recentlyCompleted: PracticeSessionStats[];
    /** Overall progress metrics */
    overallProgress: LearnerOverallProgress;
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