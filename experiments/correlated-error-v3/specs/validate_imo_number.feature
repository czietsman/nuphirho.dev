Feature: IMO number validation per IMO Resolution A.1079(28)

  Scenario: Valid IMO number with descending weights 7 to 2
    Given an IMO number "IMO9074729"
    When I validate the IMO number
    Then the IMO number should be valid

  Scenario: Another valid IMO number
    Given an IMO number "IMO9000027"
    When I validate the IMO number
    Then the IMO number should be valid

  Scenario: Invalid check digit is rejected
    Given an IMO number "IMO9074720"
    When I validate the IMO number
    Then the IMO number should be invalid

  Scenario: IMO number without prefix
    Given an IMO number "9074729"
    When I validate the IMO number
    Then the IMO number should be valid
