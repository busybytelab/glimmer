import pb from '$lib/pocketbase';
import type { Account } from '$lib/types';

class AccountService {
    /**
     * Gets the account information for the currently logged in user
     * @returns The account information for the current user
     * @throws Error if user is not logged in or account cannot be found
     */
    async getAccount(): Promise<Account> {
        const currentUser = pb.authStore.model;
        if (!currentUser) {
            throw new Error('You must be logged in to access account information');
        }

        try {
            const accounts = await pb.collection('accounts').getList(1, 1, {
                filter: `owner.id = "${currentUser.id}"`
            });
            
            if (!accounts || accounts.items.length === 0) {
                throw new Error('No account found for current user');
            }

            return accounts.items[0] as unknown as Account;
        } catch (err) {
            console.error('Failed to get account:', err);
            throw new Error('Failed to get account information');
        }
    }
}

export const accountService = new AccountService(); 