<script lang="ts">
	import ProgressCard from '$lib/components/progress.svelte';
	import StreakTracker from '$lib/components/streak-tracker.svelte';
	import Hydration from '$lib/components/hydration.svelte';
	import RecentActivity from '$lib/components/activity.svelte';
	import ThemeToggle from '$lib/components/theme-toggle.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const liquidIntake = data.liquidIntake;
	const hydrationGoal = data.hydrationGoal;

	const intakeEntries = liquidIntake.entries;

	// calculate total intake for today
	const today = new Date();
	today.setHours(0, 0, 0, 0);

	const todayIntake = intakeEntries
		.filter((entry) => {
			const entryDate = new Date(entry.time);
			const startOfDay = new Date(today);
			startOfDay.setHours(0, 0, 0, 0);
			const endOfDay = new Date(today);
			endOfDay.setHours(23, 59, 59, 999);
			return entryDate >= startOfDay && entryDate <= endOfDay;
		})
		.reduce((sum, entry) => sum + entry.volumeInLiter * 1000, 0); // convert to ml

	// latest hydration goal
	const currentGoal =
		hydrationGoal.entries[hydrationGoal.entries.length - 1]?.volumeInLiter * 1000 || 2000; // default to 2000ml if no goal set

	// calculate percentage
	const percentage = Math.min(Math.round((todayIntake / currentGoal) * 100), 100);

	// calculate streaks
	const calculateStreaks = () => {
		const entriesByDate = new Map<string, number>();

		// Group entries by date and sum their volumes
		intakeEntries.forEach((entry) => {
			const date = new Date(entry.time).toDateString();
			const volume = entry.volumeInLiter * 1000; // Convert to ml
			entriesByDate.set(date, (entriesByDate.get(date) || 0) + volume);
		});

		let currentStreak = 0;
		let personalBest = 0;
		let tempStreak = 0;

		const oldestEntryDate = Array.from(entriesByDate.keys())[0];

		// check if we have any entries
		if (intakeEntries.length === 0) {
			return { currentStreak: 0, personalBest: 0 };
		}

		// oldest date in entry at midnight
		// this holds the current date we're evaluating in the loop
		let currentDatePtr = new Date(oldestEntryDate);
		currentDatePtr.setHours(0, 0, 0, 0);

		// loop through every day from the oldest entry to the current date
		while (currentDatePtr.getTime() < today.getTime()) {
			const dateStr = currentDatePtr.toDateString();
			const dayIntake = entriesByDate.get(dateStr) || 0;

			if (dayIntake >= currentGoal) {
				tempStreak++;
				if (currentDatePtr.getTime() === today.getTime()) {
					currentStreak = tempStreak;
				}
				personalBest = Math.max(personalBest, tempStreak);
			} else {
				if (currentDatePtr.getTime() === today.getTime()) {
					currentStreak = 0;
				}
			}

			currentDatePtr.setDate(currentDatePtr.getDate() + 1);
		}

		return { currentStreak, personalBest };
	};

	const { currentStreak, personalBest } = calculateStreaks();
</script>

<div class="min-h-screen bg-gray-50 transition-colors dark:!bg-gray-900">
	<div class="mx-auto max-w-7xl px-6 py-8">
		<header class="mb-8 flex items-center justify-between">
			<div class="flex">
				<ThemeToggle />
			</div>
		</header>

		<div class="grid grid-cols-1 gap-8 lg:grid-cols-12">
			<div class="space-y-8 lg:col-span-4">
				<ProgressCard {percentage} currentIntake={todayIntake} goal={currentGoal} />

				<div class="grid grid-cols-2 gap-4">
					<StreakTracker
						title="CURRENT STREAK"
						days={currentStreak}
						iconColor="bg-orange-500"
						iconName="flame"
					/>
					<StreakTracker
						title="PERSONAL BEST"
						days={personalBest}
						iconColor="bg-amber-500"
						iconName="trophy"
					/>
				</div>
			</div>

			<div class="space-y-8 lg:col-span-8">
				<Hydration
					entries={intakeEntries.map((entry) => ({
						id: entry.dateCreated,
						amount: entry.volumeInLiter,
						timestamp: new Date(entry.time)
					}))}
				/>

				<RecentActivity
					entries={intakeEntries
						.filter((entry) => {
							const entryDate = new Date(entry.time);
							return entryDate.toDateString() === today.toDateString();
						})
						.map((entry) => ({
							id: entry.dateCreated,
							amount: entry.volumeInLiter,
							timestamp: new Date(entry.time)
						}))}
				/>
			</div>
		</div>
	</div>
</div>
