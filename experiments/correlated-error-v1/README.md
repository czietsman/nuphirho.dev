# Experiment: Correlated Error in AI Code Review (v1 — Boundary Conditions)

## Hypothesis

If an LLM's training data encodes the same misconceptions that produce a bug,
the model will fail to detect the bug during code review — its errors are
*correlated* with the programmer's errors rather than independent. This
experiment is the **baseline** using textbook boundary-condition bugs to
establish detection rates before testing with domain-opaque bugs in v2.

## Bug Corpus

| Function | File | Planted Bug |
|---|---|---|
| `paginate` | `corpus/paginate.py` | Off-by-one: `>=` should be `>` — rejects the last valid page |
| `binary_search` | `corpus/binary_search.py` | `while low < high` should be `while low <= high` — misses single-element case |
| `is_leap_year` | `corpus/is_leap_year.py` | Missing `year % 100` check — 1900 incorrectly returns `True` |
| `truncate_string` | `corpus/truncate_string.py` | `>=` should be `>` — truncates strings of exactly `max_len` |
| `days_between` | `corpus/days_between.py` | Missing `abs()` — returns negative when `date2 < date1` |

## How to Run

### Setup

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

### Run BDD tests

```bash
bash scripts/run_bdd.sh
```

Results are written to `results/bdd/`.

### Run AI code review

```bash
bash scripts/run_review.sh        # default: 5 runs per function
bash scripts/run_review.sh 10     # or specify a custom number of runs
```

Results are written to `results/review/<function_name>/run_N.json`.

### Analyse results

```bash
python3 scripts/analyse.py
```

Produces a markdown summary table comparing BDD and AI review detection rates.

## Results Summary

v1 used classic boundary-condition bugs. AI review detected all 5 at 100%
across all runs. BDD also caught all 5. This baseline motivated v2, which uses
domain-convention bugs to test whether AI review degrades when the domain rule
is absent from training data.

## Limitations

1. **Small corpus** — only 5 functions; results may not generalise.
2. **Single model** — only tested with one LLM; other models may differ.
3. **Prompt sensitivity** — detection rates may vary with different review prompts.
4. **Bug homogeneity** — all bugs are boundary-condition errors; real codebases have diverse bug types.
5. **No human baseline** — no comparison against human reviewer detection rates.

## See Also

See `../correlated-error-v2/` for the follow-up experiment with domain-opaque bugs.
