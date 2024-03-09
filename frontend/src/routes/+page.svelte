<script lang="ts">
        console.log("before");
    export async function load({ fetch }: { fetch: (input: RequestInfo, init?: RequestInit) => Promise<Response> }) {
        console.log("inside");
      const response = await fetch("http:/localhost:8080/api/pods");
      const text = await response.text(); // Get the raw response text
      console.log("Raw response:", text);
      const podsInfo = await response.json(); // Parse the response as JSON
      console.log(podsInfo);
      return { props: { podsInfo } }; // Return podsInfo for Svelte to use
    }
    export let podsInfo: Array<{name: string, timeScanned: string}> = [];
    console.log(podsInfo);
    console.log("after");
    console.log("Script executed");
    fetch('http://localhost:8080/api/pods')
        .then(response => response.text())
        .then(text => console.log("Direct fetch response:", text))
        .catch(error => console.error("Fetch error:", error));
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