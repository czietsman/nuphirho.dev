Feature: Binary search on a sorted list

  Scenario: Target is the first element
    Given a sorted list [1, 3, 5, 7, 9]
    When I search for 1
    Then the result should be index 0

  Scenario: Target is in the middle
    Given a sorted list [1, 3, 5, 7, 9]
    When I search for 5
    Then the result should be index 2

  Scenario: Target is not in the list
    Given a sorted list [1, 3, 5, 7, 9]
    When I search for 4
    Then the result should be index -1

  Scenario: Single element list where target is present
    Given a sorted list [42]
    When I search for 42
    Then the result should be index 0

  Scenario: Empty list
    Given a sorted list []
    When I search for 1
    Then the result should be index -1
