---
title: "Who should approve a change?"
slug: "who-should-approve-a-change"
tags: [ai-governance, organisational-design]
publish_date: 2026-11-20
cover_image_prompt: |
  A full explainer infographic in the style of a colourful, cheerful, hand-drawn line-art poster on light cream paper, square format. Title at the top in friendly hand-lettered text: "Who Should Approve A Change?". On the left, a small round table labelled "Change board" with a stack of forms and a tag "Every change, same form". On the right, a pipeline of rounded boxes labelled "Peer review", "Tests", "Rollback", with a tag "Evidence the pipeline already produces". Above both, a signpost labelled "Executives: set the conditions" pointing to three tags: "A clear stance", "Guardrails", "Competent deciders". A closing hand-lettered line at the bottom: "Approve the risky few. Equip everyone else." Bright coral, turquoise, sunny yellow, leaf green, sky blue, and warm orange accents, with dark legible lettering, rounded boxes, lively arrows, and friendly icons. Happy, welcoming, and energetic; not institutional, clinical, or ominous. No people, no photographic elements, no logos, and no watermark. Keep all text well inside the edges.
linkedin_url:
---

Most change approval asks the wrong person, and leans on the wrong safety mechanism. Approval was never what kept systems safe. Fast detection and recovery are.

A board keeps a log of change events, and what it needs is a graph of dependencies. The act that creates a dependency, pointing one system at another's address, is almost never filed as a change. Who would list every page where an address was added, so that a board could know? Should a board hold that detail at all? The real answer is an up to date dependency graph, written by the system, with tests that check a dependency when the thing it depends on changes, not only when the artefact that contains it is updated.

Where you can, make being wrong cheap and easy to detect. Parallel change lets old and new run side by side. Where a platform cannot do that, the change deserves extra care, whoever approves it. If you are building a platform, make changes reversible, and when someone adds a custom address, let the old one keep working.

The research and the frameworks point the same way. The 2019 State of DevOps report found that external change approval, such as a change advisory board, was associated with worse delivery performance, with respondents 2.6 times more likely to be low performers, and found no evidence that it lowered change failure rates. That is survey evidence of association, from one edition. ITIL 4 says changes should be assessed by "people who are able to understand the risks and the expected benefits", and notes that high-velocity organisations commonly decentralise approval to peer review. DORA's 2026 report adds that as AI raises code volume, change approval has to scale with it, which takes the board as given instead of asking whether it is the right mechanism.

So the executive role changes, from approving each change to setting the conditions changes run under: a clear stance, guardrails, and the competence of the people deciding. The EU AI Act's Article 4 requires providers and deployers of any AI system to take measures to ensure a sufficient level of AI literacy among their staff. DORA names a clear and communicated AI stance as one of seven capabilities that amplify AI's benefit, and finds that more frequent rollback amplifies AI's effect on team performance.

The board's attention then goes to the changes that are genuinely risky or cross-cutting, using the evidence the pipeline already produces. Everything else is approved by the people who can understand it, and made cheap to undo.

Which of your approvals would you hand back to the team, and which deserve the board's full attention?
