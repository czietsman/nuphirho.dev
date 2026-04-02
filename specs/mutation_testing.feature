Feature: Mutation testing workflow
  As the validation pipeline
  I want mutation testing to run with caching
  So that high-consequence packages are checked in CI without avoidable cold starts

  Scenario: Validate PR installs the pinned go-mutesting tool
    Given the validation workflow configuration
    Then the validation workflow installs "github.com/avito-tech/go-mutesting/cmd/go-mutesting@v0.0.0-20251226130216-48d0401f00fb"

  Scenario: Validate PR enables Go caching for mutation testing
    Given the validation workflow configuration
    Then the validation workflow enables Go dependency caching

  Scenario: Validate PR runs mutation testing against frontmatter
    Given the validation workflow configuration
    Then the validation workflow runs mutation testing for "./internal/frontmatter"
