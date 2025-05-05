<script lang="ts">
	export let title: string;
	export let description: string = '';
</script>

<div class="bg-white shadow overflow-hidden sm:rounded-lg border border-gray-200 mb-6">
	<div class="px-4 py-5 sm:px-6 bg-gray-50">
		<h3 class="text-lg leading-6 font-medium text-gray-900">{title}</h3>
		{#if description}
<p class="mt-1 max-w-2xl text-sm text-gray-500 mx-auto">{description}</p>
		{/if}
	</div>
	<div class="border-t border-gray-200">
		<slot />
	</div>
</div> 