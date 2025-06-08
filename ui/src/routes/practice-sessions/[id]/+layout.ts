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

        throw redirect(302, `/practice-sessions/${sessionId}/instructor`);
        // TODO: for learner view, we need to refactor, move this under leaner route
        //throw redirect(302, `/practice-sessions/${sessionId}/learner`);
        
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