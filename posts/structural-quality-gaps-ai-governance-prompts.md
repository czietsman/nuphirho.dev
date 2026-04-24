---
title: "Structural Quality Gaps in AI Governance Prompts"
slug: "structural-quality-gaps-ai-governance-prompts"
draft: false
tags: [ai-governance, promptq, research]
---
Most AI agent governance documents in production are structurally incomplete. That is not an observation. It is an empirical finding.

A second paper from the nuphirho.dev research programme is now on arXiv. "Structural Quality Gaps in Practitioner AI Governance Prompts: An Empirical Study Using a Five-Principle Evaluation Framework" (arXiv:2604.21090) applies a five-principle evaluation framework to 34 public AGENTS.md files sourced from GitHub. Three independent language model evaluators scored each document. Thirty-seven percent scored below the structural completeness threshold.

The five principles the framework evaluates are drawn from computability theory, proof theory, and Bayesian epistemology. Together they ask whether a governance document defines a decidable problem: does it specify success criteria, embed an assessment schema, require external verification of factual claims, scope what the agent will not do, and constrain the output format? A document that fails these tests will produce coherent output. Whether that output is correct depends entirely on what the document did not say, and there is no mechanism inside the document to catch what was left out.

The central finding is an artefact classification gap. The same file format, AGENTS.md, is being used for three architecturally incompatible purposes: task orchestration, behavioural governance, and architectural specification. There is no consensus. Different teams are solving different problems with the same artefact and calling it the same thing. That is a tractable requirements engineering problem with no current solution.

This paper is the empirical foundation for the PromptQ framework, which evaluates governance documents at authorship time. The framework is in development at promptq.ai.

arXiv:2604.21090 | [doi.org/10.48550/arXiv.2604.21090](https://doi.org/10.48550/arXiv.2604.21090)
