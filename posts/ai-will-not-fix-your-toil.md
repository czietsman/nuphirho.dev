---
title: "AI will not fix your toil by default. You have to aim it."
slug: "ai-will-not-fix-your-toil"
tags: [ai-governance, software-engineering]
publish_date: 2026-11-16
cover_image_prompt: |
  A full explainer infographic in the style of a colourful, cheerful, hand-drawn line-art poster on light cream paper. Title at the top in friendly hand-lettered text: "Aim It At The Toil". On the left, a friendly robot labelled "AI, default settings" happily speeds up a conveyor of code boxes while a pile labelled "Meetings, waiting, gates" stays exactly the same size. On the right, the same robot wears a small compass and points at that pile, with three tidy labelled streams: "Changes in flight", "Gates and comments", and "Plan, apply, watch, roll back". A small caveat tag reads "One year of survey data and one person's workflow". A closing hand-lettered line at the bottom: "Look at where the time goes. Point the tools there." Bright coral, turquoise, sunny yellow, leaf green, sky blue, and warm orange accents, with dark legible lettering, rounded boxes, lively arrows, and friendly icons. Happy, welcoming, and energetic; not institutional, clinical, or ominous. No people, no photographic elements, no logos, and no watermark.
linkedin_url:
---

DORA's 2024 report found that AI adoption does not, by default, reduce toil. A 25% increase in an individual's AI adoption went with about 2.6% less time spent on work people consider valuable, while time on toilsome work was essentially unaffected. Flow, job satisfaction, and productivity all rose. DORA's own conclusion is that AI has not yet helped people avoid the drudgery of meetings and bureaucracy, and it offers what it calls a vacuum hypothesis for why.

That matches where my time goes now. Planning and brainstorming up front lets me batch work from end to end. While agents implement and verify, I switch to another project and start another agent. The hard part is sequencing. Small batches mean proving each step, and a change that spans systems touches many of them, so each piece waits on its own gate: a review someone has to give, a comment that arrives without notice, a pipeline that fails, or a manual approval, including the plan step on infrastructure as code. Comments get lost in email.

I chose not to chase colleagues or ask them to change how they work. I built a tool instead that shows every change in flight, every pipeline, every gate, and every question waiting for me. AI maps the sequence of steps. An agent checks that each infrastructure plan contains only the intended change, applies it, watches it, and rolls it back if it does not take effect, while I judge whether to proceed. DORA's 2025 model found that frequent rollback amplifies AI's benefit to team performance, which is the habit I lean on.

Cleanup closes the loop. Leftover artefacts confuse colleagues and other AI, which then build on things nobody should rely on. Leaving the system in a better state is what lets the next person, or the next agent, avoid being sent down old paths.

My hypothesis, and it is only that, is that small batches multiply the number of concurrent changes, each blocked on its own gate, and a change that spans systems behaves like a transaction with nobody coordinating it. That overhead is what the tool is aimed at. I do not know of research that measures it.

AI will not aim itself at your toil. On DORA's 2024 evidence, the default effect is faster code and the same meetings, bureaucracy, and waiting. Someone has to look at where the time goes and point the tools there.

There are two places to point them. The first is your own workflow, and that is where I always start, because it is where I can prove a change for myself. The second is the team's and the organisation's workflow. I work on that too. The stakes are higher there, so the reward is higher, but it is slower, and I find it easier to argue for something I have already proven at home.

Where does the time go in your week now, and which part of it have you pointed AI at?
