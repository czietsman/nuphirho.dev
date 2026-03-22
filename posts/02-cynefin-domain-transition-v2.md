---
title: "From Complex to Complicated: What Executable Specifications Actually Do"
slug: "executable-specifications-cynefin-domain-transition"
draft: false
tags: ["cynefin", "specifications", "complexity", "ai", "software-engineering"]
allow_emdash: true
series: "The Specification as Quality Gate"
---

*This is Part 2 of a four-part series, "The Specification as Quality Gate." [Part 1](/echo-chamber-ai-code-review-correlated-error) developed the correlated error hypothesis. This post grounds the specification-first argument in complexity science. Parts 3 and 4 will follow.*

Ask an engineering team why they write tests, and most will say: to catch
bugs. Ask why they write specifications, and the answers are less coherent.
Documentation. Compliance. Handover. The process said so.

None of those answers explain what a specification actually does to the
nature of the problem being solved. Dave Snowden's Cynefin framework offers
the clearest account. The argument here is that executable specifications
do not simply improve quality; they change the class of problem you are
working on.

---

## Two Kinds of Problem

Cynefin distinguishes problem domains by the relationship between cause
and effect. Two are relevant here.

In the complicated domain, cause and effect are knowable through analysis.
A qualified expert can study the system, apply the right methods, and
arrive at a good answer. The problem is hard but it is tractable. Engineers
building against a well-specified API are working in the complicated domain.
The system's behaviour is knowable if you do the analysis.

In the complex domain, cause and effect are only knowable in hindsight.
The system responds to interventions in ways you cannot fully predict
before they happen. You probe, observe what emerges, and respond to what
you find. Raising a child is the canonical example. But so is asking an
AI agent to implement a feature from a two-sentence description. The
output is emergent. The only way to know if it is right is to observe
what it does.

This is not a question of difficulty. It is a question of ontological
character. Complicated problems are hard. Complex problems are of a
different kind.

---

## What Separates the Domains

Snowden distinguishes the domains not just by cause-and-effect
relationships but by the constraints that govern system behaviour.

In the complicated domain, governing constraints apply. These are the
rules, standards, and boundaries that define acceptable behaviour. They
do not specify every action, but they bound the solution space tightly
enough that analysis can identify a good approach.

In the complex domain, enabling constraints apply. These are looser.
They allow a system to function and create conditions for patterns to
emerge, but the behaviour of the system is not derivable from the
constraints alone. You must observe it.

An AI agent operating from a vague natural language prompt operates under
enabling constraints at best. The prompt allows the agent to function and
creates conditions for code to emerge, but edge cases, boundary conditions,
and architectural choices are resolved by the model's priors, not by the
stated requirements.

An AI agent operating from a precise executable specification is in a
different situation. A BDD scenario makes a specific causal claim: given
this precondition, when this action occurs, then this outcome is required.
That claim is verifiable before hindsight. The agent cannot legally produce
code that fails the scenario. When the pipeline runs, the question "does
the system do what it should?" has a deterministic answer. Cause and effect
are knowable through analysis of the specification itself, which is
exactly what defines the complicated domain.

Writing executable specifications converts enabling constraints into
governing constraints. The problem does not change. The constraint type
does. And with it, the domain.

---

## A Note on Terminology

An earlier draft of this argument used "exaptation" for this transition.
That was wrong.

Exaptation in Snowden's framework describes radical repurposing: taking
something designed for one function and finding it serves a different
function in a different context. It is an innovation mechanism operating
within the complex domain, not a mechanism for moving between domains.
The word comes from evolutionary biology, where Gould and Vrba used it
in 1982 to describe traits coopted for functions different from the ones
they evolved for.

Constraint transformation is the accurate description. The act of writing
an executable specification converts the type of constraint governing the
system's behaviour, which changes the nature of the problem.

---

## What AI Changes About the Economics

For most of software engineering history, writing precise executable
specifications cost more relative to writing the code they described.
BDD scenarios require domain knowledge, collaboration, and iteration.
Mutation testing requires tooling investment. Contract tests require
discipline about API boundaries. Many teams concluded the benefit did not
justify the cost.

Two things have shifted this calculation.

The DORA 2026 report, based on 1,110 Google engineer responses, found
that higher AI adoption correlates with higher throughput and higher
instability simultaneously. Time saved generating code is re-spent
auditing it. The bottleneck has moved from writing code to knowing what
to write and verifying that what was written is correct. Specifications
are no longer overhead relative to implementation. They are the scarce
resource.

AI also reduces the mechanical cost of expressing intent in executable
form. The evidence here is early but directional. Fonseca, Lima, and
Faria (arXiv:2510.18861, 2025) measured end-to-end Gherkin scenario
generation from JIRA tickets on a production mobile application at BMW
and found that practitioners reported time savings often amounting to a
full developer-day per feature, with the automated generation completing
in under five minutes. In a separate quasi-experiment, Hassani,
Sabetzadeh, and Amyot (arXiv:2508.20744v2, 2026) had software
engineering practitioners evaluate LLM-generated Gherkin specifications
from legal text; 91.7% of ratings on the time savings criterion fell in
the top two categories. Neither study is a controlled experiment. A
2025 systematic review of 74 LLM-for-requirements studies (Zadenoori
et al., arXiv:2509.11446) found that studies are predominantly evaluated
in controlled environments using output quality metrics, with limited
industry use, consistent with specification authoring cost being an
underexplored measurement target in the literature. The evidence is
emerging, not settled.

What both studies show consistently is that AI shifts the human effort
from authoring to reviewing: from expressing intent to validating
that the expressed intent is accurate. The intent, the domain knowledge,
and the judgment that the scenarios accurately describe what the system
should do remain the human's responsibility. The mechanical work of
expressing that intent in executable form has become substantially
cheaper.

The domain transition from complex to complicated is now economically
viable at scale in a way it was not five years ago. The claim is not
that AI makes specification automatic. It is that the economics have
shifted enough to make specification-first development the rational
default rather than the disciplined exception.

---

## Where Khononov's Observation Stops

In May 2025, Vlad Khononov argued on LinkedIn that LLMs turn complex
problems into complicated ones. The direction is right. The mechanism
is incomplete.

Khononov attributes the domain transition to LLMs broadly. The argument
here is more specific: the LLM is not what makes the transition stick.
The specification is. An LLM operating without a specification produces
output that sits in the complex domain because the agent's behaviour is
still emergent and only knowable in hindsight. An LLM operating within
a specification-governed verification pipeline produces output that can
be assessed deterministically.

Without the specification, the problem bounces back to complex every time
the agent drifts, which it will. The governing constraint is what holds
the transition.

---

## A Likely Objection

A Cynefin-literate reader will raise a fair challenge: software systems
are deterministic at the execution level. Given the same inputs and state,
they produce the same outputs. If cause and effect are always knowable in
principle, the complex framing is a category error. Software development
is always in the complicated domain.

The response is that the complexity resides in the problem space, not
the execution. What users need, which edge cases matter, how requirements
will evolve. These are genuinely complex. They are knowable only in
hindsight, they respond unpredictably to interventions, and they exhibit
the emergent properties that Snowden places in the complex domain. An
executable specification narrows the problem space by articulating
requirements precisely enough to be analysed. The domain transition
occurs in the problem space, not in the implementation.

This is not an endorsed Cynefin position. Dave Snowden is actively working
on the relationship between AI and Cynefin domains as of early 2026,
following practitioner discussions, and has not published conclusions.
The argument here is an application of Cynefin vocabulary, checked
carefully against canonical definitions. Readers familiar with the
framework are invited to challenge it.

---

## What Changes in Practice

Teams that operate AI agents without specifications are navigating a
complex domain with complicated-domain tools. They apply analysis and
best practices to a system whose behaviour is emergent. They will get
some things right through pattern matching and experience, and they will
miss systematically at the boundaries, because boundaries are exactly
where emergent behaviour is most unpredictable.

Teams operating AI agents within a specification-governed verification
pipeline have moved the problem. They are not analysing emergent
behaviour. They are verifying compliance with stated constraints. The
failure mode is specification incompleteness, not emergent
unpredictability. Incompleteness is solvable. A specification gets less
incomplete over time. Emergent unpredictability does not.

These are not the same problem with different amounts of effort applied.
They are different problems, requiring different investment, different
staffing, and different success measures.

---

*The Cynefin framework was developed by Dave Snowden and described in
Snowden, D. and Boone, M. (2007), "A Leader's Framework for Decision
Making," Harvard Business Review, November 2007. The enabling and
governing constraints terminology is from Snowden's Cynefin wiki:
cynefin.io/wiki/Constraints; this language is not in the 2007 HBR
paper, which focuses on cause-and-effect relationships by domain.
Khononov's observation on LLMs and Cynefin domains was made on LinkedIn
in May 2025. Exaptation definition: Gould, S.J. and Vrba, E.S. (1982),
"Exaptation: a missing term in the science of form," Paleobiology 8(1).
DORA 2026: "Balancing AI Tensions: Moving from AI adoption to effective
SDLC use," dora.dev, March 2026. Fonseca, P.L., Lima, B., and Faria,
J.P. (2025), "Streamlining Acceptance Test Generation for Mobile
Applications Through Large Language Models: An Industrial Case Study,"
arXiv:2510.18861. Hassani, S., Sabetzadeh, M., and Amyot, D. (2026),
"From Law to Gherkin: A Human-Centred Quasi-Experiment on the Quality
of LLM-Generated Behavioural Specifications from Food-Safety
Regulations," arXiv:2508.20744v2 (updated March 2026). Zadenoori,
M.A., Dabrowski, J., Alhoshan, W., Zhao, L., and Ferrari, A. (2025),
"Large Language Models (LLMs) for Requirements Engineering (RE): A
Systematic Literature Review," arXiv:2509.11446.*

---

*Series: The Specification as Quality Gate*
*[Part 1: The Echo Chamber in Your Pipeline](/echo-chamber-ai-code-review-correlated-error)*
*Part 2: From Complex to Complicated (this post)*
*Part 3: What Specifications Cannot Catch, coming soon*
*Part 4: The Practitioner Paper, coming soon*
