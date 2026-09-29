import { describe, expect, it } from 'vitest';
import { GET } from './+server.js';

describe('sitemap', () => {
	it('includes every public top-level page', async () => {
		const response = GET({} as Parameters<typeof GET>[0]);
		const body = await response.text();

		for (const route of ['/business-card', '/papers', '/privacy', '/cookies']) {
			expect(body).toContain(`<loc>https://nuphirho.dev${route}</loc>`);
		}
	});
});
