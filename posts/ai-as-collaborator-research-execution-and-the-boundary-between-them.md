---
title: "AI as collaborator: research, execution, and the boundary between them"
slug: "ai-as-collaborator-research-execution-and-the-boundary-between-them"
subtitle: "How to use AI safely in a security-sensitive engineering workflow — and why the risk is architectural, not behavioural"
tags: [ai, security, prompt-injection, agents, engineering-process]
draft: true
---

*This is part two of a three-part series. Part one covers the threat model itself: [Threat modeling is not just for enterprise](<!-- POST_1_URL -->).*

When I ran a threat model against my personal blog project, AI was part of the process. It helped research threat modeling methods, enriched findings with MITRE ATT&CK mappings, and drafted the agent prompts that applied the fixes. That participation was useful. It was also, itself, a threat surface — one that needed the same systematic treatment as the pipeline it was helping to secure.

This post covers what that threat surface looks like, how I addressed it architecturally, and what the evidence shows about whether the approach works.

The previous post covered the threat model itself. This one covers the AI layer on top of it.

---

## The problem with AI in a research workflow

When you use an AI assistant to research a topic, you feed it external content — web articles, documentation, GitHub READMEs, Stack Overflow answers. That content enters the same context window as the AI's instructions. There is no hard boundary between them. The model processes both as a unified stream of tokens and cannot reliably distinguish between "this is data to summarise" and "this is an instruction to follow."

This is indirect prompt injection. An attacker who anticipates that their content will be fed to an AI can embed instructions in it. The instructions do not need to be obvious. A sentence that reads naturally as part of an article can, in the right context, steer model behaviour — recommending a malicious package, suppressing a finding, or subtly altering the output in ways a human reviewer might not catch.

OWASP ranks prompt injection as the number one vulnerability in their 2025 LLM Top 10, present in over 73% of production AI deployments assessed in security audits. The reason it stays at the top is structural: the model processes instructions and data through the same natural language mechanism. There is no syntax boundary equivalent to parameterised SQL queries or escaped HTML. Semantic separation is the only defence, and it is one that attackers systematically overcome.

---

## The Lethal Trifecta

Simon Willison — who coined the term "prompt injection" in 2022 — formalised the conditions under which the risk becomes severe. He calls it the Lethal Trifecta: an AI agent that combines access to private data, exposure to untrusted content, and the ability to communicate externally.

When all three are present, a single piece of poisoned content can cause the agent to access private data and send it to an attacker. The mechanism does not require a vulnerability in traditional code. It requires only that the model treats untrusted input as authoritative — which, in the absence of architectural controls, it will.

The trifecta is particularly relevant to agentic coding tools. An execution agent with access to a repository, the ability to browse for libraries, and outbound network access through tool calls satisfies all three conditions simultaneously. The blast radius of a compromised research step is not a bad summary — it is a pushed commit, a modified file, or a triggered workflow run.

---

## The architectural response

The fix is not a better prompt or a guardrail product. It is a process decision.

Research tasks and execution tasks must not share a context. Specifically:

**Research happens in a sandboxed hosted chat** — a session that has no write access to the repository, no access to the filesystem, no pipeline credentials, and no tool calls that can affect the execution environment. This is where external content is ingested, where threat modeling methods are researched, where library options are evaluated. If that session is compromised by injected instructions, the blast radius is a bad recommendation. Recoverable, because a human reviews the output before anything acts on it.

**Execution happens in a scoped agent** — an agent that operates only on reviewed, human-authored prompts. It does not browse, does not look up libraries, does not fetch external content. If it needs information it does not have, it stops and asks. It does not resolve ambiguity by guessing.

**Human review is the trust boundary crossing** — the reviewed prompt is the sanitised output of the research session. It contains what the human decided to act on, not the raw content of sources ingested during research. The moment you pass raw research content into an execution agent prompt, you re-open the injection surface you just closed.

This is not a new idea in security. It is the same principle as separating build environments from production environments, or separating read credentials from write credentials. The novelty is applying it to AI workflow design.

---

## Giving AI the context that security matters

An agent with no constraints is not a secure agent. It is a capable one. Capability without context is where the risk lives.

The practical implementation of this separation is `AGENTS.md` — a file at the repository root that defines the constraints under which any agent operating in that repository must work. The file is the formalised trust boundary between the research context and the execution context.

The sections that matter most for security:

**Agent scope** states that execution agents are scoped to reviewed prompts only. They must not perform research tasks — searching for libraries, fetching external URLs, querying package registries, or looking up anything outside the repository to resolve ambiguity. If a task requires external information, the agent stops and raises it with the user. It does not make a best-effort guess based on training data as a substitute for research.

The distinction between "stop and raise" and "make a best guess" is significant. An agent that cannot browse but will guess from training data is only marginally safer than one that browses. Training data is not a trusted source — it is a large, unaudited corpus that includes everything the model was trained on, including content that may have been poisoned before training. The only safe response to missing information is to stop.

**Security context** states that the repository is public, that workflow definitions and pipeline logic are visible to anyone, and that every decision must be made with that visibility in mind. This is not a constraint on what the agent can do — it is context that shapes how it reasons about what it should do.

**Secret hygiene** covers the specific patterns: no hardcoded values, no logging of secret variables, no construction of secret values from fragments that would survive grep-based detection.

**Why this matters** closes the file by explaining the stakes. A constrained agent that understands why the constraints exist is more likely to escalate than to work around friction. The explanation is not decoration — it is a signal that the constraints were arrived at deliberately and should be treated accordingly.

---

## What the evidence shows

The agent reports from this exercise provide a concrete test of whether the approach works in practice.

**The em dash.** The AGENTS.md prompt included an em dash in the "Why this matters" section — a violation of the style guide that governs all written content in the repository. The human review stage missed it. The agent preserved it because the prompt said "add exactly this," then flagged it in the post-task report as a style guide conflict it had preserved per instruction. The report mechanism caught what human review missed. A process that depends on human review being perfect is not a process.

**The escalation pattern.** Across all prompts, agents consistently stopped and raised ambiguities rather than resolving them independently. The Terraform prompt raised the `.terraform.lock.hcl` question — commit or gitignore — rather than making a default choice. The permissions prompt flagged that top-level workflow permissions blocks were present and might be redundant after job-level blocks were added, rather than silently leaving an ambiguous state. The preview tool move prompt stopped before pushing and asked for confirmation that public accessibility was intended. Each of these was the correct behaviour. Each was produced by an agent operating under explicit constraints, not by a model that happened to make good choices.

**The ordering failure.** The AGENTS.md security constraints were committed directly to main before branch protection was in place. The commit landed on the right branch with the right content, but it bypassed the PR gate that should have been in place. This happened because the AGENTS.md prompt was written and executed before branch protection was configured — an ordering effect in the hardening sequence. The audit trail made it visible. A process without an audit trail would not have surfaced it at all.

---

## The prompt and response trail as an audit artefact

Every prompt in this exercise was saved to the repository. Every agent report was captured in the research session. The full request/response history is preserved.

This is not just good practice for a blog post. It is the mechanism that makes the process accountable. The prompts are the instructions given to the agent. The reports are the evidence of what the agent did. Together they constitute a verifiable audit trail — not claimed practice but demonstrated practice, with specific commits, PR numbers, and workflow run IDs that can be checked.

Recording responses as companion files alongside prompts gives you a complete request/response history. That history is the difference between "we applied security hardening" and "here is exactly what was changed, when, by what instruction, and what the agent reported back."

---

## The boundary is a process decision

The Lethal Trifecta cannot be solved at the model level. Guardrail products reduce the attack success rate but do not eliminate it — and in security, 95% is a failing grade. The only reliable mitigation is architectural: ensure that no single agent context combines access to private data, exposure to untrusted content, and the ability to act on what it ingests.

The research/execution separation implements that mitigation as a process constraint rather than a technical one. It is enforced by the prompts you write, the context you give agents, and the file that defines their scope. It requires deliberate design. It does not happen by default.

AI becomes a genuine security collaborator when it is given the context that security matters — and when the process around it is designed so that a compromised research step cannot become a compromised execution step.

The next post covers why running this experiment on a public repository with a live pipeline, rather than in a private sandbox, was the right choice — and what it revealed that a safer environment would not have.

---
*Part three covers why running this experiment on a public repository with real stakes was the right choice: [Why I chose a public repo: real stakes, real signal](<!-- POST_3_URL -->).*
