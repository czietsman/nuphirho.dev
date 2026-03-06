Feature: Publish to Hashnode
  As the blog author
  I want posts pushed to main to be published on Hashnode
  So that blog.nuphirho.dev always reflects the latest content

  Background:
    Given the Hashnode API token is stored in GitHub Secrets
    And the Hashnode publication ID is stored in GitHub Secrets

  Scenario: New post is published to Hashnode
    Given a markdown file "posts/my-new-post.md" is pushed to the main branch
    And the file contains valid front matter with title, subtitle, tags, and slug
    When the publish pipeline runs
    Then the post is created on Hashnode via the GraphQL API
    And the canonical URL is set to "https://blog.nuphirho.dev/<slug>"
    And the pipeline reports success with the published URL

  Scenario: Existing post is updated on Hashnode
    Given a markdown file "posts/my-existing-post.md" is modified on the main branch
    And the post already exists on Hashnode with a matching slug
    When the publish pipeline runs
    Then the existing post on Hashnode is updated with the new content
    And the canonical URL remains "https://blog.nuphirho.dev/<slug>"

  Scenario: Post with missing front matter is rejected
    Given a markdown file "posts/bad-post.md" is pushed to the main branch
    And the file is missing required front matter fields
    When the publish pipeline runs
    Then the pipeline fails with a clear error indicating the missing fields
    And no post is created or updated on Hashnode

  Scenario: Draft post is created as unpublished on Hashnode
    Given a markdown file "posts/draft-post.md" is pushed to the main branch
    And the front matter contains "draft: true"
    When the publish pipeline runs
    Then the post is created on Hashnode as a draft via the GraphQL API
    And the post is not publicly visible on Hashnode
    And the pipeline reports the post was pushed as a draft

  Scenario: Draft post is published when draft flag is removed
    Given a markdown file "posts/draft-post.md" is modified on the main branch
    And the front matter no longer contains "draft: true"
    And the post was previously pushed as a draft to Hashnode
    When the publish pipeline runs
    Then the post is published on Hashnode with the updated content
    And the canonical URL is set to "https://blog.nuphirho.dev/<slug>"
    And the pipeline reports success with the published URL
