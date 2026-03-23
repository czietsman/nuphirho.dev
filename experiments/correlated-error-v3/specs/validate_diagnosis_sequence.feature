Feature: Diagnosis sequence validation per ICD-10-CM Guidelines Section I.C.20.a

  Scenario: Valid diagnosis sequence with no external cause codes
    Given a diagnosis sequence of ["S72001A", "M80051A"]
    When I validate the diagnosis sequence
    Then the sequence should be valid

  Scenario: External cause code V00-Y99 is rejected as principal diagnosis
    Given a diagnosis sequence of ["V031XXA", "S61001A"]
    When I validate the diagnosis sequence
    Then the sequence should be invalid
    And the errors should include "external cause code cannot be the principal diagnosis"

  Scenario: External cause code is accepted in secondary position
    Given a diagnosis sequence of ["S61001A", "V031XXA"]
    When I validate the diagnosis sequence
    Then the sequence should be valid

  Scenario: External cause code Y93 is rejected as principal diagnosis
    Given a diagnosis sequence of ["Y9389", "S72001A"]
    When I validate the diagnosis sequence
    Then the sequence should be invalid
    And the errors should include "external cause code cannot be the principal diagnosis"
