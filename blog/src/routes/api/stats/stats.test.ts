import { describe, expect, it } from 'vitest';
import { POST } from './+server';

// A stand-in for the BLOG_ANALYTICS KV namespace.
function fakeKV(initial: Record<string, string> = {}) {
	const store = new Map(Object.entries(initial));
	return {
		store,
		get: async (key: string) => store.get(key) ?? null,
		put: async (key: string, value: string) => {
			store.set(key, value);
		}
	};
}

async function post(host: string, kv: ReturnType<typeof fakeKV>, path = '/some-post') {
	const url = new URL(`https://${host}/api/stats`);
	const request = new Request(url, { method: 'POST', body: JSON.stringify({ path }) });
	const response = await POST({ request, url, platform: { env: { BLOG_ANALYTICS: kv } } } as never);
	return (await response.json()) as { count?: number; error?: string };
}

describe('/api/stats', () => {
	it('counts a visit on the canonical host', async () => {
		const kv = fakeKV({ 'visits:/some-post': '41' });
		expect(await post('blog.nuphirho.dev', kv)).toEqual({ count: 42 });
		expect(kv.store.get('visits:/some-post')).toBe('42');
	});

	// The Pages alias and every preview deployment serve the same app. They
	// must not move the production count, which is what anyone could do by
	// posting to nuphirho-blog.pages.dev.
	it('reports without counting on any other host', async () => {
		const kv = fakeKV({ 'visits:/some-post': '41' });
		for (const host of ['nuphirho-blog.pages.dev', 'abc123.nuphirho-blog.pages.dev', 'localhost:5173']) {
			expect(await post(host, kv)).toEqual({ count: 41 });
		}
		expect(kv.store.get('visits:/some-post')).toBe('41');
	});

	it('still rejects an invalid path before touching the store', async () => {
		const kv = fakeKV();
		expect(await post('blog.nuphirho.dev', kv, '../etc')).toEqual({ error: 'invalid path' });
		expect(kv.store.size).toBe(0);
	});
});
