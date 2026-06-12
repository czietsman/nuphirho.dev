import type { RequestHandler } from './$types';
import { getAllPosts } from '$lib/posts';

export const prerender = true;

const BASE = 'https://blog.nuphirho.dev';
const TITLE = 'nuphirho';
const DESCRIPTION = 'Technical writing by Christo Zietsman on software engineering, AI, security, and systems thinking.';

function esc(s: string): string {
	return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

export const GET: RequestHandler = async () => {
	const posts = await getAllPosts();

	const items = posts.map((p) => {
		const url = `${BASE}/${p.slug}`;
		const desc = p.subtitle ? `\n    <description>${esc(p.subtitle)}</description>` : '';
		return [
			'  <item>',
			`    <title>${esc(p.title)}</title>`,
			`    <link>${url}</link>`,
			`    <guid isPermaLink="true">${url}</guid>`,
			`    <pubDate>${new Date(p.publishDate).toUTCString()}</pubDate>${desc}`,
			'  </item>'
		].join('\n');
	});

	const body = [
		'<?xml version="1.0" encoding="UTF-8"?>',
		'<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">',
		'  <channel>',
		`    <title>${TITLE}</title>`,
		`    <link>${BASE}</link>`,
		`    <description>${DESCRIPTION}</description>`,
		'    <language>en-gb</language>',
		`    <atom:link href="${BASE}/rss.xml" rel="self" type="application/rss+xml" />`,
		...items,
		'  </channel>',
		'</rss>'
	].join('\n');

	return new Response(body, {
		headers: {
			'content-type': 'application/rss+xml; charset=utf-8',
			'cache-control': 'max-age=3600'
		}
	});
};
