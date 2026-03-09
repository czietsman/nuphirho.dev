Feature: CLI wrapper
  As the pipeline operator
  I want a CLI that parses flags, loads configuration, and delegates to the pipeline
  So that the publish flow can be invoked from the command line and CI

  # --- Probe ---

  Scenario: Probe flag calls probe and exits without files
    Given valid credentials via environment variables
    When the CLI runs with "--probe"
    Then ProbeAll is called
    And the exit code is 0

  # --- Dry-run ---

  Scenario: Dry-run flag produces JSON to stdout
    Given valid credentials via environment variables
    And a post file "my-post.md" with:
      | title | My Post |
      | slug  | my-post |
      | tags  | go      |
    When the CLI runs with "--dry-run"
    Then stdout contains valid JSON
    And the exit code is 0

  Scenario: Dry-run with no post files exits with empty results
    Given valid credentials via environment variables
    When the CLI runs with "--dry-run"
    Then stdout contains "[]"
    And the exit code is 0

  # --- Publish ---

  Scenario: Valid post file publishes successfully
    Given valid credentials via environment variables
    And a post file "my-post.md" with:
      | title | My Post |
      | slug  | my-post |
      | tags  | go      |
    When the CLI runs
    Then the exit code is 0

  Scenario: Validation failure exits with code 1
    Given valid credentials via environment variables
    And a post file "bad.md" with:
      | slug | bad |
      | tags | go  |
    When the CLI runs
    Then the exit code is 1

  Scenario: Publish failure exits with code 2
    Given valid credentials via environment variables
    And a post file "fail.md" with:
      | title | Fail |
      | slug  | fail |
      | tags  | go   |
    And Hashnode publish will fail with "server error"
    When the CLI runs
    Then the exit code is 2

  # --- Credential handling ---

  Scenario: Missing Hashnode token exits with clear error
    When the CLI runs
    Then stderr contains "Hashnode token is required"
    And the exit code is 1

  Scenario: Environment variables provide token fallback
    Given HASHNODE_TOKEN is set to "test-token"
    And HASHNODE_PUBLICATION_ID is set to "test-pub-id"
    And DEVTO_API_KEY is set to "test-api-key"
    And a post file "my-post.md" with:
      | title | My Post |
      | slug  | my-post |
      | tags  | go      |
    When the CLI runs
    Then the exit code is 0

  # --- Tags file ---

  Scenario: Tags file defaults to tags.json in working directory
    Given valid credentials via environment variables
    And a post file "my-post.md" with:
      | title | My Post |
      | slug  | my-post |
      | tags  | go      |
    When the CLI runs without --tags-file
    Then the exit code is 0

  # --- Unknown flag ---

  Scenario: Unknown flag exits with usage message
    When the CLI runs with "--unknown-flag"
    Then stderr contains "flag provided but not defined"
    And the exit code is 1
