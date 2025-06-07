import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export const isAuthenticated = writable(false);
export const isLoading = writable(false);
export const error = writable<string | null>(null);
export const isAuthLoading = writable(true);

// Get initial theme value safely
function getInitialTheme(): 'light' | 'dark' {
  if (!browser) {
    return 'light'; // Default for SSR
  }
  
  // Check localStorage first
  const storedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null;
  if (storedTheme) {
    return storedTheme;
  }
  
  // Fall back to system preference
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

// Theme store
export const theme = writable<'light' | 'dark'>(getInitialTheme());

// Update the DOM when theme changes - only in browser
if (browser) {
  theme.subscribe(value => {
    document.documentElement.classList.toggle('dark', value === 'dark');
    localStorage.setItem('theme', value);
  });
} 