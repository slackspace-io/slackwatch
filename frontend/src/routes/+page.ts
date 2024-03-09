import { BE_BASE_URL } from "../config";

	console.log('before');
	export async function load({
		fetch
	}: {
		fetch: (input: RequestInfo, init?: RequestInit) => Promise<Response>;
	}) {
    const response = await fetch(`${BE_BASE_URL}api/pods`);
    const podsInfo = await response.json(); // Parse the response as JSON
    return { podsInfo };
	}
