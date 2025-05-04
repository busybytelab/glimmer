import { writable } from 'svelte/store';
import type { Instructor, Learner } from './types';

export const user = writable<Instructor | Learner | null>(null);
export const isAuthenticated = writable(false);
export const isLoading = writable(false);
export const error = writable<string | null>(null);
export const isAuthLoading = writable(true); 