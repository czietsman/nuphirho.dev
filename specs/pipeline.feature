Feature: Pipeline orchestrator
  As the publish binary
  I want to validate, publish, and cross-post posts in sequence
  So that the entire pipeline is coordinated from a single entry point

  Background:
    Given a tag glossary:
      """
      {"ai":{"devto":"ai"},"go":{"devto":"go"},"testing":{"devto":"testing"},"bdd":{"devto":"bdd"},"security":{"devto":"security"},"devops":{"devto":"devops"}}
      """

  # --- Happy path ---

  Scenario: Valid post publishes to Hashnode and cross-posts to Dev.to
    Given a post file "posts/my-post.md" with:
      | title | My Post    |
      | slug  | my-post    |
      | tags  | go,testing |
    When the pipeline runs
    Then Hashnode publish is called with slug "my-post"
    And Dev.to cross-post is called with slug "my-post"
    And the summary contains "my-post"
    And the exit code is 0

  # --- Draft ---

  Scenario: Draft post is skipped
    Given a post file "posts/draft-post.md" with:
      | title | Draft Post |
      | slug  | draft-post |
      | tags  | go         |
      | draft | true       |
    When the pipeline runs
    Then Hashnode publish is not called
    And Dev.to cross-post is not called
    And the summary contains "skipped (draft)"
    And the exit code is 0

  # --- Validation ---

  Scenario: Validation failure exits with code 1 and no API calls
    Given a post file "posts/bad-post.md" with:
      | slug | bad-post |
      | tags | go       |
    When the pipeline runs
    Then the exit code is 1
    And Hashnode publish is not called
    And Dev.to cross-post is not called

  # --- Publish failures ---

  Scenario: Hashnode failure exits with code 2 and Dev.to is not attempted
    Given a post file "posts/hn-fail.md" with:
      | title | HN Fail |
      | slug  | hn-fail |
      | tags  | go      |
    And Hashnode publish will fail with "connection refused"
    When the pipeline runs
    Then the exit code is 2
    And Dev.to cross-post is not called

  Scenario: Dev.to failure exits with code 2 and Hashnode result is preserved
    Given a post file "posts/dt-fail.md" with:
      | title | DT Fail |
      | slug  | dt-fail |
      | tags  | go      |
    And Dev.to cross-post will fail with "rate limited"
    When the pipeline runs
    Then the exit code is 2
    And Hashnode publish is called with slug "dt-fail"
    And the summary contains "Hashnode"
    And the summary contains "dt-fail"

  # --- Invalid tags ---

  Scenario: Unmapped hyphenated tag skips Dev.to with warning
    Given a post file "posts/bad-tags.md" with:
      | title | Bad Tags                     |
      | slug  | bad-tags                     |
      | tags  | go,software-engineering      |
    When the pipeline runs
    Then Hashnode publish is called with slug "bad-tags"
    And Dev.to cross-post is not called
    And the summary contains "invalid Dev.to tag"
    And the summary contains "software-engineering"
    And the exit code is 0

  # --- Series ---

  Scenario: Series is resolved and passed to Hashnode and Dev.to
    Given a post file "posts/series-post.md" with:
      | title  | Series Post              |
      | slug   | series-post              |
      | tags   | go                       |
      | series | Process Over Technology  |
    When the pipeline runs
    Then Hashnode publish is called with slug "series-post"
    And the Hashnode call for "series-post" includes series ID "series-resolved"
    And Dev.to cross-post is called with slug "series-post"
    And the Dev.to call for "series-post" includes series "Process Over Technology"
    And the exit code is 0

  Scenario: Post without series does not resolve series
    Given a post file "posts/no-series.md" with:
      | title | No Series |
      | slug  | no-series |
      | tags  | go        |
    When the pipeline runs
    Then Hashnode publish is called with slug "no-series"
    And the Hashnode call for "no-series" does not include series
    And the exit code is 0

  Scenario: Multiple posts in same series resolve once
    Given a post file "posts/series-one.md" with:
      | title  | Series One               |
      | slug   | series-one               |
      | tags   | go                       |
      | series | Process Over Technology  |
    And a post file "posts/series-two.md" with:
      | title  | Series Two               |
      | slug   | series-two               |
      | tags   | go                       |
      | series | Process Over Technology  |
    When the pipeline runs
    Then Hashnode publish is called with slug "series-one"
    And Hashnode publish is called with slug "series-two"
    And the series resolver is called once for "Process Over Technology"
    And the exit code is 0

  Scenario: Series resolution failure exits with code 2
    Given a post file "posts/bad-series.md" with:
      | title  | Bad Series              |
      | slug   | bad-series              |
      | tags   | go                      |
      | series | Failing Series          |
    And series resolution will fail with "series API error"
    When the pipeline runs
    Then the exit code is 2
    And Hashnode publish is not called
    And the summary contains "series"

  # --- Dry-run ---

  Scenario: Dry-run produces JSON with both platform results
    Given a post file "posts/dry-post.md" with:
      | title   | Dry Post   |
      | slug    | dry-post   |
      | tags    | go,testing |
      | content | Text.\n\n%[https://twitter.com/example/status/123]\n\nEnd. |
    And dry-run mode is enabled
    When the pipeline runs
    Then the output is valid JSON
    And the JSON result for "posts/dry-post.md" has hashnode action "publish"
    And the JSON result for "posts/dry-post.md" has devto action "create"
    And the JSON result for "posts/dry-post.md" has devto embeds_converted 1
    And the JSON result for "posts/dry-post.md" has devto published true
    And the exit code is 0

  Scenario: Dry-run for draft post shows skip action
    Given a post file "posts/dry-draft.md" with:
      | title | Dry Draft |
      | slug  | dry-draft |
      | tags  | go        |
      | draft | true      |
    And dry-run mode is enabled
    When the pipeline runs
    Then the output is valid JSON
    And the JSON result for "posts/dry-draft.md" has hashnode action "skip"
    And the JSON result for "posts/dry-draft.md" has devto action "skip"
    And the exit code is 0

  # --- Multiple files ---

  Scenario: Multiple files are each processed independently
    Given a post file "posts/post-one.md" with:
      | title | Post One |
      | slug  | post-one |
      | tags  | go       |
    And a post file "posts/post-two.md" with:
      | title | Post Two |
      | slug  | post-two |
      | tags  | testing  |
    When the pipeline runs
    Then Hashnode publish is called with slug "post-one"
    And Hashnode publish is called with slug "post-two"
    And Dev.to cross-post is called with slug "post-one"
    And Dev.to cross-post is called with slug "post-two"
    And the summary contains "post-one"
    And the summary contains "post-two"
    And the exit code is 0

  Scenario: Second file failure preserves first file result
    Given a post file "posts/good-post.md" with:
      | title | Good Post |
      | slug  | good-post |
      | tags  | go        |
    And a post file "posts/bad-hn.md" with:
      | title | Bad HN |
      | slug  | bad-hn |
      | tags  | go     |
    And Hashnode publish will fail for slug "bad-hn" with "server error"
    When the pipeline runs
    Then Hashnode publish is called with slug "good-post"
    And Dev.to cross-post is called with slug "good-post"
    And the summary contains "good-post"
    And the summary contains "published"
    And the exit code is 2

  # --- Warnings in summary ---

  Scenario: Tag truncation warning appears in summary
    Given a post file "posts/many-tags.md" with:
      | title | Many Tags                         |
      | slug  | many-tags                         |
      | tags  | go,testing,bdd,security,devops,ai |
    When the pipeline runs
    Then the summary contains "tags dropped"
    And the exit code is 0

  Scenario: Em dash warning appears in summary
    Given a post file "posts/em-dash.md" with:
      | title   | Em Dash Post                          |
      | slug    | em-dash                               |
      | tags    | go                                    |
      | content | Some text with an — em dash. |
    When the pipeline runs
    Then the summary contains "em dash"
    And the exit code is 0

  # --- Probe ---

  Scenario: Probe flag calls ProbeAll and exits before publish
    Given a post file "posts/some-post.md" with:
      | title | Some Post |
      | slug  | some-post |
      | tags  | go        |
    When the pipeline runs with --probe
    Then ProbeAll is called
    And Hashnode publish is not called
    And Dev.to cross-post is not called
