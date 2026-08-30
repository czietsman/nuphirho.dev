---
title: "Where aviation's framework breaks down, and where it does not"
cover_image_prompt: |
  A full explainer infographic in the style of a colourful, cheerful, hand-drawn line-art poster on light cream paper. Title at the top in friendly hand-lettered text: "Where aviation's framework breaks down, and where it does not". Draw a split diagram. On the left, "AI system" leads to "Statistical behaviour" and a broken traceability arrow labelled "Framework breaks down". On the right, "Governance document" leads to "Static human-authored artefact" and three checkable cards: "Linkage", "Validity boundary", and "Evidence architecture". Add a source box: "DO-178C; arXiv:2606.25120". A closing hand-lettered line at the bottom: "System limits do not erase document-level controls." Bright coral, turquoise, sunny yellow, leaf green, sky blue, and warm orange accents, with dark legible lettering, rounded boxes, lively arrows, and friendly icons. Happy, welcoming, and energetic; not institutional, clinical, or ominous. No people, no photographic elements, no logos, and no watermark.
slug: "aviation-li-2-separation"
series: "aviation-specification-completeness"
series_part: 2
tags: [ai-governance, research]
publish_date: 2026-07-09
linkedin_url: https://www.linkedin.com/posts/christo-zietsman_the-aviation-certification-community-has-activity-7480871615862411266-0B0k
---

The aviation certification community has been clear that its frameworks break down for AI systems. DO-178C, the standard governing aviation software certification for thirty years, requires bidirectional traceability: every requirement traces forward to a test and backward to a source requirement. A language model's behaviour is statistical rather than deterministic. You cannot trace an inference to a requirement. The framework breaks down.

This is settled. The standard-setting bodies have said so.

What has not been addressed is a different question.

Aviation's framework breaks down at the system level because AI systems are non-deterministic. A governance document, the system prompt or policy document specifying what the system is authorised to do, is not the AI system. It is a static artefact. It was written by a human at a specific time. Its structural properties can be evaluated and required independently of the stochastic system it governs.

The question of whether aviation's structural requirements transfer to AI governance documents is not the same question as whether they transfer to AI systems. The first is settled. The second has not been examined.

A new paper makes this separation explicit. It argues that three structural requirements from aviation certification, structured governance linkage, context-bounded validity, and an objective evidence architecture, are transferable to natural language governance documents for the same reason they were required of software requirements documents: both are specifications that a conforming agent must be derivable from, and both can fail in ways that are structural, identifiable, and consequential at authorship time.

The aviation framework fails when applied to AI systems. It does not fail when applied to the documents that govern them.

That separation is the paper's primary conceptual move.
