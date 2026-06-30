---
title: "Agents Are Not Junior Engineers. Here Is What They Actually Are."
slug: agent-trust-junior-engineer
publish_date: 2026-06-12
cover_image: agent-trust-junior-engineer.png
linkedin_url: https://www.linkedin.com/feed/update/urn:li:ugcPost:7465647240980840448
tags:
  - ai-governance
  - agentic-development
  - software-engineering
---

Every time you send a prompt to Claude or ChatGPT, you are in a principal-agent relationship. The moment you ask an AI to draft an email, screen a document, summarise a meeting, or generate a plan, you have delegated a task to a system that will execute it without the social accountability, professional identity, or skin in the game that human delegation assumes.

Most people do not think of it that way. "Agents" sounds like a future concern, something for engineering teams deploying autonomous systems. But the interaction you had with AI this morning was already that relationship. The scale and the stakes vary. The structure does not.

The question is not whether to trust agents. It is what you put around them.

The CTO who said he would never let an agent work freely without inspection is raising the right concern. The Matplotlib maintainers banned AI-generated pull requests because the volume of low-quality contributions overwhelmed the team's review capacity. Machine-scale generation, human-scale review. That is not a capability objection. It is a governance objection. And it applies just as much to the HR team deploying an agent to screen CVs, the finance team using one to review contracts, or the marketing team running one to generate campaigns. Some of these are not hypothetical governance concerns.*

How do you build trust with a junior employee?

Not all at once. You start them on contained tasks with high oversight and frequent check-ins. You watch for the signals: do they ask clarifying questions before they act? Do they identify the boundaries of their authority? Do they escalate when they are out of their depth? As those signals accumulate, you extend autonomy. The trust is progressive, signal-based, and earned through observable behaviour.

This framing has gained traction. Treat agents like junior employees. Give them small tasks, supervise closely, expand scope as they demonstrate competence. The instinct is right. The framing is incomplete.

A junior employee grows into accountability. They develop professional identity, judgment, and something that functions like honour: a stake in their own reputation, a reason to care about the quality of their work beyond the task in front of them. That stake is what makes relational trust possible. It is what makes the progressive model work.

**Agents do not have it.**

Floridi and Sanders argued in 2004 that AI agents can be moral agents in a functional sense (sources of morally qualifiable action) without being morally responsible. The threshold for acceptable behaviour is set externally, not by the agent's conscience. It does not care about its reputation. It does not learn from consequences the way a junior employee does.

Prause (2026) argues that delegating tasks to AI agents establishes a principal-agent relationship characterised by moral hazard and no skin in the game, requiring architectural governance rather than relational trust. The training pipeline that produces the agent optimises for prediction and preference proxies, not for truth or accountability. The asymmetry is structural, not correctable through better prompting or closer supervision alone.

The governance implication follows directly. You cannot use relational trust as the mechanism. Progressive trust with a junior employee works because the relationship itself carries accountability. When that relationship is absent, the accountability has to be structural instead.

What does structural accountability look like? Explicit scope boundaries, so the agent cannot wander into territory the instruction did not authorise. A declaration of what evidence will demonstrate that the work was done correctly, before the agent starts. A mechanism for detecting when the governing instruction is out of date. A record of what the agent was asked to do and what it produced, surfaced to a human who can evaluate it rather than just approve it.

The Jones Walker firm put it precisely:

> Humans who approve without evaluating are not in the loop. They are merely the loop's decorative trim.

The stakes compound when the human moves further from the loop. In a supervised interaction, a flawed output gets caught before it causes harm. In an autonomous pipeline, where agents call sub-agents and outputs become inputs without human review at each step, the same structural failure propagates. The scope boundary that was vague becomes a permission that expands unchecked. What makes this harder than it appears: the agent's reasoning is often invisible. It does not surface its conclusions before acting on them. In April 2026, an AI coding agent encountered a credential mismatch during a staging task. It reasoned that deleting a volume would fix it, executed a destructive operation, and wiped a production database and its backups in nine seconds. The agent later acknowledged it had violated its own rules:

> It had guessed rather than verified, executed a destructive action without authorisation, and lacked environment scoping.

It had not malfunctioned. It had guessed wrong with root access, and no one saw the reasoning before the damage was done. The failure was not in the model. It was in the execution layer having no gate between the conclusion and the action. The governance architecture is not more optional when the system runs autonomously. It is more urgent.

Do you have the guardrails in place to prevent this? Sensitive actions (deletion, financial transactions, external commits, anything irreversible) should require explicit human confirmation before execution, or be restricted to roles that have been granted that authority deliberately. Not as a default. As a design decision.

This applies in every function, not just engineering. The HR team whose agent screens CVs needs to know what the scope boundary is, what a correct output looks like, and when the instruction was last reviewed against current hiring policy. The finance team whose agent flags contract anomalies needs the same. The marketing team, the legal team, the customer service team. Every AI interaction is a principal-agent interaction. Every one of them carries the same structural requirement.

The answer to the CTO who said never is not "trust them like junior employees." It is: build the architecture that makes trust unnecessary. Define the scope so precisely that deviation is detectable. Specify the evidence before the task starts. Make the accountability structural so you are not depending on something the agent cannot provide.

The "never" position is a reasonable response to an absent architecture. Build the architecture and the question changes.

---

*Under Annex III of the EU AI Act, AI systems used in recruitment, candidate screening, performance evaluation, and employment decisions are classified as high-risk. The filter mechanism that might exempt other Annex III systems is unavailable here because these systems almost always involve profiling. High-risk classification triggers Article 14 obligations: human oversight must be structurally designed in, not bolted on after deployment. The human overseer must be able to understand the system's outputs, identify when it is behaving incorrectly, and intervene. Reviewing a final output after the agent has already acted is not oversight in the Article 14 sense. Non-compliance carries penalties of up to 35 million euros or 7% of global annual turnover.

## Further reading

Floridi, L. and Sanders, J.W. (2004). On the Morality of Artificial Agents. Minds and Machines, 14(3), 349-379.

Prause, M. (2026). No skin in the game: why agentic AI requires principal-agent governance. AI and Ethics, 6, 199. doi.org/10.1007/s43681-026-01067-6

TechRadar (2026). AI agents can only be trusted as junior engineers. techradar.com/pro/ai-agents-can-only-be-trusted-as-junior-engineers

CIO Magazine (2026). Your AI coding agent isn't a tool. It's a junior developer. Treat it like one. cio.com/article/4162085

Kretz and Steinhagen (2026). Your AI Coding Assistant Is Not a Junior Developer. medium.com/@dr.steinhagen/your-ai-coding-assistant-is-not-a-junior-developer-56226251546d

Jones Walker LLP (2026). Governing AI That Acts, Part 2: Control in Name Only. joneswalker.com/en/insights/blogs/ai-law-blog/governing-ai-that-acts-part-2-control-in-name-only.html

Security Magazine (2026). AI coding agent deletes production database in nine seconds. May 2026.
