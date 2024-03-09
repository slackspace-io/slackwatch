<script lang="ts">
    export let podsInfo: Array<{ name: string; timeScanned: string }> = [];
</script>
<main>
	<h1>Pods Information</h1>
	{#each podsInfo as { name, timeScanned }}
		<div>
			<p>Image Name: {name}</p>
			<p>Time Scanned: {timeScanned}</p>
		</div>
	{/each}
</main>
