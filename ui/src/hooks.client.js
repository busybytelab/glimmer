// This file contains SvelteKit client hooks
// It can be used to initialize client-side libraries or handle client-side events

/** @type {import('@sveltejs/kit').HandleClientError} */
export function handleError({ error, event }) {
  console.error('Client error:', error);
  return {
    message: 'An unexpected error occurred. Please try again.'
  };
} 