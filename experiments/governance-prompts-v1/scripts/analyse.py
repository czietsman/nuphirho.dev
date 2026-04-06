#!/usr/bin/env python3
"""Analyse governance-prompts-v1 scoring results and produce summary tables."""

import json
import os
import sys

PROJECT_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
MODELS = ["claude", "codex", "gemini"]
PRINCIPLES = ["P1", "P2", "P3", "P4", "P5"]

VERSION = sys.argv[1] if len(sys.argv) > 1 else "v1"


def scores_dir(model):
    return os.path.join(PROJECT_DIR, "results", f"scores-{VERSION}", model)


def summary_dir():
    return os.path.join(PROJECT_DIR, "results", "summary")


def load_scores(model):
    """Return dict of {basename: parsed_json} for a model."""
    d = scores_dir(model)
    if not os.path.isdir(d):
        return {}
    results = {}
    for f in sorted(os.listdir(d)):
        if not f.endswith(".json"):
            continue
        base = f.replace(".json", "")
        if "_run" in base:
            base = base.rsplit("_run", 1)[0]
        path = os.path.join(d, f)
        try:
            with open(path) as fh:
                data = json.load(fh)
            results[base] = data
        except (json.JSONDecodeError, IOError):
            pass
    return results


def active_models():
    """Return models that have at least one score file."""
    return [m for m in MODELS if load_scores(m)]


def write_scores_by_principle(models_data):
    """Write score distribution per principle across all files and models."""
    out = summary_dir()
    lines = ["# Scores by Principle\n"]

    header = "| Principle | Score |"
    sep = "|---|---|"
    for model in models_data:
        header += f" {model.capitalize()} (n) |"
        sep += "---|"
    lines.append(header)
    lines.append(sep)

    for p in PRINCIPLES:
        for score_val in [0, 0.5, 1]:
            row = f"| {p} | {score_val} |"
            for model, data in models_data.items():
                count = sum(
                    1 for d in data.values()
                    if d.get("scores", {}).get(p) == score_val
                )
                row += f" {count} |"
            lines.append(row)

    # Mean per principle per model
    lines.append("")
    lines.append("## Mean score per principle")
    lines.append("")
    header = "| Principle |"
    sep = "|---|"
    for model in models_data:
        header += f" {model.capitalize()} |"
        sep += "---|"
    lines.append(header)
    lines.append(sep)

    for p in PRINCIPLES:
        row = f"| {p} |"
        for model, data in models_data.items():
            vals = [d.get("scores", {}).get(p, 0) for d in data.values()]
            mean = sum(vals) / len(vals) if vals else 0
            row += f" {mean:.2f} |"
        lines.append(row)

    with open(os.path.join(out, "scores_by_principle.md"), "w") as f:
        f.write("\n".join(lines) + "\n")


def write_scores_by_repo(models_data):
    """Write total score per repo per model."""
    out = summary_dir()
    lines = ["# Scores by Repository\n"]

    header = "| # | Repository |"
    sep = "|---|---|"
    for model in models_data:
        header += f" {model.capitalize()} |"
        sep += "---|"
    header += " Mean | Range |"
    sep += "---|---|"
    lines.append(header)
    lines.append(sep)

    # Collect all repo names across models
    all_repos = set()
    for data in models_data.values():
        all_repos.update(data.keys())

    for i, repo in enumerate(sorted(all_repos), 1):
        row = f"| {i} | {repo} |"
        totals = []
        for model, data in models_data.items():
            if repo in data:
                total = data[repo].get("total", 0)
                totals.append(total)
                row += f" {total} |"
            else:
                row += " -- |"

        if totals:
            mean = sum(totals) / len(totals)
            rng = max(totals) - min(totals)
            row += f" {mean:.2f} | {rng:.1f} |"
        else:
            row += " -- | -- |"
        lines.append(row)

    with open(os.path.join(out, "scores_by_repo.md"), "w") as f:
        f.write("\n".join(lines) + "\n")


def write_divergence(models_data):
    """Write cross-model divergence report."""
    out = summary_dir()
    lines = ["# Cross-Model Divergence Report\n"]

    model_names = list(models_data.keys())
    if len(model_names) < 2:
        lines.append("Divergence requires at least two models. Skipping.\n")
        with open(os.path.join(out, "divergence.md"), "w") as f:
            f.write("\n".join(lines) + "\n")
        return

    # Repos with divergence > 1
    lines.append("## Repositories with total score divergence > 1\n")
    header = "| Repository |"
    sep = "|---|"
    for model in model_names:
        header += f" {model.capitalize()} |"
        sep += "---|"
    header += " Range |"
    sep += "---|"
    lines.append(header)
    lines.append(sep)

    all_repos = set()
    for data in models_data.values():
        all_repos.update(data.keys())

    divergent_count = 0
    for repo in sorted(all_repos):
        totals = []
        row = f"| {repo} |"
        for model in model_names:
            data = models_data[model]
            if repo in data:
                total = data[repo].get("total", 0)
                totals.append(total)
                row += f" {total} |"
            else:
                row += " -- |"
        if len(totals) >= 2:
            rng = max(totals) - min(totals)
            row += f" {rng:.1f} |"
            if rng > 1:
                lines.append(row)
                divergent_count += 1
        else:
            row += " -- |"

    if divergent_count == 0:
        lines.append("| (none) |" + " -- |" * (len(model_names) + 1))

    # Per-principle agreement
    lines.append("")
    lines.append("## Per-principle agreement rates\n")
    lines.append("Agreement = both models score within 0.5 of each other.\n")

    pairs = []
    for i, m1 in enumerate(model_names):
        for m2 in model_names[i + 1:]:
            pairs.append((m1, m2))

    header = "| Principle |"
    sep = "|---|"
    for m1, m2 in pairs:
        header += f" {m1.capitalize()}-{m2.capitalize()} |"
        sep += "---|"
    lines.append(header)
    lines.append(sep)

    for p in PRINCIPLES:
        row = f"| {p} |"
        for m1, m2 in pairs:
            d1 = models_data[m1]
            d2 = models_data[m2]
            common = set(d1.keys()) & set(d2.keys())
            if not common:
                row += " -- |"
                continue
            agree = sum(
                1 for repo in common
                if abs(d1[repo].get("scores", {}).get(p, 0) -
                       d2[repo].get("scores", {}).get(p, 0)) <= 0.5
            )
            pct = agree / len(common) * 100
            row += f" {pct:.0f}% |"
        lines.append(row)

    with open(os.path.join(out, "divergence.md"), "w") as f:
        f.write("\n".join(lines) + "\n")


def write_paper_tables(models_data):
    """Write tables formatted for Section 4.3 of the paper."""
    out = summary_dir()
    lines = ["# Paper Tables -- Section 4.3\n"]

    model_names = list(models_data.keys())

    # Table 1: Score distribution per principle (counts and percentages)
    lines.append("## Table 1: Score distribution per principle\n")
    for p in PRINCIPLES:
        lines.append(f"### {p}\n")
        header = "| Score |"
        sep = "|---|"
        for model in model_names:
            header += f" {model.capitalize()} (n) | {model.capitalize()} (%) |"
            sep += "---|---|"
        lines.append(header)
        lines.append(sep)

        for score_val in [0, 0.5, 1]:
            row = f"| {score_val} |"
            for model in model_names:
                data = models_data[model]
                total_files = len(data)
                count = sum(
                    1 for d in data.values()
                    if d.get("scores", {}).get(p) == score_val
                )
                pct = count / total_files * 100 if total_files > 0 else 0
                row += f" {count} | {pct:.0f}% |"
            lines.append(row)
        lines.append("")

    # Table 2: Mean total score
    lines.append("## Table 2: Mean total score\n")
    header = "| Metric |"
    sep = "|---|"
    for model in model_names:
        header += f" {model.capitalize()} |"
        sep += "---|"
    header += " Overall |"
    sep += "---|"
    lines.append(header)
    lines.append(sep)

    all_totals = []
    row_mean = "| Mean total |"
    row_median = "| Median total |"
    row_below = "| Below 2.5 (%) |"
    row_below_count = "| Below 2.5 (n) |"

    for model in model_names:
        data = models_data[model]
        totals = sorted(d.get("total", 0) for d in data.values())
        all_totals.extend(totals)
        n = len(totals)
        mean = sum(totals) / n if n > 0 else 0
        median = totals[n // 2] if n > 0 else 0
        below = sum(1 for t in totals if t < 2.5)
        below_pct = below / n * 100 if n > 0 else 0
        row_mean += f" {mean:.2f} |"
        row_median += f" {median:.1f} |"
        row_below += f" {below_pct:.0f}% |"
        row_below_count += f" {below} |"

    if all_totals:
        overall_mean = sum(all_totals) / len(all_totals)
        all_sorted = sorted(all_totals)
        overall_median = all_sorted[len(all_sorted) // 2]
        overall_below = sum(1 for t in all_totals if t < 2.5)
        overall_below_pct = overall_below / len(all_totals) * 100
        row_mean += f" {overall_mean:.2f} |"
        row_median += f" {overall_median:.1f} |"
        row_below += f" {overall_below_pct:.0f}% |"
        row_below_count += f" {overall_below} |"
    else:
        row_mean += " -- |"
        row_median += " -- |"
        row_below += " -- |"
        row_below_count += " -- |"

    lines.append(row_mean)
    lines.append(row_median)
    lines.append(row_below)
    lines.append(row_below_count)

    # Table 3: Highest and lowest scoring repositories
    lines.append("")
    lines.append("## Table 3: Highest and lowest scoring repositories\n")

    # Compute overall mean total per repo
    all_repos = set()
    for data in models_data.values():
        all_repos.update(data.keys())

    repo_means = {}
    for repo in all_repos:
        totals = []
        for data in models_data.values():
            if repo in data:
                totals.append(data[repo].get("total", 0))
        if totals:
            repo_means[repo] = sum(totals) / len(totals)

    sorted_repos = sorted(repo_means.items(), key=lambda x: x[1], reverse=True)

    lines.append("### Top 5\n")
    lines.append("| Repository | Mean total |")
    lines.append("|---|---|")
    for repo, mean in sorted_repos[:5]:
        lines.append(f"| {repo} | {mean:.2f} |")

    lines.append("")
    lines.append("### Bottom 5\n")
    lines.append("| Repository | Mean total |")
    lines.append("|---|---|")
    for repo, mean in sorted_repos[-5:]:
        lines.append(f"| {repo} | {mean:.2f} |")

    # Table 4: Weakest and strongest principles
    lines.append("")
    lines.append("## Table 4: Principle strength ranking\n")
    lines.append("| Rank | Principle | Overall mean |")
    lines.append("|---|---|---|")

    principle_means = {}
    for p in PRINCIPLES:
        vals = []
        for data in models_data.values():
            for d in data.values():
                vals.append(d.get("scores", {}).get(p, 0))
        principle_means[p] = sum(vals) / len(vals) if vals else 0

    ranked = sorted(principle_means.items(), key=lambda x: x[1], reverse=True)
    for i, (p, mean) in enumerate(ranked, 1):
        lines.append(f"| {i} | {p} | {mean:.2f} |")

    with open(os.path.join(out, "paper_tables.md"), "w") as f:
        f.write("\n".join(lines) + "\n")


def main():
    os.makedirs(summary_dir(), exist_ok=True)

    found = active_models()
    if not found:
        print("No scoring results found. Nothing to analyse.")
        return

    models_data = {m: load_scores(m) for m in found}

    print(f"Models with results: {', '.join(found)}")
    for model, data in models_data.items():
        print(f"  {model}: {len(data)} files scored")

    write_scores_by_principle(models_data)
    print("Wrote scores_by_principle.md")

    write_scores_by_repo(models_data)
    print("Wrote scores_by_repo.md")

    if len(found) >= 2:
        write_divergence(models_data)
        print("Wrote divergence.md")
    else:
        # Write placeholder so BDD specs that check for file existence
        # only fail when divergence is expected (2+ models).
        path = os.path.join(summary_dir(), "divergence.md")
        if not os.path.exists(path):
            with open(path, "w") as f:
                f.write("# Cross-Model Divergence Report\n\n"
                        "Divergence requires at least two models. Skipping.\n")
            print("Wrote divergence.md (placeholder -- single model)")

    write_paper_tables(models_data)
    print("Wrote paper_tables.md")

    print("Analysis complete.")


if __name__ == "__main__":
    main()
