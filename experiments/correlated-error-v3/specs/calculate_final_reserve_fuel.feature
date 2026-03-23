Feature: Final reserve fuel per ICAO Annex 6 Part I Section 4.3.6.3

  Scenario: Turbojet aircraft requires 30 minutes of reserve fuel
    Given a "jet" engine with a fuel burn rate of 2400 kg per hour
    When I calculate the final reserve fuel
    Then the reserve fuel should be 1200.00 kg

  Scenario: Piston aircraft requires 45 minutes of reserve fuel
    Given a "piston" engine with a fuel burn rate of 80 kg per hour
    When I calculate the final reserve fuel
    Then the reserve fuel should be 60.00 kg

  Scenario: Turbojet at lower burn rate still uses 30 minutes
    Given a "jet" engine with a fuel burn rate of 1800 kg per hour
    When I calculate the final reserve fuel
    Then the reserve fuel should be 900.00 kg

  Scenario: Piston at higher burn rate still uses 45 minutes
    Given a "piston" engine with a fuel burn rate of 120 kg per hour
    When I calculate the final reserve fuel
    Then the reserve fuel should be 90.00 kg
