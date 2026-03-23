# Experiment: Correlated Error in AI Code Review (v3 -- Cross-Family Panel)

## Hypothesis

The correlated error hypothesis applies where domain conventions are absent from training data. An AI reviewer without an external specification cannot identify a bug that is only wrong relative to a published domain convention -- not because the model lacks reasoning capability, but because the information required to catch the bug does not exist in the code or the docstring.

V1 tested classic boundary conditions: all models caught all bugs at 100%. V2 tested domain-convention violations with a same-family reviewer (Claude reviewing Claude-generated code): interpolate_rate was missed at 0%, calculate_dilution at 80%.

V3 extends the test to a cross-family reviewer panel: four distinct models review the same corpus independently, without any external specification. The hypothesis predicts detection rates will vary by domain opacity and by model training distribution.

## Models Under Test

| Model | CLI tool | Version | Family |
|---|---|---|---|
| Claude | Claude Code | 2.1.81 | Anthropic |
| Codex | codex-cli | 0.116.0 (gpt-5.4) | OpenAI |
| Gemini | gemini | 0.34.0 | Google |
| Amazon Q | q | 1.18.1 (claude-sonnet-4.5) | AWS / Anthropic |

Note: Amazon Q identifies as claude-sonnet-4.5 internally. This makes Claude and Q same-family (both Anthropic), which is worth noting in results interpretation. Codex and Gemini are genuinely cross-family.

## Methodology

The experiment has two independent conditions:

**Condition A -- AI review without specification.** Each buggy function from the corpus is submitted to each of the four models via their CLI tools with a standardised review prompt. No domain context, no specification, no hints. Five runs per function per model (100 total Condition A runs). Detection is recorded as pass if the reviewer identifies the planted bug.

**Condition B -- BDD pipeline.** Pre-written BDD scenarios targeting the exact domain convention run via behave. One run per function, deterministic.

The analysis script (`scripts/analyse.py`) compares results from both conditions against the ground truth to produce a summary table showing, for each function, whether BDD caught the bug and what fraction of AI review runs per model detected it.

## Bug Corpus

| Function | Domain | Convention | Planted Bug |
|---|---|---|---|
| `calculate_final_reserve_fuel` | Aviation | ICAO Annex 6 Part I Section 4.3.6.3 | Reserve durations swapped: jet assigned 45 min, piston 30 min (should be jet 30, piston 45) |
| `get_gas_day` | Energy trading | NAESB WGQ 1.3.2 / FERC Order 698 | Gas day starts at midnight instead of 09:00 Central Time |
| `validate_diagnosis_sequence` | Healthcare coding | ICD-10-CM Guidelines Section I.C.20.a | Allows V00-Y99 external cause codes as principal diagnosis |
| `validate_imo_number` | Maritime | IMO Resolution A.1079(28) | Check digit uses ascending weights 2-7 instead of descending 7-2 |
| `electricity_cost` | Utility billing | Standard inclining block tariff | Flat rate per tier instead of progressive application across tiers |

Corpus was designed collaboratively across six models (Grok 4.20, Gemini, DeepSeek-V3, Qwen3.5, Claude Opus 4.6, ChatGPT GPT-5.3). Amazon Q Developer declined to participate in corpus generation on safety grounds. Full candidate sets are in the experiments parent directory.

## How to Run

```bash
# Install dependencies
pip install -r requirements.txt

# Run BDD tests (Condition B)
./scripts/run_bdd.sh

# Run AI code review (Condition A) -- requires the relevant CLI tool
./scripts/run_review.sh claude corpus/calculate_final_reserve_fuel.py
./scripts/run_review.sh codex corpus/calculate_final_reserve_fuel.py
./scripts/run_review.sh gemini corpus/calculate_final_reserve_fuel.py
./scripts/run_review.sh q corpus/calculate_final_reserve_fuel.py

# Analyse results
python scripts/analyse.py
```

## Success Criteria

The experiment confirms the hypothesis if detection rates vary systematically by domain opacity across models, and if the BDD pipeline catches all five bugs while at least two models miss at least two bugs without the specification.

Either result is worth publishing.

## Limitations

1. Same-family overlap: Claude (Anthropic) and Q (claude-sonnet-4.5, Anthropic) share
   a training distribution. This partially replicates the same-family condition from v2.
   Codex (OpenAI) and Gemini (Google) are genuinely cross-family.
2. Planted bug corpus. Bugs were designed to test the hypothesis, not sampled from a
   natural defect distribution.
3. Bugs designed by models that are not in the reviewer panel (Claude Opus 4.6,
   Qwen3.5, ChatGPT GPT-5.3). The reviewer models (Claude Code 2.1.81, Codex,
   Gemini, Q) did not design the corpus, providing a degree of separation.
4. Keyword matching for detection is imprecise. Manual review of outputs required to
   confirm genuine detections.
5. Five functions is a small corpus. Results are directional, not statistically significant.
6. Amazon Q was excluded from corpus generation on safety grounds and used as a
   reviewer only. This is noted as a finding in the experiment write-up.
