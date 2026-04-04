Feature: Content validation
  As the blog author
  I want posts to be validated before publishing
  So that published content meets the style guide standards

  Scenario: Front matter is validated
    Given a markdown file in "posts/" is pushed to the main branch
    Then the front matter must contain "title" as a non-empty string
    And the front matter must contain "slug" as a lowercase hyphenated string
    And the front matter must contain "tags" as a non-empty list
    And the front matter may contain "subtitle" as a string
    And the front matter may contain "draft" as a boolean

  Scenario: No secrets in post content
    Given a markdown file in "posts/" is pushed to the main branch
    When the content validation step runs
    Then the post content is scanned for patterns matching API keys, tokens, and credentials
    And the pipeline fails if any potential secrets are detected

  Scenario: Style guide compliance check
    Given a markdown file in "posts/" is pushed to the main branch
    When the content validation step runs
    Then the post is checked for em dash usage
    And the post is checked for emoji usage
    And em dash violations fail validation
    And emoji violations are reported as warnings

  Scenario: All local posts are reconciled against publishing targets
    Given multiple markdown files exist in "posts/"
    When the publish pipeline runs
    Then all local posts are validated for publication eligibility
    And posts missing from a publishing target are published there
    And unchanged posts are not updated
