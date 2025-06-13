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

<div class="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-900 px-4 sm:px-6">
	<div class="max-w-md w-full space-y-8 p-6 sm:p-8 bg-white dark:bg-gray-800 rounded-xl shadow-md">
		<div class="text-center">
			<div class="flex items-center justify-center">
				<img src={glimmerLogoUrl} alt="Glimmer Logo" class="h-12 w-12 sm:h-14 sm:w-14 mr-3" />
				<h1 class="text-2xl sm:text-3xl font-bold text-primary dark:text-white">Glimmer</h1>
			</div>
			<div class="mt-8 space-y-6">
				<div class="space-y-4">
					<div class="flex justify-center">
						<div class="h-12 w-12 sm:h-14 sm:w-14 text-secondary">
							<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 19v-8.93a2 2 0 01.89-1.664l7-4.666a2 2 0 012.22 0l7 4.666A2 2 0 0121 10.07V19M3 19a2 2 0 002 2h14a2 2 0 002-2M3 19l6.75-4.5M21 19l-6.75-4.5M3 10l6.75 4.5M21 10l-6.75 4.5m0 0l-1.14.76a2 2 0 01-2.22 0l-1.14-.76" />
							</svg>
						</div>
					</div>
					<h2 class="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">Check your email</h2>
					<p class="text-base sm:text-lg text-gray-600 dark:text-gray-300">
						We've sent a verification link to your email address. Please check your inbox and click the link to verify your account.
					</p>
					<p class="text-sm text-gray-500 dark:text-gray-400">
						Didn't receive the email? Check your spam folder or
						<button
							class="text-secondary hover:text-emerald-500 focus:outline-none focus:underline font-medium"
							on:click={() => window.location.reload()}
						>
							click here to resend
						</button>
					</p>
				</div>
				<div class="pt-4">
					<button
						on:click={() => goto('/login')}
						class="w-full flex justify-center py-3 sm:py-4 px-4 border-2 border-gray-300 dark:border-gray-600 text-base sm:text-lg font-semibold rounded-lg text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition-all duration-200 transform hover:scale-[1.02] shadow-sm hover:shadow-md"
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

	.text-secondary {
		color: var(--secondary);
	}
</style> 