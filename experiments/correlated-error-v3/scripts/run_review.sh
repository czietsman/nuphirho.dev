#!/usr/bin/env bash
# Usage: ./scripts/run_review.sh <model> <function_file> [runs]
# model: claude | codex | gemini | q

set -euo pipefail

MODEL=$1
FUNCTION_FILE=$2
RUNS=${3:-5}

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

FUNCTION_NAME=$(basename "$FUNCTION_FILE" .py)
RESULTS_DIR="results/review/${MODEL}/${FUNCTION_NAME}"
mkdir -p "$RESULTS_DIR"

PROMPT_TEMPLATE=$(cat prompts/review.txt)
SOURCE=$(cat "$FUNCTION_FILE")
FULL_PROMPT="${PROMPT_TEMPLATE/\{function_source\}/$SOURCE}"

for i in $(seq 1 "$RUNS"); do
  echo "[$MODEL] Run $i of $RUNS for $FUNCTION_NAME"

  case "$MODEL" in
    claude)
      RAW=$(echo "$FULL_PROMPT" | claude -p 2>/dev/null)
      # Strip markdown fences
      RESULT=$(echo "$RAW" | sed 's/^```json//;s/^```//' | tr -d '`')
      ;;
    codex)
      RAW=$(codex exec "$FULL_PROMPT" 2>/dev/null)
      # Try to extract line after "codex" label; fall back to raw output
      LABEL_MATCH=$(echo "$RAW" | grep -A1 "^codex$" | tail -1 || true)
      if [ -n "$LABEL_MATCH" ]; then
        RESULT="$LABEL_MATCH"
      else
        RESULT="$RAW"
      fi
      ;;
    gemini)
      # Use -p flag for non-interactive mode; run from /tmp to prevent filesystem access
      RAW=$(cd /tmp && gemini -p "$FULL_PROMPT" -s 2>/dev/null)
      RESULT="$RAW"
      ;;
    q)
      RAW=$(q chat --non-interactive "$FULL_PROMPT" 2>/dev/null)
      # Strip ANSI escape codes, then filter banner and echo lines
      STRIPPED=$(echo "$RAW" | sed 's/\x1b\[[0-9;]*m//g')
      RESULT=$(echo "$STRIPPED" | grep -v "^>\|^(\|^/\|━\|╭\|╰\|│\|You are chatting\|Did you know\|ctrl\|settings\|checkpoints" || true)
      ;;
    *)
      echo "Unknown model: $MODEL"
      exit 1
      ;;
  esac

  echo "$RESULT" > "$RESULTS_DIR/run_${i}.json"
  sleep 2
done
