<script lang="ts">
	export const podsInfo: Array<{ image: string; timeScanned: string }> = [];
	export let data;
</script>

<main>
	<h1>Pods Information</h1>
	{#each data.podsInfo as { podName, containerName, image, timeScanned }}
		<div class="pod-details">
			<h3>Container: {containerName}</h3>
			<p>Image: {image}</p>
			<p>Pod: {podName}</p>
			<time datetime="{timeScanned}">Scanned: {timeScanned}</time>
		</div>
		<style>
			.pod-details {
				background-color: #e3f2fd;
				border-left: 5px solid #2196f3;
				padding: 10px 20px;
				margin: 10px 0;
				border-radius: 4px;
			}
			.pod-details h3 {
				margin: 0;
				color: #0d47a1;
				font-size: 1.1rem;
			}
			.pod-details p, .pod-details time {
				margin: 5px 0;
				color: #333;
				font-size: 0.9rem;
			}
		</style>
		<div>
		</div>
	{/each}
</main>
