Feature: Medium cross-post (manual)
  As the blog author
  I want clear guidance on manually cross-posting to Medium
  So that content reaches the Medium audience without an automated API

  Background:
    Given Medium API tokens are not available for new integrations

  Scenario: Pipeline outputs Medium import instructions
    Given a post has been successfully published to Hashnode
    When the publish pipeline completes
    Then the pipeline output includes the Medium import URL
    And the output includes the instruction to use "Import a story" on Medium
    And the output includes the canonical URL to import

  Scenario: Draft posts do not generate Medium instructions
    Given a post was skipped because it is a draft
    When the publish pipeline completes
    Then no Medium import instructions are included in the output
