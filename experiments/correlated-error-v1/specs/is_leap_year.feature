Feature: Leap year determination

  Scenario: Common year
    Given the year 2023
    When I check if it is a leap year
    Then the leap year result should be False

  Scenario: Year divisible by 4
    Given the year 2024
    When I check if it is a leap year
    Then the leap year result should be True

  Scenario: Century year not divisible by 400
    Given the year 1900
    When I check if it is a leap year
    Then the leap year result should be False

  Scenario: Century year divisible by 400
    Given the year 2000
    When I check if it is a leap year
    Then the leap year result should be True

  Scenario: Another century year not divisible by 400
    Given the year 1800
    When I check if it is a leap year
    Then the leap year result should be False
