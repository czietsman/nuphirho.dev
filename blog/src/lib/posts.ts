import { readdir, readFile } from 'node:fs/promises';
import { resolve } from 'node:path';
import matter from 'gray-matter';
import { marked } from 'marked';

export interface PostMeta {
	title: string;
	slug: string;
	subtitle?: string;
	tags: string[];
	series?: string;
	publishDate: string; // ISO date string
	brief?: string;
}

export interface Post extends PostMeta {
	html: string;
}

const POSTS_DIR = resolve('..', 'posts');

const EMBED_RE = /^%\[(.+?)\]$/gm;

function convertEmbeds(md: string): string {
	return md.replace(EMBED_RE, (_match, url) =>
		`<figure class="embed"><a href="${url}" rel="noopener noreferrer" target="_blank">${url}</a></figure>`
	);
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function buildMeta(data: Record<string, any>, slug: string): PostMeta {
	return {
		title: String(data.title ?? ''),
		slug,
		subtitle: data.subtitle ? String(data.subtitle) : undefined,
		tags: Array.isArray(data.tags) ? data.tags.map(String) : [],
		series: data.series ? String(data.series) : undefined,
		publishDate: String(data.publish_date),
	};
}

async function parsePost(filename: string): Promise<PostMeta | null> {
	const raw = await readFile(resolve(POSTS_DIR, filename), 'utf-8');
	const { data } = matter(raw);

	if (data.draft === true || !data.publish_date) return null;

	return buildMeta(data, String(data.slug ?? filename.replace(/\.md$/, '')));
}

export async function getAllPosts(): Promise<PostMeta[]> {
	const files = await readdir(POSTS_DIR);
	const mds = files.filter((f) => f.endsWith('.md'));
	const posts = await Promise.all(mds.map(parsePost));
	return (posts.filter(Boolean) as PostMeta[]).sort(
		(a, b) => new Date(b.publishDate).getTime() - new Date(a.publishDate).getTime()
	);
}

export async function getPost(slug: string): Promise<Post | null> {
	const files = await readdir(POSTS_DIR);
	const mds = files.filter((f) => f.endsWith('.md'));

	for (const file of mds) {
		const raw = await readFile(resolve(POSTS_DIR, file), 'utf-8');
		const { data, content } = matter(raw);

		if (data.draft === true || !data.publish_date) continue;

		const postSlug = String(data.slug ?? file.replace(/\.md$/, ''));
		if (postSlug !== slug) continue;

		return {
			...buildMeta(data, postSlug),
			html: marked.parse(convertEmbeds(content)) as string,
		};
	}

	return null;
}

export function formatDate(iso: string): string {
	const d = new Date(iso);
	return d.toLocaleDateString('en-GB', { day: 'numeric', month: 'long', year: 'numeric' });
}
