import { describe, it, expect } from 'vitest';
import { papers, type ExternalPaper, type ProgrammePaper } from './papers.js';

describe('papers', () => {
	it('contains three programme papers and one external prior work', () => {
		expect(papers).toHaveLength(4);
		expect(papers.filter((paper) => paper.kind === 'programme')).toHaveLength(3);
		expect(papers.filter((paper) => paper.kind === 'external')).toHaveLength(1);
	});

	it('lists programme papers in reverse chronological order', () => {
		const dates = papers
			.filter((paper): paper is ProgrammePaper => paper.kind === 'programme')
			.map((paper) => new Date(paper.date).getTime());
		expect(dates).toEqual([...dates].sort((a, b) => b - a));
	});

	it('each programme paper has a slug, arXiv URL, summary, and date', () => {
		for (const p of papers.filter((paper): paper is ProgrammePaper => paper.kind === 'programme')) {
			expect(p.slug.length).toBeGreaterThan(0);
			expect(p.title.length).toBeGreaterThan(0);
			expect(p.arxivUrl).toMatch(/^https:\/\/arxiv\.org\/abs\//);
			expect(p.summary.length).toBeGreaterThan(0);
			expect(p.date.length).toBeGreaterThan(0);
		}
	});

	it('contains the SUNScholar thesis as external prior work', () => {
		const thesis = papers.find(
			(paper): paper is ExternalPaper => paper.kind === 'external' && paper.slug === 'hierarchical-boundary-element-solver',
		);
		expect(thesis).toBeDefined();
		expect(thesis?.title).toBe('A hierarchical linear elastic boundary element solver for lenticular ore bodies');
		expect(thesis?.author).toBe('Christiaan Abraham Zietsman');
		expect(thesis?.degree).toBe('MSc (Mathematical Sciences. Applied Mathematics)');
		expect(thesis?.institution).toBe('University of Stellenbosch');
		expect(thesis?.year).toBe(2007);
		expect(thesis?.externalUrl).toBe('https://scholar.sun.ac.za/items/fbb527fb-e56a-4821-be51-3fec94929774');
		expect(thesis?.abstract).toContain('linear elastic static stress boundary element solver');
	});

	it('models programme and external papers as distinct records', () => {
		const programme: ProgrammePaper = papers[0];
		const external: ExternalPaper = {
			kind: 'external',
			slug: 'prior-work',
			title: 'Prior work',
			externalUrl: 'https://example.com/prior-work',
			abstract: 'An abstract.',
			author: 'An Author',
			institution: 'An Institution',
			degree: 'A Degree',
			year: 2000,
		};

		expect(programme.kind).toBe('programme');
		expect(programme.arxivUrl).toMatch(/^https:\/\/arxiv\.org\/abs\//);
		expect(external.kind).toBe('external');
		expect(external.externalUrl).toBe('https://example.com/prior-work');
	});

	it('contains the specification completeness paper with correct date', () => {
		const found = papers.find((p) => p.arxivUrl === 'https://arxiv.org/abs/2606.25120');
		expect(found).toBeDefined();
		expect(found?.date).toBe('23 Jun 2026');
		expect(found?.slug).toBe('aviation-specification-completeness');
	});

	it('contains the structural quality gaps paper with correct date', () => {
		const found = papers.find((p) => p.arxivUrl === 'https://arxiv.org/abs/2604.21090');
		expect(found).toBeDefined();
		expect(found?.date).toBe('22 Apr 2026');
		expect(found?.slug).toBe('structural-quality-gaps-agents-md');
	});

	it('contains the specification as quality gate paper with correct date', () => {
		const found = papers.find((p) => p.arxivUrl === 'https://arxiv.org/abs/2603.25773');
		expect(found).toBeDefined();
		expect(found?.date).toBe('26 Mar 2026');
		expect(found?.slug).toBe('specification-as-quality-gate');
	});
});
