import { getHydrationGoal, getLiquidIntake, getUserInfo } from '$lib/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	const liquidIntake = await getLiquidIntake({});

	const hydrationGoal = await getHydrationGoal({
		index: 'time',
		viewFrom: 'right'
	});

	const userInfo = await getUserInfo();

	return {
		liquidIntake,
		hydrationGoal,
		userInfo
	};
};
