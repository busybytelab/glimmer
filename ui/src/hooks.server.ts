import { redirect, type Handle } from '@sveltejs/kit';
import { building } from '$app/environment';
import { authService } from '$lib/services/auth';

export const handle: Handle = async ({ event, resolve }) => {
    // Skip auth check during build
    if (building) {
        return await resolve(event);
    }

    const path = event.url.pathname;

    // Allow access to public routes
    if (authService.isPublicRoute(path)) {
        return await resolve(event);
    }

    // Check for authentication
    const token = event.request.headers.get('cookie')?.match(/pb_auth_token=([^;]+)/)?.[1];
    if (!token) {
        throw redirect(303, '/login');
    }

    return await resolve(event);
}; 