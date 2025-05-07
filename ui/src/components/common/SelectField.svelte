<script lang="ts">
	export let id: string;
	export let label: string;
	export let value: string | undefined;
	export let disabled: boolean = false;
	export let required: boolean = false;
	export let cols: string = 'col-span-6 sm:col-span-3';
	export let inline: boolean = false;
</script>

<div class={`${cols} mb-4 ${inline ? 'flex items-center gap-2' : ''}`}>
	<label for={id} class="block text-sm font-medium text-gray-700 text-left ${inline ? '' : 'mb-1'}">
		{label} {#if required}<span class="text-red-500">*</span>{/if}
	</label>
	<select 
		name={id} 
		{id} 
		bind:value
		{required}
		{disabled}
		class="mt-1 block w-full pl-3 pr-10 py-2 border border-gray-300 rounded-md text-sm shadow-sm placeholder-gray-400
		focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500
		dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600 dark:placeholder-gray-500
		disabled:bg-gray-50 disabled:text-gray-500 disabled:border-gray-200 disabled:shadow-none
		{disabled ? '' : 'hover:border-gray-400'} ${inline ? 'mt-0' : ''}
		appearance-none bg-no-repeat bg-[length:1.5em_1.5em] bg-[right_0.5rem_center] bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2020%2020%22%3E%3Cpath%20stroke%3D%22%236B7280%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%221.5%22%20d%3D%22m6%208%204%204%204-4%22%2F%3E%3C%2Fsvg%3E')]"
	>
		<slot />
	</select>
</div> 