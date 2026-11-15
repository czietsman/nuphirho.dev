---
title: "Can the person who approves this understand it, and stop it?"
slug: "ai-act-oversight-test"
tags: [ai-governance, organisational-design]
publish_date: 2026-11-18
cover_image_prompt: |
  A full explainer infographic in the style of a colourful, cheerful, hand-drawn line-art poster on light cream paper. Title at the top in friendly hand-lettered text: "The Oversight Test". A large friendly question card in the centre reads "Can the person who approves this understand it? Can they stop it?" with two bright tick boxes beneath it labelled "Understanding" and "Authority". On the left, a stack of paper forms funnels into a small dark box labelled "Black hole". On the right, a tidy flow of rounded boxes labelled "Pull request", "Gate results", and "Deployment plan" feeds into a smiling board table labelled "Board sees the evidence". A small caveat tag reads "A standard to measure against, not a rule about approvals". A closing hand-lettered line at the bottom: "Keep the purpose. Drop the duplicate step." Bright coral, turquoise, sunny yellow, leaf green, sky blue, and warm orange accents, with dark legible lettering, rounded boxes, lively arrows, and friendly icons. Happy, welcoming, and energetic; not institutional, clinical, or ominous. No people, no photographic elements, no logos, and no watermark.
linkedin_url:
---

Many organisations run on rules they invented, built on their own sense of what control requires. The clearest case I know is an approval chain where the person signing cannot understand what they are signing. I have watched a senior executive approve a ticket with no idea what it was about.

The EU AI Act says what oversight means. For high-risk AI systems, Article 26(2) requires that oversight go to people "who have the necessary competence, training and authority, as well as the necessary support". Those duties apply to high-risk AI systems only, but they make a good standard to hold any approval up against. So here is the test. Can the person who approves this understand what they are approving, and do they have the authority to stop it? If not, the approval is theatre.

A change advisory board form is a good example. What it is for is real: what is changing, what could break, who is affected, how to roll back. But to the person making the change it is a black hole, the people who receive it usually lack the context to judge it, and it duplicates a record that already exists, a reviewed, gated pull request.

The deeper point is that approval was never the safety mechanism. Fast detection and recovery are. We once changed a support portal's address from the vendor's to a custom one, which was the right call. The vendor had not documented how the switch would behave, the platform could not serve the old and new addresses at once, and the old address was configured in other systems, which only became obvious after the fact. No review could have rated that risk. We noticed quickly and responded quickly, and that is what made a breaking change a fixable issue.

Take the test to your own approvals. Which of them would pass?
