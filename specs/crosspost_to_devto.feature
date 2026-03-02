Feature: Cross-post to Dev.to
  As the blog author
  I want published posts to be automatically cross-posted to Dev.to
  So that content reaches the Dev.to developer audience

  Background:
    Given the Dev.to API key is stored in GitHub Secrets

  Scenario: New post is cross-posted to Dev.to
    Given a post has been successfully published to Hashnode
    When the cross-post pipeline runs
    Then the post is created on Dev.to via the REST API
    And the canonical URL is set to "https://nuphirho.dev/<slug>"
    And the post is published (not saved as draft)

  Scenario: Existing post is updated on Dev.to
    Given a post has been updated on Hashnode
    And the post already exists on Dev.to
    When the cross-post pipeline runs
    Then the existing post on Dev.to is updated with the new content
    And the canonical URL remains "https://nuphirho.dev/<slug>"

  Scenario: Draft posts are not cross-posted
    Given a post was skipped on Hashnode because it is a draft
    When the cross-post pipeline runs
    Then no post is created or updated on Dev.to

  Scenario: Cross-post failure does not block the pipeline
    Given a post has been successfully published to Hashnode
    And the Dev.to API returns an error
    When the cross-post pipeline runs
    Then the pipeline reports the Dev.to failure as a warning
    And the overall pipeline does not fail
