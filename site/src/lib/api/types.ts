export interface MakeApiRequestProps {
	/** The route to request on the API */
	route: string;

	/** The request method to be used */
	requestMethod: 'GET' | 'POST';

	/** Extra data to be submitted either in the body or as query params */
	additionalData?: object;

	/** Headers to be sent with the request */
	headers?: object;
}

export interface MakeApiRequestResponse<T> {
	/** The data returned from the Api */
	data: T;

	/** The status of the response */
	status: number;
}

export interface GetLiquidIntakeProps {
	startTime?: string;
	endTime?: string;
	index?: string;
}

export interface GetLiquidIntakeResponse {
	entries: {
		dateCreated: string;
		source: string;
		time: string;
		type: string;
		volumeInLiter: number;
	}[];
}

export interface GetHydrationGoalProps {
	index: string;
	viewFrom: string;
}

export interface GetHydrationGoalResponse {
	entries: {
		time: string;
		type: string;
		volumeInLiter: number;
	}[];
}
