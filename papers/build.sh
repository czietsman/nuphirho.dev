#!/bin/bash
# build.sh -- Build specification-as-quality-gate.pdf from source
#
# Requirements: see README.md for system dependencies
# Usage: ./build.sh [--clean]
#
# The script must be run from the papers/ directory, or the PAPER variable
# adjusted to point to the correct path. arxiv.sty must be present in the
# same directory as the .tex file.

set -e

PAPER="specification-as-quality-gate"
TEX_FILE="${PAPER}.tex"

# Clean mode
if [[ "$1" == "--clean" ]]; then
  echo "Cleaning auxiliary files..."
  rm -f "${PAPER}.aux" "${PAPER}.log" "${PAPER}.out" "${PAPER}.toc"
  echo "Done."
  exit 0
fi

# Check for required files
if [[ ! -f "${TEX_FILE}" ]]; then
  echo "Error: ${TEX_FILE} not found. Run this script from the papers/ directory."
  exit 1
fi

if [[ ! -f "arxiv.sty" ]]; then
  echo "Error: arxiv.sty not found. It must be present in the same directory as ${TEX_FILE}."
  exit 1
fi

# Check for pdflatex
if ! command -v pdflatex &> /dev/null; then
  echo "Error: pdflatex not found. Install texlive-latex-base."
  exit 1
fi

echo "Building ${PAPER}.pdf..."

# First pass -- generates .aux file for cross-references
pdflatex -interaction=nonstopmode "${TEX_FILE}"

# Second pass -- resolves references
pdflatex -interaction=nonstopmode "${TEX_FILE}"

echo ""
echo "Build complete: ${PAPER}.pdf"
