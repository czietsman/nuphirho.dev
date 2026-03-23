#!/usr/bin/env python3
"""Analyse experiment results and produce a summary table."""

import json
import os
import sys

PROJECT_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
MODELS = ["claude", "codex", "gemini", "q"]


def load_ground_truth():
    with open(os.path.join(PROJECT_DIR, "ground_truth.json")) as f:
        return json.load(f)


def check_bdd_result(function_name):
    """Check if BDD caught the bug (i.e., test failed)."""
    result_file = os.path.join(PROJECT_DIR, "results", "bdd", f"{function_name}.txt")
    if not os.path.exists(result_file):
        return "NOT RUN"
    with open(result_file) as f:
        content = f.read()
    # behave reports "failed" when a scenario fails
    if "failed" in content.lower():
        return "FAIL (caught)"
    return "PASS (missed)"


def check_review_result(model, function_name, keywords):
    """Check how many review runs detected the planted bug."""
    review_dir = os.path.join(PROJECT_DIR, "results", "review", model, function_name)
    if not os.path.exists(review_dir):
        return 0, 0, 0

    runs = sorted(
        f for f in os.listdir(review_dir)
        if f.startswith("run_") and f.endswith(".json")
    )
    total = len(runs)
    detected = 0
    parse_failures = 0

    for run_file in runs:
        filepath = os.path.join(review_dir, run_file)
        try:
            with open(filepath) as f:
                content = f.read()

            # Try to parse as JSON to check for parse failures
            stripped = content.strip()
            try:
                json.loads(stripped)
            except (json.JSONDecodeError, ValueError):
                parse_failures += 1

            # Check if any detection keyword appears in the output
            content_lower = content.lower()
            if any(kw.lower() in content_lower for kw in keywords):
                detected += 1
        except IOError:
            parse_failures += 1

    return detected, total, parse_failures


def main():
    gt = load_ground_truth()

    # Header
    header = "| Function | BDD |"
    separator = "|---|---|"
    for model in MODELS:
        header += f" {model.capitalize()} |"
        separator += "---|"
    header += " Notes |"
    separator += "---|"

    print(header)
    print(separator)

    for func in gt["functions"]:
        name = func["name"]
        keywords = func["detection_keywords"]

        bdd = check_bdd_result(name)
        row = f"| {name} | {bdd} |"

        notes = []
        for model in MODELS:
            detected, total, parse_failures = check_review_result(model, name, keywords)

            if total > 0:
                rate = f"{detected}/{total}"
            else:
                rate = "NOT RUN"

            if parse_failures > 0:
                notes.append(f"{model}: {parse_failures} JSON parse failure(s)")

            row += f" {rate} |"

        row += f" {'; '.join(notes) if notes else ''} |"
        print(row)


if __name__ == "__main__":
    main()
