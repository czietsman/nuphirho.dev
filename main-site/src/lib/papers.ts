interface PaperBase {
	kind: 'programme' | 'external';
	slug: string;
	title: string;
}

export interface ProgrammePaper extends PaperBase {
	kind: 'programme';
	date: string;
	arxivUrl: string;
	summary: string;
}

export interface ExternalPaper extends PaperBase {
	kind: 'external';
	externalUrl: string;
	abstract: string;
	author: string;
	institution: string;
	degree: string;
	year: number;
}

export type Paper = ProgrammePaper | ExternalPaper;

export const papers: Paper[] = [
	{
		kind: 'programme',
		slug: 'aviation-specification-completeness',
		date: '23 Jun 2026',
		title: 'Fifty Years of Specification Completeness: What Aviation Certification Tells AI Governance About Epoch Limits, Proof Surfaces, and the Structural Gap',
		arxivUrl: 'https://arxiv.org/abs/2606.25120',
		summary: 'Maps three structural requirements that aviation software certification has enforced since 1992 onto AI governance documents, showing that no existing framework requires these properties of individual governance artefacts.',
	},
	{
		kind: 'programme',
		slug: 'structural-quality-gaps-agents-md',
		date: '22 Apr 2026',
		title: 'Structural Quality Gaps in Practitioner AI Governance Prompts: An Empirical Study Using a Five-Principle Evaluation Framework',
		arxivUrl: 'https://arxiv.org/abs/2604.21090',
		summary: 'Introduces a five-principle framework for evaluating AI governance prompt quality and applies it empirically to 34 real governance files, finding 37% fall below the structural completeness threshold.',
	},
	{
		kind: 'programme',
		slug: 'specification-as-quality-gate',
		date: '26 Mar 2026',
		title: 'The Specification as Quality Gate: Three Hypotheses on AI-Assisted Code Review',
		arxivUrl: 'https://arxiv.org/abs/2603.25773',
		summary: 'Argues that AI code review is structurally circular without an executable specification, and proposes three testable hypotheses showing how specifications break the circularity and enable reliable AI assistance at scale.',
	},
	{
		kind: 'external',
		slug: 'hierarchical-boundary-element-solver',
		title: 'A hierarchical linear elastic boundary element solver for lenticular ore bodies',
		externalUrl: 'https://scholar.sun.ac.za/items/fbb527fb-e56a-4821-be51-3fec94929774',
		author: 'Christiaan Abraham Zietsman',
		institution: 'University of Stellenbosch',
		degree: 'MSc (Mathematical Sciences. Applied Mathematics)',
		year: 2007,
		abstract: "South Africa is involved in huge mining operations deep in the earth's crust. Stresses induced by these mining operations may cause seismic events or rockbursts to occur, which could damage infrastructure and put miners' lives at risk. The effect of different mining layouts are modelled and used by engineers to make design decisions. The frequency at which models are updated and integrated with the decision making process is not optimal. These large mining layouts can not be modelled adequately using domain methods, but they are particularly well suited for the boundary element method (BEM).\n\nThis work focuses on the theory and background needed for creating a linear elastic static stress boundary element solver suited to South African mining layouts. It starts with linear elastic theory and subsequently describes the physical continuum, governing equations and the fundamental solutions which are an integral part of the BEM. Kelvin's solution cannot be applied to crack-like excavations, therefore the displacement discontinuity kernels, which are very well suited to model fractures, are derived. The derivation is approached from both the direct and indirect BEM's perspectives. The problem is cast as a boundary integral equation which can be solved using the BEM. Some of the different specializations of the BEM are discussed. The major drawback of the BEM is that it produces a dense influence matrix which quickly becomes intractable on desktop computers. Generally a mining layout requires a large amount of boundary elements, even for coarse discretization, therefore different techniques of representing the influence matrix are discussed, which, combined with an iterative solver like GMRES or Bi-CG, allows solving linear elastic static stress models.",
	},
];
