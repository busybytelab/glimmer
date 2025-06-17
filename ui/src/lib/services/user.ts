import pb from '$lib/pocketbase';
import type { RegistrationForm, User } from '$lib/types';

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

	async resendVerificationEmail(email?: string): Promise<void> {
		if (email) {
			await pb.collection('users').requestVerification(email);
			console.log('Verification email sent');
		} else {
			console.error('No email provided');
		}
	}

	async getCurrentUser(): Promise<User> {
		if (!pb.authStore.record?.id) {
			throw new Error('No user found');
		}
		const user = await pb.collection('users').getOne(pb.authStore.record?.id);
		return user;
	}

	async updateUser(user: User): Promise<User> {
		const updatedUser = await pb.collection('users').update(user.id, user);
		return updatedUser;
	}
}

export const userService = new UserService(); 