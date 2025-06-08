<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	// Use the public URL instead of importing the asset
	const glimmerLogoUrl = '/glimmer.svg';

	onMount(() => {
		// If no email parameter is provided, redirect to login
		if (!$page.url.searchParams.has('email')) {
			goto('/login');
		}
	});
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-900">
	<div class="max-w-sm w-full space-y-6 p-6 bg-white dark:bg-gray-800 rounded-xl shadow-md">
		<div class="text-center">
			<div class="flex items-center justify-center">
				<img src={glimmerLogoUrl} alt="Glimmer Logo" class="h-10 w-10 mr-3" />
				<h1 class="text-2xl font-bold text-primary dark:text-white">Glimmer</h1>
			</div>
			<div class="mt-6 space-y-4">
				<h2 class="text-xl font-semibold text-gray-900 dark:text-white">Check your email</h2>
				<p class="text-gray-600 dark:text-gray-300">
					We've sent a verification link to your email address. Please check your inbox and click the link to verify your account.
				</p>
				<div class="pt-4">
					<button
						on:click={() => goto('/login')}
						class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-secondary hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-secondary transition-colors duration-200"
					>
						Return to login
					</button>
				</div>
			</div>
		</div>
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
</style> 