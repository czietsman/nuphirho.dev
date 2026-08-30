---
title: "IEC 61508 and the governance document that never expires"
cover_image_prompt: >-
  A full explainer infographic in the style of a colourful, cheerful, hand-drawn line-art poster on light cream paper. Title at the top in friendly hand-lettered text: "IEC 61508 and the governance document that never expires". Draw a validity card labelled "Safety case" connected to three change triggers: "Component change", "Architecture change", and "Requirement change"; each trigger points to "Revalidate". Beside it, show an AI governance card continuing past the same triggers with a question mark. Add callouts: "Validity depends on assumptions", "The burden sits with the operator", and "Default to invalidation when assumptions change". Include a small source box: "IEC 61508". A closing hand-lettered line at the bottom: "Governance needs declared staleness triggers." Bright coral, turquoise, sunny yellow, leaf green, sky blue, and warm orange accents, with dark legible lettering, rounded boxes, lively arrows, and friendly refresh and document icons. Happy, welcoming, and energetic; not institutional, clinical, or ominous. No people, no photographic elements, no logos, and no watermark.
slug: regulated-industry-3-iec61508-epoch-limits
publish_date: 2026-06-24
linkedin_url: https://www.linkedin.com/posts/christo-zietsman_iec-61508-is-the-foundational-standard-for-activity-7475435798046371842-gAK2
tags:
  - ai-governance
  - research
---

IEC 61508 is the foundational standard for functional safety in industrial systems. It has been refined over fifty years. Version by version, it has accumulated answers to a single question: when does a system's safety case stop being valid?

The answer is not temporal. It is structural. The safety case expires when the assumptions underlying it no longer hold. A component change, an architecture modification, a change to the safety requirements themselves: any of these can require full re-entry into the safety lifecycle. The standard specifies three categories of modification and what each one demands by way of revalidation. The burden of proof sits with the operator. The default is invalidation.

AI governance documents have no equivalent. A governance document written to govern one version of a system does not automatically expire when the system is upgraded. It continues to assert authority over a context it may no longer correctly describe. IEC 61508 has a name for this: an invalid safety case. AI governance has no name for it at all.

Safety engineers have been solving this problem since the 1970s. They did not arrive at the answer quickly. They arrived at it after incidents and investigations that showed what happens when a safety case outlives its assumptions. The documentation of that process is public, precise, and available to anyone who wants to learn from it.

The AI governance field is not learning from it. The governance documents going into production today carry no staleness declaration, no event-based trigger, no re-evaluation threshold. They assert permanent validity without stating the assumptions that validity depends on.

That is the precise failure mode IEC 61508 was built to prevent.

PromptQ scores governance documents on whether they declare their own staleness conditions. Most do not. The waitlist is at promptq.ai.
