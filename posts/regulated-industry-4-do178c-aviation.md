---
title: "Instructions versus specifications"
cover_image_prompt: >-
  A full explainer infographic in the style of a colourful, cheerful, hand-drawn line-art poster on light cream paper. Title at the top in friendly hand-lettered text: "Instructions versus specifications". Draw a side-by-side comparison: "Instruction" points only to "Do this", while "Specification" flows through "Done looks like", "Verify it", and "Show evidence". Add callouts: "Intent is not a quality gate", "Traceability makes claims testable", and "Probabilistic systems limit the analogy". Include a small source box: "DO-178C". A closing hand-lettered line at the bottom: "Define what done means and how you will know." Bright coral, turquoise, sunny yellow, leaf green, sky blue, and warm orange accents, with dark legible lettering, rounded boxes, lively arrows, and friendly instruction-card and checklist icons. Happy, welcoming, and energetic; not institutional, clinical, or ominous. No people, no photographic elements, no logos, and no watermark.
slug: regulated-industry-4-do178c-aviation
publish_date: 2026-07-01
linkedin_url: https://www.linkedin.com/posts/christo-zietsman_aviation-software-has-been-certified-against-activity-7477972513012998144-Z1R-
tags:
  - ai-governance
  - research
  - software-engineering
---

Aviation software has been certified against a standard called DO-178C since 1992. Every line of code in a commercial aircraft must be provably derivable from a requirement. Every requirement must trace forward to a test. Every test result must trace back to a requirement. Extraneous code has no traceable requirement and is a certification finding. Undocumented behaviour has no traceable test and is a certification finding.

The specification is not a document that describes the software. It is a first-class artefact the software must be demonstrably derived from.

AI governance prompts are not specifications in this sense. They are instructions. The difference matters more than it sounds.

An instruction tells an agent what to do. A specification defines what done looks like, how it will be verified, and what evidence demonstrates the definition was met. A governance prompt that says "behave helpfully and safely" is an instruction. A specification would say: here is what helpful and safe looks like, here is how you verify it, here is what is in and out of scope, here is how you weight conflicting evidence, and here is when this document expires. The first is common. The second is rare.

DO-178C has required the second for every line of code in every commercial aircraft for over thirty years. The aviation industry did not arrive at that standard by preference. It arrived there because the cost of incomplete specifications was paid in lives and investigated in accident reports.

AI governance does not yet have that history. But the structural lesson does not require repeating the accidents. It is already written down.

The analogy is not exact. DO-178C governs deterministic software. AI systems are probabilistic. The traceability model does not transfer directly. But the foundational principle does: a governance document that does not specify what done looks like and how you would know if it had not been achieved is not a governance document. It is a statement of intent.

PromptQ is the instrument for closing that gap at authorship time. The waitlist is at promptq.ai.
