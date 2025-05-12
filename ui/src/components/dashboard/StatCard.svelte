<script lang="ts">
	export let title: string;
	export let value: number | string;
	export let isLoading: boolean = false;
	export let icon: string;
</script>

<div class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
	<div class="p-5">
		<div class="flex items-center">
			<div class="flex-shrink-0">
				<svg class="h-6 w-6 text-gray-400 dark:text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					{@html icon}
				</svg>
			</div>
			<div class="ml-5 w-0 flex-1">
				<dl>
					<dt class="text-sm font-medium text-gray-500 dark:text-gray-300 truncate">
						{title}
					</dt>
					<dd class="flex items-baseline">
						<div class="text-2xl font-semibold text-gray-900 dark:text-white">
							{#if isLoading}
								<svg class="animate-spin h-5 w-5 text-gray-400 dark:text-gray-300" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
									<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
									<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
								</svg>
							{:else}
								{value}
							{/if}
						</div>
					</dd>
				</dl>
			</div>
		</div>
	</div>
</div> 