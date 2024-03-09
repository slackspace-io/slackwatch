<script lang="ts">
	export let podsInfo: Array<{ name: string; timeScanned: string }> = [];
	export let data;
	console.log(podsInfo);
	console.log(data);
</script>

<main>
	<h1>Pods Information</h1>
	{#each data.podsInfo as { name, timeScanned }}
		<div>
			<article class="pod-info">
				<h2>{name}</h2>
				<time datetime="{timeScanned}">{timeScanned}</time>
			</article>
			<style>
				.pod-info {
					background-color: #f0f0f0;
					border-radius: 8px;
					padding: 20px;
					margin-bottom: 10px;
					box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
				}
				.pod-info h2 {
					color: #333;
					font-size: 1.2rem;
				}
				.pod-info time {
					font-size: 0.9rem;
					color: #666;
				}
			</style>
		</div>
	{/each}
</main>
