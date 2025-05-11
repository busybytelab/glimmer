<script lang="ts">
    import type { PracticeItem } from '$lib/types';
    export let item: PracticeItem;
    export let index: number;
</script>

<h4 class="text-md font-medium text-gray-900 mb-2">
    Question {index + 1}{#if item?.id} (ID: {item.id}){/if}
</h4> 