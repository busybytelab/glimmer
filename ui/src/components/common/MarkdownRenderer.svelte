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