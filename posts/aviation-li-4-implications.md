---
title: "What aviation's discipline implies for AI governance practice"
slug: "aviation-li-4-implications"
series: "aviation-specification-completeness"
series_part: 4
tags: [ai-governance, research]
publish_date: 2026-07-13
linkedin_url:
---

Aviation's software certification tradition has three operational lessons for AI governance practice. None requires replicating the aviation machinery. Each requires a decision about defaults.

The first is to treat governance documents as certification artefacts. DO-178C treats the requirements document as a first-class artefact with quality criteria, completeness requirements, and traceability obligations. The code must be provably derivable from it. The same standard applies to the system prompt, AGENTS.md file, or policy document governing an AI system. It is not an operational instruction. It is a specification. Its structural properties are evaluable and should be required.

The second is to invert the burden of proof. Aviation's default for a qualified tool is invalidity when context changes; the deployer must demonstrate continuing validity. The AI governance default is the opposite: the document remains authoritative until something goes wrong. Inverting that default means requiring a staleness declaration in every governance document: the named conditions under which the document requires revalidation before continuing to govern the system. Not a review schedule. A validity boundary.

The third is to specify evidence architecture at authorship time. DO-178C requires that evidence categories and objectives be defined during development, not assembled retrospectively at audit. The AI governance equivalent is a proof surface: a statement, in the document itself, of what evidence is required to demonstrate compliance with its governance claims, what makes that evidence sufficient, and how deviations are reported. A governance document that makes compliance claims without specifying how they can be assessed provides no basis for independent oversight.

These three are not novel requirements. They are the conditions for the governance claims in the document to be governable. A document that does not meet them is not a governance instrument. It is an assertion.

New paper at arXiv:2606.25120. Companion empirical paper at arXiv:2604.21090.
