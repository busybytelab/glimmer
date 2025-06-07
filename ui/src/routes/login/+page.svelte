<script lang="ts">
	import { isAuthenticated, error } from '$lib/stores';
	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';
	import { saveAuthToken } from '$lib/auth';
	import { page } from '$app/stores';
	import ErrorAlert from '../../components/common/ErrorAlert.svelte';
	import LoadingSpinner from '../../components/common/LoadingSpinner.svelte';
	// Use the public URL instead of importing the asset
	const glimmerLogoUrl = '/glimmer.svg';

	interface LoginForm {
		email: string;
		password: string;
		rememberMe: boolean;
	}

	let form: LoginForm = {
		email: '',
		password: '',
		rememberMe: false
	};

	let isSubmitting = false;
	let formErrors: Partial<Record<keyof LoginForm, string>> = {};
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

		return isValid;
	}

	async function handleSubmit() {
		if (!validateForm()) return;

		isSubmitting = true;
		error.set(null);

		try {
			await pb.collection('users').authWithPassword(form.email, form.password);
			// Use the centralized auth utility to save the token if rememberMe is checked
			saveAuthToken(form.rememberMe);
			isAuthenticated.set(true);
			
			// Redirect to returnUrl if it exists, otherwise to role selection
			const redirectPath = returnUrl ? decodeURIComponent(returnUrl) : '/select-role';
			window.location.href = redirectPath;
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Login failed. Please check your credentials.';
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
			<p class="mt-1 text-gray-600 dark:text-gray-300 text-base">Learning Helper for Kids</p>
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
						autocomplete="current-password"
						required
						bind:value={form.password}
						class="block w-full px-3 py-2 rounded-b-md border {formErrors.password ? 'border-red-500 dark:border-red-400' : 'border-gray-300 dark:border-gray-600'} placeholder-gray-400 dark:placeholder-gray-500 text-gray-900 dark:text-gray-100 bg-gray-100 dark:bg-gray-700 focus:outline-none focus:ring-secondary focus:border-secondary sm:text-sm"
						placeholder="Password"
						disabled={isSubmitting}
						aria-invalid={!!formErrors.password}
						aria-describedby={formErrors.password ? 'password-error' : undefined}
					/>
					{#if formErrors.password}
						<p id="password-error" class="mt-1 text-xs text-red-600 dark:text-red-300">{formErrors.password}</p>
					{/if}
				</div>
			</div>
			<div class="flex items-center justify-between">
				<div class="flex items-center">
					<input
						id="remember-me"
						name="remember-me"
						type="checkbox"
						bind:checked={form.rememberMe}
						class="h-4 w-4 text-secondary focus:ring-secondary border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700"
						disabled={isSubmitting}
					/>
					<label for="remember-me" class="ml-2 block text-xs text-gray-900 dark:text-gray-200">
						Remember me
					</label>
				</div>
				<div class="text-xs">
					<a href="/forgot-password" class="font-medium text-secondary hover:text-secondary focus:outline-none focus:underline dark:text-blue-400">
						Forgot your password?
					</a>
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
						Signing in...
					{:else}
						Sign in
					{/if}
				</button>
			</div>
			<div class="text-center text-sm">
				<p class="text-gray-600 dark:text-gray-400">
					Don't have an account?
					<a href="/register" class="font-medium text-secondary hover:text-blue-600 focus:outline-none focus:underline dark:text-blue-400">
						Create one
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