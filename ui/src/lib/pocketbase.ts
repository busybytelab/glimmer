/**
 * PocketBase Client Instance
 * 
 * This file initializes and exports the application's PocketBase client.
 * It provides a typed instance with collection methods for all models.
 * 
 * USAGE GUIDELINES:
 * - Import this instance when you need to interact with the database
 * - Example: import pb from '$lib/pocketbase';
 * - Use pb.collection('collection_name') to access a specific collection
 * - The client is already typed to provide proper types for each collection
 * 
 * DO NOT:
 * - Initialize new PocketBase instances elsewhere
 * - Add type definitions here (use pocketbase-types.ts instead)
 * - Add service methods here (create separate service files)
 */

import PocketBase from 'pocketbase';
import type { PocketBaseCollections } from './pocketbase-types';

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
const pb = new PocketBase(getPocketBaseUrl()) as PocketBaseCollections;

// Configure PocketBase to use cookies
pb.authStore.onChange(() => {
  // This ensures cookies are properly set when auth state changes
  document.cookie = `pb_auth_token=${pb.authStore.token || ''}; path=/; SameSite=Lax`;
});

export default pb; 