import { describe, it, expect } from 'vitest';
import { marked } from 'marked';
import { isPostVisible, resolveCoverImage, buildMeta, renderMarkdown } from './posts.js';

describe('isPostVisible', () => {
	const today = '2026-06-25';

	it('excludes draft posts', () => {
		expect(isPostVisible({ draft: true, publish_date: '2026-01-01' }, today)).toBe(false);
	});

	it('excludes posts without a publish_date', () => {
		expect(isPostVisible({ draft: false, publish_date: undefined }, today)).toBe(false);
	});

	it('excludes posts with a future publish_date', () => {
		expect(isPostVisible({ draft: false, publish_date: '2026-12-31' }, today)).toBe(false);
	});

	it('includes posts published today', () => {
		expect(isPostVisible({ draft: false, publish_date: '2026-06-25' }, today)).toBe(true);
	});

	it('includes posts published in the past', () => {
		expect(isPostVisible({ draft: false, publish_date: '2026-01-01' }, today)).toBe(true);
	});

	it('handles Date objects from gray-matter YAML parsing', () => {
		expect(isPostVisible({ draft: false, publish_date: new Date('2026-12-31') }, today)).toBe(false);
		expect(isPostVisible({ draft: false, publish_date: new Date('2026-01-01') }, today)).toBe(true);
	});
});

describe('resolveCoverImage', () => {
	it('returns undefined when no cover image is set', () => {
		expect(resolveCoverImage(undefined)).toBeUndefined();
		expect(resolveCoverImage('')).toBeUndefined();
	});

	it('converts a bare relative filename to an absolute blog URL', () => {
		expect(resolveCoverImage('the-governance-document-that-never-expires.png')).toBe(
			'https://blog.nuphirho.dev/the-governance-document-that-never-expires.png'
		);
	});

	it('strips a leading slash before joining to the origin', () => {
		expect(resolveCoverImage('/hero.png')).toBe('https://blog.nuphirho.dev/hero.png');
	});

	it('leaves an absolute http(s) URL unchanged', () => {
		expect(resolveCoverImage('https://cdn.example.com/hero.png')).toBe(
			'https://cdn.example.com/hero.png'
		);
		expect(resolveCoverImage('http://cdn.example.com/hero.png')).toBe(
			'http://cdn.example.com/hero.png'
		);
	});
});

describe('buildMeta', () => {
	it('exposes linkedinUrl when linkedin_url frontmatter is set', () => {
		const meta = buildMeta(
			{ title: 'Title', publish_date: '2026-01-01', linkedin_url: 'https://www.linkedin.com/feed/update/urn:li:share:123' },
			'slug',
			'content'
		);
		expect(meta.linkedinUrl).toBe('https://www.linkedin.com/feed/update/urn:li:share:123');
	});

	it('leaves linkedinUrl undefined when linkedin_url frontmatter is absent', () => {
		const meta = buildMeta({ title: 'Title', publish_date: '2026-01-01' }, 'slug', 'content');
		expect(meta.linkedinUrl).toBeUndefined();
	});

	it('leaves linkedinUrl undefined when linkedin_url frontmatter is empty', () => {
		const meta = buildMeta({ title: 'Title', publish_date: '2026-01-01', linkedin_url: '' }, 'slug', 'content');
		expect(meta.linkedinUrl).toBeUndefined();
	});
});

describe('heading rendering', () => {
	it('renders inline markdown in the heading and derives a slug id from its text', () => {
		expect(marked.parse('## Hello `World`')).toBe('<h2 id="hello-world">Hello <code>World</code></h2>\n');
	});

	it('decodes numeric entities before slugifying', () => {
		expect(marked.parse("### Don't panic")).toBe('<h3 id="dont-panic">Don&#39;t panic</h3>\n');
	});
});

describe('renderMarkdown', () => {
	it('converts an embed line into a linked figure', () => {
		const { html } = renderMarkdown('%[https://twitter.com/example/status/1]');
		expect(html).toContain(
			'<figure class="embed"><a href="https://twitter.com/example/status/1" rel="noopener noreferrer" target="_blank">https://twitter.com/example/status/1</a></figure>'
		);
	});

	it('removes a manually written series navigation paragraph', () => {
		const { html } = renderMarkdown('*Series: Part 1 of 3*\n\nBody text.');
		expect(html).not.toContain('Series:');
		expect(html).toContain('<p>Body text.</p>');
	});

	it('highlights fenced code in a known language', () => {
		const { html } = renderMarkdown('```typescript\nconst x = 1;\n```');
		expect(html).toContain('<code class="hljs language-typescript">');
		expect(html).toContain('<span class="hljs-keyword">const</span>');
	});

	it('escapes fenced code in an unknown language instead of failing', () => {
		const { html } = renderMarkdown('```nosuchlanguage\n<script>alert(1)</script>\n```');
		expect(html).toContain('&lt;script&gt;');
		expect(html).not.toContain('<script>');
	});

	it('lists only second-level headings in the table of contents', () => {
		const { toc } = renderMarkdown('# Title\n\n## First\n\n### Nested\n\n## Second');
		expect(toc).toEqual([
			{ id: 'first', text: 'First' },
			{ id: 'second', text: 'Second' }
		]);
	});

	it('gives every table of contents entry an id that matches a rendered heading anchor', () => {
		const { html, toc } = renderMarkdown(
			[
				'## Using `npm ci` in CI',
				'## Read [the docs](https://example.com) first',
				'## *Emphasis* and **strong**',
				"## Don't panic",
				'## Q&A'
			].join('\n\n')
		);
		const anchors = [...html.matchAll(/<h2 id="([^"]*)">/g)].map((m) => m[1]);
		expect(toc.map((entry) => entry.id)).toEqual(anchors);
	});

	it('shows table of contents text without markdown syntax or HTML entities', () => {
		const { toc } = renderMarkdown(
			['## Using `npm ci` in CI', '## Read [the docs](https://example.com) first', "## Don't panic", '## Q&A'].join('\n\n')
		);
		expect(toc.map((entry) => entry.text)).toEqual([
			'Using npm ci in CI',
			'Read the docs first',
			"Don't panic",
			'Q&A'
		]);
	});
});
