Feature: Progressive marginal tax calculation

  Scenario: Income within the first bracket
    Given tax brackets [(0, 0.10), (10000, 0.20), (40000, 0.30)]
    And an income of 5000
    When I calculate the tax
    Then the tax should be 500.00

  Scenario: Income spanning two brackets
    Given tax brackets [(0, 0.10), (10000, 0.20), (40000, 0.30)]
    And an income of 25000
    When I calculate the tax
    Then the tax should be 4000.00

  Scenario: Income spanning all three brackets
    Given tax brackets [(0, 0.10), (10000, 0.20), (40000, 0.30)]
    And an income of 60000
    When I calculate the tax
    Then the tax should be 13000.00
