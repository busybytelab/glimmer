import pb from '$lib/pocketbase';
import type { RegistrationForm } from '$lib/types';

class UserService {
	/**
	 * Register a new user
	 * @param form The registration form data
	 * @throws Error if registration fails
	 */
	async register(form: RegistrationForm): Promise<void> {
		try {
			await pb.collection('users').create({
				email: form.email,
				password: form.password,
				passwordConfirm: form.passwordConfirm
			});
		} catch (err) {
			// Handle specific PocketBase errors
			if (err instanceof Error) {
				if (err.message.includes('email')) {
					throw new Error('This email is already registered');
				}
				if (err.message.includes('password')) {
					throw new Error('Password requirements not met');
				}
			}
			throw new Error('Registration failed. Please try again.');
		}
	}
}

export const userService = new UserService(); 