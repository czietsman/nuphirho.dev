Feature: Governance prompts corpus capture

  Background:
    Given GITHUB_TOKEN is set in the environment
    And the corpus version is "v1"

  Scenario: All repositories are retrieved
    When the capture tool runs
    Then 34 agent files exist in experiments/governance-prompts-v1/agents/
    And the manifest contains 34 rows
    And no row in the manifest has status "Not found"

  Scenario: All commit SHAs are present
    When the capture tool runs
    Then every agent file has a commit SHA in its front matter
    And no commit SHA equals "unknown"
    And every SHA is exactly 40 characters

  Scenario: File content is not truncated
    When the capture tool runs
    Then every agent file has content beyond the front matter
    And the line count in the front matter matches the actual file content

  Scenario: Manifest is well-formed
    When the capture tool runs
    Then manifest.md exists at experiments/governance-prompts-v1/manifest.md
    And the manifest contains a header with retrieval date
    And the manifest table has columns: #, Repository, Short SHA, Full SHA, Lines, Status

  Scenario: Support files are present
    When the capture tool runs
    Then README.md exists at experiments/governance-prompts-v1/README.md
    And GOVERNANCE_PROMPTS.md exists at experiments/GOVERNANCE_PROMPTS.md
