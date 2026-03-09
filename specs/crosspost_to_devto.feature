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
    And the canonical URL is set to "https://blog.nuphirho.dev/<slug>"
    And the post is published (not saved as draft)

  Scenario: Existing post is updated on Dev.to
    Given a post has been updated on Hashnode
    And the post already exists on Dev.to
    When the cross-post pipeline runs
    Then the existing post on Dev.to is updated with the new content
    And the canonical URL remains "https://blog.nuphirho.dev/<slug>"

  Scenario: Draft post is created as unpublished on Dev.to
    Given a post has been created on Hashnode as a draft
    And the front matter contains "draft: true"
    When the cross-post pipeline runs
    Then the post is created on Dev.to via the REST API with "published: false"
    And the canonical URL is set to "https://blog.nuphirho.dev/<slug>"

  Scenario: Draft post is published when draft flag is removed
    Given a post has been published on Hashnode after removing the draft flag
    And the post already exists on Dev.to as unpublished
    When the cross-post pipeline runs
    Then the existing post on Dev.to is updated with "published: true"
    And the canonical URL remains "https://blog.nuphirho.dev/<slug>"

  Scenario: Hashnode embeds are converted to Dev.to Liquid tags
    Given a post contains Hashnode embed syntax "%[https://twitter.com/user/status/123]"
    When the cross-post pipeline runs
    Then the embed is converted to Dev.to Liquid tag syntax "{% embed https://twitter.com/user/status/123 %}"
    And no Hashnode embed syntax remains in the Dev.to article body

  Scenario: Cross-post failure does not block the pipeline
    Given a post has been successfully published to Hashnode
    And the Dev.to API returns an error
    When the cross-post pipeline runs
    Then the pipeline reports the Dev.to failure as a warning
    And the overall pipeline does not fail
