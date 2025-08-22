import { getHydrationGoal, getLiquidIntake } from '$lib/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	const liquidIntake = await getLiquidIntake({});

	const hydrationGoal = await getHydrationGoal({
		index: 'time',
		viewFrom: 'right'
	});

	return {
		liquidIntake,
		hydrationGoal
	};
};
