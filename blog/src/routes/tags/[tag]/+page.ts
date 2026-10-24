import type { EntryGenerator } from './$types';
import { getAllPosts } from '#lib/posts.ts';

export const prerender = true;

export const entries: EntryGenerator = async () => {
	const posts = await getAllPosts();
	const tags = new Set(posts.flatMap((p) => p.tags));
	return [...tags].map((tag) => ({ tag }));
};
