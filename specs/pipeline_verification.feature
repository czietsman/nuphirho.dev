Feature: Pipeline verification before deployment
  As the repository owner
  I want every change to an app or its dependencies verified by CI
  So that a dependency update cannot reach main without its tests running

  Scenario: A blog change runs the blog unit tests before building
    Given a pull request changes files under blog/
    When the blog deployment workflow runs
    Then the blog unit tests run after dependencies are installed
    And the blog is built only after the unit tests pass

  Scenario: A dependabot blog pull request does not attempt a preview deploy
    Given dependabot opens a pull request that changes files under blog/
    When the blog deployment workflow runs
    Then the blog unit tests and build run
    And the Cloudflare Pages preview deploy is skipped

  Scenario: A main site pull request is tested and built without deploying
    Given a pull request changes files under main-site/
    When the main site workflow runs
    Then the main site unit tests run after dependencies are installed
    And the main site is built
    And nothing is deployed to Cloudflare Pages

  Scenario: A main site change on main is deployed
    Given a push to main changes files under main-site/
    When the main site workflow runs
    Then the main site unit tests run
    And the main site is deployed to the nuphirho-main Cloudflare Pages project

  Scenario: A Terraform pull request whose plan never ran gets no plan comment
    Given a Terraform pull request whose init step fails
    When the Terraform workflow runs
    Then no plan comment is posted to the pull request
    And the workflow reports the init failure
