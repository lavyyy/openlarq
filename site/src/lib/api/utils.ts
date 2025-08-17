import { API_URL } from '$env/static/private';
import axios from 'axios';
import type { MakeApiRequestProps, MakeApiRequestResponse } from './types';

export const makeApiRequest = async <T>(
	props: MakeApiRequestProps
): Promise<MakeApiRequestResponse<T>> => {
	const method = props.requestMethod;
	const route = props.route;
	const headers = props.headers;
	const data = props.additionalData;

	if (method === 'POST') {
		const response = await axios.post(`${API_URL}${route}`, data, { headers });

		return { data: response.data, status: response.status };
	}

	const response = await axios.get(`${API_URL}/${route}`, {
		params: data,
		headers
	});

	return { data: response.data, status: response.status };
};
