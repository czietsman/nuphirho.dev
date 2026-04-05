Feature: Tag validation at PR time

  Background:
    Given a tag glossary:
      """
      {
        "ai-governance": { "devto": "aigovernance" },
        "software-verification": { "devto": "softwareverification" }
      }
      """

  Scenario: Post with all tags mapped passes validation
    Given a post with tags ["ai-governance", "epistemology"]
    When tag validation runs
    Then validation passes

  Scenario: Post with unmapped hyphenated tag fails validation
    Given a post with tags ["new-tag-with-hyphens"]
    When tag validation runs
    Then validation fails
    And the error identifies "new-tag-with-hyphens" as unmapped

  Scenario: Post with tag containing no hyphens needs no mapping
    Given a post with tags ["epistemology"]
    When tag validation runs
    Then validation passes

  Scenario: Draft posts are not validated
    Given a draft post with tags ["new-tag-with-hyphens"]
    When tag validation runs
    Then validation is skipped

  Scenario: Multiple unmapped tags are all reported
    Given a post with tags ["bad-tag-one", "bad-tag-two", "goodtag"]
    When tag validation runs
    Then validation fails
    And the error identifies "bad-tag-one" as unmapped
    And the error identifies "bad-tag-two" as unmapped
