import { describe, it, expect } from 'vitest';
import { papers } from './papers.js';

describe('papers', () => {
	it('contains exactly three papers', () => {
		expect(papers).toHaveLength(3);
	});

	it('lists papers in reverse chronological order', () => {
		const dates = papers.map((p) => new Date(p.date).getTime());
		expect(dates).toEqual([...dates].sort((a, b) => b - a));
	});

	it('each paper has a non-empty title, url, summary, and date', () => {
		for (const p of papers) {
			expect(p.title.length).toBeGreaterThan(0);
			expect(p.url.length).toBeGreaterThan(0);
			expect(p.summary.length).toBeGreaterThan(0);
			expect(p.date.length).toBeGreaterThan(0);
		}
	});

	it('all urls point to arxiv', () => {
		for (const p of papers) {
			expect(p.url).toMatch(/^https:\/\/arxiv\.org\/abs\//);
		}
	});

	it('contains the specification completeness paper with correct date', () => {
		const found = papers.find((p) => p.url === 'https://arxiv.org/abs/2606.25120');
		expect(found).toBeDefined();
		expect(found?.date).toBe('23 Jun 2026');
	});

	it('contains the structural quality gaps paper with correct date', () => {
		const found = papers.find((p) => p.url === 'https://arxiv.org/abs/2604.21090');
		expect(found).toBeDefined();
		expect(found?.date).toBe('22 Apr 2026');
	});

	it('contains the specification as quality gate paper with correct date', () => {
		const found = papers.find((p) => p.url === 'https://arxiv.org/abs/2603.25773');
		expect(found).toBeDefined();
		expect(found?.date).toBe('26 Mar 2026');
	});
});
