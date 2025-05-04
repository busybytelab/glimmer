<script lang="ts">
	import { isAuthenticated, error } from '$lib/stores';
	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';

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
			if (form.rememberMe) {
				// Store auth token in localStorage for persistence
				localStorage.setItem('authToken', pb.authStore.token);
			}
			isAuthenticated.set(true);
			// Use reload instead of history manipulation to ensure proper auth state
			window.location.href = '/dashboard';
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Login failed. Please check your credentials.';
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

<div class="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-900">
	<div class="max-w-sm w-full space-y-6 p-6 bg-white dark:bg-gray-800 rounded-xl shadow-md">
		<div class="text-center">
			<div class="flex items-center justify-center">
				<img src="/glimmer.svg" alt="Glimmer Logo" class="h-10 w-10 mr-3" />
				<h1 class="text-2xl font-bold text-primary dark:text-white">Glimmer</h1>
			</div>
			<p class="mt-1 text-gray-600 dark:text-gray-300 text-base">Learning Helper for Kids</p>
		</div>
		<form on:submit|preventDefault={handleSubmit} class="space-y-4">
			{#if $error}
				<div class="bg-red-100 dark:bg-red-900 border border-red-400 dark:border-red-700 text-red-700 dark:text-red-200 px-3 py-2 rounded relative text-sm" role="alert">
					<strong class="font-bold">Error!</strong>
					<span class="block sm:inline">{$error}</span>
				</div>
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
						<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Signing in...
					{:else}
						Sign in
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>

<style>
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