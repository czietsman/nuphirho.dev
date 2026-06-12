import type { RequestHandler } from './$types';

export const prerender = true;

const BASE = 'https://nuphirho.dev';
const ROUTES = ['/', '/about', '/novel-findings', '/roadmap', '/words-of-meaning'];

export const GET: RequestHandler = () => {
	const urls = ROUTES.map((r) => `  <url><loc>${BASE}${r}</loc></url>`);

	const body = [
		'<?xml version="1.0" encoding="UTF-8"?>',
		'<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">',
		...urls,
		'</urlset>'
	].join('\n');

	return new Response(body, {
		headers: {
			'content-type': 'application/xml; charset=utf-8',
			'cache-control': 'max-age=3600'
		}
	});
};
