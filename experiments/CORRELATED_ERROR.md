# Correlated Error in AI Code Review

## Overview

This experiment tests whether AI code review shares systematic blind spots with AI code generation. A deterministic BDD pipeline encodes domain conventions as executable specifications and catches defects that same-family AI review misses.

Two iterations were run, progressing from simple to complex:

## v1 — Boundary Conditions (`correlated-error-v1/`)

Five Python functions with classic textbook boundary-condition bugs: off-by-one errors, incorrect loop termination, missing `abs()`, and similar.

**Result:** AI review detected all 5 bugs at 100% across all runs. BDD also caught all 5. These bugs are too well-represented in training data to test the correlated error hypothesis.

## v2 — Domain-Convention Violations (`correlated-error-v2/`)

Five Python functions with domain-opaque bugs: ISDA day count conventions, marginal tax calculation, aviation maintenance logic (OR vs AND), VC option pool shuffle, and log-linear rate interpolation. Each function has neutral docstrings that do not reveal the domain rule — the bug is only visible with external domain knowledge.

**Result:** BDD caught all 5 bugs. AI review detection varied by domain opacity:

| Function | BDD | AI review | Domain |
|---|---|---|---|
| `prorate_premium` | caught | 5/5 (100%) | Insurance day count |
| `apply_tiered_tax` | caught | 5/5 (100%) | Tax brackets |
| `schedule_maintenance` | caught | 5/5 (100%) | Aviation maintenance |
| `calculate_dilution` | caught | 4/5 (80%) | VC term sheets |
| `interpolate_rate` | caught | 0/5 (0%) | Fixed income |

The gradient in detection rate correlates with domain obscurity in training data. `interpolate_rate` (log-linear convention) was completely missed. `calculate_dilution` (option pool shuffle) was partially detected. The others contained code-level signals detectable without domain knowledge.

## Key finding

AI review degrades where domain conventions are absent from training data. BDD is invariant to this — executable specifications encode the convention directly, regardless of whether the reviewer has seen it before.

## Experimental design note

v2 initially used docstrings that stated the domain convention explicitly (e.g., "ISDA actual/actual", "log-linear interpolation"). With those docstrings, AI review caught all 5 at 100% — it was comparing code to its own docstring, not applying domain knowledge. The docstrings were stripped to neutral descriptions for the final run. This intermediate result is itself informative: the specification-as-docstring is equivalent to having a BDD spec, and its removal is what exposes the correlated error.
