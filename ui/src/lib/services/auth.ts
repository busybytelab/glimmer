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

  /**
   * Refreshes the authentication token using PocketBase's authRefresh method.
   * @param customFetch Optional fetch function to use (provided by SvelteKit)
   * @returns Promise that resolves when the token is refreshed
   * @throws Error if the refresh fails
   */
  async refreshAuthToken(customFetch?: typeof fetch): Promise<void> {
    try {
      if (customFetch) {
        const originalSend = pb.send;
        pb.send = async <T = any>(path: string, options: any = {}) => {
          const url = pb.buildURL(path);
          // Ensure headers are properly set
          options.headers = {
            ...options.headers,
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${this.getAuthToken()}`
          };
          const response = await customFetch(url, options);
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          return response.json() as Promise<T>;
        };
        await pb.collection('users').authRefresh();
        pb.send = originalSend;
      } else {
        await pb.collection('users').authRefresh();
      }
      this.saveAuthToken();
    } catch (error) {
      this.clearAuthToken();
      throw error;
    }
  }

  /**
   * Authenticates a user with email and password.
   * @param email User's email address
   * @param password User's password
   * @param rememberMe Whether to persist the auth token
   * @returns Promise that resolves when authentication is successful
   * @throws Error if authentication fails
   */
  async authWithPassword(email: string, password: string, rememberMe: boolean = true): Promise<void> {
    try {
      await pb.collection('users').authWithPassword(email, password);
      this.saveAuthToken(rememberMe);
    } catch (error) {
      this.clearAuthToken();
      // handle error: The request doesn't satisfy the collection requirements to authenticate
      throw error;
    }
  }
}

// Export a singleton instance
export const authService = new AuthService(); 