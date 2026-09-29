import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';
import { handle } from './hooks.server.js';

describe('public response security', () => {
	it('sets the browser security baseline', async () => {
		const response = await handle({
			event: {} as never,
			resolve: async () => new Response('ok')
		});

		expect(response.headers.get('strict-transport-security')).toContain('max-age=63072000');
		expect(response.headers.get('content-security-policy')).toContain("frame-ancestors 'none'");
		expect(response.headers.get('content-security-policy')).toContain('https://fonts.googleapis.com');
		expect(response.headers.get('content-security-policy')).toContain('https://fonts.gstatic.com');
		expect(response.headers.get('permissions-policy')).toBe('camera=(), microphone=(), geolocation=()');
	});

	it('publishes the same baseline for static responses', () => {
		const headers = readFileSync('_headers', 'utf8');

		expect(headers).toContain('Strict-Transport-Security: max-age=63072000; includeSubDomains; preload');
		expect(headers).toContain("frame-ancestors 'none'");
		expect(headers).toContain('X-Frame-Options: DENY');
		expect(headers).toContain('Permissions-Policy: camera=(), microphone=(), geolocation=()');
	});

	it('publishes a current security contact', () => {
		const security = readFileSync('static/.well-known/security.txt', 'utf8');

		expect(security).toContain('Contact: mailto:contact@nuphirho.dev');
		expect(security).toContain('Expires: 2027-09-29T00:00:00Z');
		expect(security).toContain('Canonical: https://blog.nuphirho.dev/.well-known/security.txt');
	});
});
