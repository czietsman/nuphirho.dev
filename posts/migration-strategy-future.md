---
title: "The migration strategy of the future"
cover_image_prompt: >-
  A full explainer infographic in the style of a colourful, cheerful, hand-drawn line-art poster on light cream paper. Title at the top in friendly hand-lettered text: "The migration strategy of the future". Draw a flow from "Running legacy behaviour" to "Executable specifications" to two target boxes labelled "Vue" and "React", both returning to a gate labelled "Run the specs". Add callouts: "Capture behaviour first", "Code becomes a reference", and "Failures show what is missing". Add a caveat: "Correctness is bounded by the specifications". A closing hand-lettered line at the bottom: "The migration strategy is not a rewrite. It is a specification." Bright coral, turquoise, sunny yellow, leaf green, sky blue, and warm orange accents, with dark legible lettering, rounded boxes, lively arrows, and friendly app-window and checklist icons. Happy, welcoming, and energetic; not institutional, clinical, or ominous. No people, no photographic elements, no logos, and no watermark.
slug: migration-strategy-future
publish_date: 2026-05-13
linkedin_url: https://www.linkedin.com/feed/update/urn:li:share:7433583304777490433
tags:
  - software-engineering
  - specification-driven-development
  - ai-assisted-development
---

I had an insight today that changes how I think about platform migrations.

I have been building a proof of concept with microfrontends, hosting Auth0 components in a Vue shell that pulls in our Vue-based platform natively and embeds our Angular-based platform alongside it. Both covered end to end by behavioural specifications, including the shell-to-microfrontend boundary.

Then it landed: if I can completely cover an existing application's behaviour with executable specs (every user flow, every edge case, every interaction), I do not need the original codebase anymore. The specs become the portable truth.

Hand those specs to AI. Say "implement this in Vue" or "implement this in React" or whatever the target stack demands. Run the specs against the output. If they pass, the implementation is correct. If they fail, the spec tells you exactly what is wrong.

The original application becomes a reference, not a dependency. The framework mismatch between legacy and target disappears. There is no line-by-line porting. No translation errors. No six-month rewrite project with a dozen engineers arguing about architecture.

The cost of this? Tokens. And the cost of tokens is in freefall.

Every enterprise I have ever worked in has at least one legacy platform nobody wants to touch. Too complex. Too risky. Too expensive to rewrite. Too critical to fail. Those platforms are not trapped by their code. They are trapped by the absence of specifications that capture what the code actually does.

Write the specs. Verify them against the running system. Once they pass, the behaviour is captured. And once the behaviour is captured, it can be reimplemented in any stack, by any team, in a fraction of the time and cost of a traditional rewrite.

The migration strategy of the future is not a rewrite. It is a specification.
