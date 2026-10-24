---
title: "When an agent writes the code, I am the other approver"
slug: "caio-i-am-the-other-approver"
tags: [ai-governance, software-engineering]
series: "The Chief AI Officer's Duties"
publish_date: 2026-10-29
cover_image_prompt: |
  A full explainer infographic in the style of a colourful, cheerful, hand-drawn line-art poster on light cream paper. Title at the top in friendly hand-lettered text: "One Second Person, Two Jobs". Two lanes side by side. The left lane is labelled "Capability: you cannot hide your own change" and shows a gate with a padlock between a small robot labelled "Author" and a branch labelled "Main". The right lane is labelled "Scepticism: a second mind asks why" and shows a friendly magnifying glass held over a checklist labelled "Intent, tests, review bot, scans". Above both lanes, a small dial runs from "Low risk: gates carry the weight" to "High risk: a person who can say no". A closing hand-lettered line at the bottom: "A rubber stamp looks like a control and is not one." Bright coral, turquoise, sunny yellow, leaf green, sky blue, and warm orange accents, with dark legible lettering, rounded boxes, lively arrows, and friendly icons. Happy, welcoming, and energetic; not institutional, clinical, or ominous. No people, no photographic elements, no logos, and no watermark.
linkedin_url:
---

When an agent writes the code, the risk moves to the approval. A signature on a change nobody understood looks like a control and does not work as one, which makes it worse than no second reviewer at all.

I am the second approver on what agents write for me. The process stays: pull requests, gates, best practice. What changes is what my approval looks at, and how much weight it carries.

The EU AI Act gives the cleanest definition of that weight. For high-risk AI systems, it asks that oversight go to people with the competence, training, and authority to do it, and that they stay aware of over-relying on the system's output, which it names automation bias. The reflexive sign-off is a named failure.

Those provisions cover high-risk systems only, and that is the useful part: oversight is attached where the stakes justify the cost. They govern deployed AI systems, not code review, so I borrow them by analogy. The lesson is to match the weight of the approver to the risk.

That gives me two tiers, and I keep the first small on purpose. Every human handoff adds the longest lead time, so the tier has to be surgical: security-sensitive, irreversible, shared, and regulated changes. There a person reads the change, with the understanding and the authority to refuse it. Learning the code belongs outside the review cycle. Everywhere else the weight moves to the specification and the gates. Intent is agreed upfront, tests, a review bot, and vulnerability scanning check correctness, and I decide whether the result is acceptable. That is closer to a product manager's job than a code reviewer's, and a role I have held before.

Separation of duties does two jobs that a second person happens to bundle: it stops whoever makes a change from hiding it, and it adds a second mind to challenge the judgement. With an agent as the author, the first holds because the decision to merge stays with me. An agent never merges its own change unreviewed. After review, or for routine housekeeping, it merges on my instruction.

The second job needs a check from somewhere other than the authoring agent. Studies have found language model judges favouring their own outputs, and models struggling to correct their own reasoning without outside feedback. That rules out the same context marking its own work, not an agent reviewing. In practice I run a separate review skill on most changes and a different reviewer on merge requests, and each lens catches things the others missed.

Two limits. Automated checks only restrict capability if nobody, administrators included, can merge past them. And because I both specify the intent and approve the result, a mistake in the ask passes both. On team projects a human reviewer catches that. On solo work it rests on my judgement.
