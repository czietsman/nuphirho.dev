import { n as getAllPosts, r as getPost } from "../../../chunks/posts.js";
import { error } from "@sveltejs/kit";
//#region src/routes/[slug]/+page.server.ts
var prerender = true;
var entries = async () => {
	return (await getAllPosts()).map((p) => ({ slug: p.slug }));
};
var load = async ({ params }) => {
	const post = await getPost(params.slug);
	if (!post) error(404, "Post not found");
	return { post };
};
//#endregion
export { entries, load, prerender };
