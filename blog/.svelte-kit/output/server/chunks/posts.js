import { readFile, readdir } from "node:fs/promises";
import { resolve } from "node:path";
import matter from "gray-matter";
//#region src/lib/posts.ts
var POSTS_DIR = resolve("..", "posts");
var EMBED_RE = /^%\[(.+?)\]$/gm;
function convertEmbeds(md) {
	return md.replace(EMBED_RE, (_match, url) => `<figure class="embed"><a href="${url}" rel="noopener noreferrer" target="_blank">${url}</a></figure>`);
}
async function parsePost(filename) {
	const { data } = matter(await readFile(resolve(POSTS_DIR, filename), "utf-8"));
	if (data.draft === true || !data.publish_date) return null;
	return {
		title: String(data.title ?? ""),
		slug: String(data.slug ?? filename.replace(/\.md$/, "")),
		subtitle: data.subtitle ? String(data.subtitle) : void 0,
		tags: Array.isArray(data.tags) ? data.tags.map(String) : [],
		series: data.series ? String(data.series) : void 0,
		publishDate: String(data.publish_date)
	};
}
async function getAllPosts() {
	const mds = (await readdir(POSTS_DIR)).filter((f) => f.endsWith(".md"));
	return (await Promise.all(mds.map(parsePost))).filter(Boolean).sort((a, b) => new Date(b.publishDate).getTime() - new Date(a.publishDate).getTime());
}
async function getPost(slug) {
	const mds = (await readdir(POSTS_DIR)).filter((f) => f.endsWith(".md"));
	for (const file of mds) {
		const { data, content } = matter(await readFile(resolve(POSTS_DIR, file), "utf-8"));
		if (data.draft === true || !data.publish_date) continue;
		const postSlug = String(data.slug ?? file.replace(/\.md$/, ""));
		if (postSlug !== slug) continue;
		const { marked } = await import("marked");
		const html = marked.parse(convertEmbeds(content));
		return {
			title: String(data.title ?? ""),
			slug: postSlug,
			subtitle: data.subtitle ? String(data.subtitle) : void 0,
			tags: Array.isArray(data.tags) ? data.tags.map(String) : [],
			series: data.series ? String(data.series) : void 0,
			publishDate: String(data.publish_date),
			brief: void 0,
			html
		};
	}
	return null;
}
function formatDate(iso) {
	return new Date(iso).toLocaleDateString("en-GB", {
		day: "numeric",
		month: "long",
		year: "numeric"
	});
}
//#endregion
export { getAllPosts as n, getPost as r, formatDate as t };
