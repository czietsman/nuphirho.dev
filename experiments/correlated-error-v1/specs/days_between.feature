Feature: Calculate days between two dates

  Scenario: Dates in chronological order
    Given date1 is "2024-01-01" and date2 is "2024-01-10"
    When I calculate the days between them
    Then the days result should be 9

  Scenario: Dates in reverse order
    Given date1 is "2024-01-10" and date2 is "2024-01-01"
    When I calculate the days between them
    Then the days result should be 9

  Scenario: Same date
    Given date1 is "2024-06-15" and date2 is "2024-06-15"
    When I calculate the days between them
    Then the days result should be 0

  Scenario: Dates spanning a year boundary
    Given date1 is "2023-12-31" and date2 is "2024-01-01"
    When I calculate the days between them
    Then the days result should be 1
