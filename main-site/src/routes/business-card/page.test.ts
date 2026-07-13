import { render } from 'svelte/server';
import { describe, expect, it } from 'vitest';
import Page from './+page.svelte';

describe('business card page', () => {
	it('renders contact, QR, native share, and contact download controls', () => {
		const body = render(Page).body;

		expect(body).toContain('Christo Zietsman');
		expect(body).toContain('nuphirho.dev');
		expect(body).toContain('0823266576');
		expect(body).toContain('href="https://www.linkedin.com/in/christo-zietsman"');
		expect(body).toContain('Scan for LinkedIn');
		expect(body).toContain('Share card');
		expect(body).toContain('Quick Share / AirDrop');
		expect(body).toContain('href="/business-card/contact.vcf"');
		expect(body).toContain('Add contact');
		expect(body).toContain('aria-label="Paperless card"');
	});
});
