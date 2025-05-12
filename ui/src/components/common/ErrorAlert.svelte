<script lang="ts">
	export let message: string;
	export let title: string = 'Error!';
</script>

<div class="bg-red-100 dark:bg-red-900/30 border border-red-400 dark:border-red-800 text-red-700 dark:text-red-300 px-4 py-3 rounded relative" role="alert">
	<strong class="font-bold">{title}</strong>
	<span class="block sm:inline">{message}</span>
</div> 