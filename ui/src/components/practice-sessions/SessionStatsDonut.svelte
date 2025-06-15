<script lang="ts">
    import type { PracticeSessionStats } from '$lib/types';

    /**
     * The session statistics to display
     */
    export let stats: PracticeSessionStats | null = null;

    /**
     * Whether to show the legend below the chart
     */
    export let showLegend: boolean = true;

    // Calculate circle segments (each segment is a percentage of the total)
    $: remainingPercent = stats ? ((stats.total_items - stats.answered_items) / stats.total_items) * 100 : 100;
    $: correctPercent = stats ? ((stats.answered_items - stats.wrong_answers_count) / stats.total_items) * 100 : 0;
    $: wrongPercent = stats ? (stats.wrong_answers_count / stats.total_items) * 100 : 0;

    // SVG circle properties
    const radius = 40;
    const circumference = 2 * Math.PI * radius;
    const strokeWidth = 12;

    // Calculate stroke-dasharray and stroke-dashoffset for each segment
    $: segments = stats ? [
        {
            color: 'stroke-slate-200 dark:stroke-slate-700', // Not answered - light gray
            percent: remainingPercent,
            offset: 0
        },
        {
            color: 'stroke-emerald-500 dark:stroke-emerald-400', // Correct - green
            percent: correctPercent,
            offset: (remainingPercent * circumference) / 100
        },
        {
            color: 'stroke-rose-500 dark:stroke-rose-400', // Wrong - red
            percent: wrongPercent,
            offset: ((remainingPercent + correctPercent) * circumference) / 100
        }
    ] : [];
</script>

{#if stats}
    <div class="relative w-32 h-32">
        <!-- Donut Chart -->
        <svg class="w-full h-full -rotate-90" viewBox="0 0 100 100">
            <!-- Background circle -->
            <circle
                cx="50"
                cy="50"
                r={radius}
                fill="none"
                class="stroke-slate-100 dark:stroke-slate-800"
                stroke-width={strokeWidth}
            />
            
            <!-- Progress segments -->
            {#each segments as segment}
                {#if segment.percent > 0}
                    <circle
                        cx="50"
                        cy="50"
                        r={radius}
                        fill="none"
                        class={segment.color}
                        stroke-width={strokeWidth}
                        stroke-dasharray={`${(segment.percent * circumference) / 100} ${circumference}`}
                        stroke-dashoffset={-segment.offset}
                        stroke-linecap="round"
                    />
                {/if}
            {/each}

            <!-- Score text in the center (rotated back to normal) -->
            <g transform="rotate(90 50 50)">
                <text
                    x="50"
                    y="50"
                    text-anchor="middle"
                    dominant-baseline="middle"
                    class="fill-gray-900 dark:fill-white"
                >
                    <tspan class="font-bold" style="font-size: 28px">{stats.total_score}</tspan><!--
                    --><tspan style="font-size: 16px">%</tspan>
                </text>
            </g>
        </svg>

        <!-- Legend -->
        {#if showLegend}
            <div class="absolute -bottom-6 left-1/2 -translate-x-1/2 flex items-center gap-2 text-[10px]">
                <div class="flex items-center gap-1">
                    <div class="w-2 h-2 rounded-full bg-emerald-500 dark:bg-emerald-400"></div>
                    <span class="text-gray-600 dark:text-gray-400">Correct</span>
                </div>
                <div class="flex items-center gap-1">
                    <div class="w-2 h-2 rounded-full bg-rose-500 dark:bg-rose-400"></div>
                    <span class="text-gray-600 dark:text-gray-400">Wrong</span>
                </div>
            </div>
        {/if}
    </div>
{/if}