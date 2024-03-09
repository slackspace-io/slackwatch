	console.log('before');
	export async function load({
		fetch
	}: {
		fetch: (input: RequestInfo, init?: RequestInit) => Promise<Response>;
	}) {
		console.log('inside');
		const response = await fetch('http://localhost:8080/api/pods');
		const text = await response.text(); // Get the raw response text
		console.log('Raw response:', text);
		const podsInfo = await response.json(); // Parse the response as JSON
		console.log(podsInfo);
		return { props: { podsInfo } }; // Return podsInfo for Svelte to use
	}
