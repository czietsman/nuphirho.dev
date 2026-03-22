# Experiment: Correlated Error in AI Code Review

## Hypothesis

AI code review shares systematic blind spots with AI code generation. When an LLM generates code containing a particular class of bug, the same (or a similar) LLM reviewing that code will fail to detect the bug at a higher rate than expected. Behaviour-Driven Development (BDD) with concrete domain-specific scenarios can catch bugs that AI review misses, because executable specifications encode external conventions that the model has no built-in knowledge of.

The bugs in this corpus are domain-convention bugs: the code is internally consistent, idiomatic, and produces plausible output. The defect is only visible with knowledge of an external specification or industry convention (ISDA day count rules, marginal tax bracket mechanics, aviation maintenance logic, VC term sheet conventions, fixed income interpolation methods). This makes them harder for an LLM to detect than syntactic or boundary-condition errors.

## Methodology

The experiment has two independent conditions:

**Condition A -- AI Code Review.** Each buggy function from the corpus is submitted to an LLM (via `claude -p`) with a standardised review prompt. The prompt asks the model to identify bugs and return structured JSON. Each function is reviewed multiple times (default 5) to measure detection consistency. A review "detects" the bug if its output contains any of the detection keywords listed in `ground_truth.json`.

**Condition B -- BDD Specification.** Each buggy function is exercised by a Gherkin feature file containing scenarios that cover both normal operation and the specific domain rule that the planted bug violates. The `behave` test runner executes these scenarios against the buggy implementations. A failing scenario means BDD caught the bug.

The analysis script (`scripts/analyse.py`) compares results from both conditions against the ground truth to produce a summary table showing, for each function, whether BDD caught the bug and what fraction of AI review runs detected it.

## Bug Corpus

| Function | File | Planted Bug |
|---|---|---|
| `prorate_premium` | `corpus/prorate_premium.py` | Uses fixed 365 divisor instead of 366 for leap years per ISDA actual/actual day count convention |
| `apply_tiered_tax` | `corpus/apply_tiered_tax.py` | Assignment (`=`) instead of accumulation (`+=`) in tax loop -- only keeps the last bracket's contribution |
| `schedule_maintenance` | `corpus/schedule_maintenance.py` | Uses AND instead of OR for flight hours and cycles check -- requires both to exceed limits instead of either |
| `calculate_dilution` | `corpus/calculate_dilution.py` | Option pool calculated on post-money shares instead of pre-money -- understates founder dilution per option pool shuffle convention |
| `interpolate_rate` | `corpus/interpolate_rate.py` | Uses linear interpolation instead of log-linear as required by market convention for rate curves |

## How to Run

```bash
# Install dependencies
pip install -r requirements.txt

# Run BDD tests (Condition B)
./scripts/run_bdd.sh

# Run AI code review (Condition A) -- requires claude CLI
./scripts/run_review.sh        # default: 5 runs per function
./scripts/run_review.sh 10     # or specify a custom number

# Analyse results
python scripts/analyse.py
```

## Success Criteria

The hypothesis is supported if BDD catches all five bugs (5/5 scenarios fail) while AI review misses at least some of them across its runs. A stronger result would show AI review consistently missing specific bug categories, suggesting systematic blind spots rather than random variation.

## Limitations

The corpus is small by design -- five functions with one planted bug each. The bugs were chosen for domain obscurity rather than being sampled from naturally occurring defects. Each bug violates an external convention (financial, aviation, venture capital) that is unlikely to be well-represented in general training data. This is sufficient to demonstrate the experimental methodology but is not large enough to draw statistically robust conclusions about AI review accuracy in general.

All five bugs were hand-planted by the experiment author rather than being naturally occurring defects produced by an LLM. This means the bugs may not be representative of the errors an AI code generator would actually introduce. The correlation between generation and review blind spots can only be fully tested with bugs that originate from the same family of models being used for review.

The detection-keyword approach used to score AI review output is a rough heuristic. A review response might identify the correct bug using different terminology than the keywords in `ground_truth.json`, leading to a false negative in scoring. Conversely, a response might mention a keyword in passing without actually identifying the bug, leading to a false positive. Manual review of the raw outputs is recommended to validate the automated scores.

The BDD scenarios were written by someone who already knew where each bug was, which means the specifications are optimally targeted at the planted defects. In a real development workflow, the person writing specifications would not have advance knowledge of bugs, and their scenarios might not cover the exact boundary conditions needed to trigger them. This gives BDD an unfair advantage in this experiment.

Only a single LLM (Claude, via the CLI) is used for the review condition. Different models may have different blind spots, and a model might perform better or worse depending on the specific prompt wording, temperature, or system prompt. The results should not be generalised to all AI code review tools without further experimentation across multiple models and prompting strategies.
