---
title: "The policy made it to the developer. The knowledge did not."
cover_image_prompt: |
  A full explainer infographic in the style of a colourful, cheerful, hand-drawn line-art poster on light cream paper. Title at the top in friendly hand-lettered text: "The policy made it to the developer. The knowledge did not.". Draw a flow from "Compliance knowledge" through a blunt document labelled "Handle PII accordingly" to "Developer", "Integration", and "Product architecture". Show missing context falling out between the first two steps. Beside it, show a better document with cards for "Data element", "Sensitivity", "Who may see it", and "When". Add callouts: "AI implements what is written" and "The gap begins at authorship". A closing hand-lettered line at the bottom: "Move the judgement into the delivery artefact." Bright coral, turquoise, sunny yellow, leaf green, sky blue, and warm orange accents, with dark legible lettering, rounded boxes, lively arrows, and friendly icons. Happy, welcoming, and energetic; not institutional, clinical, or ominous. No people, no photographic elements, no logos, and no watermark.
slug: governance-delivery-gap
publish_date: 2026-05-15
tags:
  - ai-governance
  - software-engineering
linkedin_url: https://www.linkedin.com/posts/christo-zietsman_think-about-the-sign-up-flow-for-almost-any-activity-7460940283082432512-JL-q
---

Think about the sign-up flow for almost any product. Full name, email address, telephone number. Simple data. But the moment you ask what you can display, to whom, and under what circumstances, you are already in contested territory.

Can the MSP see the end user's email address? Can the distributor? Can your own support team? What about the user's mailbox: is that user data, organisational data, or both? At each layer of the stack the question changes, and the policy usually says the same thing: this is PII, handle it accordingly.

That answer is not wrong. It is incomplete. And the developer who receives it has no choice but to build to the most restrictive interpretation, because the policy gave them no basis for anything else.

This is not an AI problem. I have navigated it directly, in products built before AI was anywhere near the delivery chain. A policy written to satisfy a compliance audit treats a telephone number the same as a financial record. The developer treats them the same. The integration treats them the same. Every layer of the system encodes the same bluntness, and by the time you realise the product cannot do what the business needs, the architecture is already there.

The compliance officer who wrote the policy knows the difference. They have the context, the legal training, the risk judgement. They know that a full name and email in a B2B SaaS product is not the same exposure as a medical record. But that knowledge lives in the compliance function. It does not travel into the delivery chain. The policy is what travels.

This is the governance-to-delivery gap. The distance between what the compliance officer knows and what the delivery artefacts say, filled by interpretation that is invisible to the developer and absent from the document.

AI makes it worse, not by creating the problem but by removing what little buffer remained. A senior developer with domain knowledge could read the policy and apply proportionality judgement the document did not contain. An AI agent cannot. It implements what is written. The blunt policy becomes blunt enforcement at the speed of code generation, across every service, every integration, every data flow.

The fix is not better enforcement. It is better authorship. The same policy document could have described each data element in detail: what it is, how sensitive it is, who is permitted to see it and under what circumstances. That work costs the compliance officer an afternoon. It saves the developer weeks, and the support team months. A governance document written at that level of detail is not a more complicated document. It is a more useful one. It costs less to implement, less to test, less to maintain, and less to explain to every support engineer who fields a ticket because the system cannot do something the customer expected.

The problem was always at authorship time. It just used to be survivable.
