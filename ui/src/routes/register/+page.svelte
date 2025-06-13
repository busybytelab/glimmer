<script lang="ts">
	import { error } from '$lib/stores';
	import { onMount } from 'svelte';
	import ErrorAlert from '$components/common/ErrorAlert.svelte';
	import LoadingSpinner from '$components/common/LoadingSpinner.svelte';
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
			// Redirect to verify-email page with email parameter
			window.location.href = `/verify-email?email=${encodeURIComponent(form.email)}`;
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Registration failed. Please try again.';
			error.set(errorMessage);
			formErrors.email = errorMessage;
		} finally {
			isSubmitting = false;
		}
	}

	onMount(() => {
		// Clear any existing errors when component mounts
		error.set(null);
	});
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-900 px-4 sm:px-6">
	<div class="max-w-md w-full space-y-8 p-6 sm:p-8 bg-white dark:bg-gray-800 rounded-xl shadow-md">
		<div class="text-center">
			<div class="flex items-center justify-center">
				<img src={glimmerLogoUrl} alt="Glimmer Logo" class="h-12 w-12 sm:h-14 sm:w-14 mr-3" />
				<h1 class="text-2xl sm:text-3xl font-bold text-primary dark:text-white">Glimmer</h1>
			</div>
			<p class="mt-2 text-gray-600 dark:text-gray-300 text-base sm:text-lg">Create your account</p>
		</div>
		<form on:submit|preventDefault={handleSubmit} class="space-y-6">
			{#if $error}
				<ErrorAlert message={$error} />
			{/if}
			<div class="space-y-4">
				<div>
					<label for="email" class="block text-sm sm:text-base font-medium text-gray-700 dark:text-gray-300 mb-1">Email address</label>
					<input
						id="email"
						name="email"
						type="email"
						autocomplete="email"
						required
						bind:value={form.email}
						class="block w-full px-4 py-3 sm:py-4 rounded-lg border {formErrors.email ? 'border-red-500 dark:border-red-400' : 'border-gray-300 dark:border-gray-600'} placeholder-gray-400 dark:placeholder-gray-500 text-gray-900 dark:text-gray-100 bg-gray-100 dark:bg-gray-700 focus:outline-none focus:ring-secondary focus:border-secondary text-base"
						placeholder="Enter your email"
						disabled={isSubmitting}
						aria-invalid={!!formErrors.email}
						aria-describedby={formErrors.email ? 'email-error' : undefined}
					/>
					{#if formErrors.email}
						<p id="email-error" class="mt-2 text-sm text-red-600 dark:text-red-300">{formErrors.email}</p>
					{/if}
				</div>
				<div>
					<label for="password" class="block text-sm sm:text-base font-medium text-gray-700 dark:text-gray-300 mb-1">Password</label>
					<input
						id="password"
						name="password"
						type="password"
						autocomplete="new-password"
						required
						bind:value={form.password}
						class="block w-full px-4 py-3 sm:py-4 rounded-lg border {formErrors.password ? 'border-red-500 dark:border-red-400' : 'border-gray-300 dark:border-gray-600'} placeholder-gray-400 dark:placeholder-gray-500 text-gray-900 dark:text-gray-100 bg-gray-100 dark:bg-gray-700 focus:outline-none focus:ring-secondary focus:border-secondary text-base"
						placeholder="Create a password"
						disabled={isSubmitting}
						aria-invalid={!!formErrors.password}
						aria-describedby={formErrors.password ? 'password-error' : undefined}
					/>
					{#if formErrors.password}
						<p id="password-error" class="mt-2 text-sm text-red-600 dark:text-red-300">{formErrors.password}</p>
					{/if}
				</div>
				<div>
					<label for="passwordConfirm" class="block text-sm sm:text-base font-medium text-gray-700 dark:text-gray-300 mb-1">Confirm Password</label>
					<input
						id="passwordConfirm"
						name="passwordConfirm"
						type="password"
						autocomplete="new-password"
						required
						bind:value={form.passwordConfirm}
						class="block w-full px-4 py-3 sm:py-4 rounded-lg border {formErrors.passwordConfirm ? 'border-red-500 dark:border-red-400' : 'border-gray-300 dark:border-gray-600'} placeholder-gray-400 dark:placeholder-gray-500 text-gray-900 dark:text-gray-100 bg-gray-100 dark:bg-gray-700 focus:outline-none focus:ring-secondary focus:border-secondary text-base"
						placeholder="Confirm your password"
						disabled={isSubmitting}
						aria-invalid={!!formErrors.passwordConfirm}
						aria-describedby={formErrors.passwordConfirm ? 'password-confirm-error' : undefined}
					/>
					{#if formErrors.passwordConfirm}
						<p id="password-confirm-error" class="mt-2 text-sm text-red-600 dark:text-red-300">{formErrors.passwordConfirm}</p>
					{/if}
				</div>
			</div>
			<div class="pt-2">
				<button
					type="submit"
					disabled={isSubmitting}
					class="w-full flex justify-center py-3 sm:py-4 px-4 border border-transparent text-base sm:text-lg font-semibold rounded-lg text-white bg-secondary hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-secondary disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 transform hover:scale-[1.02] shadow-sm hover:shadow-md"
				>
					{#if isSubmitting}
						<div class="w-6 h-6 mr-3">
							<LoadingSpinner size="sm" color="white" />
						</div>
						Creating account...
					{:else}
						Create account
					{/if}
				</button>
			</div>
			<!-- Password requirements hint -->
			<div class="text-sm text-gray-500 dark:text-gray-400 space-y-1 mt-4">
				<p class="font-medium">Password requirements:</p>
				<ul class="list-disc list-inside">
					<li>At least 6 characters long</li>
					<li>Passwords must match</li>
				</ul>
			</div>
			<!-- Improved sign in section with better visibility -->
			<div class="relative mt-8">
				<div class="absolute inset-0 flex items-center">
					<div class="w-full border-t border-gray-300 dark:border-gray-600"></div>
				</div>
				<div class="relative flex justify-center text-base sm:text-lg">
					<span class="px-4 bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-400 font-medium">Already have an account?</span>
				</div>
			</div>
			<div class="mt-6">
				<a
					href="/login"
					class="w-full flex justify-center py-3 sm:py-4 px-4 border-2 border-gray-300 dark:border-gray-600 text-base sm:text-lg font-semibold rounded-lg text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition-all duration-200 transform hover:scale-[1.02] shadow-sm hover:shadow-md"
				>
					Sign in to your account
				</a>
				<p class="mt-4 text-center text-sm text-gray-500 dark:text-gray-400">
					By creating an account, you agree to our
					<a href="/terms" class="text-secondary hover:text-emerald-500 focus:outline-none focus:underline dark:text-blue-400">
						Terms & Conditions
					</a>
					and
					<a href="/privacy" class="text-secondary hover:text-emerald-500 focus:outline-none focus:underline dark:text-blue-400">
						Privacy Policy
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