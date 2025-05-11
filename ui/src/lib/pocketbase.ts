import PocketBase from 'pocketbase';
import type { TypedPocketBase } from './types';

// Get the PocketBase URL from environment variables
// In development, it's typically http://localhost:8787
// In production, it should be the same origin as the UI
const getPocketBaseUrl = () => {
  // Use Vite's environment variables
  // VITE_* variables are exposed to the client
  const envUrl = import.meta.env.VITE_POCKETBASE_URL;
  
  if (envUrl) {
    return envUrl;
  }
  
  // Check if we're in a browser environment
  if (typeof window !== 'undefined') {
    // Always use the same origin as the UI to ensure cookies work properly
    return window.location.origin;
  }
  
  // Default to localhost for server-side rendering
  return 'http://localhost:8787';
};

// Create a typed PocketBase client instance
const pb = new PocketBase(getPocketBaseUrl()) as TypedPocketBase;

// Configure PocketBase to use cookies
pb.authStore.onChange(() => {
  // This ensures cookies are properly set when auth state changes
  document.cookie = `pb_auth_token=${pb.authStore.token || ''}; path=/; SameSite=Lax`;
});

export default pb; 