#!/usr/bin/env bash
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"
mkdir -p results/bdd
for FEATURE_FILE in specs/*.feature; do
  FEATURE_NAME=$(basename "$FEATURE_FILE" .feature)
  echo "Running BDD for $FEATURE_NAME"
  python3 -m behave "$FEATURE_FILE" --no-capture 2>&1 | tee "results/bdd/${FEATURE_NAME}.txt"
done
