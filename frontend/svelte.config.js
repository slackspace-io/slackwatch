import adapter from '@sveltejs/adapter-node';

import sveltePreprocess from 'svelte-preprocess';

// If you're using vitePreprocess (mentioned in the provided snippet),
// ensure you're importing it correctly. Typically, for SvelteKit, sveltePreprocess is what's needed.
// The vite-plugin-svelte is automatically used by SvelteKit, so you usually don't need to manually import it.

/** @type {import('@sveltejs/kit').Config} */
const config = {
    // Consult https://kit.svelte.dev/docs/integrations#preprocessors
    // for more information about preprocessors
    preprocess: sveltePreprocess({
        typescript: {
            tsconfigFile: './tsconfig.json', // Adjusted path
        },
    }),

    kit: {
        // adapter-auto only supports some environments, see https://kit.svelte.dev/docs/adapter-auto for a list.
        // If your environment is not supported or you settled on a specific environment, switch out the adapter.
        // See https://kit.svelte.dev/docs/adapters for more information about adapters.
        adapter: adapter()
    }
};

export default config;
