# governance-prompts-v1

Structural quality evaluation of 34 AGENTS.md governance prompts against
five PromptQ principles. Part of the RE@Next! 2026 paper submission.

- **Version:** v1
- **Captured:** 2026-04-05
- **Tool:** capture-governance-prompts-corpus
- **Corpus:** 34 repositories (see manifest.md)

## Experiment

Each AGENTS.md file in agents/ is scored by multiple LLMs against five
principles (P1--P5). Cross-model divergence is recorded as a finding.

The five principles are defined in ground_truth.json:

- **P1: Success Definition** -- completion criteria
- **P2: Assessment Rubric** -- self-evaluation criteria
- **P3: Scope Boundary** -- permitted/prohibited actions
- **P4: Data Classification** -- content sensitivity handling
- **P5: Quality Gate** -- verification before acceptance

## Scoring

```bash
go run ./cmd/score-governance-prompts --model claude
go run ./cmd/score-governance-prompts --model all
```

Run from the repository root. The tool reads from
`experiments/governance-prompts-v1/` by default.

The prompt template is in `prompts/` (gitignored).

## How to Run

### Prerequisites

```bash
pip install -r requirements.txt
```

### Run BDD specs

```bash
cd experiments/governance-prompts-v1
python3 -m behave
```

BDD specs verify output structure and completeness only. They do not verify
that scores are correct (there is no oracle).

### Analysis

```bash
python3 scripts/analyse.py
```

Produces four files in results/summary/:
- scores_by_principle.md
- scores_by_repo.md
- divergence.md (requires 2+ models)
- paper_tables.md (formatted for Section 4.3)

## Data Policy

Raw scores and the scoring prompt are gitignored.
Only summary outputs (results/summary/) are committed to the repository.

## Results archiving

Raw scores are gitignored from this repository.
If `EXPERIMENT_ARCHIVE_DIR` is set, the scoring tool copies results to
the archive directory. Commit there separately.
