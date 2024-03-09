<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchPods } from '$lib/api';

	let podsInfo: Array<{name: string, timeScanned: string}> = [];

	onMount(async () => {
		try {
			podsInfo = await fetchPods();
		} catch (error) {
			console.error('Failed to fetch pods:', error);
			podsInfo = [{name: 'Error fetching pods', timeScanned: ''}];
		}
	});
</script>

<main>
	<h1>Pods Information</h1>
	{#each podsInfo as {name, timeScanned}}
	    <div>
	        <p>Image Name: {name}</p>
	        <p>Time Scanned: {timeScanned}</p>
	    </div>
	{/each}
</main>
