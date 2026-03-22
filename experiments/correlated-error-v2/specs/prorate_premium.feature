Feature: Insurance premium proration using ISDA actual/actual convention

  Scenario: Full year in a standard year
    Given an annual premium of 1200.00
    When I prorate from "2023-01-01" to "2023-12-31"
    Then the prorated premium should be 1196.71

  Scenario: Quarter in a standard year
    Given an annual premium of 3650.00
    When I prorate from "2023-01-01" to "2023-04-01"
    Then the prorated premium should be 900.00

  Scenario: Full year in a leap year
    Given an annual premium of 1200.00
    When I prorate from "2024-01-01" to "2024-12-31"
    Then the prorated premium should be 1196.72

  Scenario: Period in a leap year uses 366 divisor
    Given an annual premium of 3660.00
    When I prorate from "2024-01-01" to "2024-04-01"
    Then the prorated premium should be 910.00
