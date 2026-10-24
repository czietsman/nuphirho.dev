import { svelte } from '@sveltejs/vite-plugin-svelte';
import { svelteTesting } from '@testing-library/svelte/vite';
import { defineConfig } from 'vitest/config';
import { fileURLToPath } from 'node:url';

export default defineConfig({
	plugins: [svelte(), svelteTesting()],
	resolve: {
		alias: {
			'cloudflare:workers': fileURLToPath(new URL('./vitest-cloudflare-workers.ts', import.meta.url)),
		},
	},
	test: {
		environment: 'jsdom',
		include: ['src/**/*.{test,spec}.{js,ts}'],
		setupFiles: ['./vitest-setup.ts'],
		deps: {
			optimizer: {
				web: { enabled: false },
			},
		},
	},
});
