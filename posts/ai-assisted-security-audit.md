---
title: "AI-Assisted Security Audit"
slug: "ai-assisted-security-audit"
draft: false
tags: [security, ai, architecture, software-engineering, process]
---

Last year I started working on a large, unfamiliar codebase. Go, new CI/CD pipelines, new deployment patterns. My job was integration work, not security. I was using AI to accelerate my understanding of the system, reading through services, tracing data flows, building a mental model of how the pieces fit together.

Then I started asking different questions.

Not "show me the code for this endpoint" but "should this role be able to perform that action?" Not "how does authentication work?" but "what happens if this token is presented to a service it was not issued for?" The shift was subtle but the consequences were not. That first question, about role permissions, led to a real finding that we have since fixed.

## Why outcome-based questions work

The questions that uncovered the issues were not on any security checklist. They came from twenty years of building and breaking systems, from having seen enough architectures to know where the seams are. The AI did not generate those questions. I did. But the AI made it possible to answer them across an entire codebase in minutes rather than weeks.

Most conversations about AI-assisted development focus on writing code faster. The value that compounds is pairing AI's traversal speed with experienced judgement. The AI can scan every file, trace every dependency, and test every boundary. But it needs someone who knows which questions to ask.

When I asked about role-based access, the AI showed me every service that checked permissions, every place a role was assigned, and every path where a request could reach a protected resource. It then helped me construct a proof-of-concept that demonstrated the issue. What would have taken days of manual code review took minutes of directed conversation.

## From code comprehension to security posture

I did not set out to do a security audit. I was trying to understand a system. But when you combine architectural intuition with an AI that can traverse an entire codebase in seconds, you start seeing things that would take weeks to uncover through manual review.

Within a day I had a clear picture of the security posture across development, staging, and production environments. Not a theoretical report or a list of hypothetical risks. Working proof-of-concepts that demonstrated real issues, followed by merge requests to resolve them. A prioritised security roadmap for a codebase I had never seen before, in under a day.

That claim sounds strong, and it should be qualified. The speed came from a specific combination of factors: the codebase was well-structured with clear service boundaries, the AI tooling was mature enough to handle large codebases, and the person asking the questions had two decades of experience knowing where to look. A junior engineer with the same tools would have found different things. Fewer in some areas, more in others, because my pattern matching has its own blind spots. The AI amplifies whatever expertise you bring. It does not replace it.

## What made AI useful for security

Security auditing has a property that makes it suited to AI assistance: the attack surface is defined by what the code does, not what the developer intended it to do. A code reviewer looks at a function and asks "does this do what the author meant?" A security reviewer looks at the same function and asks "what else could this do?"

That second question requires traversing the codebase for breadth rather than depth. You need to know every path that reaches a function, every input that can influence its behaviour, every assumption that upstream code makes about downstream validation. AI does this without getting tired of tracing call chains, without losing track of which services talk to which. It can hold the entire dependency graph in context and answer questions about it on demand.

The division of labour was consistent throughout: AI for breadth and speed, human for judgement and prioritisation. I decided which findings mattered, which were theoretical risks versus exploitable issues, and how to sequence remediation. The AI did the reconnaissance. I did the analysis.

## The systems view

This experience connects to a theme running through this series. The verification process I described in my earlier post on the trust barrier, the idea that you build confidence in AI output through measurable process rather than blind trust, showed up in this work. I verified each finding with a working proof-of-concept before reporting it. I tested each fix against the original exploit. Following that process is what made the findings trustworthy, not the AI's confidence in its own output. Without the proof-of-concepts, the findings would have been suggestions. With them, they were evidence.

The systems thinking I explored when writing about the transition from craftsman to factory in AI-assisted development appeared here too. This audit was not a sequence of isolated code reviews. It was a systematic traversal of a system, asking architectural questions about how components interact, where trust boundaries exist, and what happens when assumptions break. Taking that systems view is what turned an interesting AI experiment into an actionable security roadmap.

## The process gap

Every codebase has issues waiting to be found. The pattern generalises: AI traversal speed paired with experienced judgement finds things that neither does alone. The AI without the experience generates noise. The experience without the AI takes weeks instead of hours.

But finding the problems is half the story. Fixing them at scale requires a level of process maturity that most teams have not built. Dan Shapiro's maturity levels for AI-assisted development describe a progression where Level 4 and 5 represent teams with verification processes and architectural constraints around AI output. Without that maturity, fixing everything uncovered in a thorough AI-assisted audit will take far longer than it should. The same process discipline that lets you find the problems is what lets you fix them.

The security gaps are real. The process gap is bigger. The teams that start building these processes now will have a compounding advantage. The teams that wait will find themselves further behind than they expect.

We need to sharpen the axe today.
