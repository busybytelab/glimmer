import pb from '$lib/pocketbase';

/**
 * Service for handling user role-related operations
 */
export const rolesService = {
    /**
     * Checks if the current user is an instructor
     * @returns {Promise<boolean>} True if the user is an instructor, false otherwise
     */
    async isInstructor(): Promise<boolean> {
        try {
            const authData = pb.authStore.model;
            if (!authData) {
                return false;
            }

            try {
                const instructorRecord = await pb.collection('instructors').getFirstListItem(`user="${authData.id}"`);
                return !!instructorRecord;
            } catch {
                return false;
            }
        } catch (err) {
            console.error('Failed to check user role:', err);
            return false;
        }
    }
}; 