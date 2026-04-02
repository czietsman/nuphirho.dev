Feature: Repository validation
  As the repository maintainer
  I want repository validation rules to be explicit and enforceable
  So that tests, workflows, and documentation remain traceable and trustworthy

  Scenario: Repository tests declare BDD traceability annotations
    Given the repository contains Go test functions
    Then each test function declares the BDD feature or scenario that backs it

  Scenario: Repository test annotations reference existing BDD specifications
    Given a Go test function declares a BDD feature or scenario reference
    Then the referenced feature file exists
    And the referenced feature or scenario exists in that file

  Scenario: PR validation enforces BDD traceability checks
    Given the pull request validation workflow runs
    Then it runs the repository BDD traceability checks

  Scenario: README documents BDD traceability enforcement
    Given repository validation enforces BDD traceability checks
    Then the README documents the BDD traceability rule

  Scenario: Workflow action pins are checked in repository validation
    Given repository validation runs
    Then it checks that workflow action revisions remain pinned

  Scenario: Terraform workflow comments the real plan output
    Given the Terraform workflow comments on a pull request
    Then it reads the generated plan output from plan.txt

  Scenario: README describes draft posts as skipped
    Given the publishing pipeline skips draft posts
    Then the README describes drafts as skipped

  Scenario: PR validation runs mutation testing
    Given pull request validation runs
    Then it runs mutation testing for internal/frontmatter

  Scenario: README documents mutation testing in validation
    Given pull request validation includes mutation testing
    Then the README documents that mutation testing runs in validation

  Scenario: The repository Go version supports the pinned mutation tool
    Given the validation workflow installs a pinned mutation tool
    Then go.mod declares a compatible Go version

  Scenario: The publish workflow sends daily notifications only when needed
    Given the scheduled publish workflow runs
    Then it prepares and sends a Telegram notification only when there is something to report

  Scenario: README documents daily publish notifications
    Given the scheduled publish workflow sends daily notifications
    Then the README documents the daily notification behaviour

  Scenario: Dependency review prompt templates are present and documented
    Given the repository contains dependency review prompt templates
    Then the prompt template files exist
    And the README documents them
