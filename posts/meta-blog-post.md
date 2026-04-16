---
title: "Treating Ideas as Releasable Software"
slug: "treating-ideas-as-releasable-software"
draft: false
tags: [ai, process, scientific-method]
publish_date: 2026-03-25
---

I recently published a four-part series on AI code review, correlated errors, and specification-driven verification. The research behind it took longer than the writing. This post is about that research process, not the conclusions it produced.

The argument is simple: published ideas have the same failure modes as shipped software. They go out with known gaps because the author ran out of time or confidence. They make claims the evidence does not support. They patch weaknesses with emphasis rather than fixing them. They pass the author's own review because the author wrote them.

I wanted to see what happens when you apply engineering rigour to content development. Not writing advice. Not style guides. Actual process: hypotheses, verification, quality gates, and honest assessment of what the evidence does and does not support. What follows is past tense. The series is published. The process is complete.

Here is what that looked like.

## Hypothesis formation through dialogue

The three core ideas in the series did not start as outlines. They emerged from research conversations where claims were made, challenged, refined, and scoped. The correlated error hypothesis started as an observation about AI code review tools. The Cynefin framing started as a question about complexity theory. The residual defect taxonomy started as an honest admission that executable specifications cannot catch everything.

The important moves happened early. Separating structural claims from empirical ones. Identifying where an argument was circular at the level of definitions rather than genuinely supported. Recognising when an observation needed a qualifier rather than a bolder statement. These are the same moves you make in a design review, and they matter just as much in content.

## Science officer review

Before any drafting, the hypotheses were stress-tested using a dedicated challenge role. The brief was explicit: do not validate, try to break them.

The results were specific. The correlated error claim needed a bridging sentence connecting ensemble machine learning theory to LLM pipelines. The Cynefin framing needed the determinism objection engaged directly, not sidestepped. The economic argument needed data, not assertion. The taxonomy needed "to the best of our knowledge" before the novelty claim.

Each of these is a gap that a validation-oriented review would have missed. A reviewer looking for confirmation finds confirmation. A reviewer tasked with finding weaknesses finds weaknesses. The role matters more than the competence of the reviewer.

## Targeted literature search

Two targeted research briefs were sent to Grok, each scoped with specific context about the hypotheses and explicit instructions about what to look for and what to ignore. The first brief asked for academic papers on correlated errors in LLM pipelines, Cynefin applied to specifications, and the limits of specification. The second asked specifically for evidence on specification authoring cost, because a Toulmin assessment had identified the economic argument as assertion-dense rather than evidence-dense.

These were not broad searches. The briefs included the current state of the argument, the specific gaps, and instructions to look outside the mainstream SDD and vibe coding conversation. The papers that closed the last Toulmin gap, Fonseca et al. and Hassani et al. on specification authoring cost, came back from the second brief. They were found because the search was scoped by the review process, not by general curiosity.

The research brief for each search was precise about what it was looking for and what it was not. This turns out to matter enormously. A broad search for "AI code review" returns hundreds of vendor marketing pieces and opinion posts. A search for "correlated errors in LLM ensembles applied to code generation pipelines" returns three papers, each directly relevant.

The SGCR paper required a correction before citation. The 90.9% figure is the adoption rate among developers who used the spec-grounded framework, not a defect detection improvement. That distinction matters. One version supports the claim. The other overstates it.

One paper from the initial search (arXiv:2603.00311) turned out to be about regex engines, not metamorphic testing in the relevant sense. It was removed. The cost of including a bad citation is higher than the cost of having one fewer reference. Precision in the search brief determined the quality of what came back. Garbage in, garbage out applies to research as much as it applies to data pipelines.

## Citation verification

Every cited paper was verified against its original abstract before inclusion. Author lists, key figures, and specific findings were checked. The 90.9% correction, a 37.6% performance loss figure, the "popularity trap" terminology, each was confirmed in the source before being used.

Three numbers from an intermediary research return required PDF-level verification, which was completed against the original papers before publishing. The distinction between what was confirmed from the abstract and what came from an intermediary's interpretation was tracked explicitly. I know which claims I can stand behind and which ones still need a primary source check.

This is tedious. It is also the difference between a post that survives scrutiny and one that does not.

## Toulmin framework as quality gate

Five of Toulmin's six elements (claim, data, warrant, rebuttal, qualifier, with backing folded into the data and warrant assessment) were applied to each piece before drafting the final versions. This is not a writing technique. It is a structural verification tool.

The warrant failures were the most important findings. The economic argument was assertion-dense rather than evidence-dense. The constraint transformation claim was circular at the level of definitions. A Category D claim used "arguably" where a citation or an explicit provisionality flag was needed.

Toulmin did not find these problems. Applying it systematically did. The framework forces you to ask "what connects my evidence to my claim?" for every claim in the piece. When the answer is "it feels right" or "it is obvious," you have found a gap.

## Stop-slop scoring

I use a scoring framework (Hardik Pandya's stop-slop skill) that evaluates prose across five dimensions: directness, rhythm, trust in the reader, authenticity, and density. Each is scored 1 to 10. Below 35 out of 50, the piece goes back for revision.

One post scored 30 on the first pass. The specific failures: a formulaic "often described as X, but actually Y" opener, repetitive staccato paragraph endings, pre-emptive hedging that signalled a weak argument rather than honest qualification, and an unsupported economic claim presented with confident language.

Each failure had a specific fix. The scoring made the prioritisation tractable. Without a framework, "this doesn't feel right" is the feedback. With one, "the rhythm score is 5 because three consecutive paragraphs end with punchy one-liners" is the feedback. The second version is actionable.

## Revision and re-scoring

The revised pieces were scored again. The 30 moved to 38. All four pieces passed the threshold. The one remaining gap in the Toulmin assessment, the economic argument, was closed by a targeted literature search that was run specifically because the review identified the gap.

One late-stage catch is worth describing because it illustrates how the review process works. The correlated error argument, as originally phrased, left a gap that an argumentative reader could exploit: "but you used different AI models for generation and review, so the correlation breaks and your claim falls apart." The argument does not fall apart, but the text needed a sentence distinguishing model diversity from ground truth. Using a different model is not the same as using an external specification. Diversity reduces correlation. It does not eliminate it. Only an external reference, a specification that defines intent independently of the code, breaks the circular validation. That gap was found during review, reasoned through, and closed with a single sentence. Drafting missed it. The process caught it.

The loop closed. The gap in the argument drove the research. The research filled the gap. The revised piece passed both the structural and the editorial quality gates.

## What this process is and is not

This is expensive for a single blog post sharing a personal experience. It is not expensive for a practitioner publication making novel claims in an active research area and citing more than a dozen empirical sources.

The discipline is not the process itself. It is knowing which level of process a given piece warrants and applying it without shortcuts when the stakes are high. A post about my experience with a security audit needs authenticity and clarity. A post proposing a novel defect taxonomy and connecting it to the oracle problem in formal testing theory needs every claim verified, every citation checked, and every logical gap identified before publication.

The same principle applies to software. Not every service needs mutation testing. Not every feature needs a formal specification. But the ones that matter, the ones where failure has consequences, need the full process. The skill is knowing which is which.

AI made this process possible at a pace that would have been impractical without it. The literature search, the cross-referencing, the citation verification, the structural analysis. All of it happened faster than I could have done manually. But the decisions about what to research, which arguments hold, which caveats matter, and when the piece is ready to publish, those were mine.

The thinking is mine. The rigour is a collaboration. And the output is still what matters most.
