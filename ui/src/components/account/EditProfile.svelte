<script lang="ts">
	import { error } from '$lib/stores';
	import { createEventDispatcher } from 'svelte';
	import ErrorAlert from '../common/ErrorAlert.svelte';
	import FormField from '../common/FormField.svelte';
	import FormButton from '../common/FormButton.svelte';
	import FormSection from '../common/FormSection.svelte';

	export let profile: { name: string; email: string } | null = null;
	export let isLoading: boolean = false;

	let name = '';
	let email = '';

	$: if (profile) {
		name = profile.name;
		email = profile.email;
	}

	const dispatch = createEventDispatcher();

	function handleSubmit() {
		dispatch('save', { name });
	}
</script>

{#if $error}
	<div class="px-6 py-4">
		<ErrorAlert message={$error} />
	</div>
{/if}

<div class="px-6 pb-6">
	<form on:submit|preventDefault={handleSubmit}>
		<div class="py-4">
			<div class="grid grid-cols-6 gap-6">
				<FormField 
					id="name"
					label="Name"
					bind:value={name}
					disabled={isLoading}
					required={true}
				/>

				<FormField 
					id="email"
					label="Email address (cannot be changed)"
					type="email"
					value={email}
					cols="col-span-6 sm:col-span-4"
					disabled={true}
					required={true}
				/>
			</div>
		</div>
		<div class="flex justify-end mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
			<FormButton 
				type="submit" 
				{isLoading}
				loadingText="Saving..."
			>
				Save
			</FormButton>
		</div>
	</form>
</div> 