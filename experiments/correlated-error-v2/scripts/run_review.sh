#!/usr/bin/env bash
# Usage: ./scripts/run_review.sh [runs_per_function]
# Runs Claude CLI review for each function in corpus/

RUNS=${1:-5}
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

for FUNCTION_FILE in corpus/*.py; do
  FUNCTION_NAME=$(basename "$FUNCTION_FILE" .py)
  RESULTS_DIR="results/review/${FUNCTION_NAME}"
  mkdir -p "$RESULTS_DIR"

  PROMPT=$(cat prompts/review.txt)
  SOURCE=$(cat "$FUNCTION_FILE")
  FULL_PROMPT="${PROMPT/\{function_source\}/$SOURCE}"

  for i in $(seq 1 "$RUNS"); do
    echo "Run $i of $RUNS for $FUNCTION_NAME"
    echo "$FULL_PROMPT" | claude -p > "$RESULTS_DIR/run_${i}.json"
    sleep 2
  done
done
