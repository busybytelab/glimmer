import type { Writable } from 'svelte/store';
import type { Learner } from './types';

// TODO: check usage of user and see if we can remove it
export const user: Writable<Learner | null>;
export const isAuthenticated: Writable<boolean>;
export const isLoading: Writable<boolean>;
export const error: Writable<string | null>;
export const isAuthLoading: Writable<boolean>; 