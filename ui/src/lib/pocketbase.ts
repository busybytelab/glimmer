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

// Store the original send method
const originalSend = pb.send;

/**
 * Checks if an error is due to PocketBase auto-cancellation
 * @param err - The error to check
 * @returns true if the error is due to auto-cancellation
 */
function isAutoCancelledErrorInternal(err: any): boolean {
  // Check for the isAbort property that PocketBase sets on auto-cancelled requests
  if (err && err.isAbort === true) {
    return true;
  }

  // Also check for error messages that indicate auto-cancellation
  if (err && err.message && typeof err.message === 'string') {
    const message = err.message.toLowerCase();
    return (
      message.includes('autocancelled') ||
      message.includes('auto-cancelled') ||
      message.includes('request was autocancelled') ||
      message.includes('the request was autocancelled')
    );
  }

  return false;
}

// Override the send method to globally handle auto-cancellation.
pb.send = async function <T = any>(path: string, options: any = {}): Promise<T> {
  try {
    // Pass all requests through to the original send method
    return await originalSend.call(this, path, options) as T;
  } catch (err) {
    // Check if the error is due to auto-cancellation
    if (isAutoCancelledErrorInternal(err)) {
      // This is an expected auto-cancellation. We can "swallow" this error
      // by returning a promise that never resolves. This prevents the error
      // from propagating to the console as an unhandled promise rejection.
      // The new request that superseded this one will proceed as normal.
      console.log(`Request cancelled: ${path}`);
      return new Promise<T>(() => {}); // A promise that never resolves
    }

    // For any other type of error, re-throw it to be handled by the caller.
    throw err;
  }
};


export default pb; 