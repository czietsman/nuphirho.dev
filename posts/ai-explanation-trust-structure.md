---
title: "Not all AI explanations build trust equally"
slug: "ai-explanation-trust-structure"
tags: [ai-governance, research]
publish_date: 2026-08-14
linkedin_url:
---

When you present an AI output to someone (a colleague, a client, a stakeholder) you almost certainly add some kind of explanation. "The model says this because...", "it flagged this because the pattern matches...", "most systems like this would give the same answer." You are trying to help them evaluate the output and decide how much to rely on it.

The type of explanation you give predicts the degree of trust it produces, and the relationship runs against most assumptions about persuasion.

A study published in the Proceedings of the ACM on Human-Computer Interaction in November 2024 tested this directly. Saumya Pareek and colleagues at the University of Melbourne designed an experiment using an AI credibility assessment tool (a system that evaluates whether a claim is credible) and varied the type of explanation provided alongside the AI's output. The goal was to measure which explanation types produced the most trust and appropriate reliance in human users.

The framework they used distinguishes four types of conceptualisation validations, a term drawn from Jaccard and Jacoby's foundational work on epistemic frameworks. A consensual validation tells the user that the claim is widely accepted or that most sources agree. An expert validation attributes the judgment to a recognised authority or specialist knowledge. An internal (logical) validation explains the reasoning: here is why the conclusion follows from the evidence. An empirical validation offers observable, verifiable evidence that supports the conclusion.

The finding worth taking in is about which of these is weakest: consensual. The explanations that appealed to broad agreement or majority acceptance produced the lowest trust gains. The other three (expert attribution, logical reasoning, and empirical evidence) were roughly twice as influential. And among those three, there was not a large gap. Empirical evidence was generally strongest, but the main story is the distance between consensus and everything else.

---

This is counter-intuitive if your mental model of persuasion is that popularity signals quality. In many domains, knowing that something is widely accepted is genuinely useful information: it suggests the claim has been tested by many people and survived. But in the specific context of AI-assisted decision-making, that signal appears to carry much less weight than evidence that the AI's conclusion can be explained structurally or grounded in observable facts.

One possible explanation: consensual validation does not tell you why the AI reached its conclusion, or whether the reasoning would hold in this specific case. It tells you that the output agrees with something widely believed. For a user who is already uncertain about whether to trust the AI's judgment, "most sources agree" does not add the kind of resolution that "here is the evidence and here is the reasoning" adds.

The second finding that matters for practice is the backfire effect among users who already distrusted the AI system. For those users, providing explanations in some conditions made things worse rather than better. The researchers describe this as a boomerang effect: an explanation that is intended to build trust can, for an already-distrustful user, read as further evidence that the system is trying to persuade rather than inform, reinforcing rather than reducing resistance.

More explanation is not always better. For a user who has already formed a negative view of the AI system, layering additional justification onto outputs can deepen the mistrust rather than addressing it. The explanation lands differently depending on the trust baseline it meets.

---

The distinction the paper draws between trust and reliance is also worth keeping. Trust is the user's attitude toward the system. Reliance is whether they actually follow its recommendations. These are related but not identical, and the paper treats them separately. An explanation can increase expressed trust without increasing appropriate reliance, and it can produce reliance changes that do not track the system's actual accuracy.

The target is appropriate reliance: using the AI when it is likely to be right and exercising independent judgment when it is likely to be wrong. Not maximum trust, and not maximum reliance.

---

Leading with consensus as your primary justification for trusting an AI output ("this is what most systems would say," "this is the standard interpretation," "the industry broadly agrees") is the least effective of the available approaches. It is also the easiest default. It does not require knowing the reasoning, engaging with the evidence, or explaining the specific basis for the conclusion. It is the explanation you reach for when you do not have the others.

The stronger explanation is structural or empirical: here is what the AI looked at, here is how the conclusion follows from what it found, here is the observable evidence it is drawing on. That is also the explanation that requires the most from the person providing it. You need to understand the reasoning well enough to communicate it, rather than pointing to the output and noting that it agrees with the consensus.

For users who start from distrust, the implication is less about explanation type and more about approach: the goal is not to explain your way through resistance. Building appropriate reliance in a sceptical user requires addressing the underlying concern rather than adding justification. More persuasion when the user does not trust the persuader can produce the opposite of the intended effect.

---

The paper is about a credibility assessment AI used in a specific experimental context, and I want to be measured about how far the findings generalise. The four validation types (consensual, expert, internal (logical), and empirical) are drawn from Jaccard and Jacoby's framework for understanding how people evaluate claims generally, and Pareek and colleagues operationalised them specifically for AI explanations. The experimental context is one domain, one task type, and one participant population.

What I take from it is not a rulebook but a reorientation of a question I was already asking. When I explain an AI output, what am I providing? Am I giving the recipient something that lets them evaluate the conclusion independently, or am I providing a social signal about what other people or authorities believe? The data suggests the former is more likely to produce trust that is actually grounded in the output's quality. The latter is more available, and considerably less useful.
