import "../../chunks/internal.js";
import { $ as getContext, B as escape_html, c as unsubscribe_stores, o as store_get, z as attr } from "../../chunks/dev.js";
import "../../chunks/client.js";
//#region node_modules/@sveltejs/kit/src/runtime/app/stores.js
/**
* A function that returns all of the contextual stores. On the server, this must be called during component initialization.
* Only use this if you need to defer store subscription until after the component has mounted, for some reason.
*
* @deprecated Use `$app/state` instead (requires Svelte 5, [see docs for more info](https://svelte.dev/docs/kit/migrating-to-sveltekit-2#SvelteKit-2.12:-$app-stores-deprecated))
*/
var getStores = () => {
	const stores$1 = getContext("__svelte__");
	return {
		/** @type {typeof page} */
		page: { subscribe: stores$1.page.subscribe },
		/** @type {typeof navigating} */
		navigating: { subscribe: stores$1.navigating.subscribe },
		/** @type {typeof updated} */
		updated: stores$1.updated
	};
};
/**
* A readable store whose value contains page data.
*
* On the server, this store can only be subscribed to during component initialization. In the browser, it can be subscribed to at any time.
*
* @deprecated Use `page` from `$app/state` instead (requires Svelte 5, [see docs for more info](https://svelte.dev/docs/kit/migrating-to-sveltekit-2#SvelteKit-2.12:-$app-stores-deprecated))
* @type {import('svelte/store').Readable<import('@sveltejs/kit').Page>}
*/
var page = { subscribe(fn) {
	return getStores().page.subscribe(fn);
} };
//#endregion
//#region src/routes/+layout.svelte
function _layout($$renderer, $$props) {
	$$renderer.component(($$renderer) => {
		var $$store_subs;
		let { children } = $$props;
		$$renderer.push(`<a href="#main" class="skip-link">Skip to content</a> <header class="site-header"><div class="container"><a href="https://nuphirho.dev" class="site-name" aria-label="nuphirho home">nuphirho</a> <nav class="site-nav" aria-label="Main navigation"><a href="/about"${attr("aria-current", store_get($$store_subs ??= {}, "$page", page).url.pathname === "/about" ? "page" : void 0)}>This Blog</a> <a href="https://nuphirho.dev/words-of-meaning">Words</a> <a href="/"${attr("aria-current", store_get($$store_subs ??= {}, "$page", page).url.pathname === "/" ? "page" : void 0)}>Posts</a> <button id="theme-toggle" class="theme-toggle" type="button"${attr("aria-label", "Switch to dark theme")}>${escape_html("Dark")}</button></nav></div></header> <main id="main" class="site-main">`);
		children($$renderer);
		$$renderer.push(`<!----></main> <footer class="site-footer"><div class="container"><span>© 2026 Christo Zietsman</span> <span class="footer-links"><a href="https://doi.org/10.48550/arXiv.2603.25773">arXiv Paper 1</a> <a href="https://doi.org/10.48550/arXiv.2604.21090">arXiv Paper 2</a> <a href="https://github.com/czietsman">GitHub</a> <a href="https://www.linkedin.com/in/christo-zietsman/">LinkedIn</a> <a href="https://nuphirho.dev/privacy">Privacy</a> <a href="https://nuphirho.dev/cookies">Cookies</a></span></div></footer>`);
		if ($$store_subs) unsubscribe_stores($$store_subs);
	});
}
//#endregion
export { _layout as default };
