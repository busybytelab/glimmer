/**
 * PocketBase Type Definitions
 * 
 * This file contains type definitions for PocketBase client operations.
 * It provides strongly typed interfaces for database operations and authentication.
 * 
 * USAGE GUIDELINES:
 * - Import these types when creating services or utilities that work with PocketBase
 * - Use RecordOperations<T> for components that only need data operations
 * - Use AuthOperations for components that only need auth operations
 * - Use RecordService<T> for components that need both
 * - Extend these interfaces if you need to add custom methods
 * 
 * DO NOT ADD:
 * - Data model interfaces (add to types.ts instead)
 * - UI component types (add to types.ts or component files)
 * - Service implementation logic (create separate service files)
 */

import PocketBase from 'pocketbase';
import type {
    PracticeItem,
    PracticeResult,
    PracticeTopic,
    PracticeSession,
    Account,
    User,
    Instructor,
    Learner
} from './types';

/**
 * Core record operations interface for data manipulation
 */
export interface RecordOperations<T = any> {
    // Core data methods
    getFullList<R = T>(options?: any): Promise<R[]>;
    getList<R = T>(page?: number, perPage?: number, options?: any): Promise<{ page: number; perPage: number; totalItems: number; totalPages: number; items: R[] }>;
    getOne<R = T>(id: string, options?: any): Promise<R>;
    getFirstListItem<R = T>(filter: string, options?: any): Promise<R>;
    create<R = T>(data: any, options?: any): Promise<R>;
    update<R = T>(id: string, data: any, options?: any): Promise<R>;
    delete(id: string, options?: any): Promise<boolean>;
}

/**
 * Authentication operations interface
 */
export interface AuthOperations {
    authWithPassword(usernameOrEmail: string, password: string, options?: any): Promise<any>;
    authRefresh(options?: any): Promise<any>;
    requestPasswordReset(email: string, options?: any): Promise<any>;
    confirmPasswordReset(resetToken: string, newPassword: string, options?: any): Promise<any>;
}

/**
 * Complete service interface combining record and auth operations
 */
export interface RecordService<T = any> extends RecordOperations<T>, AuthOperations {}

/**
 * Type for the PocketBase client with strongly typed collections
 */
export type PocketBaseCollections = Omit<PocketBase, 'collection'> & {
    collection(idOrName: string): RecordService; 
    collection(idOrName: 'practice_items'): RecordService<PracticeItem>;
    collection(idOrName: 'practice_topics'): RecordService<PracticeTopic>;
    collection(idOrName: 'practice_sessions'): RecordService<PracticeSession>;
    collection(idOrName: 'practice_results'): RecordService<PracticeResult>;
    collection(idOrName: 'accounts'): RecordService<Account>;
    collection(idOrName: 'users'): RecordService<User>;
    collection(idOrName: 'instructors'): RecordService<Instructor>;
    collection(idOrName: 'learners'): RecordService<Learner>;
} 