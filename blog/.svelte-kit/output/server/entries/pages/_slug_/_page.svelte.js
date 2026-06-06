import { B as escape_html, i as head, l as html, n as derived, r as ensure_array_like, z as attr } from "../../../chunks/dev.js";
import { t as formatDate } from "../../../chunks/posts.js";
//#region src/routes/[slug]/+page.svelte
function _page($$renderer, $$props) {
	$$renderer.component(($$renderer) => {
		let { data } = $$props;
		const post = derived(() => data.post);
		head("jot9ci", $$renderer, ($$renderer) => {
			$$renderer.title(($$renderer) => {
				$$renderer.push(`<title>${escape_html(post().title)} — nuphirho</title>`);
			});
			if (post().subtitle) {
				$$renderer.push("<!--[0-->");
				$$renderer.push(`<meta name="description"${attr("content", post().subtitle)}/>`);
			} else $$renderer.push("<!--[-1-->");
			$$renderer.push(`<!--]-->`);
		});
		$$renderer.push(`<div class="container"><article><header class="post-header"><h1 class="post-title">${escape_html(post().title)}</h1> `);
		if (post().subtitle) {
			$$renderer.push("<!--[0-->");
			$$renderer.push(`<p class="post-subtitle">${escape_html(post().subtitle)}</p>`);
		} else $$renderer.push("<!--[-1-->");
		$$renderer.push(`<!--]--> <div class="post-meta"><time${attr("datetime", post().publishDate)}>${escape_html(formatDate(post().publishDate))}</time> `);
		if (post().series) {
			$$renderer.push("<!--[0-->");
			$$renderer.push(`<span>${escape_html(post().series)}</span>`);
		} else $$renderer.push("<!--[-1-->");
		$$renderer.push(`<!--]--> `);
		if (post().tags.length) {
			$$renderer.push("<!--[0-->");
			$$renderer.push(`<div class="tags"><!--[-->`);
			const each_array = ensure_array_like(post().tags);
			for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
				let tag = each_array[$$index];
				$$renderer.push(`<span class="tag">${escape_html(tag)}</span>`);
			}
			$$renderer.push(`<!--]--></div>`);
		} else $$renderer.push("<!--[-1-->");
		$$renderer.push(`<!--]--></div></header> <div class="post-content">${html(post().html)}</div></article> <a href="/" class="back-link">← All posts</a></div>`);
	});
}
//#endregion
export { _page as default };
