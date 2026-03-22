Feature: Venture capital dilution with option pool shuffle

  Scenario: Standard Series A with 10% option pool
    Given 9000000 existing shares
    And a 10% option pool
    And an investment of 2000000 at 10000000 pre-money valuation
    When I calculate the dilution
    Then founder ownership should be 75.00%
    And investor ownership should be 16.67%
    And pool ownership should be 8.33%

  Scenario: Series A with 20% option pool
    Given 8000000 existing shares
    And a 20% option pool
    And an investment of 5000000 at 10000000 pre-money valuation
    When I calculate the dilution
    Then founder ownership should be 53.33%
    And investor ownership should be 33.33%
    And pool ownership should be 13.33%
