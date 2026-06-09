import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ request, platform }) => {
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
	await kv.put(key, String(current + 1));
	return json({ count: current + 1 });
};
