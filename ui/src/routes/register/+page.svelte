<script lang="ts">
	import { error } from '$lib/stores';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import ErrorAlert from '../../components/common/ErrorAlert.svelte';
	import LoadingSpinner from '../../components/common/LoadingSpinner.svelte';
	import { userService } from '$lib/services/user';

	// Use the public URL instead of importing the asset
	const glimmerLogoUrl = '/glimmer.svg';

	interface RegistrationForm {
		email: string;
		password: string;
		passwordConfirm: string;
	}

	let form: RegistrationForm = {
		email: '',
		password: '',
		passwordConfirm: ''
	};

	let isSubmitting = false;
	let formErrors: Partial<Record<keyof RegistrationForm, string>> = {};
	let returnUrl: string | null = null;

	function validateForm(): boolean {
		formErrors = {};
		let isValid = true;

		if (!form.email) {
			formErrors.email = 'Email is required';
			isValid = false;
		} else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
			formErrors.email = 'Please enter a valid email address';
			isValid = false;
		}

		if (!form.password) {
			formErrors.password = 'Password is required';
			isValid = false;
		} else if (form.password.length < 6) {
			formErrors.password = 'Password must be at least 6 characters';
			isValid = false;
		}

		if (!form.passwordConfirm) {
			formErrors.passwordConfirm = 'Please confirm your password';
			isValid = false;
		} else if (form.password !== form.passwordConfirm) {
			formErrors.passwordConfirm = 'Passwords do not match';
			isValid = false;
		}

		return isValid;
	}

	async function handleSubmit() {
		if (!validateForm()) return;

		isSubmitting = true;
		error.set(null);

		try {
			await userService.register(form);
			// Redirect to login page with success message
			window.location.href = '/login?registered=true';
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Registration failed. Please try again.';
			error.set(errorMessage);
			formErrors.email = errorMessage;
		} finally {
			isSubmitting = false;
		}
	}

	onMount(() => {
		// Get returnUrl from query parameters
		returnUrl = $page.url.searchParams.get('returnUrl');
		// Clear any existing errors when component mounts
		error.set(null);
	});
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-900">
	<div class="max-w-sm w-full space-y-6 p-6 bg-white dark:bg-gray-800 rounded-xl shadow-md">
		<div class="text-center">
			<div class="flex items-center justify-center">
				<img src={glimmerLogoUrl} alt="Glimmer Logo" class="h-10 w-10 mr-3" />
				<h1 class="text-2xl font-bold text-primary dark:text-white">Glimmer</h1>
			</div>
			<p class="mt-1 text-gray-600 dark:text-gray-300 text-base">Create your account</p>
		</div>
		<form on:submit|preventDefault={handleSubmit} class="space-y-4">
			{#if $error}
				<ErrorAlert message={$error} />
			{/if}
			<div class="space-y-2">
				<div>
					<label for="email" class="sr-only">Email address</label>
					<input
						id="email"
						name="email"
						type="email"
						autocomplete="email"
						required
						bind:value={form.email}
						class="block w-full px-3 py-2 rounded-t-md border {formErrors.email ? 'border-red-500 dark:border-red-400' : 'border-gray-300 dark:border-gray-600'} placeholder-gray-400 dark:placeholder-gray-500 text-gray-900 dark:text-gray-100 bg-gray-100 dark:bg-gray-700 focus:outline-none focus:ring-secondary focus:border-secondary sm:text-sm"
						placeholder="Email address"
						disabled={isSubmitting}
						aria-invalid={!!formErrors.email}
						aria-describedby={formErrors.email ? 'email-error' : undefined}
					/>
					{#if formErrors.email}
						<p id="email-error" class="mt-1 text-xs text-red-600 dark:text-red-300">{formErrors.email}</p>
					{/if}
				</div>
				<div>
					<label for="password" class="sr-only">Password</label>
					<input
						id="password"
						name="password"
						type="password"
						autocomplete="new-password"
						required
						bind:value={form.password}
						class="block w-full px-3 py-2 border {formErrors.password ? 'border-red-500 dark:border-red-400' : 'border-gray-300 dark:border-gray-600'} placeholder-gray-400 dark:placeholder-gray-500 text-gray-900 dark:text-gray-100 bg-gray-100 dark:bg-gray-700 focus:outline-none focus:ring-secondary focus:border-secondary sm:text-sm"
						placeholder="Password"
						disabled={isSubmitting}
						aria-invalid={!!formErrors.password}
						aria-describedby={formErrors.password ? 'password-error' : undefined}
					/>
					{#if formErrors.password}
						<p id="password-error" class="mt-1 text-xs text-red-600 dark:text-red-300">{formErrors.password}</p>
					{/if}
				</div>
				<div>
					<label for="passwordConfirm" class="sr-only">Confirm Password</label>
					<input
						id="passwordConfirm"
						name="passwordConfirm"
						type="password"
						autocomplete="new-password"
						required
						bind:value={form.passwordConfirm}
						class="block w-full px-3 py-2 rounded-b-md border {formErrors.passwordConfirm ? 'border-red-500 dark:border-red-400' : 'border-gray-300 dark:border-gray-600'} placeholder-gray-400 dark:placeholder-gray-500 text-gray-900 dark:text-gray-100 bg-gray-100 dark:bg-gray-700 focus:outline-none focus:ring-secondary focus:border-secondary sm:text-sm"
						placeholder="Confirm Password"
						disabled={isSubmitting}
						aria-invalid={!!formErrors.passwordConfirm}
						aria-describedby={formErrors.passwordConfirm ? 'password-confirm-error' : undefined}
					/>
					{#if formErrors.passwordConfirm}
						<p id="password-confirm-error" class="mt-1 text-xs text-red-600 dark:text-red-300">{formErrors.passwordConfirm}</p>
					{/if}
				</div>
			</div>
			<div>
				<button
					type="submit"
					disabled={isSubmitting}
					class="w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-secondary hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-secondary disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
				>
					{#if isSubmitting}
						<div class="w-5 h-5 mr-3">
							<LoadingSpinner size="sm" color="white" />
						</div>
						Creating account...
					{:else}
						Create account
					{/if}
				</button>
			</div>
			<div class="text-center text-sm">
				<p class="text-gray-600 dark:text-gray-400">
					Already have an account?
					<a href="/login" class="font-medium text-secondary hover:text-blue-600 focus:outline-none focus:underline dark:text-blue-400">
						Sign in
					</a>
				</p>
			</div>
		</form>
	</div>
</div>

<style lang="postcss">
	:global(html) {
		--primary: #2c3e50;
		--secondary: #3498db;
	}

	.text-primary {
		color: var(--primary);
	}

	.bg-secondary {
		background-color: var(--secondary);
	}

	.focus\:ring-secondary:focus {
		--tw-ring-color: var(--secondary);
	}

	.focus\:border-secondary:focus {
		border-color: var(--secondary);
	}
</style> 