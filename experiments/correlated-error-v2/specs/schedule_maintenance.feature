Feature: Aircraft maintenance scheduling

  Scenario: Both hours and cycles exceed limits
    Given last service was "2024-01-15"
    And the aircraft has 520 flight hours and 310 cycles
    When I check maintenance status on "2024-06-15"
    Then maintenance should be due
    And the reasons should include "flight_hours"
    And the reasons should include "cycles"

  Scenario: Only flight hours exceed limit
    Given last service was "2024-01-15"
    And the aircraft has 520 flight hours and 200 cycles
    When I check maintenance status on "2024-06-15"
    Then maintenance should be due
    And the reasons should include "flight_hours"

  Scenario: Only cycles exceed limit
    Given last service was "2024-01-15"
    And the aircraft has 400 flight hours and 310 cycles
    When I check maintenance status on "2024-06-15"
    Then maintenance should be due
    And the reasons should include "cycles"

  Scenario: Calendar limit exceeded
    Given last service was "2023-01-01"
    And the aircraft has 100 flight hours and 50 cycles
    When I check maintenance status on "2024-06-15"
    Then maintenance should be due
    And the reasons should include "calendar"

  Scenario: No thresholds exceeded
    Given last service was "2024-06-01"
    And the aircraft has 100 flight hours and 50 cycles
    When I check maintenance status on "2024-06-15"
    Then maintenance should not be due
