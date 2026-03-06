#!/usr/bin/env bash
#
# Unit tests for the secret detection grep patterns used in .husky/pre-push.
# Run with: bash tests/test-secret-patterns.sh

set -euo pipefail

PASS=0
FAIL=0

# -- Patterns (must mirror .husky/pre-push) ------------------------------------

PATTERNS=(
  'AKIA[0-9A-Z]{16}'
  'ghp_[A-Za-z0-9]{36}'
  'github_pat_[A-Za-z0-9_]{22,}'
  'gho_[A-Za-z0-9]{36}'
  'ghu_[A-Za-z0-9]{36}'
  'ghr_[A-Za-z0-9]{36}'
  'ghs_[A-Za-z0-9]{36}'
  '-----BEGIN [A-Z ]*PRIVATE KEY-----'
  '(CLOUDFLARE_API_TOKEN|CLOUDFLARE_ACCOUNT_ID|HASHNODE_TOKEN|HASHNODE_PUBLICATION_ID|DEVTO_API_KEY|R2_ACCESS_KEY_ID|R2_SECRET_ACCESS_KEY)[[:space:]]*[=:][[:space:]]*[A-Za-z0-9_/+.=-]{20,}'
  '(api_key|api_secret|apikey|secret_key|access_key|auth_token|password|passwd|private_key|token)[[:space:]]*[=:][[:space:]]*['\''"]?[A-Za-z0-9_/+.=-]{20,}'
)

COMBINED=""
for p in "${PATTERNS[@]}"; do
  if [ -z "$COMBINED" ]; then
    COMBINED="$p"
  else
    COMBINED="$COMBINED|$p"
  fi
done

# -- Helpers -------------------------------------------------------------------

expect_match() {
  local label="$1"
  local input="$2"
  if echo "$input" | grep -Eq "$COMBINED"; then
    PASS=$((PASS + 1))
  else
    FAIL=$((FAIL + 1))
    echo "FAIL (should match): $label"
    echo "  input: $input"
  fi
}

expect_no_match() {
  local label="$1"
  local input="$2"
  if echo "$input" | grep -Eq "$COMBINED"; then
    FAIL=$((FAIL + 1))
    echo "FAIL (should not match): $label"
    echo "  input: $input"
  else
    PASS=$((PASS + 1))
  fi
}

# -- Test cases: should match --------------------------------------------------

expect_match "AWS access key" \
  "AKIAIOSFODNN7EXAMPLE"

expect_match "AWS key in assignment" \
  "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE"

expect_match "GitHub PAT (classic)" \
  "ghp_ABCDEFghijklmnopqrstuvwxyz0123456789"

expect_match "GitHub fine-grained PAT" \
  "github_pat_ABCDE12345fghij67890klmnop_extra"

expect_match "GitHub OAuth token" \
  "gho_ABCDEFghijklmnopqrstuvwxyz0123456789"

expect_match "GitHub user-to-server token" \
  "ghu_ABCDEFghijklmnopqrstuvwxyz0123456789"

expect_match "GitHub refresh token" \
  "ghr_ABCDEFghijklmnopqrstuvwxyz0123456789"

expect_match "GitHub server-to-server token" \
  "ghs_ABCDEFghijklmnopqrstuvwxyz0123456789"

expect_match "RSA private key header" \
  "-----BEGIN RSA PRIVATE KEY-----"

expect_match "EC private key header" \
  "-----BEGIN EC PRIVATE KEY-----"

expect_match "Generic private key header" \
  "-----BEGIN PRIVATE KEY-----"

expect_match "Cloudflare API token assignment" \
  "CLOUDFLARE_API_TOKEN=v1.0abcdef1234567890abcdef1234567890"

expect_match "Hashnode token assignment" \
  "HASHNODE_TOKEN=abcdefghij1234567890abcdefghij"

expect_match "Dev.to API key assignment" \
  "DEVTO_API_KEY=abcdefghij1234567890abcdefghij"

expect_match "R2 access key assignment" \
  "R2_ACCESS_KEY_ID=ABCDEFGHIJ1234567890ABCDEFGHIJ"

expect_match "R2 secret key assignment" \
  "R2_SECRET_ACCESS_KEY=abcdefghij1234567890abcdefghij"

expect_match "Generic api_key assignment (double quotes)" \
  'api_key = "abcdefghij1234567890abcdefghij"'

expect_match "Generic token assignment (single quotes)" \
  "token = 'abcdefghij1234567890abcdefghij'"

expect_match "Generic password assignment" \
  "password=abcdefghijklmnopqrstuvwxyz"

expect_match "Generic secret_key with colon" \
  "secret_key: abcdefghij1234567890abcdefghij"

# -- Test cases: should not match ----------------------------------------------

expect_no_match "Short string (not a secret)" \
  "token = abc"

expect_no_match "AKIA prefix with too few chars" \
  "AKIA1234567890"

expect_no_match "Normal code variable" \
  "const maxRetries = 3"

expect_no_match "URL without secrets" \
  "https://api.example.com/v1/users"

expect_no_match "Comment about tokens" \
  "# This function validates tokens"

expect_no_match "Git SHA (40 hex chars but no keyword)" \
  "abc123def456789012345678901234567890abcd"

expect_no_match "Package name with ghp prefix but wrong length" \
  "ghp_short"

expect_no_match "Markdown text" \
  "The deployment uses Cloudflare DNS with automatic SSL."

expect_no_match "Empty string" \
  ""

expect_no_match "Keyword with short value" \
  "password=short"

# -- Summary -------------------------------------------------------------------

echo ""
echo "Results: $PASS passed, $FAIL failed"

if [ "$FAIL" -gt 0 ]; then
  exit 1
fi

echo "All tests passed."
