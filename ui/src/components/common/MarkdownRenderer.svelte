<!-- ui/src/components/common/MarkdownRenderer.svelte -->
<script lang="ts">
  import { marked } from 'marked';
  import DOMPurify from 'dompurify';
  import { onMount } from 'svelte';

  /**
   * The markdown content to render
   */
  export let content: string;

  /**
   * Whether to sanitize the HTML output
   * @default true
   */
  export let sanitize: boolean = true;

  /**
   * Additional CSS classes to apply to the rendered content
   */
  export let className: string = '';

  let container: HTMLElement;

  onMount(() => {
    // Configure marked options
    marked.setOptions({
      breaks: true, // Convert line breaks to <br>
      gfm: true, // GitHub Flavored Markdown
    });
  });

  $: {
    if (container) {
      // Use marked.parse() which returns a string directly
      const html = marked.parse(content, { async: false });
      container.innerHTML = sanitize ? DOMPurify.sanitize(html) : html;
    }
  }
</script>

<div 
  bind:this={container} 
  class="prose dark:prose-invert max-w-none {className}"
  data-testid="markdown-content"
></div>

<style>
  /* Base styles for markdown content */
  :global(.prose) {
    @apply text-gray-700;
  }

  :global(.prose pre) {
    @apply bg-gray-50 dark:bg-gray-800 rounded-lg p-4 overflow-x-auto;
  }

  :global(.prose code) {
    @apply bg-gray-50 dark:bg-gray-800 px-1 py-0.5 rounded text-sm;
  }

  :global(.prose pre code) {
    @apply bg-transparent p-0 text-gray-800 dark:text-gray-200;
  }

  :global(.prose a) {
    @apply text-indigo-600 dark:text-indigo-400 hover:text-indigo-800 dark:hover:text-indigo-300;
  }

  :global(.prose blockquote) {
    @apply border-l-4 border-gray-200 dark:border-gray-700 pl-4 italic;
  }

  :global(.prose table) {
    @apply border-collapse w-full;
  }

  :global(.prose th) {
    @apply border border-gray-200 dark:border-gray-700 px-4 py-2 bg-gray-50 dark:bg-gray-800;
  }

  :global(.prose td) {
    @apply border border-gray-200 dark:border-gray-700 px-4 py-2;
  }

  /* Fix for dark mode text color */
  :global(.dark .prose) {
    @apply text-gray-300;
  }

  /* Fix for dark mode heading color */
  :global(.dark .prose h1),
  :global(.dark .prose h2),
  :global(.dark .prose h3),
  :global(.dark .prose h4),
  :global(.dark .prose h5),
  :global(.dark .prose h6) {
    @apply text-blue-200;
  }

  /* Fix for dark mode list items and strong tags */
  :global(.dark .prose li),
  :global(.dark .prose strong),
  :global(.dark .prose b),
  :global(.dark .prose dt),
  :global(.dark .prose dd) {
    @apply text-blue-100;
  }

  /* Fix for dark mode list markers */
  :global(.dark .prose ul li::before) {
    @apply bg-blue-200;
  }

  :global(.dark .prose ol li::before) {
    @apply text-blue-200;
  }

  /* Light mode heading color */
  :global(.prose h1),
  :global(.prose h2),
  :global(.prose h3),
  :global(.prose h4),
  :global(.prose h5),
  :global(.prose h6) {
    @apply text-gray-900;
  }

  /* Light mode list items and strong tags */
  :global(.prose li),
  :global(.prose strong),
  :global(.prose b),
  :global(.prose dt),
  :global(.prose dd) {
    @apply text-gray-800;
  }

  /* Light mode list markers */
  :global(.prose ul li::before) {
    @apply bg-gray-400;
  }

  :global(.prose ol li::before) {
    @apply text-gray-400;
  }
</style> 