import { buildVCard } from '../card.js';

export function GET(): Response {
	return new Response(buildVCard(), {
		headers: {
			'content-type': 'text/vcard; charset=utf-8',
			'content-disposition': 'attachment; filename="christo-zietsman.vcf"',
		},
	});
}
