<script lang="ts">
	import { user } from '$lib/stores';
	
	export let compact: boolean = false;
	export let onLogout: () => void;
</script>

<div class="flex-shrink-0 w-full group block">
	<div class="flex items-center">
		<div class="ml-3">
			<a 
				href="#profile" 
				class={`${compact ? 'text-sm' : 'text-base'} font-medium text-gray-700 group-hover:text-gray-900 hover:underline cursor-pointer`}
			>
				{$user?.user?.name || 'User'}
			</a>
			<button
				on:click={onLogout}
				class={`${compact ? 'text-xs' : 'text-sm'} font-medium text-gray-500 group-hover:text-gray-700 ml-0 block mt-1`}
			>
				Sign out
			</button>
		</div>
	</div>
</div> 