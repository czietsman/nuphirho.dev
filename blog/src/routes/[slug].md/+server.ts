import type { RequestHandler, EntryGenerator } from './$types';
import { getAllPosts, getPostMarkdown } from '$lib/posts';
import { error } from '@sveltejs/kit';

export const prerender = true;

export const entries: EntryGenerator = async () => {
	const posts = await getAllPosts();
	return posts.map((p) => ({ slug: p.slug }));
};

export const GET: RequestHandler = async ({ params }) => {
	const md = await getPostMarkdown(params.slug);
	if (md === null) error(404, 'Not found');
	return new Response(md, {
		headers: {
			'content-type': 'text/markdown; charset=utf-8',
			'cache-control': 'max-age=3600'
		}
	});
};
