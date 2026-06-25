---
title: "Gold In My Citation Framework"
slug: promptq-d1d5-use-case
publish_date: 2026-06-08
tags:
  - ai-governance
  - promptq
  - research
---

Three months ago I had no citation verification framework at all. I did not know I needed one yet.

I knew AI made things up and I have read enough news about lawyers getting caught out where they did not vet the referenced cases. I had also run into this when planning out features and got API endpoints which never existed. AI models fabricate facts and citations with enough plausibility to pass a casual check. Confident title, plausible author, real-looking DOI. The paper does not exist. Or it exists but says something entirely different.

I realised this early, because I was preparing to publish a paper. There was no way I would include a bad citation in something going on the public record. So I built a framework to catch them.

The framework itself went through multiple rounds. Every new fabrication pattern I found was a reason to re-evaluate and find a solution. The citations in the published papers were manually checked on top of the framework. When something is going on the public record, good enough is not good enough.

### The framework

The framework consists of a few checks.

1. **Resolve** checks the URL loads.
2. **Trust** checks the source is credible.
3. **Match** checks the title is correct.
4. **Describe** checks the content is as described.
5. **Verify** checks the specific claim is actually in the document.

That framework was already a game changer. Citations that would have slipped through now get caught. The research became more reliable. And then something interesting happened.

Today I applied PromptQ to the framework itself.

PromptQ is a structural evaluation tool for governance documents. The citation verification framework is a governance document. Applying PromptQ to it was the scientific method applied to my own tools: here is the hypothesis, here is the instrument, now turn the instrument on itself and see what happens.

It scored 43%.

> Finding a gap in your own tools is a gold nugget. It means the instrument is working.

That result was not discouraging. It was exciting. Finding a gap in your own tools is a gold nugget. It means the instrument is working. It means there is unexplored terrain right in your own backyard. And it means you can measure the improvement.

Three gaps stood out.

The first was pass/fail logic. What happens when **Resolve** passes but **Verify** fails? The source exists but the claim attributed to it is not in it. That is the most dangerous pattern in citation verification: a real source wrapping a fabricated claim. The original framework had nothing to say about it. The remediated version makes it explicit: a **Verify** failure is Overreach, regardless of the **Trust** score. A credible source that does not contain the attributed claim fails.

The second was scope boundary. Three edge cases walked into the research and the framework had nothing to say about any of them. Paywalled sources: does a landing page count as a **Resolve** pass? Post-verification sources: a paper published after the finding's verification date satisfies **Resolve**, **Trust**, and **Match** but is not prior art. Non-English sources: do **Describe** and **Verify** require translation? In each case, the researcher was left to decide. And researcher decisions at the edges of a verification framework are exactly where confirmation bias lives.

The third was re-evaluation triggers. A verification framework with no expiry date will silently become wrong, or miss the opportunity to be improved. AI model citation behaviour changes. New fabrication patterns emerge. The original framework had no mechanism to detect that it had become inadequate. The remediated version declares three trigger types: event-based, evidence-based, and time-based.

> Before: 43%. After: 86%.

### Scorecard

**Original framework: 3.0/7**

| Principle | Score | Notes |
|---|---|---|
| P1 Success Definition | 0.5 | Dimensions defined; no pass/fail logic or verdicts |
| P2 Assessment Rubric | 0.5 | Partial criteria; no 0/0.5/1 scale |
| P3 Scope Boundary | 0 | No negative boundary; three edge cases undefined |
| P4 Data Classification | 0.5 | Source types implicit; no classification table |
| P5 Quality Gate | 0.5 | Framework is a gate; no precedence rule |
| P6 Internal Consistency | 0.5 | No precedence rule for Verify over Trust |
| P7 Contextual Currency | 0.5 | No re-evaluation triggers declared |
| **Total** | **3.0/7** | Below threshold |

**Remediated framework: 6.0/7**

| Principle | Score | Notes |
|---|---|---|
| P1 Success Definition | 1.0 | Pass/fail logic and citation verdicts explicit |
| P2 Assessment Rubric | 1.0 | 0/0.5/1 criteria defined for each dimension |
| P3 Scope Boundary | 1.0 | In scope, out of scope, and edge cases defined |
| P4 Data Classification | 0.5 | Source classification table provided; handling partial |
| P5 Quality Gate | 0.5 | Framework is a gate; no second-reviewer requirement |
| P6 Internal Consistency | 1.0 | Precedence rule defined (Verify over Trust) |
| P7 Contextual Currency | 1.0 | Three trigger types declared |
| **Total** | **6.0/7** | Above threshold |

The best part? I can now rerun the original citations through the improved framework and measure the delta. The scientific method applied: hypothesis, instrument, result, refinement, retest. That cycle is what makes a research programme rather than a pile of notes.

Our instinct when we find a gap is to remediate it and share what we found. If it happened to us, it is happening to others. Perhaps we can help everyone improve.

The lesson is not that the original framework was wrong. It was doing its job. It just had not been to the edges yet. PromptQ forced the question before the edge cases arrived and caused real damage.

The waitlist is at promptq.ai.
