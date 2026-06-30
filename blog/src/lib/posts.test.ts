import { describe, it, expect } from 'vitest';
import { isPostVisible, resolveCoverImage, buildMeta } from './posts.js';

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
