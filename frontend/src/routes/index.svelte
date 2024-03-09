<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchPods } from '$lib/api';

	let podsInfo = '';

	onMount(async () => {
		try {
			podsInfo = await fetchPods();
		} catch (error) {
			console.error('Failed to fetch pods:', error);
			podsInfo = 'Error fetching pods';
		}
	});
</script>

<main>
	<h1>Pods Information</h1>
	<pre>{podsInfo}</pre>
</main>
