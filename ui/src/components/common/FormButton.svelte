<script lang="ts">
	export let type: "button" | "submit" | "reset" = "button";
	export let isLoading: boolean = false;
	export let disabled: boolean = false;
	export let variant: "primary" | "secondary" | "danger" = "primary";
	export let loadingText: string = "Loading...";

	const variants = {
		primary: "bg-indigo-600 hover:bg-indigo-700 focus:ring-indigo-500 border-indigo-600 dark:bg-indigo-700 dark:hover:bg-indigo-800 dark:border-indigo-700",
		secondary: "bg-gray-600 hover:bg-gray-700 focus:ring-gray-500 border-gray-600 dark:bg-gray-700 dark:hover:bg-gray-800 dark:border-gray-700",
		danger: "bg-red-600 hover:bg-red-700 focus:ring-red-500 border-red-600 dark:bg-red-700 dark:hover:bg-red-800 dark:border-red-700"
	};
</script>

<button 
	{type} 
	class={`inline-flex justify-center py-2 px-5 border shadow-sm text-sm font-medium rounded-md text-white 
	${variants[variant]} 
	focus:outline-none focus:ring-2 focus:ring-offset-2 dark:focus:ring-offset-gray-800
	transition-all duration-200 ease-in-out
	disabled:opacity-50 disabled:cursor-not-allowed
	`}
	disabled={disabled || isLoading}
	on:click
>
	{#if isLoading}
		<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
			<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
			<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
		</svg>
		{loadingText}
	{:else}
		<slot>Submit</slot>
	{/if}
</button> 