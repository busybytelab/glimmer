<script lang="ts">
	export let id: string;
	export let label: string;
	export let type: string = 'text';
	export let value: string;
	export let disabled: boolean = false;
	export let required: boolean = false;
	export let placeholder: string = '';
	export let cols: string = 'col-span-6 sm:col-span-3';
</script>

<div class={`${cols} mb-4`}>
	<label for={id} class="block text-sm font-medium text-gray-700 mb-1 text-left">
		{label} {#if required}<span class="text-red-500">*</span>{/if}
	</label>
	<input 
		{type} 
		name={id} 
		{id} 
		bind:value
		{placeholder}
		{required}
		{disabled}
		class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md text-sm shadow-sm placeholder-gray-400
		focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500
		dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600 dark:placeholder-gray-500
		disabled:bg-gray-50 disabled:text-gray-500 disabled:border-gray-200 disabled:shadow-none
		{disabled ? '' : 'hover:border-gray-400'}" 
	/>
</div> 