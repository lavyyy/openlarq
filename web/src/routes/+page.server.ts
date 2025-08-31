import { getHydrationGoal, getLiquidIntake, getUserInfo } from '$lib/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	// calculate date range
	const today = new Date();
	const oneYearAgo = new Date();
	oneYearAgo.setFullYear(today.getFullYear() - 1);

	// format dates
	const endTime = today.toISOString().split('T')[0];
	const startTime = oneYearAgo.toISOString().split('T')[0];

	const [liquidIntake, hydrationGoal, userInfo] = await Promise.all([
		getLiquidIntake({
			startTime,
			endTime
		}),
		getHydrationGoal({
			index: 'time',
			viewFrom: 'right'
		}),
		getUserInfo()
	]);

	return {
		liquidIntake,
		hydrationGoal,
		userInfo
	};
};
