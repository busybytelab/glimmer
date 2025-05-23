import { rolesService } from '$lib/services/roles';
import { redirect } from '@sveltejs/kit';
import pb from '$lib/pocketbase';

export const load = async ({ params }: { params: { id: string } }) => {
    try {
        const sessionId = params.id;
        if (!sessionId) {
            return {
                loading: false,
                error: 'Invalid session ID'
            };
        }

        // Check if user is authenticated
        if (!pb.authStore.isValid) {
            // Redirect to login page with return URL
            const returnUrl = encodeURIComponent(`/practice-sessions/${sessionId}/learner`);
            throw redirect(302, `/login?returnUrl=${returnUrl}`);
        }

        // Get user role
        const isInstructor = await rolesService.isInstructor();
        
        // Redirect based on role
        if (isInstructor) {
            throw redirect(302, `/practice-sessions/${sessionId}/instructor`);
        } else {
            throw redirect(302, `/practice-sessions/${sessionId}/learner`);
        }
    } catch (err) {
        if (err instanceof Error && err.message.includes('redirect')) {
            throw err; // Re-throw redirect errors
        }
        console.error('Error in layout load:', err);
        return {
            loading: false,
            error: err instanceof Error ? err.message : 'An unexpected error occurred'
        };
    }
}; 