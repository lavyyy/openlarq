import type {
	GetHydrationGoalProps,
	GetHydrationGoalResponse,
	GetLiquidIntakeProps,
	GetLiquidIntakeResponse,
	GetUserInfoResponse
} from './types';
import { makeApiRequest } from './utils';

export const getLiquidIntake = async (
	props: GetLiquidIntakeProps
): Promise<GetLiquidIntakeResponse> => {
	const response = await makeApiRequest<GetLiquidIntakeResponse>({
		route: '/liquid-intake',
		requestMethod: 'GET',
		additionalData: {
			startTime: props.startTime,
			endTime: props.endTime,
			index: props.index
		}
	});

	return response.data;
};

export const getHydrationGoal = async (
	props: GetHydrationGoalProps
): Promise<GetHydrationGoalResponse> => {
	const response = await makeApiRequest<GetHydrationGoalResponse>({
		route: '/hydration-goal',
		requestMethod: 'GET',
		additionalData: {
			index: props.index,
			viewFrom: props.viewFrom
		}
	});

	return response.data;
};

export const getUserInfo = async (): Promise<GetUserInfoResponse> => {
	const response = await makeApiRequest<GetUserInfoResponse>({
		route: '/user-info',
		requestMethod: 'GET'
	});

	return response.data;
};
