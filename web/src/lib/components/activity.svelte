<script lang="ts">
	import { format } from 'date-fns';
	const { entries }: { entries: { id: string; amount: number; timestamp: Date }[] } = $props();

	const entriesByDate = entries.map((e) => {
		return {
			id: e.id,
			amount: e.amount * 1000,
			timestamp: e.timestamp
		};
	});

	const sortedEntries = [...entriesByDate].sort((a, b) => b.timestamp.getTime() - a.timestamp.getTime());
	const uniqueDates = [...new Map(sortedEntries.map((e) => [e.timestamp.toDateString(), e.timestamp])).values()].sort(
		(a, b) => b.getTime() - a.getTime()
	);
</script>

<div class="bg-card text-card-foreground rounded-lg border shadow-sm dark:border-gray-700">
	<div class="flex flex-col space-y-1.5 p-6 pb-2">
		<div
			class="text-navy-blue text-lg font-semibold leading-none tracking-tight dark:text-blue-400"
		>
			Recent Activity
		</div>
	</div>

	<div class="p-6 pt-0">
		{#if Object.entries(entriesByDate).length > 0}
			{#each uniqueDates as date}
				<div class="mb-4">
					<h3 class="mb-2 text-sm font-medium text-gray-500 dark:text-gray-400">
						{date.toDateString()}
					</h3>
					<div class="space-y-2">
						{#each entriesByDate
								.filter((e) => e.timestamp.toDateString() === date.toDateString())
								.sort((a, b) => b.timestamp.getTime() - a.timestamp.getTime()) as entry}
							<div
								class="flex items-center justify-between rounded-lg border border-gray-100 bg-gray-50 p-2 dark:border-gray-600 dark:bg-gray-700"
							>
								<div class="flex items-center">
									<div
										class="bg-water-medium/20 mr-2 flex h-6 w-6 items-center justify-center rounded-full dark:bg-blue-500/20"
									>
										<span class="text-water-dark text-xs dark:text-blue-400">💧</span>
									</div>
									<span class="dark:text-gray-200">{entry.amount.toFixed(0)} ml</span>
								</div>
								<span class="text-xs text-gray-400 dark:text-gray-500">
									{format(entry.timestamp, 'h:mm a')}
								</span>
							</div>
						{/each}
					</div>
				</div>
			{/each}
		{:else}
			<div class="py-6 text-center text-gray-500 dark:text-gray-400">
				No water intake recorded yet
			</div>
		{/if}
	</div>
</div>
