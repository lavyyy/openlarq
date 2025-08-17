<script lang="ts">
	import { calMonths, getCalendar, getColor, getColSpanForMonth } from './utils.js';
	import type { Props } from './utils.js';
	import { theme } from '$lib/stores/theme.js';

	let {
		data,
		onclick,
		onmouseout,
		onmouseover,
		onfocus,
		onblur,
		colors = {
			dark: ['#151b23', '#033a16', '#196c2e', '#2ea043', '#56d364'],
			light: ['#eff2f5', '#aceebb', '#4ac26b', '#2da44e', '#116329']
		},
		year = new Date().getFullYear(),
		lday = true,
		lmonth = true
	}: Props = $props();

	let { max, calendar } = $derived(getCalendar(data, year));
</script>

<div class="max-w-full overflow-x-auto overflow-y-hidden">
	<table class="relative w-max border-separate border-spacing-[3px] overflow-hidden">
		{#if lmonth}
			<thead>
				<tr class="h-[10px]">
					<td class="h-[10px]"></td>
					{#each calMonths as m}
						<td colspan={getColSpanForMonth(m)} class="h-[10px] text-[12px]">{m}</td>
					{/each}
				</tr>
			</thead>
		{/if}
		<tbody>
			{#each calendar as w, i}
				<tr class="h-[10px]">
					{#if lday}
						<td class="relative h-[10px] pr-[8px] align-top text-[12px] leading-[10px]">
							{['', 'Mon', '', 'Wed', '', 'Fri', ''][i]}
						</td>
					{/if}
					{#each w as d}
						{#if d}
							<td
								class="h-[10px] w-[10px] rounded-[2px]"
								style={`background:${getColor(colors, max, d.value, $theme)}`}
								data-date={d.date}
								data-value={d.value}
								{onclick}
								{onmouseout}
								{onmouseover}
								{onfocus}
								{onblur}
							></td>
						{:else}
							<td></td>
						{/if}
					{/each}
				</tr>
			{/each}
		</tbody>
	</table>
</div>
