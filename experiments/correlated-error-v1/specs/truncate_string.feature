Feature: Truncate a string to a maximum length

  Scenario: String shorter than max_len
    Given the string "hello" with max_len 10
    When I truncate the string
    Then the truncated result should be "hello"

  Scenario: String longer than max_len
    Given the string "hello world, this is a test" with max_len 10
    When I truncate the string
    Then the truncated result should be "hello w..."

  Scenario: String exactly max_len characters
    Given the string "hello worl" with max_len 10
    When I truncate the string
    Then the truncated result should be "hello worl"

  Scenario: String one character longer than max_len
    Given the string "hello world" with max_len 10
    When I truncate the string
    Then the truncated result should be "hello w..."
