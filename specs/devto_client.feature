Feature: Dev.to client
  As the publish pipeline
  I want to interact with the Dev.to REST API
  So that posts are cross-posted, updated, and managed correctly

  Background:
    Given a Dev.to client configured with API key "test-api-key"

  # --- Create ---

  Scenario: Create a new article
    Given no article exists with canonical URL "https://blog.nuphirho.dev/my-new-post"
    When the pipeline creates an article:
      | title     | My New Post                                    |
      | slug      | my-new-post                                    |
      | content   | Hello world.                                   |
      | tags      | go,testing                                     |
      | published | true                                           |
    Then a POST request is sent to "/api/articles"
    And the response contains article ID 12345 and URL "https://dev.to/nuphirho/my-new-post"

  Scenario: Update an existing article when edited_at is newer than the remote article
    Given an article exists with canonical URL "https://blog.nuphirho.dev/existing-post" and ID 67890
    And the article "https://blog.nuphirho.dev/existing-post" has remote edited_at "2026-03-10T09:00:00Z"
    When the pipeline creates an article:
      | title     | Updated Post                                   |
      | slug      | existing-post                                  |
      | content   | Updated content.                               |
      | tags      | go                                             |
      | published | true                                           |
      | edited_at | 2026-03-11T09:00:00Z                           |
    Then a PUT request is sent to "/api/articles/67890"
    And the response contains article URL "https://dev.to/nuphirho/existing-post"

  Scenario: Skip updating an existing article without edited_at
    Given an article exists with canonical URL "https://blog.nuphirho.dev/unchanged-post" and ID 67891
    When the pipeline creates an article:
      | title     | Same Post                                      |
      | slug      | unchanged-post                                 |
      | content   | Same content.                                  |
      | tags      | go                                             |
      | published | true                                           |
    Then no POST or PUT request is made
    And the dry-run result action is "unchanged"

  Scenario: Skip update when remote article has millisecond-precision edited_at timestamp
    Given an article exists with canonical URL "https://blog.nuphirho.dev/stable-post" and ID 50000
    And the article "https://blog.nuphirho.dev/stable-post" has remote edited_at "2026-04-11T05:00:30.123Z"
    When the pipeline creates an article:
      | title     | Stable Post                                    |
      | slug      | stable-post                                    |
      | content   | Content.                                       |
      | tags      | go                                             |
      | published | true                                           |
      | edited_at | 2026-04-10T00:00:00Z                           |
    Then no POST or PUT request is made
    And the dry-run result action is "unchanged"

  Scenario: Article not found by canonical URL falls through to create
    Given no article exists with canonical URL "https://blog.nuphirho.dev/brand-new"
    When the pipeline creates an article:
      | title     | Brand New Post                                 |
      | slug      | brand-new                                      |
      | content   | Content.                                       |
      | tags      | new                                            |
      | published | true                                           |
    Then a POST request is sent to "/api/articles"

  # --- Draft toggle ---

  Scenario: Create article as unpublished draft
    Given no article exists with canonical URL "https://blog.nuphirho.dev/draft-post"
    When the pipeline creates an article:
      | title     | Draft Post                                     |
      | slug      | draft-post                                     |
      | content   | Draft content.                                 |
      | tags      | drafts                                         |
      | published | false                                          |
    Then a POST request is sent to "/api/articles"
    And the request body has "published" set to false

  Scenario: Publish a previously unpublished article
    Given an unpublished article exists with canonical URL "https://blog.nuphirho.dev/was-draft" and ID 11111
    When the pipeline creates an article:
      | title     | Was Draft                                      |
      | slug      | was-draft                                      |
      | content   | Now published.                                 |
      | tags      | publishing                                     |
      | published | true                                           |
    Then a PUT request is sent to "/api/articles/11111"
    And the request body has "published" set to true

  # --- Embed conversion ---

  Scenario: Hashnode tweet embed is converted to Dev.to Liquid tag
    Given no article exists with canonical URL "https://blog.nuphirho.dev/the-trust-barrier"
    When the pipeline creates an article:
      | title     | The Trust Barrier                              |
      | slug      | the-trust-barrier                              |
      | content   | Some text.\n\n%[https://twitter.com/unclebobmartin/status/1962636247769530650]\n\nMore text. |
      | tags      | ai,process                                     |
      | published | false                                          |
    Then the request body content contains "{% embed https://twitter.com/unclebobmartin/status/1962636247769530650 %}"
    And the request body content does not contain "%[https://twitter.com/"

  Scenario: Multiple Hashnode embeds are all converted
    Given no article exists with canonical URL "https://blog.nuphirho.dev/multi-embed"
    When the pipeline creates an article:
      | title     | Multi Embed Post                               |
      | slug      | multi-embed                                    |
      | content   | Intro.\n\n%[https://twitter.com/unclebobmartin/status/2008879916301898134]\n\nMiddle.\n\n%[https://twitter.com/unclebobmartin/status/2019025982863069621]\n\nEnd. |
      | tags      | ai                                             |
      | published | false                                          |
    Then the request body content contains "{% embed https://twitter.com/unclebobmartin/status/2008879916301898134 %}"
    And the request body content contains "{% embed https://twitter.com/unclebobmartin/status/2019025982863069621 %}"
    And the request body content does not contain "%[https://twitter.com/"

  Scenario: Non-embed lines are not affected by conversion
    Given no article exists with canonical URL "https://blog.nuphirho.dev/no-embeds"
    When the pipeline creates an article:
      | title     | No Embeds                                      |
      | slug      | no-embeds                                      |
      | content   | A paragraph with a [link](https://example.com) in it.\n\nAnother paragraph. |
      | tags      | test                                           |
      | published | true                                           |
    Then the request body content contains "A paragraph with a [link](https://example.com) in it."
    And the request body content contains "Another paragraph."

  # --- Tag mapping ---

  Scenario: Tags are included in the request body
    Given no article exists with canonical URL "https://blog.nuphirho.dev/tagged-post"
    When the pipeline creates an article:
      | title     | Tagged Post                                    |
      | slug      | tagged-post                                    |
      | content   | Content.                                       |
      | tags      | go,testing,bdd,security                        |
      | published | true                                           |
    Then the request body tags are ["go", "testing", "bdd", "security"]

  # --- Error handling ---

  Scenario: Invalid API key returns authentication error
    Given the Dev.to API returns a 401 unauthorized error
    When the pipeline creates an article:
      | title     | Fail Post                                      |
      | slug      | fail-post                                      |
      | content   | Content.                                       |
      | tags      | test                                           |
      | published | true                                           |
    Then the error is "authentication failed"

  Scenario: Rate limit response returns rate limit error after retries exhausted
    Given the Dev.to API returns a 429 rate limit error
    When the pipeline creates an article:
      | title     | Rate Limited                                   |
      | slug      | rate-limited                                   |
      | content   | Content.                                       |
      | tags      | test                                           |
      | published | true                                           |
    Then the error is "rate limited"
    And the client retried 3 times

  Scenario: Retry succeeds after transient rate limit
    Given the Dev.to API returns 2 rate limit errors before succeeding
    And no article exists with canonical URL "https://blog.nuphirho.dev/retry-post"
    When the pipeline creates an article:
      | title     | Retry Post                                     |
      | slug      | retry-post                                     |
      | content   | Content.                                       |
      | tags      | test                                           |
      | published | true                                           |
    Then a POST request is sent to "/api/articles"
    And the client retried 2 times

  Scenario: Retry-After header is respected
    Given the Dev.to API returns 1 rate limit error with Retry-After header "5"
    And no article exists with canonical URL "https://blog.nuphirho.dev/retry-after-post"
    When the pipeline creates an article:
      | title     | Retry After Post                               |
      | slug      | retry-after-post                                |
      | content   | Content.                                       |
      | tags      | test                                           |
      | published | true                                           |
    Then a POST request is sent to "/api/articles"
    And the client retried 1 time
    And the first retry delay was 5 seconds

  Scenario: Unhandled API error returns status and body
    Given the Dev.to API returns a 422 error with body '{"error":"Title is too short"}'
    When the pipeline creates an article:
      | title     | Bad Post                                       |
      | slug      | bad-post                                       |
      | content   | Content.                                       |
      | tags      | test                                           |
      | published | true                                           |
    Then the error contains "API error (HTTP 422)"
    And the error contains "Title is too short"

  Scenario: Network failure returns connection error
    Given the Dev.to API is unreachable
    When the pipeline creates an article:
      | title     | Network Fail                                   |
      | slug      | network-fail                                   |
      | content   | Content.                                       |
      | tags      | test                                           |
      | published | true                                           |
    Then the error contains "connection"

  # --- Dry-run ---

  Scenario: Dry-run for new article builds payload without sending
    Given no article exists with canonical URL "https://blog.nuphirho.dev/dry-run-new"
    And dry-run mode is enabled
    When the pipeline creates an article:
      | title     | Dry Run New                                    |
      | slug      | dry-run-new                                    |
      | content   | Content.                                       |
      | tags      | test                                           |
      | published | true                                           |
    Then no POST or PUT request is made
    And the dry-run result action is "create"
    And the dry-run result published is true

  Scenario: Dry-run for existing article builds update payload
    Given an article exists with canonical URL "https://blog.nuphirho.dev/dry-run-update" and ID 99999
    And the article "https://blog.nuphirho.dev/dry-run-update" has remote edited_at "2026-03-10T09:00:00Z"
    And dry-run mode is enabled
    When the pipeline creates an article:
      | title     | Dry Run Update                                 |
      | slug      | dry-run-update                                 |
      | content   | Content.                                       |
      | tags      | test                                           |
      | published | true                                           |
      | edited_at | 2026-03-11T09:00:00Z                           |
    Then no POST or PUT request is made
    And the dry-run result action is "update"
    And the dry-run result existing ID is 99999

  # --- Series ---

  Scenario: Series is included in the request body
    Given no article exists with canonical URL "https://blog.nuphirho.dev/series-post"
    When the pipeline creates an article:
      | title     | Series Post                                    |
      | slug      | series-post                                    |
      | content   | Content.                                       |
      | tags      | go                                             |
      | published | true                                           |
      | series    | Process Over Technology                        |
    Then the request body has "series" set to "Process Over Technology"

  Scenario: Series is omitted when empty
    Given no article exists with canonical URL "https://blog.nuphirho.dev/no-series"
    When the pipeline creates an article:
      | title     | No Series Post                                 |
      | slug      | no-series                                      |
      | content   | Content.                                       |
      | tags      | go                                             |
      | published | true                                           |
    Then the request body does not have "series"

  # --- Dry-run ---

  Scenario: Dry-run reports embed conversion count
    Given no article exists with canonical URL "https://blog.nuphirho.dev/dry-run-embeds"
    And dry-run mode is enabled
    When the pipeline creates an article:
      | title     | Dry Run Embeds                                 |
      | slug      | dry-run-embeds                                 |
      | content   | Text.\n\n%[https://twitter.com/unclebobmartin/status/1962636247769530650]\n\n%[https://twitter.com/unclebobmartin/status/2023158252700066287]\n\nEnd. |
      | tags      | test                                           |
      | published | false                                          |
    Then the dry-run result embeds converted is 2
