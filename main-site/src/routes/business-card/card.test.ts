import { describe, expect, it, vi } from 'vitest';
import { BUSINESS_CARD, buildVCard, isBusinessCardPath, shareBusinessCard } from './card.js';

describe('public business card', () => {
	it('uses the approved contact and sharing destinations', () => {
		expect(BUSINESS_CARD).toEqual({
			name: 'Christo Zietsman',
			website: 'nuphirho.dev',
			phone: '0823266576',
			internationalPhone: '+27823266576',
			linkedInUrl: 'https://www.linkedin.com/in/christo-zietsman',
			publicUrl: 'https://nuphirho.dev/business-card',
		});
	});

	it('builds an importable vCard', () => {
		const card = buildVCard();
		expect(card).toContain('BEGIN:VCARD\r\nVERSION:4.0\r\n');
		expect(card).toContain('FN:Christo Zietsman\r\n');
		expect(card).toContain('TEL;TYPE=cell:+27823266576\r\n');
		expect(card).toContain('URL:https://nuphirho.dev\r\n');
		expect(card).toContain('X-SOCIALPROFILE;TYPE=linkedin:https://www.linkedin.com/in/christo-zietsman\r\n');
		expect(card).toMatch(/END:VCARD\r\n$/);
	});

	it('opens the native share sheet when available', async () => {
		const share = vi.fn().mockResolvedValue(undefined);
		const result = await shareBusinessCard({ share });

		expect(result).toBe('shared');
		expect(share).toHaveBeenCalledWith({
			title: 'Christo Zietsman',
			text: 'Christo Zietsman | nuphirho.dev',
			url: 'https://nuphirho.dev/business-card',
		});
	});

	it('copies the public URL when native sharing is unavailable', async () => {
		const writeText = vi.fn().mockResolvedValue(undefined);
		const result = await shareBusinessCard({ clipboard: { writeText } });

		expect(result).toBe('copied');
		expect(writeText).toHaveBeenCalledWith('https://nuphirho.dev/business-card');
	});

	it('marks only the public card route as immersive', () => {
		expect(isBusinessCardPath('/business-card')).toBe(true);
		expect(isBusinessCardPath('/business-card/')).toBe(true);
		expect(isBusinessCardPath('/about')).toBe(false);
	});
});
