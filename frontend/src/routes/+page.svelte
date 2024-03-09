<script lang="ts">
    import { fetchPods } from '$lib/api';
    
    export async function load({ }) {
        try {
        const podsInfo = await fetchPods();
        return {
            props: {
                podsInfo
            }
        };
    } catch (error) {
        console.error('Failed to fetch pods:', error);
        return {
            props: {
                podsInfo: [{name: 'Error fetching pods', timeScanned: ''}]
            }
        };
    }
}

	export let podsInfo: Array<{name: string, timeScanned: string}> = [];
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