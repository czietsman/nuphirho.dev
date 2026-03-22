---
title: "The Echo Chamber in Your Pipeline"
slug: "echo-chamber-ai-code-review-correlated-error"
draft: false
tags: ["ai", "code-review", "specifications", "bdd", "software-engineering"]
series: "The Specification as Quality Gate"
---

*This is Part 1 of a four-part series, "The Specification as Quality Gate." The series develops three hypotheses about executable specifications, AI code review, and what each is actually for. Parts 2, 3, and 4 will follow.*

When an AI coding agent generates code and a separate AI reviewer examines
it, both agents are reasoning from the same artefact: the code itself. If
neither has an external reference (a specification, a contract, a
statement of what the system is supposed to do), the reviewer has no
ground truth to compare against. It checks the code against the code. Not
against intent.

The architecture is circular. And there is now enough empirical evidence
to say precisely where and how it fails.

---

## The Mechanism

In classical ensemble learning, stacking multiple estimators improves
reliability under one condition: the estimators must fail independently.
If two classifiers share the same training distribution and the same blind
spots, combining them does not reduce error. It consolidates it. The joint
miss rate of two correlated estimators approaches the miss rate of either
one alone, not the product of both.

This holds for LLM pipelines for a specific reason: a code-generating
model and a code-reviewing model drawn from the same model family share
architecture, training corpus, and reward signal. That is the same class
of correlation the independence condition prohibits. They are not two
independent estimators. They are two samples from the same prior.

A 2025 paper (Vallecillos-Ruiz, Hort, and Moonen, arXiv:2510.21513)
studying LLM ensembles for code generation and repair named this the
"popularity trap." Models trained on similar distributions converge on
the same syntactically plausible but semantically wrong answers. Consensus
selection (the default review heuristic in most multi-agent pipelines)
filters out the minority correct solutions and amplifies the shared
error. Diversity-based selection, by contrast, recovered up to 95% of
the gain that a perfectly independent ensemble would achieve. The
implication is direct: naive stacking of homogeneous LLMs is
counterproductive.

A second paper (Mi et al., arXiv:2412.11014) examined a
researcher-then-reviser pipeline for Verilog code generation and found
that erroneous information from the first agent was accepted downstream
because the agents shared the same training distribution and lacked
adversarial diversity. Self-review repeated the original errors rather
than correcting them.

A February 2026 paper (Pappu et al., arXiv:2602.01011) went further,
showing that even heterogeneous multi-agent teams consistently fail to
match their best individual member, incurring performance losses of up to
37.6%, even when explicitly told which member is the expert. The failure
mechanism is consensus-seeking over expertise. Homogeneous copies make
it worse.

None of these papers tested a pure generator-then-reviewer pipeline
without external grounding in the exact configuration described here.
That controlled experiment has not been published. But three independent
papers from 2025 and 2026, approaching the problem from different angles,
document the same failure mode: shared training distributions produce
correlated errors, and consensus amplifies rather than corrects them.

---

## A Concrete Example

This example uses a deliberately simple case. Modern frontier models catch
classic boundary conditions reliably, as the experiments described in the
caveats section confirm this. The point here is the sequence: how a BDD
scenario makes a defect detectable before any reviewer is involved.

Consider a pagination function. A developer asks an AI coding agent to
implement it. The agent produces:

```python
def paginate(items, page, page_size):
    start = page * page_size
    end = start + page_size
    return items[start:end]
```

The agent also produces tests: one for page 0 and one for page 1 against
a ten-item list. Both pass. The implementation is correct for those cases.

The flaw is in what was not tested. When the total number of items is
exactly divisible by the page size and the caller requests the last page
by calculating it from the total count, the edge case is invisible to
the tests provided.

The review agent, given only the code and the tests, validates the
implementation against what is there. The tests pass. The reviewer
reports no issues. It has no basis to ask what was not tested, because
there is no specification defining the expected behaviour for boundary
conditions.

The BDD scenario that would have caught it:

```gherkin
Scenario: Last page when total is exactly divisible by page size
  Given a list of 10 items
  And a page size of 5
  When I request page 1
  Then I receive items 6 through 10
  And the result contains exactly 5 items
```

This scenario would have failed against the implementation before the
reviewer was ever involved. The pipeline stops at the specification, not
the review.

There is a second reason this matters beyond the review question. Even
where models reliably identify boundary conditions during focused review,
general refactoring passes introduce a different risk. An agent
refactoring for readability or performance is not focused on correctness.
Conditional checks and edge-case handling can be silently removed. The
BDD pipeline catches that regression the same way it catches the original
defect. The protection is unconditional on what the agent was trying to do.

---

## What the SGCR Paper Actually Says

A December 2025 paper from HiThink Research (arXiv:2512.17540) proposed
a Specification-Grounded Code Review framework and reported a 90.9%
improvement over a baseline LLM reviewer. This figure circulates widely
and requires a correction.

The 90.9% is an improvement in developer adoption rate of review
suggestions: 42% versus 22%. Developer adoption reflects whether
suggestions are relevant and actionable. It is not a defect detection
rate. These are different claims.

The SGCR paper validates the hypothesis partially. Without specification
grounding, baseline LLM reviewers produce generic suggestions and
hallucinated issues. Grounding in human-authored specifications filters
the noise and produces suggestions developers act on. But the paper's
architecture, combining an explicit specification-driven path and an
implicit heuristic discovery path, does not claim specifications
make review redundant. It positions them as making review better.

SGCR's implicit pathway explicitly handles discovery of issues beyond the
stated rules. This is the paper's own acknowledgement that specifications
alone are insufficient. It is consistent with the argument that an AI
review agent has a legitimate residual role, but only after the
specification pipeline has done its work.

---

## The Independence Condition

The correlated error claim has a precise boundary. Stacking any two AI
reviewers is not always counterproductive. The failure condition is
specific: estimators that share a training distribution and lack an
external reference exhibit correlated failures. The condition for genuine
benefit is diversity plus external grounding.

If the reviewer uses a different model family, a different temperature, or
a substantially different prompting strategy, errors may be partially
independent. A cross-family pipeline, Grok reviewing Claude-generated
code for instance, has more independence than a same-family pipeline.
But model diversity and external grounding are two separate conditions.
Diversity reduces correlation between estimators. It does not supply
ground truth. A cross-family reviewer without an external specification
is still checking code against code, not code against intent. It will
catch some things the generator missed. It will share the same blind
spots on anything not well-represented in either training corpus, and
it has no basis to identify what was not specified regardless of how
different its architecture is from the generator.

The ensemble literature is clear that diversity is what makes stacking
work. The specification is what makes review non-circular. Both
conditions matter. The current industry architecture typically satisfies
neither.

---

## But What About Documentation?

A reasonable objection: if the problem is that the reviewer lacks ground
truth, why not solve it by adding documentation? Write better docstrings.
Maintain a spec document. Give the reviewer the context it needs.

Our own experiment accidentally tested this. The first version of the v2
corpus included docstrings that stated the domain convention explicitly.
The `prorate_premium` docstring said "ISDA actual/actual." The
`schedule_maintenance` docstring said "ANY of these thresholds." Claude
caught every bug at 100%. Documentation in close proximity to the code
works, at least when the agent reads it.

That result is real, but it does not solve the problem. It reframes it.

The first issue is retrieval and compliance. A docstring sitting on the
same function is the best case for specification proximity. An external
policy document, a requirements file, a Confluence page is a much weaker
guarantee. The agent needs to find it, read it, and apply it faithfully.
Anyone who has watched an AI agent work at scale knows that is not a safe
assumption. Agents read context selectively. They produce plausible output
that satisfies surface checks without necessarily following the rule they
were given. I have seen this enough times firsthand to know I cannot
treat documentation-as-specification as reliable ground truth.

The second issue is drift. Docstrings are not executable. They do not
fail when the code diverges from them. A function can be refactored, a
business rule can change, an edge case can be added, and the docstring
sits there describing what the function used to do, confidently, in
present tense. Every codebase accumulates this. The reviewer checking
code against a stale docstring is not checking against intent. It is
checking against what someone intended when they first wrote it.

A BDD scenario cannot drift silently. When the system stops behaving the
way the scenario describes, the scenario fails. The build stops. The
divergence is visible immediately. That is not a property of documentation.
It is a property of executable specifications. The difference is not one
of quality or discipline. It is structural.

The objection "just write better documentation" asks for the same
discipline that produced the bad documentation in the first place, and
adds an assumption that AI agents will follow it reliably. The BDD
pipeline removes both dependencies.

---

## The Implication

An AI reviewer without an external specification is a probability estimate,
not a quality gate. It samples from the same distribution as the generator
and applies pattern matching against its training data. It will catch some
things the generator missed, particularly common vulnerability patterns
and well-documented anti-patterns that appear frequently in training data.
It will systematically miss whatever the generator systematically missed,
because they share the same prior.

A BDD scenario that fails is a falsified claim. The pipeline either passes
or it does not. There is no hallucination, no false confidence, no noise
to filter.

The architecture that follows is not "no AI review." It is: specifications
first, verification pipeline second, AI review only for the residual. The
residual is real. It includes architectural properties, structural drift,
and defect classes that resist specification, but it is a much smaller
target than the current industry framing suggests. The next post in this
series maps that residual precisely.

---

## Caveats and Open Questions

This argument rests on three empirical papers and two small contrived
experiments. Both experiments are same-family only (Claude reviewing
Claude-generated code) and use a planted bug corpus rather than a natural
defect sample. They are directional evidence, not a controlled
demonstration.

The first experiment used classic boundary-condition bugs. Claude caught
all five at 100%, which refined rather than confirmed the hypothesis.
Classic boundary conditions are pattern-recognition problems that are
dense in training data. The correlated error claim applies where the
convention is absent from training data, not where it is ubiquitous.

The second experiment used domain-convention violations with neutral
docstrings. Detection ranged from 0% to 100% depending on domain opacity.
The log-linear interpolation function, where the convention is market
practice in fixed income rather than general programming knowledge, was
missed in all five runs. The reviewer flagged a different concern instead
and declared the implementation correct. BDD caught it. The full corpus,
specifications, scripts, and results are available in the
[experiments directory](https://github.com/czietsman/nuphirho.dev/tree/main/experiments).

The independence condition also requires more empirical work. Same-family
models fail the independence requirement. How much architectural diversity (different model families, different prompting strategies, different
temperatures) is required to provide genuine signal rather than
correlated noise remains to be established.

These are open questions, not fatal weaknesses. Without an external
reference, circular review is circular. The empirical work so far points
in one direction. A controlled demonstration using a natural defect sample
and cross-family pipelines would make the argument more precise.

The authors of the papers cited here, along with anyone working on LLM
ensemble reliability, multi-agent code pipelines, or specification-grounded
review, are welcome to challenge, correct, or build on this argument.
If the hypothesis is wrong, knowing precisely where it breaks is more
useful than leaving it unchallenged. Comments are open.

---

*The papers cited in this post are: Vallecillos-Ruiz, F., Hort, M., and
Moonen, L., "Wisdom and Delusion of LLM Ensembles for Code Generation and
Repair," arXiv:2510.21513, October 2025; Mi, Z. et al., "CoopetitiveV:
Leveraging LLM-powered Coopetitive Multi-Agent Prompting for High-quality
Verilog Generation," arXiv:2412.11014, v2 June 2025; Pappu, A. et al.,
"Multi-Agent Teams Hold Experts Back," arXiv:2602.01011, February 2026;
and Wang, K., Mao, B. et al., "SGCR: A Specification-Grounded Framework
for Trustworthy LLM Code Review," arXiv:2512.17540, December 2025, v2
January 2026. Author lists and findings verified against the original
papers.*

---

*Series: The Specification as Quality Gate*
*Part 1: The Echo Chamber in Your Pipeline (this post)*
*[Part 2: From Complex to Complicated](/executable-specifications-cynefin-domain-transition)*
*Part 3: What Specifications Cannot Catch, coming soon*
*Part 4: The Practitioner Paper, coming soon*
