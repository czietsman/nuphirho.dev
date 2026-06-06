import type { PageServerLoad, EntryGenerator } from './$types';
import { getAllPosts, getPost } from '$lib/posts';
import { error } from '@sveltejs/kit';

export const prerender = true;

export const entries: EntryGenerator = async () => {
	const posts = await getAllPosts();
	return posts.map((p) => ({ slug: p.slug }));
};

export const load: PageServerLoad = async ({ params }) => {
	const post = await getPost(params.slug);
	if (!post) error(404, 'Post not found');
	return { post };
};
