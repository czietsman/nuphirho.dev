---
title: "A governance document is a snapshot, not a wall"
slug: governance-document-retrieval-gap
platform: linkedin
tags: [ai-governance, software-engineering]
d1_provenance: []
publish_date: 2026-09-09
cover_image_prompt: |
  A full explainer infographic in the style of a colourful, cheerful, hand-drawn line-art poster on light cream paper. Title at the top in friendly hand-lettered text: "A governance document is a snapshot, not a wall". A timeline begins with a document labelled "Policy written: known conditions", moves to a later retrieval box labelled "New content at runtime", and ends at a warning checkpoint labelled "Never anticipated by the document". A return arrow points to "Flag and re-evaluate" rather than pretending the original document blocks the new condition. Three compact callout boxes read "Authorship fixes a snapshot", "Retrieval moves later", and "Thorough drafting cannot close time". A small caveat box reads "TOCTOU is an analogy, not an equivalence". A closing hand-lettered line at the bottom: "Detect when reality has moved beyond the snapshot." Bright coral, turquoise, sunny yellow, leaf green, sky blue, and warm orange accents, with dark legible lettering, rounded boxes, lively arrows, and friendly icons. Happy, welcoming, and energetic; not institutional, clinical, or ominous. No people, no photographic elements, no logos, and no watermark.
stop_slop: 43/50
toulmin: Track A 6/6, Track B 6/6
notes: "Source: #587, findings/so-novel-findings-registry.md F24 (retrieval-layer governance asymmetry). TOCTOU (time-of-check to time-of-use) used explicitly as an analogy, not a formal equivalence, per the finding's own adjacency note and the issue's guidance. Distinct from the staleness/epoch-limits material in the-governance-document-that-never-expires.md: this is about retrieval-time content the document never anticipated, not about the document's validity horizon. Closing line replaced 6 September 2026 on Director instruction: 'Most governance documents have no answer to that question at all' became the more direct, second-person 'I doubt your governance documents can answer that question today,' addressing the reader rather than generalising about governance documents in the abstract. Added the missing d1_provenance field."
linkedin_url:
---

A governance document written today can only govern what its authors thought to anticipate today. An AI agent that retrieves information at the moment it acts, a search, a database lookup, a document pulled into context, can encounter material the document never accounted for, because that material did not exist yet, or nobody thought to ask about it, when the document was written.

Security engineers have a name for a version of this problem: time-of-check to time-of-use, the gap between when a permission is verified and when it is actually exercised, during which the world can change underneath it. The parallel is not exact, TOCTOU is about access control evaluated at the wrong moment, and this is about a governance document's authority over content it could never have anticipated. But the shape of the gap is the same: something was fixed in advance, something else moved at runtime, and the fixed thing has no way to know.

This is not a drafting failure. No amount of additional care at authorship time closes this gap, because the gap is structural: retrieval happens after the document is written, by definition. You cannot pre-specify a response to content that does not exist yet. What you can do is stop treating the document as a wall that blocks everything outside its scope, and start treating it as a snapshot, current as of the moment it was written, that needs a mechanism for flagging when what it governs has moved past what it anticipated.

The practical question this leaves is not "how do we write a more thorough policy." It is "how do we know when the policy has met something it was never written for." I doubt your governance documents can answer that question today.
