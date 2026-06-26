export interface Paper {
	date: string;
	title: string;
	url: string;
	summary: string;
}

export const papers: Paper[] = [
	{
		date: '23 Jun 2026',
		title: 'Fifty Years of Specification Completeness: What Aviation Certification Tells AI Governance About Epoch Limits, Proof Surfaces, and the Structural Gap',
		url: 'https://arxiv.org/abs/2606.25120',
		summary: 'Maps three structural requirements that aviation software certification has enforced since 1992 onto AI governance documents, showing that no existing framework requires these properties of individual governance artefacts.',
	},
	{
		date: '22 Apr 2026',
		title: 'Structural Quality Gaps in Practitioner AI Governance Prompts: An Empirical Study Using a Five-Principle Evaluation Framework',
		url: 'https://arxiv.org/abs/2604.21090',
		summary: 'Introduces a five-principle framework for evaluating AI governance prompt quality and applies it empirically to 34 real governance files, finding 37% fall below the structural completeness threshold.',
	},
	{
		date: '26 Mar 2026',
		title: 'The Specification as Quality Gate: Three Hypotheses on AI-Assisted Code Review',
		url: 'https://arxiv.org/abs/2603.25773',
		summary: 'Argues that AI code review is structurally circular without an executable specification, and proposes three testable hypotheses showing how specifications break the circularity and enable reliable AI assistance at scale.',
	},
];
