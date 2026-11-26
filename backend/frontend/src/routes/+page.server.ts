import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/private';

export const load: PageServerLoad = async ({ fetch }) => {
	const apiURL = env.API_URL ?? '';

	try {
		const res = await fetch(`${apiURL}/api/articles?limit=40`);
		if (!res.ok) throw new Error(`API ${res.status}`);
		const data = await res.json();
		return { articles: data.articles ?? [], error: null };
	} catch (err) {
		console.error('Failed to load articles:', err);
		return {
			articles: [],
			error: 'Could not reach the backend. Make sure it is running on port 8080.'
		};
	}
};
