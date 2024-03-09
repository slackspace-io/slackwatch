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
<script context="module" lang="ts">
        console.log("before");
    export async function load({ fetch }: { fetch: (input: RequestInfo, init?: RequestInit) => Promise<Response> }) {
        console.log("inside");
      const response = await fetch('http://localhost:8080/api/pods');
      const data = await response.json();
      return { props: { data } };
    }
    export let podsInfo: Array<{name: string, timeScanned: string}> = [];
    console.log(podsInfo);
    console.log("after");
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