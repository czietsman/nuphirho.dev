Feature: Electricity billing using progressive inclining block tariff

  Scenario: Consumption within the first tier
    Given electricity tiers [(100, 0.05), (300, 0.10), (inf, 0.15)]
    And consumption of 80 kWh
    When I calculate the electricity cost
    Then the electricity cost should be 4.00

  Scenario: Consumption spanning two tiers applies each rate progressively
    Given electricity tiers [(100, 0.05), (300, 0.10), (inf, 0.15)]
    And consumption of 200 kWh
    When I calculate the electricity cost
    Then the electricity cost should be 15.00

  Scenario: Consumption spanning all tiers applies each rate to its block
    Given electricity tiers [(100, 0.05), (300, 0.10), (inf, 0.15)]
    And consumption of 400 kWh
    When I calculate the electricity cost
    Then the electricity cost should be 40.00

  Scenario: Exact tier boundary applies lower rate to full block
    Given electricity tiers [(100, 0.05), (300, 0.10), (inf, 0.15)]
    And consumption of 100 kWh
    When I calculate the electricity cost
    Then the electricity cost should be 5.00
