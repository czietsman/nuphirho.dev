Feature: Gas day determination per NAESB WGQ 1.3.2

  Scenario: Timestamp after 09:00 Central Time belongs to the same calendar date gas day
    Given a timestamp of "2024-06-15T10:30:00-05:00"
    When I determine the gas day
    Then the gas day should be "2024-06-15"

  Scenario: Timestamp before 09:00 Central Time belongs to the previous calendar date gas day
    Given a timestamp of "2024-06-15T08:00:00-05:00"
    When I determine the gas day
    Then the gas day should be "2024-06-14"

  Scenario: Timestamp at exactly 09:00 Central Time starts a new gas day
    Given a timestamp of "2024-06-15T09:00:00-05:00"
    When I determine the gas day
    Then the gas day should be "2024-06-15"

  Scenario: Timestamp at midnight Central Time belongs to the previous gas day
    Given a timestamp of "2024-06-15T00:00:00-05:00"
    When I determine the gas day
    Then the gas day should be "2024-06-14"
