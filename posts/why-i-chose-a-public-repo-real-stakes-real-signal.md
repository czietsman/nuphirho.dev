---
title: "Why I chose a public repo: real stakes, real signal"
slug: "why-i-chose-a-public-repo-real-stakes-real-signal"
subtitle: "Simulated environments produce simulated learning. The discomfort of working in public is the point."
tags: [engineering-process, security, ai, personal-practice, devops]
draft: true
---

*This is part three of a three-part series. Part one covers the threat model: [Threat modeling is not just for enterprise](<!-- POST_1_URL -->). Part two covers the AI collaboration layer: [AI as collaborator: research, execution, and the boundary between them](<!-- POST_2_URL -->).*

When I decided to build nuphirho.dev as a managed software project — with a pipeline, infrastructure as code, BDD specifications, and a process that mirrors what I advocate professionally — the question of repository visibility was not incidental. It was a design decision.

I chose public deliberately. Not because it was the easiest option. Because it was the most honest one.

This post is about why that choice matters, what it revealed that a private sandbox would not have, and what it means to use a personal project as a real-stakes test bed for practices under professional evaluation.

The previous two posts covered the threat model and the AI collaboration layer. This one covers the experiment itself.

---

## The case against a safe environment

The instinct when learning or experimenting is to minimise risk. Use a private repo. Use a throwaway domain. Use test credentials. Keep it contained so that mistakes do not matter.

The problem with that instinct is that contained mistakes do not teach the same things as real ones. When nothing is at stake, the feedback loop is artificial. You learn whether the process works in the absence of consequences. You do not learn whether it holds under the conditions that actually apply when consequences exist.

I have spent enough time advocating for rigorous engineering practices — BDD, threat modeling, AI-assisted workflows with deliberate constraints — to know that the hardest question is not "does this work in theory?" It is "do I actually trust this enough to use it on something that matters?"

The only way to answer that question honestly is to use it on something that matters.

---

## What "real stakes" means for a personal blog

The stakes here are modest. Nobody is going to lose money if my pipeline breaks. My reputation would survive a misconfigured DNS record. The blast radius of the worst plausible incident — a compromised GitHub account leading to malicious content published under my name — is recoverable.

But "modest" is not the same as "none." The domain is real and costs money. The Hashnode account has a publishing history attached to my name. The Cloudflare configuration controls DNS for a domain I own. The GitHub repository is public, which means a secret committed accidentally is immediately public and must be treated as compromised regardless of how quickly it is removed.

These are real constraints that produce real behaviour. When I ran the threat model, I was not performing a theoretical exercise. I was asking what would actually happen to something I actually own if an adversary spent ten minutes looking at it. That question has a different quality than the same question asked about a throwaway project.

---

## What the public repo changes about the threat model

The most direct illustration of why visibility matters: when the repository is public, the threat model changes in ways that are not obvious until you model it explicitly.

An adversary looking at a private repository needs to compromise the account first. An adversary looking at a public repository can read the workflow YAML, understand which APIs the pipeline calls, identify which secret variable names are in use, and map the Terraform configuration before they have touched anything. The reconnaissance cost is zero.

That is not a reason to make the repository private. It is a reason to be deliberate about what the public visibility implies — and to run the threat model with that visibility as a given constraint rather than a problem to solve.

The STRIDE analysis surfaced this directly. Information Disclosure finding I2 — "workflow definitions expose pipeline logic publicly" — was partially mitigated by a control I had not explicitly thought about: GitHub restricts unauthenticated API access to workflow files, `AGENTS.md`, and `.secretscanignore`. The files are readable by a logged-in user, but the bar is higher than anonymous access. The threat model made that explicit. I would not have articulated it without running the exercise.

---

## The ordering failure

The most instructive moment in the hardening sequence was not a finding. It was a mistake.

Branch protection on `main` was configured as part of the hardening work. It was applied via the `gh` tool in one of the later prompts. But three commits landed on `main` directly before that protection was in place — the AGENTS.md security constraints and two commits in the `blog-nuphirho.dev` repository.

All three commits had the right content. All three had the right messages. They bypassed the PR gate that was supposed to be in place because the gate was not yet in place when they were made. The ordering of the hardening sequence created a window in which the process that was supposed to be enforced was not yet enforced.

In a private sandbox, this would have been invisible. The commits would have landed, the protection would have been applied, and the sequence would appear clean in retrospect. In a live repository with an audit trail — agent reports capturing what was done and when, commit SHAs linked to workflow runs — the gap is visible. The trail shows exactly where the process deviated from the intent.

That visibility is the value. A process you can only audit in retrospect is a process you can only fix in retrospect. The ordering failure here was low consequence. The same failure in a production pipeline with a live secret rotation sequence could be considerably less recoverable.

---

## The em dash

The agent caught something the human reviewer missed.

The AGENTS.md prompt contained an em dash — a violation of the style guide that governs all written content in the repository. The style guide is explicit: no em dashes. The prompt had been reviewed before being given to the agent. The review missed it.

The agent preserved the em dash because the prompt said "add exactly this." Then it flagged the conflict in its post-task report.

The report mechanism is easy to treat as ceremony — a structured receipt of what the agent did. The em dash shows it is not. It is a second review pass by a system that does not get tired, does not skim, and does not carry the reviewer's familiarity bias. The human who wrote the prompt and the human who reviewed it both missed the em dash because they were reading for meaning, not for typographic compliance. The agent found it because it was specifically instructed to flag style guide conflicts.

This is a small example of a general principle: the loop benefits from the agent reporting back into it. Not because the agent is more reliable than a human reviewer, but because it is differently reliable. The combination catches more than either does alone.

---

## Drafts in public

There is one more implication of the public repository worth naming directly. Posts committed to the repository as drafts are technically public from the moment of the commit. The `draft: true` front matter flag controls whether they are published to Hashnode and Dev.to — but the markdown source is readable on GitHub by anyone who looks.

I am comfortable with that. The thinking is already in the open. The process is designed to be transparent. A reader who finds a draft post in the repository before it is published is seeing exactly what they would see after publication, just earlier. There is no meaningful exposure that the draft status is protecting against.

This is a deliberate consequence of the public repository choice, not an oversight. It means the repository is not just the source of the blog — it is the blog, for anyone who wants to read it that way.

---

## What this experiment is actually for

The parallel to my day job is explicit and intentional.

At CyberSentriq, the same questions apply at a different scale. AI-assisted workflows with deliberate constraints. Threat modeling applied to systems that hold customer data and run security-critical infrastructure. The research/execution separation as a process design principle rather than a personal practice. Branch protection and PR gates as organisational controls rather than individual discipline.

The personal project is a test bed. It is where I run the practices under conditions I control completely, with stakes low enough that mistakes are recoverable but real enough that the feedback is genuine. When I advocate for these practices professionally, I am not advocating for something I have read about. I am advocating for something I have run, broken, fixed, and documented.

That is the value of real stakes. Not the risk, but the signal.

---

## The thinking is mine

I want to be direct about the role AI played in this.

AI helped research threat modeling methods. It drafted the agent prompts. It wrote the first versions of these posts. It caught the em dash. It flagged ambiguities in every prompt it was given and stopped when it did not have the information it needed to proceed safely.

The decisions were mine. The threat model findings were produced by running a structured analysis against a system I understand. The process design — the research/execution separation, the AGENTS.md constraints, the audit trail — came from thinking about what security actually requires in an AI-assisted workflow, not from delegating that thinking to the AI.

The thinking is mine. The clarity is a collaboration. But collaboration requires that the stakes are real enough to matter — because that is when you find out whether the thinking actually holds.
