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

<div class="py-6">
	<div class="max-w-3xl mx-auto px-4 sm:px-6 md:px-8">
		<h1 class="text-2xl font-semibold text-gray-900 mb-6">Edit Profile</h1>
	</div>
	
	{#if $error}
		<div class="max-w-3xl mx-auto px-4 sm:px-6 md:px-8 mb-6">
			<ErrorAlert message={$error} />
		</div>
	{/if}
	
	<div class="max-w-3xl mx-auto px-4 sm:px-6 md:px-8">
		<FormSection title="Profile Information" description="Update your personal details">
			<form on:submit|preventDefault={handleSubmit} class="p-4">
				<div class="py-4 bg-white">
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
				<div class="px-4 py-3 bg-gray-50 text-right sm:px-6 mt-4 -mx-4 -mb-4 border-t border-gray-200">
					<FormButton 
						type="submit" 
						{isLoading}
						loadingText="Saving..."
					>
						Save
					</FormButton>
				</div>
			</form>
		</FormSection>
	</div>
</div> 