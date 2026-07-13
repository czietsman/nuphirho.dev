import { describe, expect, it } from 'vitest';
import { GET } from './+server.js';

describe('business card vCard endpoint', () => {
	it('downloads Christo Zietsman as a vCard', async () => {
		const response = GET();

		expect(response.headers.get('content-type')).toBe('text/vcard; charset=utf-8');
		expect(response.headers.get('content-disposition')).toBe(
			'attachment; filename="christo-zietsman.vcf"',
		);
		expect(await response.text()).toContain('FN:Christo Zietsman');
	});
});
