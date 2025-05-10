import { redirect, type Handle } from '@sveltejs/kit';
import { building } from '$app/environment';

// List of public routes that don't require authentication
const publicRoutes = ['/login', '/forgot-password', '/reset-password'];

export const handle: Handle = async ({ event, resolve }) => {
    // Skip auth check during build
    if (building) {
        return await resolve(event);
    }

    const path = event.url.pathname;
    const isPublicRoute = publicRoutes.some(route => path.startsWith(route));

    // Get auth token from cookies
    const authCookie = event.request.headers.get('cookie')?.split(';')
        .find(c => c.trim().startsWith('pb_auth='));
    
    // Check if we have a valid auth cookie
    const isAuthenticated = !!authCookie;
    console.log('path: '+ path + ', isAuthenticated: ' + isAuthenticated + ', isPublicRoute: ' + isPublicRoute);

    // If authenticated and trying to access login page, redirect to dashboard
    if (isAuthenticated && isPublicRoute) {
        throw redirect(303, '/dashboard');
    }

    // Redirect to login if not authenticated and trying to access protected route
    if (!isAuthenticated && !isPublicRoute) {
        // Get the current URL including search params
        const returnUrl = event.url.pathname + event.url.search;
        // Only encode if there's actually a returnUrl
        const encodedReturnUrl = returnUrl ? encodeURIComponent(returnUrl) : '';
        const redirectUrl = `/login${encodedReturnUrl ? `?returnUrl=${encodedReturnUrl}` : ''}`;
        throw redirect(303, redirectUrl);
    }

    if (!isAuthenticated && isPublicRoute) {
        const returnUrl = event.url.search;
        console.log('returnUrl: ' + returnUrl);

    }    

    return await resolve(event);
}; 