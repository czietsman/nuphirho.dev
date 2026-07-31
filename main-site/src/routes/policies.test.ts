import { render } from 'svelte/server';
import { describe, expect, it } from 'vitest';
import CookiesPage from './cookies/+page.svelte';
import PrivacyPage from './privacy/+page.svelte';

describe('privacy policy', () => {
	it('describes the personal data and request metadata that may be processed', () => {
		const body = render(PrivacyPage).body;

		expect(body).toContain('31 July 2026');
		expect(body).toContain('email address');
		expect(body).toContain('Google Fonts');
		expect(body).toContain('IP address');
		expect(body).toContain('page path and an aggregate view count');
		expect(body).not.toContain('We do not collect personal data.');
	});

	it('discloses the locally stored theme preference', () => {
		const body = render(PrivacyPage).body;

		expect(body).toContain('local storage');
		expect(body).toContain('theme preference');
	});
});

describe('cookie policy', () => {
	it('distinguishes cookies, local storage, and the cookie-free visitor counter', () => {
		const body = render(CookiesPage).body;

		expect(body).toContain('31 July 2026');
		expect(body).toContain('may set');
		expect(body).toContain('Local storage');
		expect(body).toContain('theme preference');
		expect(body).toContain('does not use cookies');
	});
});
