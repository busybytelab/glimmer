import type { Writable } from 'svelte/store';
import type { Instructor, Learner } from './types';

export const user: Writable<Instructor | Learner | null>;
export const isAuthenticated: Writable<boolean>;
export const isLoading: Writable<boolean>;
export const error: Writable<string | null>;
export const isAuthLoading: Writable<boolean>; 