export interface LiquidIntakeEntry {
	amount: number;
	time: string;
}

export interface HydrationGoal {
	goal: number;
}

export interface PageData {
	liquidIntake: LiquidIntakeEntry[];
	hydrationGoal: HydrationGoal;
}
