import { n as getAllPosts } from "../../chunks/posts.js";
//#region src/routes/+page.server.ts
var prerender = true;
var load = async () => {
	return { posts: await getAllPosts() };
};
//#endregion
export { load, prerender };
