<script lang="ts">
    /**
     * Comma-separated string of tags
     */
    export let value: string = '';
    
    /**
     * Input field ID
     */
    export let id: string;
    
    /**
     * Field label
     */
    export let label: string = 'Tags';
    
    /**
     * Placeholder text
     */
    export let placeholder: string = 'Enter tags separated by commas';
    
    /**
     * Whether the field is disabled
     */
    export let disabled: boolean = false;
    
    /**
     * Whether the field is required
     */
    export let required: boolean = false;
    
    /**
     * CSS classes for the container
     */
    export let cols: string = 'col-span-6';
    
    let inputValue = value;
    
    // Parse tags from the value string
    $: tags = value.split(',').map(tag => tag.trim()).filter(Boolean);
    
    // Update the bound value when input changes
    $: value = inputValue;
    
    /**
     * Remove a tag at the specified index
     */
    function removeTag(index: number) {
        const newTags = tags.filter((_, i) => i !== index);
        inputValue = newTags.join(', ');
    }
    
    /**
     * Handle keyboard events for better UX
     */
    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.preventDefault();
            const input = event.target as HTMLInputElement;
            input.blur();
            input.focus();
        }
    }
</script>

<div class={`${cols} mb-4`}>
    <label for={id} class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1 text-left">
        {label} {#if required}<span class="text-red-500">*</span>{/if}
    </label>
    
    <div class="space-y-3">
        <input
            type="text"
            {id}
            name={id}
            bind:value={inputValue}
            {placeholder}
            {required}
            {disabled}
            on:keydown={handleKeyDown}
            class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
        />
        
        {#if tags.length > 0}
            <div class="flex flex-wrap gap-2">
                {#each tags as tag, index}
                    <span class="inline-flex items-center gap-1 bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-300 text-xs font-medium px-2.5 py-1 rounded-full">
                        <span>{tag}</span>
                        {#if !disabled}
                            <button
                                type="button"
                                class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-200 ml-1 focus:outline-none"
                                on:click={() => removeTag(index)}
                                title="Remove tag: {tag}"
                            >
                                <svg class="h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                                </svg>
                            </button>
                        {/if}
                    </span>
                {/each}
            </div>
        {/if}
        
        {#if tags.length === 0 && !value}
            <p class="text-xs text-gray-500 dark:text-gray-400">
                Enter tags separated by commas. Tags help categorize and search for practice items.
            </p>
        {/if}
    </div>
</div> 