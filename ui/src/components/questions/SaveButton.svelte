<script lang="ts">
    export let disabled = false;
    export let onClick: () => void;
</script>

<div class="mt-4">
    <button
        class="inline-flex items-center px-6 py-3 min-h-[48px] min-w-[48px] border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 active:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed"
        {disabled}
        on:click={onClick}
    >
        Save
    </button>
</div> 