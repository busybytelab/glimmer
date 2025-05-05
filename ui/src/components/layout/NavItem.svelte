<script lang="ts">
	export let href: string;
	export let text: string;
	export let icon: string;
	export let active: boolean = false;
	export let mobile: boolean = false;
	
	// Convert regular paths to hash-based routes
	const hashHref = href.startsWith('/') ? `#${href.substring(1)}` : `#${href}`;
	
	function handleClick() {
		// No need to prevent default as we want the hash navigation to work
		// The hash change will trigger the router in App.svelte
	}
</script>

<a
	href={hashHref}
	on:click={handleClick}
	class={`group flex items-center px-2 py-2 ${
		mobile ? 'text-base' : 'text-sm'
	} font-medium rounded-md ${
		active
			? 'text-gray-900 bg-gray-100'
			: 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
	}`}
>
	<svg
		class={`${mobile ? 'mr-4' : 'mr-3'} h-6 w-6 ${
			active ? 'text-gray-500' : 'text-gray-400 group-hover:text-gray-500'
		}`}
		fill="none"
		viewBox="0 0 24 24"
		stroke="currentColor"
	>
		{@html icon}
	</svg>
	{text}
</a> 