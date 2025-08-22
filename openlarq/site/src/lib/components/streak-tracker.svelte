<script lang="ts">
	const {
		title,
		days,
		iconColor,
		iconName
	}: { title: string; days: number; iconColor: string; iconName: string } = $props();
</script>

<div
	class="rounded-xl bg-white p-4 shadow-sm transition-all hover:shadow-md dark:border-gray-700 dark:bg-gray-800"
>
	<div class="mb-2 flex items-center">
		<div class={`${iconColor} mr-3 flex h-10 w-10 items-center justify-center rounded-full`}>
			{#if iconName === 'flame'}
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2">
					<path
						d="M12 2c1 8 7 12 7 16a6 6 0 11-12 0c0-4 6-8 7-16z"
						stroke-linecap="round"
						stroke-linejoin="round"
					/>
				</svg>
			{:else}
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2">
					<path
						d="M6 9H4.5a2.5 2.5 0 0 0 0 5H6m12-5h1.5a2.5 2.5 0 0 1 0 5H18M4 22h16a2 2 0 0 0 2-2V4a2 2 0 0 0-2-2H4a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2z"
						stroke-linecap="round"
						stroke-linejoin="round"
					/>
				</svg>
			{/if}
		</div>
		<div>
			<div class="text-xs text-gray-500 dark:text-gray-400">{title}</div>
			<div class="text-navy-blue text-xl font-bold dark:text-blue-400">
				{days} day{days !== 1 ? 's' : ''}
			</div>
		</div>
	</div>
</div>
