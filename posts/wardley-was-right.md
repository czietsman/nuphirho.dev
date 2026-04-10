---
title: "Wardley Was Right"
slug: wardley-was-right
publish_date: 2026-04-11
tags:
  - ai-assisted-development
  - agentic-development
  - software-engineering
---

**Simon Wardley** changed the way I think about a problem I thought I had solved.

He published a post a few weeks ago about what happens to software engineers when agents write the code. His answer: the role survives, the title probably does not. The people who matter are the ones who can maintain a chain of comprehension over a billion lines of code they did not write and cannot read. He called them AI wranglers, agentic herders. Alice, the software engineer you fired because a thought leader said you no longer needed her.

I read it and posted a comment. I was confident. Executable specifications, I said. If the spec proves itself in the pipeline, nobody needs to read the code. The human writes intent. The agents write code. The pipeline is the reviewer.

Wardley replied in four words: "agreed, up until the specification part."

He was right.

Not entirely right. But right enough that I spent the next morning following the challenge instead of defending the position. That is usually more productive.

---

The chain of comprehension problem is real. At a billion lines of code, reading is not a strategy. Something has to carry the comprehension that used to live in the engineer's head. My argument was that executable specifications are that something. A BDD scenario that runs in the pipeline is a piece of comprehension that is smaller than the code, verifiable, and does not depend on any individual's memory of what was intended.

Wardley's pushback was precise: it is not clear yet what those emerging practices are.

He was pointing at a scale problem I had underweighted. Executable specifications work within a bounded context. Within a single service, a single domain, a single team's scope of ownership, BDD scenarios that run and mutate are a mature enough practice to defend. But across a billion lines the scenario explosion problem makes it unmanageable. The number of valid scenarios grows with the product of features and their interactions. Across domain boundaries, that number is combinatorially out of reach. Nobody writes BDD scenarios for what happens when the system fails at that scale in a specific sequence. Who owns that scenario?

The organisations operating at that scale do not use specifications as the primary comprehension mechanism. The pipeline is still the reviewer. But what the pipeline is checking is not behaviour. It is properties.

That is a different thing.

---

The distinction that came out of following the challenge is between behavioural comprehension and property comprehension. A BDD specification captures behavioural intent: given this context, when this action, then this outcome. It is precise, it is owned, and it fails deterministically when the behaviour changes. A fitness function captures a system property: latency stays below this threshold, coupling between these services stays below this score, security posture does not degrade. It runs across the whole system and fails when the architecture drifts.

Both are verifiable artefacts. Neither is a human reading code. But they operate at different levels and answer different questions.

Behavioural specifications within domains. Contract tests at service boundaries. Architectural fitness functions across the system. Observability in production where the emergent behaviour is too complex to specify in advance. Each layer comprehends its own scope. The layers compose.

Wardley's craft-to-engineering-discipline shift is the right frame for what happens next. SREs did not stop caring about servers. They started caring about the system properties that matter at scale: reliability, observability, deployment safety. The shift from software engineer to whatever comes next is the same shape. Stop caring about the code. Start caring about the specifications at the domain level, the fitness functions at the system level, and the governance model that connects them.

The question I did not have an answer to when I posted that comment, and still do not have a complete answer to now, is what the governance model looks like when the agents are generating both the code and the specifications simultaneously. The comprehension hierarchy I described above assumes humans are writing the specifications that the agents implement. What happens when the agents are writing the specifications too?

That is the open question Wardley's pushback correctly identifies. The practices are co-evolving. What they settle into at a billion lines of AI-generated code, governed by AI-generated specifications, is not yet known.

What I am more confident about is the principle: comprehension has to live in verifiable artefacts, not in people's heads. What those artefacts turn out to be at each level of scale is still being worked out. But the direction is clear. You are not going back to reading code. The question is what you are going forward to.

That is the problem I am working on.
