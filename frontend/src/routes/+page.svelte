<script lang="ts">
	export const imageUpdates: Array<{ containerName: string, currentTag: string, foundAt: string, image: string, isNewer: boolean, newTag: string, podName: string, timeScanned: string }> = []
	export let data;
</script>

<main>
	<h1>Available Updates</h1>
	<p>
		<a href="/running">View Running Containers</a>
	</p>
	{#each data.imageUpdates as { containerName, currentTag, foundAt, image, isNewer, newTag, podName, timeScanned }}
		<div class="update-details">
			<h3>{containerName} ({podName})</h3>
			<p>Current Image: {image}</p>
			<p>Current Tag: {currentTag}</p>
			<p>New Tag: {newTag}</p>
			<p>Is Newer: {isNewer}</p>
			<time datetime={foundAt}>Found At: {foundAt}</time>
		</div>
	{/each}
</main>

<style>
	.update-details {
		background-color: #f0f0f0;
		border-left: 5px solid #4caf50;
		padding: 10px 20px;
		margin: 10px 0;
		border-radius: 4px;
	}
	.update-details h3 {
		margin: 0;
		color: #333;
		font-size: 1.2rem;
	}
	.update-details p,
	.update-details time {
		margin: 5px 0;
		color: #666;
		font-size: 1rem;
	}
</style>
