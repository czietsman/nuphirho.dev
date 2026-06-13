import type { PageServerLoad } from './$types';
import { getAllPosts } from '$lib/posts';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params }) => {
	const posts = await getAllPosts();
	const filtered = posts.filter((p) => p.tags.includes(params.tag));
	if (filtered.length === 0) error(404, 'Tag not found');
	return { tag: params.tag, posts: filtered };
};
