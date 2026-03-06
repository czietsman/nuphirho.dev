Feature: Pre-push secret detection
  As the repository owner
  I want secrets to be detected before code is pushed
  So that credentials are never exposed in the public git history

  Background:
    Given a pre-push hook is installed via husky
    And the hook scans added lines in the commit range for secret patterns

  Scenario: AWS access key is detected
    Given a commit adds a line containing "AKIAIOSFODNN7EXAMPLE"
    When the pre-push hook runs
    Then the push is blocked
    And the output identifies the offending file and line

  Scenario: GitHub personal access token is detected
    Given a commit adds a line containing "ghp_ABCDEFghijklmnopqrstuvwxyz0123456789"
    When the pre-push hook runs
    Then the push is blocked
    And the output identifies the offending file and line

  Scenario: GitHub fine-grained token is detected
    Given a commit adds a line containing "github_pat_ABCDE12345fghij67890klmnop"
    When the pre-push hook runs
    Then the push is blocked
    And the output identifies the offending file and line

  Scenario: PEM private key header is detected
    Given a commit adds a line containing "-----BEGIN RSA PRIVATE KEY-----"
    When the pre-push hook runs
    Then the push is blocked
    And the output identifies the offending file and line

  Scenario: Generic secret assignment is detected
    Given a commit adds a line containing "CLOUDFLARE_API_TOKEN=v1.0abc123def456ghi789jkl012mno"
    When the pre-push hook runs
    Then the push is blocked
    And the output identifies the offending file and line

  Scenario: Clean commit passes the scan
    Given a commit adds only non-secret content
    When the pre-push hook runs
    Then the push is allowed

  Scenario: Branch deletion is ignored
    Given the push is deleting a remote branch
    When the pre-push hook runs
    Then the hook exits without scanning

  Scenario: Paths in .secretscanignore are excluded
    Given a commit adds a line containing "AKIAIOSFODNN7EXAMPLE"
    And the file is listed in ".secretscanignore"
    When the pre-push hook runs
    Then the push is allowed

  Scenario: Bypass with --no-verify
    Given a commit adds a line containing a secret
    When the developer pushes with "--no-verify"
    Then the pre-push hook is skipped by git
    And the push proceeds
