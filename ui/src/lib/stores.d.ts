import type { Writable } from 'svelte/store';
import type { Learner } from './types';

// TODO: look at the usage of user and fix it
export const user: Writable<| Learner | null>;
export const isAuthenticated: Writable<boolean>;
export const isLoading: Writable<boolean>;
export const error: Writable<string | null>;
export const isAuthLoading: Writable<boolean>; 