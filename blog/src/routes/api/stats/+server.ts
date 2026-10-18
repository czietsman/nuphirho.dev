import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Only the canonical host counts. The Pages alias and every preview deployment
// serve this same endpoint, and a visit there must not move the production
// count; they get the current count back without adding to it.
const COUNTING_HOST = 'blog.nuphirho.dev';

// Counts are approximate by design. KV has no atomic increment, so two visits
// landing at the same moment can record one. At this blog's traffic that is
// rare and harmless, and not worth a separate datastore.
export const POST: RequestHandler = async ({ request, url, platform }) => {
	let path: unknown;
	try {
		({ path } = await request.json());
	} catch {
		return json({ error: 'invalid JSON' }, { status: 400 });
	}

	if (!path || typeof path !== 'string' || path.length > 500 || !/^\/[\w\-.~%/]*$/.test(path)) {
		return json({ error: 'invalid path' }, { status: 400 });
	}

	const kv = platform?.env?.BLOG_ANALYTICS;
	if (!kv) {
		return json({ count: 0 });
	}

	const key = `visits:${path}`;
	const current = parseInt((await kv.get(key)) ?? '0', 10);
	if (url.host !== COUNTING_HOST) {
		return json({ count: current });
	}
	await kv.put(key, String(current + 1));
	return json({ count: current + 1 });
};
