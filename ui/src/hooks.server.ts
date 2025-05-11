import { redirect, type Handle } from '@sveltejs/kit';
import { building } from '$app/environment';
import { isPublicRoute } from '$lib/auth';

export const handle: Handle = async ({ event, resolve }) => {
    // Skip auth check during build
    if (building) {
        return await resolve(event);
    }

    const path = event.url.pathname;
    const isPublic = isPublicRoute(path);

    // Get auth token from cookies
    const authCookie = event.request.headers.get('cookie')?.split(';')
        .find(c => {
            const trimmed = c.trim();
            if (!trimmed.startsWith('pb_auth_token=')) {
                return false;
            }
            const value = trimmed.split('=')[1];
            if (value && value !== 'null' && value.trim().length > 0) {
                console.log('authCookie', value);
                return true;
            }
            return false;
        });
    
    // Check if we have a valid auth cookie
    const isAuthenticated = !!authCookie;
    console.log('isAuthenticated', isAuthenticated ? authCookie : 'no');

    // If authenticated and trying to access login page, redirect to dashboard
    if (isAuthenticated && isPublic) {
        throw redirect(303, '/dashboard');
    }

    // Redirect to login if not authenticated and trying to access protected route
    if (!isAuthenticated && !isPublic) {
        // Get the current URL including search params
        const returnUrl = event.url.pathname + event.url.search;
        // Only encode if there's actually a returnUrl
        const encodedReturnUrl = returnUrl ? encodeURIComponent(returnUrl) : '';
        const redirectUrl = `/login${encodedReturnUrl ? `?returnUrl=${encodedReturnUrl}` : ''}`;
        throw redirect(303, redirectUrl);
    }

    return await resolve(event);
}; 