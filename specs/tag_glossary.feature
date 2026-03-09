Feature: Tag glossary
  As the publish pipeline
  I want to map canonical tags to platform-specific equivalents
  So that posts use valid tags on each platform

  Background:
    Given a tag glossary:
      """
      {
        "_comment": "Test glossary",
        "infrastructure-as-code": { "devto": "iac" },
        "spatial-thinking": { "devto": "spatialthinking" },
        "ai-assisted-development": { "devto": "aidev" },
        "threat-modeling": { "devto": "threatmodeling" }
      }
      """

  Scenario: Tag in glossary is mapped to platform-specific value
    Given canonical tags ["infrastructure-as-code", "go"]
    When tags are mapped for "devto"
    Then the mapped tags are ["iac", "go"]

  Scenario: Tag not in glossary passes through unchanged
    Given canonical tags ["go", "testing"]
    When tags are mapped for "devto"
    Then the mapped tags are ["go", "testing"]
    And no warnings are produced

  Scenario: Multiple tags are mapped correctly
    Given canonical tags ["spatial-thinking", "ai-assisted-development", "career"]
    When tags are mapped for "devto"
    Then the mapped tags are ["spatialthinking", "aidev", "career"]

  Scenario: Dev.to tags are limited to 4 with warning
    Given canonical tags ["go", "testing", "bdd", "security", "devops"]
    When tags are mapped for "devto"
    Then the mapped tags are ["go", "testing", "bdd", "security"]
    And a warning is produced: "1 tags dropped for Dev.to: devops"

  Scenario: Multiple Dev.to tags dropped produces warning listing all
    Given canonical tags ["go", "testing", "bdd", "security", "devops", "ai"]
    When tags are mapped for "devto"
    Then the mapped tags are ["go", "testing", "bdd", "security"]
    And a warning is produced: "2 tags dropped for Dev.to: devops, ai"

  Scenario: Dev.to tags with hyphens are rejected without glossary mapping
    Given canonical tags ["threat-modeling", "engineering-process"]
    When tags are validated for "devto"
    Then "threat-modeling" is valid because it has a glossary mapping
    And "engineering-process" is invalid because it contains hyphens

  Scenario: Hashnode tags pass through without restriction
    Given canonical tags ["infrastructure-as-code", "go", "testing"]
    When tags are mapped for "hashnode"
    Then the mapped tags are ["infrastructure-as-code", "go", "testing"]

  Scenario: Hashnode does not limit tag count
    Given canonical tags ["go", "testing", "bdd", "security", "devops", "ai"]
    When tags are mapped for "hashnode"
    Then the mapped tags are ["go", "testing", "bdd", "security", "devops", "ai"]

  Scenario: Comment field in glossary is ignored
    Given canonical tags ["_comment"]
    When tags are mapped for "devto"
    Then the mapped tags are ["_comment"]

  Scenario: Empty tag list produces empty result
    Given canonical tags []
    When tags are mapped for "devto"
    Then the mapped tags are []

  Scenario: Glossary with no mapping for platform passes tag through
    Given a tag glossary:
      """
      {
        "special-tag": { "hashnode": "special" }
      }
      """
    And canonical tags ["special-tag"]
    When tags are mapped for "devto"
    Then the mapped tags are ["special-tag"]

  Scenario: Dev.to tag validation rejects hyphens
    Given canonical tags ["valid", "has-hyphen"]
    When tags are validated for "devto"
    Then "valid" is valid
    And "has-hyphen" is invalid because it contains hyphens

  Scenario: Dev.to tag validation accepts alphanumeric only
    Given canonical tags ["go", "bdd", "ai", "testing123"]
    When tags are validated for "devto"
    Then all tags are valid
