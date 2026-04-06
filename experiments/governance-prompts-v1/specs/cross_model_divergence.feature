Feature: Cross-model divergence

  Scenario: Divergence report is produced when multiple models have scored
    Given scoring results exist for at least two models
    When the analysis script runs
    Then divergence.md exists in results/summary/

  Scenario: Summary tables are produced
    Given scoring results exist for at least one model
    When the analysis script runs
    Then scores_by_principle.md exists in results/summary/
    And scores_by_repo.md exists in results/summary/
    And paper_tables.md exists in results/summary/
