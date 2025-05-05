<script lang="ts">
	import SideNav from './SideNav.svelte';

	export let sidebarOpen = true;

	// Toggle sidebar
	function toggleSidebar() {
		sidebarOpen = !sidebarOpen;
	}
</script>

<div class="h-screen flex overflow-hidden bg-gray-100">
	<!-- Mobile sidebar -->
	{#if sidebarOpen}
		<div class="md:hidden fixed inset-0 flex z-40">
			<button
				class="fixed inset-0 bg-gray-600 bg-opacity-75"
				on:click={() => sidebarOpen = false}
				on:keydown={(e) => e.key === 'Escape' && (sidebarOpen = false)}
				aria-label="Close sidebar overlay"
			></button>
			
			<!-- Mobile SideNav -->
			<SideNav isOpen={true} />
		</div>
	{/if}

	<!-- Desktop sidebar -->
	<div class="hidden md:flex md:flex-shrink-0">
		<SideNav isOpen={sidebarOpen} />
	</div>

	<!-- Main content -->
	<div class="flex flex-col w-0 flex-1 overflow-hidden h-screen">
		<div class="md:hidden pl-1 pt-1 sm:pl-3 sm:pt-3">
			<button
				on:click={toggleSidebar}
				class="-ml-0.5 -mt-0.5 h-12 w-12 inline-flex items-center justify-center rounded-md text-gray-500 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-secondary"
			>
				<span class="sr-only">Open sidebar</span>
				<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
				</svg>
			</button>
		</div>
		<main class="flex-1 relative z-0 overflow-y-auto focus:outline-none h-full">
			<slot />
		</main>
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

	.focus\:ring-secondary:focus {
		--tw-ring-color: var(--secondary);
	}
</style> 