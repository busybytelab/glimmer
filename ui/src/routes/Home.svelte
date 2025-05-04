<script lang="ts">
	import { onMount } from 'svelte';

	function navigateToLogin() {
		window.history.pushState({}, '', '/login');
		window.dispatchEvent(new PopStateEvent('popstate'));
	}

	onMount(() => {
		// Set focus to the main heading for screen readers
		const heading = document.querySelector('h1');
		if (heading) {
			heading.setAttribute('tabindex', '-1');
			heading.focus();
		}
	});
</script>

<div class="min-h-screen flex flex-col items-center justify-center bg-gray-50 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8 text-center">
		<div>
			<h1 class="text-4xl font-extrabold text-primary tracking-tight">
				Welcome to Glimmer
			</h1>
			<p class="mt-2 text-lg text-gray-600">
				A self-hosted learning assistant for parents and children
			</p>
		</div>

		<div class="mt-8 space-y-4">
			<button
				on:click={navigateToLogin}
				class="group relative w-full flex justify-center py-3 px-4 border border-transparent text-lg font-medium rounded-md text-white bg-secondary hover:bg-secondary focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-secondary transition-colors duration-200"
			>
				Get Started
			</button>

			<div class="mt-6">
				<div class="relative">
					<div class="absolute inset-0 flex items-center">
						<div class="w-full border-t border-gray-300"></div>
					</div>
					<div class="relative flex justify-center text-sm">
						<span class="px-2 bg-gray-50 text-gray-500">
							Key Features
						</span>
					</div>
				</div>

				<div class="mt-6 grid grid-cols-1 gap-4">
					<div class="flex items-center space-x-3">
						<svg class="h-6 w-6 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
						</svg>
						<span class="text-gray-700">Personalized Learning Paths</span>
					</div>
					<div class="flex items-center space-x-3">
						<svg class="h-6 w-6 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
						</svg>
						<span class="text-gray-700">Progress Tracking</span>
					</div>
					<div class="flex items-center space-x-3">
						<svg class="h-6 w-6 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
						</svg>
						<span class="text-gray-700">Safe and Private</span>
					</div>
				</div>
			</div>
		</div>
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
</style> 