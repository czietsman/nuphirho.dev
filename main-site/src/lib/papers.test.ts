import { describe, it, expect } from 'vitest';
import { papers } from './papers.js';

describe('papers', () => {
	it('contains exactly three papers', () => {
		expect(papers).toHaveLength(3);
	});

	it('lists papers in reverse chronological order', () => {
		const years = papers.map((p) => p.year);
		expect(years).toEqual([...years].sort((a, b) => b - a));
	});

	it('each paper has a non-empty title, url, and summary', () => {
		for (const p of papers) {
			expect(p.title.length).toBeGreaterThan(0);
			expect(p.url.length).toBeGreaterThan(0);
			expect(p.summary.length).toBeGreaterThan(0);
		}
	});

	it('all urls point to arxiv', () => {
		for (const p of papers) {
			expect(p.url).toMatch(/^https:\/\/arxiv\.org\/abs\//);
		}
	});

	it('contains the specification completeness paper', () => {
		const found = papers.find((p) => p.url === 'https://arxiv.org/abs/2606.25120');
		expect(found).toBeDefined();
	});

	it('contains the structural quality gaps paper', () => {
		const found = papers.find((p) => p.url === 'https://arxiv.org/abs/2604.21090');
		expect(found).toBeDefined();
	});

	it('contains the specification as quality gate paper', () => {
		const found = papers.find((p) => p.url === 'https://arxiv.org/abs/2603.25773');
		expect(found).toBeDefined();
	});
});
