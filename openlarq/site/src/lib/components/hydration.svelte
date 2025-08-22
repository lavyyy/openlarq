<script lang="ts">
	import Heatmap from './heatmap/heatmap.svelte';

	const { entries }: { entries: { id: string; amount: number; timestamp: Date }[] } = $props();

	let data = $state<{ [key: string]: number }>({});
	let year = $state<number>(2025);
	let years: number[] = $state<number[]>([]);

	// update years when entries change
	$effect(() => {
		years = Array.from(new Set(entries.map((e) => new Date(e.timestamp).getFullYear()))).sort(
			(a, b) => b - a
		);
	});

	function fillMap() {
		let map: { [key: string]: number } = {};

		// process entries and aggregate by date for the selected year
		entries.forEach((entry) => {
			const date = new Date(entry.timestamp);
			if (date.getFullYear() !== year) return;
			const dateKey = `${date.getFullYear()}-${('0' + (date.getMonth() + 1)).slice(-2)}-${('0' + date.getDate()).slice(-2)}`;

			// aggregate amounts for the same date
			if (map[dateKey]) {
				map[dateKey] += entry.amount;
			} else {
				map[dateKey] = entry.amount;
			}
		});

		data = map;
	}

	// update data when entries or year change
	$effect(() => {
		fillMap();
	});
</script>

<div class="bg-card text-card-foreground rounded-lg border shadow-sm dark:border-gray-700">
	<div class="flex flex-col space-y-1.5 p-6 pb-2">
		<div
			class="text-navy-blue text-lg font-semibold leading-none tracking-tight dark:text-blue-400"
		>
			Hydration
		</div>

		<div class="mt-2">
			<select
				bind:value={year}
				class="rounded border bg-white px-2 py-1 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-200"
			>
				{#each years as y}
					<option value={y}>{y}</option>
				{/each}
			</select>
		</div>
	</div>

	<div class="p-6 pt-4">
		<Heatmap {data} {year} />
	</div>
</div>
