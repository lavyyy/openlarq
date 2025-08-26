<script lang="ts">
	const {
		percentage,
		currentIntake,
		goal
	}: { percentage: number; currentIntake: number; goal: number } = $props();

	const radius = 120;
	const circumference = 2 * Math.PI * radius;
	const strokeDashoffset = circumference - (percentage / 100) * circumference;

	// convert ml to oz for display (roughly 1 oz = 29.57 ml)
	const currentOz = (currentIntake / 29.57).toFixed(1);
	const goalOz = (goal / 29.57).toFixed(0);

	// calculate remaining amount
	const remaining = Math.max(0, goal - currentIntake);
	const remainingOz = (remaining / 29.57).toFixed(1);
</script>

<div class="bg-card text-card-foreground rounded-lg border shadow-sm dark:border-gray-700">
	<!-- Header -->
	<div class="flex flex-col space-y-1.5 p-6 pb-2">
		<div
			class="text-navy-blue text-lg font-semibold leading-none tracking-tight dark:text-blue-400"
		>
			Daily Progress
		</div>
	</div>

	<!-- Progress Circle -->
	<div class="flex flex-col items-center p-6 pt-4">
		<div class="relative inline-flex items-center justify-center">
			<svg width="280" height="280" class="-rotate-90 transform">
				<circle
					cx="140"
					cy="140"
					r={radius}
					fill="none"
					stroke="#E6EEF7"
					stroke-width="26"
					stroke-linecap="round"
					class="dark:stroke-gray-700"
				/>
				<circle
					cx="140"
					cy="140"
					r={radius}
					fill="none"
					stroke="#1a5fb4"
					stroke-width="26"
					stroke-linecap="round"
					stroke-dasharray={circumference}
					stroke-dashoffset={strokeDashoffset}
					class="transition-all duration-700 ease-in-out dark:stroke-blue-500"
				/>

				{#each [25, 50, 75] as markerPercent}
					<circle
						cx={140 + radius * Math.cos(((markerPercent / 100) * 360 - 90) * (Math.PI / 180))}
						cy={140 + radius * Math.sin(((markerPercent / 100) * 360 - 90) * (Math.PI / 180))}
						r="4"
						fill="#E6EEF7"
						class={percentage >= markerPercent
							? 'fill-water-dark dark:fill-blue-500'
							: 'fill-gray-500 dark:fill-gray-300'}
					/>
				{/each}

				{#if percentage >= 100}
					<circle
						cx={140 + radius * Math.cos(0)}
						cy={140 + radius * Math.sin(0)}
						r="6"
						fill="#FFD700"
					/>
				{/if}
			</svg>

			<div class="absolute flex flex-col items-center justify-center text-center">
				<span class="text-navy-blue text-6xl font-bold dark:text-blue-400">{percentage}%</span>
				<div class="mt-2 flex flex-col gap-1 text-xl text-gray-500 dark:text-gray-400">
					<span>{currentOz} / {goalOz} oz</span>
					{#if remaining > 0}
						<span class="text-water-dark text-lg font-medium dark:text-blue-500"
							>{remainingOz} oz left</span
						>
					{/if}
				</div>
			</div>
		</div>
	</div>
</div>
