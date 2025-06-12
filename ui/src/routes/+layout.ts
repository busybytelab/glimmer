import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { authService } from '$lib/services/auth';
import type { LayoutLoad } from './$types';

export const prerender = false;
export const ssr = false;

export const load: LayoutLoad = async ({ url, fetch }) => {
  // Skip auth check for public routes
  if (authService.isPublicRoute(url.pathname)) {
    return {
      authenticated: false
    };
  }

  // Only run auth checks in the browser
  if (browser) {
    try {
      const token = authService.getAuthToken();
      
      if (!token) {
        // No token found, redirect to login
        await goto('/login');
        return {
          authenticated: false
        };
      }

      // Try to refresh the token using the provided fetch instance
      await authService.refreshAuthToken(fetch);
      
      return {
        authenticated: true
      };
    } catch (error) {
      // Auth refresh failed, redirect to login
      await goto('/login');
      return {
        authenticated: false
      };
    }
  }

  return {
    authenticated: false
  };
}; 
