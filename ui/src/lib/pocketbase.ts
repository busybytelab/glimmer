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
  
  // Fallback to development URL if no environment variable is set
  if (import.meta.env.DEV) {
    return 'http://localhost:8787';
  }
  
  // In production, use the same origin as the UI
  return window.location.origin;
};

// Create a typed PocketBase client instance
const pb = new PocketBase(getPocketBaseUrl()) as TypedPocketBase;

export default pb; 