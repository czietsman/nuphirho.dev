import { B as escape_html, i as head, r as ensure_array_like, s as stringify, z as attr } from "../../chunks/dev.js";
import { t as formatDate } from "../../chunks/posts.js";
//#region src/routes/+page.svelte
function _page($$renderer, $$props) {
	$$renderer.component(($$renderer) => {
		let { data } = $$props;
		head("1uha8ag", $$renderer, ($$renderer) => {
			$$renderer.title(($$renderer) => {
				$$renderer.push(`<title>nuphirho — Posts</title>`);
			});
			$$renderer.push(`<meta name="description" content="Writing on software engineering, process, and AI."/>`);
		});
		$$renderer.push(`<div class="container"><div class="page-intro"><h1>Posts</h1> <p>Writing on software engineering, process, and AI.</p></div> <ul class="post-list"><!--[-->`);
		const each_array = ensure_array_like(data.posts);
		for (let $$index_1 = 0, $$length = each_array.length; $$index_1 < $$length; $$index_1++) {
			let post = each_array[$$index_1];
			$$renderer.push(`<li class="post-item"><h2 class="post-item-title"><a${attr("href", `/${stringify(post.slug)}`)}>${escape_html(post.title)}</a></h2> <div class="post-item-meta"><time${attr("datetime", post.publishDate)}>${escape_html(formatDate(post.publishDate))}</time> `);
			if (post.series) {
				$$renderer.push("<!--[0-->");
				$$renderer.push(`<span>${escape_html(post.series)}</span>`);
			} else $$renderer.push("<!--[-1-->");
			$$renderer.push(`<!--]--> `);
			if (post.tags.length) {
				$$renderer.push("<!--[0-->");
				$$renderer.push(`<div class="tags"><!--[-->`);
				const each_array_1 = ensure_array_like(post.tags);
				for (let $$index = 0, $$length = each_array_1.length; $$index < $$length; $$index++) {
					let tag = each_array_1[$$index];
					$$renderer.push(`<span class="tag">${escape_html(tag)}</span>`);
				}
				$$renderer.push(`<!--]--></div>`);
			} else $$renderer.push("<!--[-1-->");
			$$renderer.push(`<!--]--></div> `);
			if (post.subtitle) {
				$$renderer.push("<!--[0-->");
				$$renderer.push(`<p>${escape_html(post.subtitle)}</p>`);
			} else $$renderer.push("<!--[-1-->");
			$$renderer.push(`<!--]--></li>`);
		}
		$$renderer.push(`<!--]--></ul></div>`);
	});
}
//#endregion
export { _page as default };
