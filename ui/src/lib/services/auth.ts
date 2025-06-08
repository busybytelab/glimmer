import pb from '../pocketbase';

// Public routes that don't require authentication
export const PUBLIC_ROUTES = ['/login', '/forgot-password', '/reset-password', '/register', '/verify-email'] as const;

// Type for public routes
export type PublicRoute = typeof PUBLIC_ROUTES[number];

export class AuthService {
  /**
   * Gets the current authentication token from either PocketBase's authStore or localStorage.
   * @returns The authentication token or an empty string if not authenticated.
   */
  getAuthToken(): string {
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
  getCurrentUserId(): string | null {
    if (pb.authStore.isValid && pb.authStore.model) {
      return pb.authStore.model.id;
    }
    return null;
  }

  /**
   * Checks if the user is currently authenticated.
   * @returns True if authenticated, false otherwise.
   */
  isAuthenticated(): boolean {
    return this.getAuthToken() !== '';
  }

  /**
   * Saves the authentication token to localStorage for persistence between sessions.
   * @param rememberMe Whether to save the token in localStorage for persistence.
   */
  saveAuthToken(rememberMe: boolean = true): void {
    if (rememberMe && pb.authStore.isValid) {
      localStorage.setItem('authToken', pb.authStore.token);
    }
  }

  /**
   * Clears the authentication token from both PocketBase's authStore and localStorage.
   */
  clearAuthToken(): void {
    pb.authStore.clear();
    localStorage.removeItem('authToken');
  }

  /**
   * Helper function to check if a path is a public route
   */
  isPublicRoute(path: string): boolean {
    return PUBLIC_ROUTES.some(route => path.startsWith(route));
  }
}

// Export a singleton instance
export const authService = new AuthService(); 