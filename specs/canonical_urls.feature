Feature: Canonical URL management
  As the blog author
  I want all cross-posted content to reference the canonical URL on nuphirho.dev
  So that SEO authority builds on the primary domain

  Scenario: Canonical URL is set on Hashnode publication
    Given a post with slug "setting-up-the-blog" is published to Hashnode
    Then the canonical URL is "https://nuphirho.dev/setting-up-the-blog"

  Scenario: Canonical URL is set on Dev.to cross-post
    Given a post with slug "setting-up-the-blog" is cross-posted to Dev.to
    Then the canonical URL is "https://nuphirho.dev/setting-up-the-blog"

  Scenario: Canonical URL format is validated before publishing
    Given a post front matter contains slug "Setting Up The Blog"
    When the publish pipeline validates the front matter
    Then the pipeline fails with an error indicating the slug must be lowercase and hyphenated
