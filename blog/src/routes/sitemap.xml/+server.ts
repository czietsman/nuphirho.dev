import type { RequestHandler } from './$types';
import { getAllPosts } from '$lib/posts';

export const prerender = true;

const BASE = 'https://blog.nuphirho.dev';

function esc(s: string): string {
	return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
}

export const GET: RequestHandler = async () => {
	const posts = await getAllPosts();

	const staticUrls = ['/', '/about'].map(
		(path) => `  <url><loc>${BASE}${path}</loc></url>`
	);

	const postUrls = posts.map((p) => {
		const lastmod = new Date(p.publishDate).toISOString().slice(0, 10);
		return `  <url><loc>${BASE}/${esc(p.slug)}</loc><lastmod>${lastmod}</lastmod></url>`;
	});

	const tags = [...new Set(posts.flatMap((p) => p.tags))];
	const tagUrls = tags.map((t) => `  <url><loc>${BASE}/tags/${esc(t)}</loc></url>`);

	const body = [
		'<?xml version="1.0" encoding="UTF-8"?>',
		'<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">',
		...staticUrls,
		...postUrls,
		...tagUrls,
		'</urlset>'
	].join('\n');

	return new Response(body, {
		headers: {
			'content-type': 'application/xml; charset=utf-8',
			'cache-control': 'max-age=3600'
		}
	});
};
