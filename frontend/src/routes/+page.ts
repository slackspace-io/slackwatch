import config  from '../config';

console.log('before');
export async function load({
  fetch
	}: {
		fetch: (input: RequestInfo, init?: RequestInit) => Promise<Response>;
	}) {
    const response = await fetch(`${config!.baseURL}api/pods`);
    const podsInfo = await response.json(); // Parse the response as JSON
    return { podsInfo };
	}
