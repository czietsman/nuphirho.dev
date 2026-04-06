Feature: Scoring output structure

  Scenario: Score file contains all required fields
    Given a score file exists for any repository and model
    Then the file contains a "scores" object
    And the file contains a "total" field
    And the file contains a "justifications" object

  Scenario: Scores are within valid range
    Given a score file exists for any repository and model
    Then every principle score is 0, 0.5, or 1
    And the total equals the sum of the five principle scores
    And the total is between 0 and 5

  Scenario: Justifications are present for all principles
    Given a score file exists for any repository and model
    Then justifications exist for P1, P2, P3, P4, and P5
    And every justification is a non-empty string
