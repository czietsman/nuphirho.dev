Feature: Notify summary CLI
  As the scheduled publish pipeline
  I want a helper that summarises blog queue and publish outcomes
  So that Telegram messages are only sent when there is something to report

  Scenario: Tomorrow's queued post is included
    Given today's date is "2026-04-02"
    And a post file "posts/tomorrow.md" with:
      | title        | Tomorrow Post |
      | slug         | tomorrow-post |
      | tags         | go            |
      | publish_date | 2026-04-03    |
    When the notify summary CLI runs
    Then stdout contains "Queued for tomorrow: Tomorrow Post"
    And the exit code is 0

  Scenario: Today's successful publish is included
    Given today's date is "2026-04-02"
    And a post file "posts/today.md" with:
      | title        | Today Post |
      | slug         | today-post |
      | tags         | go         |
      | publish_date | 2026-04-02 |
    And the publish step exit code is 0
    When the notify summary CLI runs
    Then stdout contains "Published today: Today Post"
    And the exit code is 0

  Scenario: Today's failed publish is included
    Given today's date is "2026-04-02"
    And a post file "posts/today.md" with:
      | title        | Today Post |
      | slug         | today-post |
      | tags         | go         |
      | publish_date | 2026-04-02 |
    And the publish step exit code is 2
    When the notify summary CLI runs
    Then stdout contains "Publish failed today: Today Post"
    And the exit code is 0

  Scenario: No queue and no publish result produces no message
    Given today's date is "2026-04-02"
    And a post file "posts/later.md" with:
      | title        | Later Post |
      | slug         | later-post |
      | tags         | go         |
      | publish_date | 2026-04-04 |
    When the notify summary CLI runs
    Then stdout is empty
    And the exit code is 0
