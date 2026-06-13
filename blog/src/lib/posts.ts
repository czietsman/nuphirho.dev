import matter from 'gray-matter';
import { marked } from 'marked';
import { markedHighlight } from 'marked-highlight';
import hljs from 'highlight.js';

function slugifyHeading(text: string): string {
	return text
		.replace(/<[^>]+>/g, '')                                          // strip HTML tags
		.replace(/&#(\d+);/g, (_, n) => String.fromCharCode(Number(n)))  // decode numeric entities
		.replace(/&[a-z]+;/g, '')                                         // strip remaining named entities
		.replace(/[`*_~[\]()]/g, '')                                      // strip inline markdown
		.toLowerCase()
		.replace(/[^a-z0-9\s-]/g, '')
		.trim()
		.replace(/\s+/g, '-');
}

marked.use(markedHighlight({
	langPrefix: 'hljs language-',
	highlight(code, lang) {
		const language = hljs.getLanguage(lang) ? lang : 'plaintext';
		return hljs.highlight(code, { language }).value;
	}
}));

marked.use({
	renderer: {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		heading(this: any, text: string, depth: number) {
			const id = slugifyHeading(text);
			return `<h${depth} id="${id}">${text}</h${depth}>\n`;
		}
	}
});

export interface TocEntry {
	id: string;
	text: string;
}

export interface PostMeta {
	title: string;
	slug: string;
	subtitle?: string;
	tags: string[];
	series?: string;
	publishDate: string; // ISO date string
	brief?: string;
	readingTimeMinutes: number;
}

export interface Post extends PostMeta {
	html: string;
	toc: TocEntry[];
	prevInSeries?: { slug: string; title: string };
	nextInSeries?: { slug: string; title: string };
}

// node:fs/promises and node:path are imported lazily inside functions so they
// are never present in the Cloudflare Worker bundle at module load time.
// These functions only run during SSR prerender at build time.

const EMBED_RE = /^%\[(.+?)\]$/gm;

function convertEmbeds(md: string): string {
	return md.replace(EMBED_RE, (_match, url) =>
		`<figure class="embed"><a href="${url}" rel="noopener noreferrer" target="_blank">${url}</a></figure>`
	);
}

function readingTimeMinutes(content: string): number {
	const words = content.trim().split(/\s+/).filter(Boolean).length;
	return Math.max(1, Math.ceil(words / 200));
}

function buildToc(content: string): TocEntry[] {
	const tokens = marked.lexer(content);
	return tokens
		.filter((t): t is { type: 'heading'; depth: number; text: string; raw: string } =>
			t.type === 'heading' && (t as { depth: number }).depth === 2
		)
		.map((t) => ({
			id: slugifyHeading(t.text),
			text: t.text,
		}));
}

function stripManualSeriesNav(html: string): string {
	return html.replace(/<p><em>Series:[\s\S]*?<\/p>/g, '');
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function buildMeta(data: Record<string, any>, slug: string, content: string): PostMeta {
	return {
		title: String(data.title ?? ''),
		slug,
		subtitle: data.subtitle ? String(data.subtitle) : undefined,
		tags: Array.isArray(data.tags) ? data.tags.map(String) : [],
		series: data.series ? String(data.series) : undefined,
		publishDate: String(data.publish_date),
		readingTimeMinutes: readingTimeMinutes(content),
	};
}

async function postsDir(): Promise<string> {
	const { resolve } = await import('node:path');
	return resolve('..', 'posts');
}

async function parsePost(filename: string, dir: string): Promise<PostMeta | null> {
	const { readFile } = await import('node:fs/promises');
	const { resolve } = await import('node:path');
	const raw = await readFile(resolve(dir, filename), 'utf-8');
	const { data, content } = matter(raw);

	if (data.draft === true || !data.publish_date) return null;

	return buildMeta(data, String(data.slug ?? filename.replace(/\.md$/, '')), content);
}

export async function getAllPosts(): Promise<PostMeta[]> {
	const { readdir } = await import('node:fs/promises');
	const dir = await postsDir();
	const files = await readdir(dir);
	const mds = files.filter((f) => f.endsWith('.md'));
	const posts = await Promise.all(mds.map((f) => parsePost(f, dir)));
	return (posts.filter(Boolean) as PostMeta[]).sort(
		(a, b) => new Date(b.publishDate).getTime() - new Date(a.publishDate).getTime()
	);
}

export async function getPost(slug: string): Promise<Post | null> {
	if (!/^[\w-]+$/.test(slug)) return null;

	const { readdir, readFile } = await import('node:fs/promises');
	const { resolve } = await import('node:path');
	const dir = await postsDir();
	const files = await readdir(dir);
	const mds = files.filter((f) => f.endsWith('.md'));

	const allMeta: PostMeta[] = [];
	let target: { meta: PostMeta; content: string } | null = null;

	for (const file of mds) {
		const raw = await readFile(resolve(dir, file), 'utf-8');
		const { data, content } = matter(raw);

		if (data.draft === true || !data.publish_date) continue;

		const postSlug = String(data.slug ?? file.replace(/\.md$/, ''));
		const meta = buildMeta(data, postSlug, content);
		allMeta.push(meta);
		if (postSlug === slug) target = { meta, content };
	}

	if (!target) return null;

	let prevInSeries: { slug: string; title: string } | undefined;
	let nextInSeries: { slug: string; title: string } | undefined;

	if (target.meta.series) {
		const seriesPosts = allMeta
			.filter((p) => p.series === target!.meta.series)
			.sort((a, b) => new Date(a.publishDate).getTime() - new Date(b.publishDate).getTime());
		const idx = seriesPosts.findIndex((p) => p.slug === slug);
		if (idx > 0) prevInSeries = { slug: seriesPosts[idx - 1].slug, title: seriesPosts[idx - 1].title };
		if (idx < seriesPosts.length - 1) nextInSeries = { slug: seriesPosts[idx + 1].slug, title: seriesPosts[idx + 1].title };
	}

	return {
		...target.meta,
		html: stripManualSeriesNav(marked.parse(convertEmbeds(target.content)) as string),
		toc: buildToc(target.content),
		prevInSeries,
		nextInSeries,
	};
}

export async function getPostMarkdown(slug: string): Promise<string | null> {
	if (!/^[\w-]+$/.test(slug)) return null;

	const { readdir, readFile } = await import('node:fs/promises');
	const { resolve } = await import('node:path');
	const dir = await postsDir();
	const files = await readdir(dir);
	const mds = files.filter((f) => f.endsWith('.md'));

	for (const file of mds) {
		const raw = await readFile(resolve(dir, file), 'utf-8');
		const { data, content } = matter(raw);

		if (data.draft === true || !data.publish_date) continue;

		const postSlug = String(data.slug ?? file.replace(/\.md$/, ''));
		if (postSlug !== slug) continue;

		return content.trim();
	}

	return null;
}

export { formatDate } from './format.js';
