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

    // Base record type that matches PocketBase's structure
    export interface PocketBaseRecord {
        id: string;
        created: string;
        updated: string;
        collectionId: string;
        collectionName: string;
        expand?: Record<string, any>;
    }

    // Type for record service
    export interface RecordService<T = any> {
        getFullList(options?: any): Promise<T[]>;
        getList(page?: number, perPage?: number, options?: any): Promise<{ page: number; perPage: number; totalItems: number; totalPages: number; items: T[] }>;
        getOne(id: string, options?: any): Promise<T>;
        getFirstListItem(filter: string, options?: any): Promise<T>;
        create(data: any, options?: any): Promise<T>;
        update(id: string, data: any, options?: any): Promise<T>;
        delete(id: string, options?: any): Promise<boolean>;
    }

    // Type for the PocketBase client with collections
    export interface PocketBaseCollections extends PocketBase {
        collection(idOrName: string): RecordService; // fallback for any other collection
        collection(idOrName: 'practice_items'): RecordService<import('./types').PracticeItem>;
        collection(idOrName: 'practice_topics'): RecordService<import('./types').PracticeTopic>;
        collection(idOrName: 'practice_sessions'): RecordService<import('./types').PracticeSession>;
        collection(idOrName: 'practice_results'): RecordService<import('./types').PracticeResult>;
        collection(idOrName: 'accounts'): RecordService<import('./types').Account>;
        collection(idOrName: 'users'): RecordService<import('./types').User>;
        collection(idOrName: 'instructors'): RecordService<import('./types').Instructor>;
        collection(idOrName: 'learners'): RecordService<import('./types').Learner>;
    }
} 