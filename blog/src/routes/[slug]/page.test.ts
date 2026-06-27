import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/svelte';
import Page from './+page.svelte';

const basePost = {
	title: 'A post with a hero',
	slug: 'a-post-with-a-hero',
	subtitle: 'A subtitle',
	tags: ['ai'],
	series: undefined,
	publishDate: '2026-06-30',
	readingTimeMinutes: 5,
	html: '<p>Body</p>',
	toc: [],
};

describe('post page cover image', () => {
	it('renders the cover image as a hero when cover_image is set', () => {
		render(Page, {
			data: { post: { ...basePost, coverImage: 'https://example.com/hero.png' } },
		});
		const img = screen.getByRole('img');
		expect(img).toHaveAttribute('src', 'https://example.com/hero.png');
	});

	it('renders no image when cover_image is absent', () => {
		render(Page, { data: { post: { ...basePost } } });
		expect(screen.queryByRole('img')).toBeNull();
	});

	it('never renders cover_post (LinkedIn-only teaser) content', () => {
		// cover_post stays in frontmatter and must never appear in the blog HTML.
		render(Page, {
			data: { post: { ...basePost, coverPost: 'LINKEDIN ONLY TEASER TEXT' } },
		});
		expect(screen.queryByText(/LINKEDIN ONLY TEASER TEXT/)).toBeNull();
	});
});

describe('post page social image meta', () => {
	it('emits og:image and twitter:image when a cover image is set', () => {
		render(Page, {
			data: { post: { ...basePost, coverImage: 'https://blog.nuphirho.dev/hero.png' } },
		});
		expect(document.head.querySelector('meta[property="og:image"]')).toHaveAttribute(
			'content',
			'https://blog.nuphirho.dev/hero.png'
		);
		expect(document.head.querySelector('meta[name="twitter:image"]')).toHaveAttribute(
			'content',
			'https://blog.nuphirho.dev/hero.png'
		);
		expect(document.head.querySelector('meta[name="twitter:card"]')).toHaveAttribute(
			'content',
			'summary_large_image'
		);
	});

	it('emits no og:image when there is no cover image', () => {
		render(Page, { data: { post: { ...basePost } } });
		expect(document.head.querySelector('meta[property="og:image"]')).toBeNull();
	});
});
