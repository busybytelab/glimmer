import pb from './pocketbase';

/**
 * Gets the current authentication token from either PocketBase's authStore or localStorage.
 * @returns The authentication token or an empty string if not authenticated.
 */
export function getAuthToken(): string {
  // Check if PocketBase has a valid token in memory
  if (pb.authStore.isValid && typeof pb.authStore.token === 'string' && pb.authStore.token !== '') {
    return pb.authStore.token;
  }
  
  // Fallback to localStorage
  const storedToken = localStorage.getItem('authToken');
  return storedToken !== null ? storedToken : '';
}

/**
 * Gets the current user ID from PocketBase's authStore.
 * @returns The user ID or null if not authenticated.
 */
export function getCurrentUserId(): string | null {
  if (pb.authStore.isValid && pb.authStore.model) {
    return pb.authStore.model.id;
  }
  return null;
}

/**
 * Checks if the user is currently authenticated.
 * @returns True if authenticated, false otherwise.
 */
export function isAuthenticated(): boolean {
  return getAuthToken() !== '';
}

/**
 * Saves the authentication token to localStorage for persistence between sessions.
 * @param rememberMe Whether to save the token in localStorage for persistence.
 */
export function saveAuthToken(rememberMe: boolean = true): void {
  if (rememberMe && pb.authStore.isValid) {
    localStorage.setItem('authToken', pb.authStore.token);
  }
}

/**
 * Clears the authentication token from both PocketBase's authStore and localStorage.
 */
export function clearAuthToken(): void {
  pb.authStore.clear();
  localStorage.removeItem('authToken');
}

// Public routes that don't require authentication
export const PUBLIC_ROUTES = ['/login', '/forgot-password', '/reset-password'] as const;

// Type for public routes
export type PublicRoute = typeof PUBLIC_ROUTES[number];

// Helper function to check if a path is a public route
export function isPublicRoute(path: string): boolean {
    return PUBLIC_ROUTES.some(route => path.startsWith(route));
} 