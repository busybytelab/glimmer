import { writable } from 'svelte/store';
import type { Instructor, Learner } from './types';

export const user = writable<Instructor | Learner | null>(null);
export const isAuthenticated = writable(false);
export const isLoading = writable(false);
export const error = writable<string | null>(null);
export const isAuthLoading = writable(true);

// Theme store
export const theme = writable<'light' | 'dark'>(
  typeof localStorage !== 'undefined' && localStorage.getItem('theme') 
    ? localStorage.getItem('theme') as 'light' | 'dark'
    : window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
);

// Update the DOM when theme changes
if (typeof document !== 'undefined') {
  theme.subscribe(value => {
    if (value === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
    // Save to localStorage
    localStorage.setItem('theme', value);
  });
} 