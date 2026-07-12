import { render } from 'svelte/server';
import { describe, expect, it } from 'vitest';
import type { ExternalPaper, ProgrammePaper } from '$lib/papers.js';
import Page from './+page.svelte';

const programmePaper: ProgrammePaper = {
	kind: 'programme',
	slug: 'programme-paper',
	date: '1 Jul 2026',
	title: 'A programme paper',
	arxivUrl: 'https://arxiv.org/abs/2607.00001',
	summary: 'The programme paper summary.',
};

const externalPaper: ExternalPaper = {
	kind: 'external',
	slug: 'external-paper',
	title: 'Earlier external work',
	externalUrl: 'https://example.com/external-paper',
	abstract: 'The external abstract.',
	author: 'A Researcher',
	institution: 'A University',
	degree: 'MSc',
	year: 2001,
};

const errata = {
	'programme-paper': {
		slug: 'programme-paper',
		linkedinUrl: 'https://www.linkedin.com/posts/example',
		vnextMarkdownUrl: 'https://example.com/programme-paper-vnext.md',
		vnextPdfUrl: 'https://example.com/programme-paper-vnext.pdf',
		errata: [
			{
				date: '2 Jul 2026',
				location: 'Section 2',
				correction: 'Corrected the stated limit.',
				reason: 'The original value was transposed.',
				type: 'factual',
			},
		],
	},
};

function renderPage(
	paperEntries: Array<ProgrammePaper | ExternalPaper>,
	paperErrata: Record<string, (typeof errata)['programme-paper']> = {},
) {
	return render(Page, { props: { paperEntries, paperErrata } }).body;
}

describe('papers page', () => {
	it('renders errata collapsed with its Wayback note and maintained vnext links', () => {
		const body = renderPage([programmePaper], errata);

		expect(body).toContain('<details');
		expect(body).not.toContain('<details open');
		expect(body).toContain('Errata');
		expect(body).toContain('2 Jul 2026');
		expect(body).toContain('Section 2');
		expect(body).toContain('Corrected the stated limit.');
		expect(body).toContain('The original value was transposed.');
		expect(body).toContain('discoverable through web archives');
		expect(body).toContain('Maintained next version (Markdown)');
		expect(body).toContain('Maintained next version (PDF)');
	});

	it('omits errata and its Wayback note when a programme paper has none', () => {
		const body = renderPage([programmePaper]);

		expect(body).not.toContain('Errata');
		expect(body).not.toContain('discoverable through web archives');
	});

	it('renders LinkedIn only when the errata record provides a URL', () => {
		const withLinkedIn = renderPage([programmePaper], errata);
		const withoutLinkedIn = renderPage([programmePaper], {
			'programme-paper': { ...errata['programme-paper'], linkedinUrl: null },
		});

		expect(withLinkedIn).toContain('Discuss on LinkedIn');
		expect(withoutLinkedIn).not.toContain('Discuss on LinkedIn');
	});

	it('renders external work with link-only controls', () => {
		const body = renderPage([externalPaper], errata);

		expect(body).toContain('Prior work');
		expect(body).toContain('Earlier external work');
		expect(body).toContain('A Researcher');
		expect(body).toContain('A University');
		expect(body).toContain('MSc');
		expect(body).toContain('2001');
		expect(body).toContain('The external abstract.');
		expect(body).toContain('View on SUNScholar');
		expect(body).not.toContain('arXiv');
		expect(body).not.toContain('Errata');
		expect(body).not.toContain('Maintained next version');
		expect(body).not.toContain('Discuss on LinkedIn');
	});
});
