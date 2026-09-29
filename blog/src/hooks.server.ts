import type { Handle } from '@sveltejs/kit';

const CONTENT_SECURITY_POLICY = [
	"default-src 'self'",
	"script-src 'self' 'unsafe-inline'",
	"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
	"font-src 'self' https://fonts.gstatic.com",
	"img-src 'self' data:",
	"connect-src 'self'",
	"object-src 'none'",
	"base-uri 'self'",
	"frame-ancestors 'none'",
	"form-action 'self'",
	'upgrade-insecure-requests'
].join('; ');

const SECURITY_HEADERS = {
	'Content-Security-Policy': CONTENT_SECURITY_POLICY,
	'Permissions-Policy': 'camera=(), microphone=(), geolocation=()',
	'Referrer-Policy': 'strict-origin-when-cross-origin',
	'Strict-Transport-Security': 'max-age=63072000; includeSubDomains; preload',
	'X-Content-Type-Options': 'nosniff',
	'X-Frame-Options': 'DENY'
};

export const handle: Handle = async ({ event, resolve }) => {
	const response = await resolve(event);

	for (const [name, value] of Object.entries(SECURITY_HEADERS)) {
		response.headers.set(name, value);
	}

	return response;
};
