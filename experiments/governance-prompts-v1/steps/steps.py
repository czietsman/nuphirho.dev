import json
import os
import subprocess
import sys

PROJECT_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0, PROJECT_DIR)

from behave import given, when, then

MODELS = ["claude", "codex", "gemini"]
PRINCIPLES = ["P1", "P2", "P3", "P4", "P5"]
VALID_SCORES = [0, 0.5, 1]
SCORING_VERSION = os.environ.get("SCORING_VERSION", "v1")


def scores_dir(model):
    return os.path.join(PROJECT_DIR, "results", f"scores-{SCORING_VERSION}", model)


def agents_dir():
    return os.path.join(PROJECT_DIR, "agents")


def summary_dir():
    return os.path.join(PROJECT_DIR, "results", "summary")


def find_any_score_file():
    """Return the path and parsed content of the first score file found."""
    for model in MODELS:
        d = scores_dir(model)
        if not os.path.isdir(d):
            continue
        for f in sorted(os.listdir(d)):
            if f.endswith(".json"):
                path = os.path.join(d, f)
                with open(path) as fh:
                    return path, json.load(fh)
    return None, None


def load_all_score_files(model):
    """Return a list of (filename, parsed_json) for a model."""
    d = scores_dir(model)
    if not os.path.isdir(d):
        return []
    results = []
    for f in sorted(os.listdir(d)):
        if f.endswith(".json"):
            with open(os.path.join(d, f)) as fh:
                results.append((f, json.load(fh)))
    return results


def models_with_results():
    """Return list of models that have at least one score file."""
    found = []
    for model in MODELS:
        d = scores_dir(model)
        if os.path.isdir(d) and any(f.endswith(".json") for f in os.listdir(d)):
            found.append(model)
    return found


# ---------- scoring_output.feature ----------

@given("a score file exists for any repository and model")
def step_score_file_exists(context):
    path, data = find_any_score_file()
    assert path is not None, "No score files found in results/scores/"
    context.score_path = path
    context.score_data = data


@then('the file contains a "scores" object')
def step_has_scores(context):
    assert "scores" in context.score_data, "Missing 'scores' key"
    assert isinstance(context.score_data["scores"], dict), "'scores' is not an object"


@then('the file contains a "total" field')
def step_has_total(context):
    assert "total" in context.score_data, "Missing 'total' key"


@then('the file contains a "justifications" object')
def step_has_justifications(context):
    assert "justifications" in context.score_data, "Missing 'justifications' key"
    assert isinstance(context.score_data["justifications"], dict), "'justifications' is not an object"


@then("every principle score is 0, 0.5, or 1")
def step_valid_scores(context):
    scores = context.score_data["scores"]
    for p in PRINCIPLES:
        assert p in scores, f"Missing principle {p}"
        assert scores[p] in VALID_SCORES, f"{p} score {scores[p]} not in {VALID_SCORES}"


@then("the total equals the sum of the five principle scores")
def step_total_sum(context):
    scores = context.score_data["scores"]
    expected = sum(scores[p] for p in PRINCIPLES)
    actual = context.score_data["total"]
    assert abs(actual - expected) < 0.01, f"Total {actual} != sum {expected}"


@then("the total is between 0 and 5")
def step_total_range(context):
    total = context.score_data["total"]
    assert 0 <= total <= 5, f"Total {total} out of range [0, 5]"


@then("justifications exist for P1, P2, P3, P4, and P5")
def step_all_justifications(context):
    justs = context.score_data["justifications"]
    for p in PRINCIPLES:
        assert p in justs, f"Missing justification for {p}"


@then("every justification is a non-empty string")
def step_non_empty_justifications(context):
    justs = context.score_data["justifications"]
    for p in PRINCIPLES:
        val = justs.get(p, "")
        assert isinstance(val, str) and len(val.strip()) > 0, \
            f"Justification for {p} is empty or not a string"


# ---------- scoring_completeness.feature ----------

@given('scoring results exist for model "{model}"')
def step_results_for_model(context, model):
    context.model = model
    d = scores_dir(model)
    assert os.path.isdir(d), f"No results directory for {model}"
    context.score_files = load_all_score_files(model)
    assert len(context.score_files) > 0, f"No score files for {model}"


@then("34 score files exist for that model")
def step_34_files(context):
    count = len(context.score_files)
    assert count == 34, f"Expected 34 score files, got {count}"


def _strip_suffix(name, suffix):
    """Strip a single trailing suffix, not all occurrences."""
    if name.endswith(suffix):
        return name[:-len(suffix)]
    return name


@then("every agent file in the corpus has a corresponding score file")
def step_corpus_coverage(context):
    agent_names = set()
    for f in os.listdir(agents_dir()):
        if f.endswith(".md"):
            agent_names.add(_strip_suffix(f, ".md"))
    score_names = set()
    for fname, _ in context.score_files:
        # Handle both name.json and name_runN.json
        base = _strip_suffix(fname, ".json")
        if "_run" in base:
            base = base.rsplit("_run", 1)[0]
        score_names.add(base)
    missing = agent_names - score_names
    assert len(missing) == 0, f"Missing scores for: {missing}"


@then("no score file exists without a corresponding agent file")
def step_no_orphan_scores(context):
    agent_names = set()
    for f in os.listdir(agents_dir()):
        if f.endswith(".md"):
            agent_names.add(_strip_suffix(f, ".md"))
    for fname, _ in context.score_files:
        base = _strip_suffix(fname, ".json")
        if "_run" in base:
            base = base.rsplit("_run", 1)[0]
        assert base in agent_names, f"Orphan score file: {fname}"


# ---------- cross_model_divergence.feature ----------

@given("scoring results exist for at least two models")
def step_at_least_two_models(context):
    found = models_with_results()
    assert len(found) >= 2, f"Only {len(found)} model(s) have results"
    context.active_models = found


@given("scoring results exist for at least one model")
def step_at_least_one_model(context):
    found = models_with_results()
    assert len(found) >= 1, "No models have results"
    context.active_models = found


@when("the analysis script runs")
def step_run_analysis(context):
    script = os.path.join(PROJECT_DIR, "scripts", "analyse.py")
    result = subprocess.run(
        [sys.executable, script],
        cwd=PROJECT_DIR,
        capture_output=True,
        text=True,
    )
    context.analysis_exit_code = result.returncode
    context.analysis_stdout = result.stdout
    context.analysis_stderr = result.stderr
    assert result.returncode == 0, \
        f"analyse.py failed (exit {result.returncode}): {result.stderr}"


@then("divergence.md exists in results/summary/")
def step_divergence_exists(context):
    path = os.path.join(summary_dir(), "divergence.md")
    assert os.path.isfile(path), f"divergence.md not found at {path}"


@then("scores_by_principle.md exists in results/summary/")
def step_scores_by_principle_exists(context):
    path = os.path.join(summary_dir(), "scores_by_principle.md")
    assert os.path.isfile(path), f"scores_by_principle.md not found at {path}"


@then("scores_by_repo.md exists in results/summary/")
def step_scores_by_repo_exists(context):
    path = os.path.join(summary_dir(), "scores_by_repo.md")
    assert os.path.isfile(path), f"scores_by_repo.md not found at {path}"


@then("paper_tables.md exists in results/summary/")
def step_paper_tables_exists(context):
    path = os.path.join(summary_dir(), "paper_tables.md")
    assert os.path.isfile(path), f"paper_tables.md not found at {path}"
