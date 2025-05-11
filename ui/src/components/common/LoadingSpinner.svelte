<script lang="ts">
  export let size: 'sm' | 'md' | 'lg' = 'md';
  export let color: string = 'primary';
  
  const sizeMap = {
    sm: 'h-8 w-8',
    md: 'h-12 w-12',
    lg: 'h-32 w-32'
  };
  
  const borderClass = color === 'primary' 
    ? 'border-primary' 
    : color === 'gray' 
      ? 'border-gray-900 dark:border-white' 
      : `border-${color}`;
</script>

<div class="animate-spin rounded-full border-t-2 border-b-2 {borderClass} {sizeMap[size]}"></div> 