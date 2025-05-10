import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import path from 'path';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapter({
			pages: 'dist',
			assets: 'dist',
			fallback: 'index.html',
			precompress: false,
			strict: true
		}),
		files: {
			routes: 'src/routes',
			lib: 'src/lib',
			assets: 'src/assets'
		},
		paths: {
			base: ''
		},
		alias: {
			$lib: path.resolve('./src/lib'),
			$components: path.resolve('./src/components')
		},
		appDir: 'app',
		prerender: {
			handleHttpError: ({ path, referrer, message }) => {
				if (message.includes('Not Found')) {
					return;
				}
				throw new Error(message);
			}
		},
		// Optimize for single bundle
		inlineStyleThreshold: 0, // Inline all styles
		version: {
			name: Date.now().toString() // Force cache busting
		}
	},
	preprocess: vitePreprocess({
		postcss: {
			configFilePath: './postcss.config.cjs'
		}
	})
};

export default config;
