---
title: "What Specifications Cannot Catch: A Proposed Taxonomy of the Residual"
slug: "what-specifications-cannot-catch-residual-defect-taxonomy"
draft: false
tags: ["specifications", "testing", "bdd", "software-engineering", "ai", "formal-methods"]
series: "The Specification as Quality Gate"
---

*This is Part 3 of a four-part series, "The Specification as Quality Gate." [Part 1](/echo-chamber-ai-code-review-correlated-error) developed the correlated error hypothesis. [Part 2](/executable-specifications-cynefin-domain-transition) grounded the argument in complexity science. This post maps what executable specifications cannot catch. Part 4 will follow.*

The argument for specification-driven development is strong. BDD scenarios
that pass or fail are not probability estimates. Mutation testing that
finds no survivors is a verified claim about the coverage of the
verification layer itself. Contract probes that confirm boundary behaviour
are deterministic gates, not heuristics. The case for putting
specifications before code, and treating the pipeline as the reviewer,
is well-founded.

"Well-founded" is not the same as "complete." Any argument that executable
specifications make AI code review redundant needs to account honestly for
what specifications cannot catch. That accounting is not a concession. It
is what makes the rest of the argument credible.

What follows is a proposed taxonomy of the defect classes that lie outside
the reach of executable specifications. To the best of our knowledge, no
prior taxonomy in the testing or formal methods literature organises
defects by specifiability rather than by severity or recovery type. The
reason to organise by specifiability is practical: severity tells you how
bad a defect is. Specifiability tells you which tool to reach for. Those
are different questions, and conflating them is why teams end up applying
AI review to problems where it is circular and missing the problems where
it would genuinely help. It should be treated as a working framework and
challenged accordingly.

---

## The Foundation: The Oracle Problem

Before the categories, a theoretical grounding that gives the taxonomy
its shape.

The oracle problem in software testing asks: how do you know whether a
test outcome is correct? For a test to be meaningful, you need an oracle,
a ground truth against which to judge the system's output. In most
testing, the oracle is implicit: the developer's understanding of what
the system should do. In specification-driven development, the oracle is
the specification.

Barr, Harman, McMinn, Shahbaz, and Yoo, in their 2015 survey in IEEE
Transactions on Software Engineering, established that complete oracles
are theoretically impossible for most real-world software. Even a formally
correct specification, internally consistent and precisely expressed,
cannot fully specify the correct output for all possible inputs in all
possible contexts.

This is not a practical limitation that better process will eventually
overcome. It is a theoretical result. There exists a class of defects for
which no specification, however precise, provides an oracle. That class
defines the permanent boundary of what specification-driven verification
can achieve.

Everything in the taxonomy that follows should be read against this
foundation.

---

## Category A: Theoretically Specifiable, Not Yet Specified

The largest and most important category. These are defects that a
specification could have caught if the scenario had been written. The
gap is a process failure, not a theoretical limitation.

The clearest examples are domain-convention violations: code that is
internally consistent and idiomatic, but wrong relative to a rule that
is not inferrable from the code alone. An interest rate interpolation
function that uses linear interpolation when the market convention for
this curve type is log-linear looks correct to any reviewer without
the specification. The code is clean. The logic is sound. Only the BDD
scenario that states the correct output for a specific tenor gap makes
the violation visible.

This is not a hypothetical. A small experiment run alongside the paper
this post is part of tested exactly this case. Claude reviewed the
function five times without a specification and missed the bug in every
run, flagging a different concern instead. The BDD scenario caught it
immediately. The full experiment is at
[correlated-error-v2](https://github.com/czietsman/nuphirho.dev/tree/main/experiments/correlated-error-v2).

Classic boundary conditions (off-by-one errors, loop termination
edge cases, leap year rules) also belong here in principle. But
frontier models catch those reliably from pattern recognition alone,
because they are dense in training data. The specification still matters
for these: it prevents regressions when a refactoring agent removes a
guard clause it considers noise. But the correlated error problem bites
hardest at the domain-specific end of Category A, where neither the
generator nor the reviewer has the convention in its prior.

Other examples: error handling paths that were not enumerated in
requirements, input validation for edge cases the specification author
did not consider, state machine transitions for states that were not
mapped.

Category A is the primary target for the specification-first argument.
AI-assisted specification drafting, property-based exploration, and
mutation testing can reduce it systematically over time. The residual
in this category is not permanent. It shrinks as specification discipline
matures.

An AI review agent operating in Category A without an external
specification shares the same blind spots as the generator: both
drew from the same prior, and whatever the prior underweights, both
miss.

---

## Category B: Specifiable in Principle, Economically Impractical

Some defect classes are theoretically specifiable but require verification
effort that exceeds the value of specifying them. A function that takes
three independent boolean parameters has eight input combinations. Ten
parameters gives 1,024. Twenty gives over a million. Every combination's
boundary conditions could be specified, but nobody does this, for good
reason.

This boundary is moving. Property-based testing frameworks generate inputs
systematically and explore combinatorial spaces that manual scenario
authoring cannot reach. The economics of Category B defects are shifting
as these tools mature and AI-assisted property-based exploration becomes
practical.

Category B is where AI review can provide legitimate signal, not by
finding what the specification missed, but by sampling the input space
more broadly than the specified scenarios cover. This is heuristic, not
deterministic, but it is genuine signal rather than circular reasoning,
provided the reviewer draws from a genuinely different prior than the
generator.

---

## Category C: Inherently Unspecifiable From Pre-Execution Context

Some defect classes depend on properties of the running system, the
environment, or the interaction between components that only manifest at
runtime. Timing-dependent race conditions, behaviour under partial network
failure, memory behaviour under sustained load, and performance degradation
under specific hardware configurations cannot be expressed as a BDD
scenario that runs before deployment. The conditions that trigger them
do not exist at specification time.

These are not unknown unknowns in the Category A sense. They are knowable
properties, observable and measurable after the fact. The issue is
temporal, not conceptual.

The boundary of Category C is moving. Santos, Pimentel, Rocha, and
Soares (Software, 2024) explicitly encoded ISO 25010 performance
efficiency requirements as BDD scenarios. Maaz et al. (arXiv:2510.09907,
2025) used LLM-assisted property-based testing to surface concurrency and
numerical precision bugs that previously required runtime observation to
detect. Some defects classified as inherently unspecifiable five years
ago are specifiable today with the right tooling.

For the defects that remain in Category C, the right tool is runtime
verification infrastructure. This includes ML-based anomaly detection,
APM tooling with learned behavioural baselines, observability platforms,
chaos engineering, and load testing. This is a mature field predating
LLMs that applies machine learning to operational data rather than to
source code. Neither a BDD pipeline nor an AI code review agent
reaches these defects. They require the system to be running.

---

## Category D: Structural and Architectural Properties

Code can pass every BDD scenario, survive every mutation, and satisfy
every contract probe, while simultaneously introducing coupling between
modules that will compound over time, violating layer boundaries the
architecture was designed to enforce, or drifting from the intended design
in ways that only become visible months later.

These are relational properties of the codebase as a whole, not properties
of individual components in isolation. They resist behavioural specification
because they concern how the code is structured relative to the intended
design, not what the code does.

But resist does not mean unspecifiable. Architectural rules are specifiable
once a human articulates them. A rule such as "the web layer may not import
from the data layer" is a constraint that can be expressed, enforced, and
tested deterministically. Tools like ArchUnit, Dependency Cruiser, and
NDepend enforce dependency rules as automated checks. Contract testing
frameworks like Pact verify that service boundaries behave as agreed
between producers and consumers. These are deterministic gates, not
heuristic opinions, and they belong in the pipeline alongside BDD
scenarios. General delivery rigour (clear module boundaries, explicit
interface definitions, disciplined use of dependency injection) reduces
the surface area of Category D by making architectural intent concrete
rather than implicit.

The residual in Category D, after tooling and contract testing are in
place, is the unarticulated architectural intent that has not yet been
expressed as a rule. Drift from a design decision that was never written
down. Coupling that violates a pattern that exists only in someone's head.
A half-completed migration to a new architectural pattern, where some
modules use the new approach and others still use the old one: no
automated rule catches that because the correct answer depends on whether
the team intended to complete the migration or abandoned it. Dead patterns
accumulate the same way: abstractions introduced for a use case that no
longer exists, base classes with a single subclass, interfaces designed
for extensibility that never came. These are structural noise that only
becomes visible to something reasoning about the whole codebase over time.

This is where an AI review agent with access to the full codebase and
architectural context adds genuine, non-circular signal. It can identify
incomplete migrations, flag dead abstractions, and surface inconsistencies
that no automated rule yet captures. The role is analogous to an expert
architect reviewing for structural coherence: the agent advises, the human
decides whether to complete the migration, remove the dead pattern, or
codify the observation as a new enforceable rule. That last step closes
the loop back into tooling, which is where the category belongs once the
intent has been made explicit.

This is the least empirically grounded category in the taxonomy: no
controlled study has yet isolated AI architectural review as a distinct
defect class separate from what specification pipelines and architectural
tooling catch. It is the most provisional claim here, and should be
treated accordingly.

Category D is the permanent legitimate home for AI review in a
specification-driven pipeline, not as a substitute for architectural
tooling and contract testing, but as the complement that operates where
rules have not yet been written.

---

## Category E: Specification Defects

The oracle problem survey makes this category unavoidable. A specification
that is internally consistent, precisely expressed, and correctly
implemented can still describe a system that does not do what users need.

Requirements elicitation is not a solved problem. Domain experts disagree.
Users articulate needs in terms of their current mental model, which may
not reflect what they actually want. Business rules change after the
specification was written. Edge cases that matter in production were not
considered during discovery.

No verification pipeline catches Category E defects, because the pipeline
verifies conformance to the specification. If the specification is wrong,
the pipeline confirms the wrong thing. An AI review agent without external
user feedback cannot catch Category E either, because the agent has no way
to know the specification does not match user intent.

The right tool for Category E is not in the engineering pipeline at all.
It is user testing, feedback loops, observability of actual usage, and
the design thinking practices that surface unstated assumptions during
requirements elicitation. These are human processes. AI can assist with
them, but not replace them.

Category E is not a process failure. It is a theoretical limitation.
Even a perfect specification process will leave some Category E defects,
because the oracle for user intent is inherently incomplete.

---

## What the Taxonomy Implies

Each category points to a different tool.

Category A is the target for specification discipline. It shrinks as
teams invest in BDD coverage, mutation testing, and AI-assisted
specification drafting. An AI review agent here is a temporary patch for
process immaturity, not a permanent fixture.

Category B is the target for property-based testing and combinatorial
exploration. AI adds value here only if it brings genuine sampling
diversity, meaning a different prior than the generator.

Category C is the target for runtime verification infrastructure:
ML-based anomaly detection, observability tooling, chaos engineering,
and load testing. This is a mature field with its own tooling. Neither
specifications nor AI code reviewers substitute for it.

Category D is the target for architectural tooling and contract testing
first: ArchUnit, Dependency Cruiser, Pact, and similar deterministic
enforcement mechanisms for rules that have been articulated. AI review
is the complement for the unarticulated residual: drift from design
intent that has not yet been codified as an enforceable rule. The agent
advises, the human decides, and the observation ideally becomes a new
rule that closes back into tooling.

Category E is the reminder that the specification is never complete. It
gets less incomplete over time. The feedback loop from production back
to requirements is a human loop, and cannot be automated away.

---

## An Invitation

This taxonomy is proposed, not established. The categories emerged from
reasoning through what executable specifications can and cannot express,
grounded in the oracle problem literature and tested against empirical
findings from the 2025-2026 AI code review research.

If the taxonomy is wrong, the corrections are welcome and will improve the
argument. If it is incomplete, the missing categories are interesting. The
most likely challenge is that the Category A and Category D boundary is
blurrier than described: architectural properties may become specifiable
as tooling matures, which would reclassify parts of Category D as
Category A over time. That boundary question is worth watching.

The goal is not to defend the taxonomy. It is to have a precise enough
framework to stop treating AI code review as a monolithic practice and
start reasoning about which kinds of defects each tool in the pipeline
is actually positioned to find.

If you work in software testing, formal methods, or AI-assisted
development and the taxonomy is wrong in ways not anticipated here,
the comments are the right place to say so. A taxonomy that provokes
rigorous challenge is more useful than one that goes unchallenged.
Corrections will be published with attribution.

---

*This post draws on: Barr, E.T., Harman, M., McMinn, P., Shahbaz, M.,
and Yoo, S. (2015), "The Oracle Problem in Software Testing: A Survey,"
IEEE Transactions on Software Engineering, 41(5), 507-525,
doi:10.1109/TSE.2014.2372785; Santos, S., Pimentel, T., Rocha, F., and
Soares, M.S. (2024), "Using Behaviour-Driven Development (BDD) for
Non-Functional Requirements," Software, 3(3), 271-283,
doi:10.3390/software3030014; and Maaz, M. et al. (2025), "Agentic
Property-Based Testing: Finding Bugs Across the Python Ecosystem,"
arXiv:2510.09907.*

---

*Series: The Specification as Quality Gate*
*[Part 1: The Echo Chamber in Your Pipeline](/echo-chamber-ai-code-review-correlated-error)*
*[Part 2: From Complex to Complicated](/executable-specifications-cynefin-domain-transition)*
*Part 3: What Specifications Cannot Catch (this post)*
*Part 4: The Practitioner Paper, coming soon*
