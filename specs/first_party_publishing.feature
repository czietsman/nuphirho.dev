Feature: Publish only to the first-party blog
  As the blog owner
  I want posts deployed only to blog.nuphirho.dev
  So that publishing does not depend on external content platforms

  Scenario: A post change deploys the first-party blog
    Given a post in the repository changes
    When the blog deployment workflow runs
    Then the SvelteKit blog is built from the repository posts
    And the blog is deployed to the nuphirho-blog Cloudflare Pages project
    And no external content platform is called

  Scenario: A scheduled post becomes visible
    Given a non-draft post has a publish date of today
    When the daily blog deployment workflow runs at 05:00 UTC
    Then the SvelteKit blog includes the post
    And the blog is deployed to the nuphirho-blog Cloudflare Pages project

  Scenario: A production deployment reports its status
    Given a production blog deployment was attempted
    When the blog deployment workflow finishes
    Then Telegram receives the final job status
    And the message includes the GitHub Actions run URL

  Scenario: A pull-request preview does not send a Telegram notification
    Given a pull-request blog preview was attempted
    When the blog deployment workflow finishes
    Then Telegram is not called
