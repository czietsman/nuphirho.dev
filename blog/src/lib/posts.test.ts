import { describe, it, expect } from 'vitest';
import { isPostVisible } from './posts.js';

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
