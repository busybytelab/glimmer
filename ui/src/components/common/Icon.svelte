<script lang="ts">
    import type { IconType } from '$lib/types';
    import { IconTypeMap } from '$lib/types';

    export let type: IconType;
    export let class_name: string = '';
</script>

<svg 
    xmlns="http://www.w3.org/2000/svg" 
    fill="none" 
    viewBox="0 0 24 24" 
    stroke-width="1.5" 
    stroke="currentColor"
    class={class_name}
>
    <path stroke-linecap="round" stroke-linejoin="round" d={IconTypeMap[type]} />
</svg> 