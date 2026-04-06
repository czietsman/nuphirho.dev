Feature: Scoring completeness

  Scenario: All 34 repositories are scored for a model
    Given scoring results exist for model "claude"
    Then 34 score files exist for that model

  Scenario: Score files match the corpus
    Given scoring results exist for model "claude"
    Then every agent file in the corpus has a corresponding score file
    And no score file exists without a corresponding agent file
