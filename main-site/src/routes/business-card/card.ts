export const BUSINESS_CARD = {
	name: 'Christo Zietsman',
	website: 'nuphirho.dev',
	phone: '0823266576',
	internationalPhone: '+27823266576',
	linkedInUrl: 'https://www.linkedin.com/in/christo-zietsman',
	publicUrl: 'https://nuphirho.dev/business-card',
} as const;

type SharePlatform = {
	share?: (data: ShareData) => Promise<void>;
	clipboard?: { writeText: (text: string) => Promise<void> };
};

export function isBusinessCardPath(pathname: string): boolean {
	return pathname.replace(/\/$/, '') === '/business-card';
}

export function buildVCard(): string {
	return [
		'BEGIN:VCARD',
		'VERSION:4.0',
		'N:Zietsman;Christo;;;',
		'FN:Christo Zietsman',
		`TEL;TYPE=cell:${BUSINESS_CARD.internationalPhone}`,
		`URL:https://${BUSINESS_CARD.website}`,
		`X-SOCIALPROFILE;TYPE=linkedin:${BUSINESS_CARD.linkedInUrl}`,
		'END:VCARD',
		'',
	].join('\r\n');
}

export async function shareBusinessCard(platform: SharePlatform): Promise<'shared' | 'copied'> {
	if (platform.share) {
		await platform.share({
			title: BUSINESS_CARD.name,
			text: `${BUSINESS_CARD.name} | ${BUSINESS_CARD.website}`,
			url: BUSINESS_CARD.publicUrl,
		});
		return 'shared';
	}
	if (!platform.clipboard) throw new Error('Sharing is not supported');
	await platform.clipboard.writeText(BUSINESS_CARD.publicUrl);
	return 'copied';
}
