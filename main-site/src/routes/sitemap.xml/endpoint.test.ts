import { describe, expect, it } from 'vitest';
import { GET } from './+server.js';

describe('sitemap', () => {
	it('includes the public business card', async () => {
		const response = GET({} as Parameters<typeof GET>[0]);
		expect(await response.text()).toContain(
			'<loc>https://nuphirho.dev/business-card</loc>',
		);
	});
});
