// Date formatting for the blog. Kept free of node:* and node-only
// dependencies so it can be imported into client-rendered components without
// pulling the post-loading data layer (gray-matter) into the browser bundle.

export function formatDate(iso: string): string {
	const d = new Date(iso);
	return d.toLocaleDateString('en-GB', { day: 'numeric', month: 'long', year: 'numeric' });
}
