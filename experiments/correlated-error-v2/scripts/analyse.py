#!/usr/bin/env python3
"""Analyse experiment results and produce a summary table."""

import json
import os
import re
import sys

PROJECT_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


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


def check_review_result(function_name, keywords):
    """Check how many review runs detected the planted bug."""
    review_dir = os.path.join(PROJECT_DIR, "results", "review", function_name)
    if not os.path.exists(review_dir):
        return 0, 0

    runs = sorted(f for f in os.listdir(review_dir) if f.startswith("run_") and f.endswith(".json"))
    total = len(runs)
    detected = 0

    for run_file in runs:
        filepath = os.path.join(review_dir, run_file)
        try:
            with open(filepath) as f:
                content = f.read().lower()
            # Check if any detection keyword appears in the output
            if any(kw.lower() in content for kw in keywords):
                detected += 1
        except IOError:
            pass

    return detected, total


def main():
    gt = load_ground_truth()

    print("| Function | BDD result | AI review detections | AI review rate |")
    print("|---|---|---|---|")

    for func in gt["functions"]:
        name = func["name"]
        keywords = func["detection_keywords"]

        bdd = check_bdd_result(name)
        detected, total = check_review_result(name, keywords)

        if total > 0:
            rate = f"{detected}/{total}"
            pct = f"{detected/total*100:.0f}%"
        else:
            rate = "NOT RUN"
            pct = "N/A"

        print(f"| {name} | {bdd} | {rate} | {pct} |")


if __name__ == "__main__":
    main()
