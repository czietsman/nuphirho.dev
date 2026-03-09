Feature: Contract probes
  As the pipeline operator
  I want to validate credentials and API shapes before publishing
  So that configuration errors are caught early, not mid-pipeline

  # --- All pass ---

  Scenario: All probes pass with valid credentials
    Given valid credentials are configured for all platforms
    When the probe command runs
    Then all probes pass
    And the exit code is 0
    And the output contains "[hashnode] credentials        OK"
    And the output contains "[hashnode] publication        OK"
    And the output contains "[hashnode] draft mutation     OK"
    And the output contains "[hashnode] cleanup            OK"
    And the output contains "[devto]    credentials        OK"
    And the output contains "[devto]    article mutation   OK"
    And the output contains "[devto]    cleanup            OK"
    And the output contains "All probes passed."

  # --- Hashnode failures ---

  Scenario: Hashnode probe fails with invalid token
    Given an invalid Hashnode token is configured
    And valid Dev.to credentials are configured
    When the probe command runs
    Then the Hashnode credentials probe fails
    And the output contains "[hashnode] credentials        FAIL"
    And the exit code is 1

  Scenario: Hashnode probe fails when publication ID does not resolve
    Given a valid Hashnode token is configured
    And an unknown publication ID is configured
    And valid Dev.to credentials are configured
    When the probe command runs
    Then the Hashnode publication probe fails
    And the output contains "[hashnode] publication        FAIL"
    And the exit code is 1

  Scenario: Hashnode probe fails when API is unreachable
    Given the Hashnode API is unreachable
    And valid Dev.to credentials are configured
    When the probe command runs
    Then the Hashnode credentials probe fails
    And the output contains "[hashnode] credentials        FAIL"
    And the exit code is 1

  # --- Dev.to failures ---

  Scenario: Dev.to probe fails with invalid API key
    Given valid Hashnode credentials are configured
    And an invalid Dev.to API key is configured
    When the probe command runs
    Then the Dev.to credentials probe fails
    And the output contains "[devto]    credentials        FAIL"
    And the exit code is 1

  Scenario: Dev.to probe fails when API is unreachable
    Given valid Hashnode credentials are configured
    And the Dev.to API is unreachable
    When the probe command runs
    Then the Dev.to credentials probe fails
    And the output contains "[devto]    credentials        FAIL"
    And the exit code is 1

  # --- Cleanup ---

  Scenario: Probe cleans up draft on Hashnode
    Given valid credentials are configured for all platforms
    When the probe command runs
    Then the Hashnode draft created by the probe is deleted
    And the output contains "[hashnode] cleanup            OK"

  Scenario: Probe cleans up article on Dev.to
    Given valid credentials are configured for all platforms
    When the probe command runs
    Then the Dev.to article created by the probe is deleted
    And the output contains "[devto]    cleanup            OK"

  Scenario: Hashnode cleanup failure is a warning not an error
    Given valid Hashnode credentials are configured
    And Hashnode draft deletion will fail
    And valid Dev.to credentials are configured
    When the probe command runs
    Then the output contains "[hashnode] cleanup            WARN"
    And the exit code is 0

  Scenario: Dev.to cleanup failure is a warning not an error
    Given valid Hashnode credentials are configured
    And valid Dev.to credentials are configured
    And Dev.to article deletion will fail
    When the probe command runs
    Then the output contains "[devto]    cleanup            WARN"
    And the exit code is 0

  # --- Summary ---

  Scenario: Failed probe reports summary with failure count
    Given an invalid Hashnode token is configured
    And an invalid Dev.to API key is configured
    When the probe command runs
    Then the output contains "probes failed"
    And the exit code is 1
