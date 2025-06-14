import type { Load } from '@sveltejs/kit';

export const load: Load = ({ params, url }) => {
    return {
        id: params.id || null,
        initialPrompt: url.searchParams.get('prompt') || null
    };
}; 