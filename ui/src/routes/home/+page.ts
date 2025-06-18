import { browser } from '$app/environment';
import { redirect } from '@sveltejs/kit';
import { goto } from '$app/navigation';
import { learnersService } from '$lib/services/learners';
import type { LoadEvent } from '@sveltejs/kit';
import pb from '$lib/pocketbase';
import type { Learner } from '$lib/types';

export interface PageData {
    learners: Learner[];
    error: string | null;
}

export const load = async ({ fetch }: LoadEvent) => {
    if (browser) {
        try {
            if (!pb.authStore.isValid) {
                await goto('/login');
                return; 
            }

            const learners = await learnersService.getLearners(1, 50, fetch);
            
            if (!learners || learners.length === 0) {
                await goto('/account/dashboard');
                return;
            }

            return { learners };
        } catch (err) {
            console.error('Error in load function:', err);
            return {
                learners: [],
                error: 'Failed to load learners. Please try again.'
            };
        }
    }

    // Server-side execution (for completeness, though our app is client-side focused)
    try {
        if (!pb.authStore.isValid) {
            throw redirect(302, '/login');
        }
        const learners = await learnersService.getLearners(1, 50, fetch);
        if (!learners || learners.length === 0) {
            throw redirect(302, '/account/dashboard');
        }
        return { learners };
    } catch(err) {
        if (err instanceof Response && err.status >= 300 && err.status < 400) {
            throw err; // Re-throw redirects
        }
        return { learners: [], error: 'Failed to load learners.' };
    }
}; 