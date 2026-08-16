---
title: "AI bias is the bias you can see"
slug: "ai-bias-you-can-see"
tags: [ai-governance, research]
publish_date: 2026-08-05
format: article
subtitle: "On hiring, noise, and why the governable failure is not the one most people are watching"
cover_image_prompt: |
  A single document on a plain desk, clearly lit on one half, the other half in shadow from an unseen source. Muted palette, cool blues and greys. Still, institutional. No people, no screens, no colour accents.
cover_post: |
  The Amazon resume-screening tool that discriminated against women is the most cited AI bias story in hiring. It is also exactly what catchable failure looks like: systematic, measurable, visible in the aggregate, correctable in one place.

  The question it does not answer: what was happening in the room before the algorithm?

  Human interviewers carry their own bias. It works differently. It is not systematic in the way algorithmic bias is. It is case-by-case, varying by interviewer, by time of day, by the sequence of candidates seen before you. The same candidate, assessed by two different people or by the same person on different days, does not reliably get the same outcome. That is noise in the technical sense, and it does not show up in a disparate-impact audit.

  Algorithmic bias is macro-level: systematic, auditable, visible in the aggregate. Human bias in hiring is micro-level: idiosyncratic, distributed, and invisible to the instruments designed to catch the macro version.

  When an organisation runs a fairness audit and passes, that tells you its process has no systematic directional lean at the population level. It does not tell you that equivalent individuals are being treated consistently.

  That is a different property. And it is not being measured.

  New piece.
stop_slop: 42/50
toulmin: Track A 6/6, Track B 6/6
cover_image: ai-bias-you-can-see.png
linkedin_url: https://www.linkedin.com/pulse/ai-bias-you-can-see-christo-zietsman-behaf
---

In 2018 Amazon scrapped a resume-screening tool it had been quietly testing. The tool had learned from historical hiring data, and the historical data reflected a decade of patterns in a male-dominated industry. The result was a system that systematically downgraded applications from women. Amazon caught it, investigated it, and shut the project down.

That story is widely repeated as a warning about AI in hiring. It is, and it should be. What it demonstrates beyond the obvious is less often noted: the bias was findable. It was systematic enough to be visible in aggregate data. Once visible, it was addressable in one place. That is not how human interviewing typically works.

---

Algorithmic bias in hiring is macro: systematic, roughly consistent across cases, inherited from training data. Because it operates at the macro level, it is detectable by disparate-impact analysis. You compare outcomes across demographic groups; if the distribution is skewed enough, the skew shows up. The Amazon tool failed this test. That is the mechanism that caught it.

Human judgment in hiring fails through two mechanisms the Kahneman literature keeps distinct. The first is idiosyncratic bias: person-specific directional leans that vary by interviewer rather than being uniform across the organisation. Bertrand and Mullainathan's field experiment, sending identical CVs under names that signalled different racial backgrounds, found callback rates differed. That is bias in the strict sense: a consistent lean on a signal that should have been irrelevant. The second is noise: unwanted variability within the same person's judgment. An interviewer assessing comparable candidates across a day, or revisiting the same file at a different time, does not reliably reach the same conclusion. Neither failure mode registers in aggregate data the way systematic algorithmic bias does.

The distinction matters for measurement. Disparate-impact testing is designed to detect directional lean in aggregate outcomes. It catches the Amazon case. What it is not designed to detect is noise: unwanted variability across cases that can cancel out in group averages while still treating individuals inconsistently. If some interviewers lean one way and others lean another, and the tendencies wash out at the population level, the fairness audit passes. The individuals affected by each individual decision do not.

---

Kahneman, Sibony, and Sunstein's 2021 book Noise uses hiring as one of its primary examples of professional judgment under conditions of inconsistency. Different interviewers assessing the same candidate, or the same interviewer assessing similar candidates across the course of a day, produce variation in their conclusions that is large and largely unexplained by anything relevant to the hiring decision. That is noise, and the book is careful to distinguish it from bias: bias is a consistent lean, noise is unwanted variability. Both are problems, but they fail differently, they are measured differently, and they are not detected by the same instruments.

An organisation can run a fairness audit on its hiring process, find no evidence of systematic disparate impact, and still have a process where equivalent candidates are treated differently depending on who reviews them, when, and in what sequence. Those are two separate properties. The audit addresses one and is silent on the other.

This is not a claim that human hiring is more biased than algorithmic hiring in aggregate. The Amazon case stands as the obvious counter-evidence, and the algorithmic fairness literature is extensive. The failure modes work differently, they are measured by different instruments, and the macro-level one is the more governable of the two.

---

If you have been a candidate, a hiring manager, or both, you have almost certainly been involved in a process that would have gone differently with a different reviewer, at a different time. That is not a claim about malice. It is a claim about variability, and the research on unstructured interview reliability supports it.

The more structured the process, the more the variability is constrained. A structured interview, where every candidate is asked the same questions in the same sequence and scores are recorded before group discussion, reduces but does not eliminate inconsistency. An unstructured interview, which is still the most common format, relies on the judgment of a particular person at a particular moment. Paul Meehl's work on clinical versus statistical prediction, which has been replicated across many domains since 1954, consistently found that mechanical rules outperform expert judgment when predicting outcomes. Hiring is one of the domains where that finding has held.

When a hiring decision is framed as "not quite the right culture fit" or "just a gut feeling", what is usually being described is a conclusion that emerged from an unstructured process and cannot be fully reconstructed from the inputs. That is not always wrong: experienced judgment carries real information. But it is also the condition under which inconsistency is highest, and under which the gap between what was evaluated and what should have been evaluated is hardest to see.

The argument for a consistent, auditable process in hiring is not that it removes all bias. The Amazon case is proof that algorithmic systems can introduce their own, and the head-to-head literature finds that a system trained on biased historical decisions will reproduce that bias faithfully. The durable advantage of the algorithmic approach is not that the machine is fairer but that it is inspectable. Kleinberg and colleagues put it precisely: an algorithm forces the objective to be stated and the weights to be examined, in a way that human cognition cannot be. Observable, inspectable failure can be investigated, measured, and addressed. The reliability failure in unstructured human judgment is a documented problem in people-decisions, and it does not appear in the fairness reports.
