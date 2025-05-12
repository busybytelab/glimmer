<script lang="ts">
	export let title: string;
	export let description: string = '';
</script>

<div class="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-lg border border-gray-200 dark:border-gray-700 mb-6 w-full">
	<div class="px-4 py-5 sm:px-6 bg-gray-50 dark:bg-gray-700">
		<h3 class="text-lg leading-6 font-medium text-gray-900 dark:text-white">{title}</h3>
		{#if description}
			<p class="mt-1 text-sm text-gray-500 dark:text-gray-300">{description}</p>
		{/if}
	</div>
	<div class="border-t border-gray-200 dark:border-gray-700 w-full">
		<slot />
	</div>
</div> 