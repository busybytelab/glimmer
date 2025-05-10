<script lang="ts">
    import { onMount, afterUpdate } from 'svelte';
    import Prism from 'prismjs';
    // Import default Prism CSS
    import 'prismjs/themes/prism.css';
    // Import common languages
    import 'prismjs/components/prism-markdown';
    import 'prismjs/components/prism-javascript';
    import 'prismjs/components/prism-typescript';
    import 'prismjs/components/prism-css';
    import 'prismjs/components/prism-json';

    export let id: string;
    export let label: string;
    export let value: string;
    export let disabled: boolean = false;
    export let required: boolean = false;
    export let placeholder: string = '';
    export let minRows: number = 3;
    export let maxRows: number = 20;
    export let cols: string = 'col-span-6';
    export let language: string = ''; // Language for syntax highlighting: 'markdown', 'javascript', etc.
    
    let textAreaRef: HTMLTextAreaElement;
    let preElement: HTMLPreElement;
    let codeElement: HTMLElement;
    let expandable = false;
    let expanded = false;
    let lineHeight: number;
    
    // Update the syntax highlighting
    function updateSyntaxHighlighting() {
        if (!codeElement) return;
        
        // Handle final newlines (adding a space to force display of empty line)
        let textToHighlight = value || '';
        if (textToHighlight[textToHighlight.length - 1] === '\n') {
            textToHighlight += ' ';
        }
        
        // Escape HTML special characters to prevent HTML rendering
        textToHighlight = textToHighlight
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;');
            
        codeElement.innerHTML = textToHighlight;
        
        // Apply highlighting if language is specified
        if (language && typeof Prism !== 'undefined') {
            Prism.highlightElement(codeElement);
        }
    }
    
    // Sync scroll position between textarea and pre/code
    function syncScroll() {
        if (preElement && textAreaRef) {
            preElement.scrollTop = textAreaRef.scrollTop;
            preElement.scrollLeft = textAreaRef.scrollLeft;
        }
    }
    
    // Adjust height of textarea based on content
    function adjustHeight() {
        if (!textAreaRef) return;
        
        // Reset height to get proper scrollHeight
        textAreaRef.style.height = 'auto';
        
        // Calculate proper line height if not set already
        if (!lineHeight) {
            lineHeight = parseInt(getComputedStyle(textAreaRef).lineHeight) || 20;
        }
        
        // Calculate min/max heights
        const minHeight = minRows * lineHeight;
        const maxHeight = maxRows * lineHeight;
        
        // Get the scroll height of the content
        const scrollHeight = textAreaRef.scrollHeight;
        
        // Determine if content exceeds max height
        expandable = scrollHeight > maxHeight;
        
        // Set height based on content and constraints
        let newHeight = Math.max(minHeight, scrollHeight);
        if (!expanded) {
            newHeight = Math.min(newHeight, maxHeight);
        }
        
        textAreaRef.style.height = `${newHeight}px`;
        if (preElement) {
            preElement.style.height = `${newHeight}px`;
        }
    }
    
    // Toggle expanded state
    function toggleExpand() {
        expanded = !expanded;
        adjustHeight();
    }
    
    // Handle tab key for indentation
    function handleTab(event: KeyboardEvent) {
        if (event.key === 'Tab') {
            event.preventDefault();
            
            const start = textAreaRef.selectionStart;
            const end = textAreaRef.selectionEnd;
            
            // For multi-line selections
            if (start !== end) {
                const selectedText = value.substring(start, end);
                const lines = selectedText.split('\n');
                
                // Process each line
                if (event.shiftKey) {
                    // Un-indent: remove up to 4 spaces from the beginning of each line
                    const processedLines = lines.map(line => line.replace(/^( {1,4})/, ''));
                    
                    // Calculate removed spaces for cursor position adjustment
                    const spacesRemoved = lines.reduce((total, line) => {
                        const spaces = line.match(/^( {1,4})/);
                        return total + (spaces ? spaces[1].length : 0);
                    }, 0);
                    
                    // Update value
                    const newValue = 
                        value.substring(0, start) + 
                        processedLines.join('\n') + 
                        value.substring(end);
                    
                    value = newValue;
                    
                    // Adjust selection
                    setTimeout(() => {
                        textAreaRef.selectionStart = start;
                        textAreaRef.selectionEnd = end - spacesRemoved;
                    }, 0);
                } else {
                    // Indent: add 4 spaces at the beginning of each line
                    const processedLines = lines.map(line => '    ' + line);
                    
                    // Calculate added spaces for cursor position adjustment
                    const spacesAdded = lines.length * 4;
                    
                    // Update value
                    const newValue = 
                        value.substring(0, start) + 
                        processedLines.join('\n') + 
                        value.substring(end);
                    
                    value = newValue;
                    
                    // Adjust selection
                    setTimeout(() => {
                        textAreaRef.selectionStart = start;
                        textAreaRef.selectionEnd = end + spacesAdded;
                    }, 0);
                }
            } else {
                // For single cursor position (no selection)
                const newValue = 
                    value.substring(0, start) + 
                    '    ' + 
                    value.substring(end);
                
                value = newValue;
                
                // Adjust cursor position
                setTimeout(() => {
                    textAreaRef.selectionStart = textAreaRef.selectionEnd = start + 4;
                }, 0);
            }
        }
    }
    
    // Handle Enter key to maintain indentation
    function handleEnter(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.preventDefault();
            
            const cursorPos = textAreaRef.selectionStart;
            
            // Get the current line (up to the cursor)
            const currentLine = value.substring(0, cursorPos).split('\n').pop() || '';
            
            // Extract indentation with null check
            const indentMatch = currentLine.match(/^\s*/);
            const indentation = indentMatch ? indentMatch[0] : '';
            
            // Insert new line with indentation
            const newValue = 
                value.substring(0, cursorPos) + 
                '\n' + indentation + 
                value.substring(cursorPos);
            
            value = newValue;
            
            // Adjust cursor position
            setTimeout(() => {
                textAreaRef.selectionStart = textAreaRef.selectionEnd = 
                    cursorPos + 1 + indentation.length;
            }, 0);
        }
    }
    
    // Handle all keyboard events
    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Tab') {
            handleTab(event);
        } else if (event.key === 'Enter') {
            handleEnter(event);
        }
    }
    
    // Initialize and update on content changes
    function handleInput() {
        updateSyntaxHighlighting();
        adjustHeight();
    }
    
    onMount(() => {
        // Initial setup
        adjustHeight();
        updateSyntaxHighlighting();
        
        // Run this outside Svelte's update cycle to ensure DOM is ready
        setTimeout(() => {
            adjustHeight();
            updateSyntaxHighlighting();
        }, 0);
    });
    
    afterUpdate(() => {
        updateSyntaxHighlighting();
    });
    
    // Watch for content changes
    $: if (value !== undefined && textAreaRef) {
        setTimeout(() => {
            adjustHeight();
            updateSyntaxHighlighting();
        }, 0);
    }
</script>

<div class={`${cols} ${label ? 'mb-4' : ''}`}>
    {#if label}
        <label for={id} class="block text-sm font-medium text-gray-700 mb-1 text-left">
            {label} {#if required}<span class="text-red-500">*</span>{/if}
        </label>
    {/if}
    
    <div class="editor-container relative mt-1">
        <!-- Syntax highlighting pre/code -->
        <pre 
            bind:this={preElement}
            class="syntax-highlight-pre"
            aria-hidden="true"
        ><code 
            bind:this={codeElement}
            class={`language-${language || 'plaintext'}`}
        ></code></pre>
        
        <!-- Actual textarea for input -->
        <textarea 
            bind:this={textAreaRef}
            bind:value
            {id}
            {placeholder}
            {disabled}
            {required}
            name={id}
            class="syntax-highlight-textarea"
            spellcheck="false"
            autocomplete="off"
            on:input={handleInput}
            on:keydown={handleKeyDown}
            on:scroll={syncScroll}
        ></textarea>
        
        <!-- Expand/Collapse button -->
        {#if expandable}
            <button 
                type="button" 
                class="absolute bottom-2 right-2 text-xs font-medium px-3 py-1.5 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-colors duration-200 z-20"
                style="backdrop-filter: blur(2px); background-color: rgba(243, 244, 246, 0.85);"
                on:click={toggleExpand}
                tabindex="-1"
            >
                {expanded ? 'Collapse' : 'Expand'}
            </button>
        {/if}
    </div>
</div>

<style>
    .editor-container {
        position: relative;
        width: 100%;
        min-height: 100px;
        border: 1px solid #d1d5db;
        border-radius: 0.375rem;
        overflow: hidden;
    }
    
    .syntax-highlight-pre {
        margin: 0;
        padding: 0.75rem 1rem;
        width: 100%;
        height: 100%;
        background-color: white;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
        font-size: 0.875rem;
        line-height: 1.5;
        tab-size: 4;
        white-space: pre;
        overflow-x: auto;
        overflow-y: hidden;
    }
    
    .syntax-highlight-pre code {
        font-family: inherit;
        font-size: inherit;
        line-height: inherit;
        tab-size: inherit;
        white-space: inherit;
        display: block;
    }
    
    .syntax-highlight-textarea {
        position: absolute;
        top: 0;
        left: 0;
        margin: 0;
        padding: 0.75rem 1rem;
        width: 100%;
        height: 100%;
        border: none;
        background: transparent;
        color: transparent;
        caret-color: #1f2937;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
        font-size: 0.875rem;
        line-height: 1.5;
        tab-size: 4;
        white-space: pre;
        overflow-x: auto;
        overflow-y: hidden;
        resize: none;
        z-index: 10;
    }
    
    .syntax-highlight-textarea:focus {
        outline: none;
        box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.5);
    }
</style> 