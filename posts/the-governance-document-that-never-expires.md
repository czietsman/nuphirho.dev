---
title: "The Governance Document That Never Expires"
slug: the-governance-document-that-never-expires
publish_date: 2026-06-30
format: article
linkedin_url: https://www.linkedin.com/posts/christo-zietsman_you-wrote-a-system-prompt-or-an-agentsmd-ugcPost-7476401732361842688-zd9I/
cover_image: the-governance-document-that-never-expires.png
cover_image_prompt: |
  Create an image: A wide, dark antique aeronautical chart in landscape format, the surface aged
  and textured like a formal navigation document from an earlier era of aviation certification.
  Across the centre and right two-thirds, fine gold line work traces airways, flight path
  corridors, and route sectors: curved arcs, precise bearings, and formally bounded regions
  suggesting an en-route or qualification chart. The line work is meticulous and hierarchical,
  each boundary precisely drawn. Near the centre, one formally delineated region -- bounded by
  the same precise gold lines as the surrounding structure -- is conspicuously empty, its interior
  in deep, unbroken shadow while the flight paths around it remain illuminated and active. A warm
  amber glow emanates from the route structure at the centre, fading to deep navy at all edges.
  The empty region reads as uncharted territory within an otherwise precisely specified domain.
  Aged parchment texture throughout. Panoramic composition with negative space on the left third
  for visual balance. No text, no people, no faces. 16:9 aspect ratio.
cover_post: |
  You wrote a system prompt. Or an AGENTS.md file. Or an AI policy document.

  When does it expire?

  Not when does it get reviewed. When does it lose its authority over the system it governs.

  Most AI governance documents do not answer that question. They govern the system indefinitely as the model updates, the deployment context shifts, the task domain expands. No trigger fires. No revalidation is required.

  Aviation solved this thirty years ago. A qualification document does not persist by assumption. It persists by demonstration. The deployer bears the burden of showing the prior qualification still holds. The default is invalidity, not continuity.

  A new paper asks: does that logic apply to AI governance documents? A governance document is a static artefact with evaluable structural properties, even when the system it governs is not.

  The companion empirical paper (arXiv:2604.21090) tested the field. Across 34 real practitioner governance documents, zero declare any trigger for revalidation. Not one.

  Aviation's default: assume invalidity.
  AI governance's default: assume permanent authority.

  Full article at blog.nuphirho.dev/the-governance-document-that-never-expires

  arXiv:2606.25120
tags:
  - ai-governance
  - aviation
  - research
  - promptq
---

Most people who work with AI systems professionally have written one. A system prompt. An AGENTS.md file. A policy document stating what the AI is authorised to do, what it must refuse, and how its output will be handled. These documents are real governance instruments. They constrain real systems. Organisations deploy against them, audit against them, and cite them when things go wrong.

Almost none of them say when they expire.

Not when they should be reviewed. When they expire. When the authority they claim over the system they govern stops being valid. When the context has changed sufficiently that the document no longer accurately describes what it purports to govern.

This absence is the subject of a new paper. The paper traces it to a thirty-year tradition in aviation certification that solved the same problem for software, and asks why that solution has not been required of AI governance documents.

---

Aviation's answer to the question of when a software artefact's authority expires is precise. The aviation tool qualification standard defines three cases. Unchanged reuse: no requalification required if nothing material has changed. Operational environment change only: no requalification required if the deployer demonstrates equivalence; the burden of proof is on the deployer, not the certifier. Tool changed: requalification required, scoped by an impact analysis identifying what remains valid and what must be renewed. The default, when equivalence is not demonstrated, is invalidity. The authority of the artefact over the system does not persist by assumption. It must be actively maintained.

No AI governance document takes this position on its own validity. The document was written. It governs the system. It continues to govern the system as the model is updated, as the deployment context changes, as the task domain expands. No trigger fires. No revalidation is required. The default is implicit permanent validity.

---

Aviation's certification framework does not transfer to AI systems at the system level. The aviation community has acknowledged this clearly. The framework's bidirectional traceability requires deterministic behaviour: every requirement traces forward to a test and backward to a source requirement. A language model's behaviour is statistical rather than deterministic. You cannot trace an inference to a requirement. The framework breaks down at the system level.

The paper separates two questions the aviation community has been treating as one. Whether aviation's frameworks apply to AI systems is settled: they do not, and the community has said so. Whether aviation's frameworks apply to AI governance documents is a different question, and it has not been addressed.

A governance document is a static artefact. It was written by a human at a specific time. Its structural properties can be evaluated and required independently of the stochastic system it governs. The framework breaks down at the system level because AI systems are non-deterministic. It does not break down at the document level, because governance documents are not.

That separation is the paper's primary conceptual contribution.

---

The paper maps three structural requirements from aviation certification onto three findings about AI governance documents.

The first is governance linkage. Every claim in a governance document must trace to evidence of its satisfaction. A document that asserts the system will not disclose confidential information, without specifying how that claim can be assessed, provides no basis for compliance verification. The claim is made. The proof surface is absent.

The second is context-bounded validity. An artefact's authority depends on the stability of the context for which it was written. Aviation's default is invalidation; the deployer bears the burden of demonstrating that the prior qualification still holds. The AI governance equivalent requires a staleness declaration: a named set of conditions under which the document requires revalidation before continuing to govern the system.

The third is the structural gap. No regulatory instrument in any of the nine sectors surveyed requires AI governance documents to satisfy structural completeness criteria before deployment. Financial services, healthcare, nuclear, legal, pharmaceutical, insurance, public sector. A structured audit of governance instruments across five language jurisdictions confirms the same absence. Aviation is the sharpest instance because it has the strongest tradition of formal specification completeness. If the absence holds there, it holds everywhere.

---

A companion empirical paper (arXiv:2604.21090) evaluated a corpus of 34 real-world AI governance documents against a seven-principle structural quality framework. Every document in the corpus scores zero on Contextual Currency, the principle governing staleness declarations and epoch limits. Not a single document declares any trigger for revalidation. Zero across 34 documents, confirmed by five independent raters from two model families. Ninety-four percent of the corpus falls below the minimum structural quality threshold under the seven-principle model.

This is not a finding about poorly written documents. The corpus represents real practitioner documents from multiple sectors and language jurisdictions. The absence of staleness declarations is the norm. The default assumption in the field is that governance documents do not expire.

---

Three practical implications follow.

Governance documents should be treated as certification artefacts. The system prompt or policy document governing an AI system is a first-class artefact with evaluable structural properties, not a static operational instruction.

The burden of proof for continuing validity should be on the deployer. A governance document that makes no staleness declaration implicitly claims permanent authority over a system whose context is changing. The appropriate default is the aviation inversion: assume invalidity, demonstrate continuing validity.

Evidence architecture should be specified at authorship time. Compliance claims without evidence mechanisms are assertions, not governance. Writing the proof surface into the document when the document is written is not additional overhead. It is what makes the claims governable.

---

The paper is at arXiv:2606.25120. The companion paper establishing the empirical finding is at arXiv:2604.21090.
