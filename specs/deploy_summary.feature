Feature: Blog deployment summary
  As the blog owner
  I want the deployment notification to say which posts went live and which are queued
  So that I know what readers can see without opening the pipeline run

  Background:
    Given these posts:
      | title                | slug               | publish_date | draft |
      | Already out          | already-out        | 2026-10-06   | false |
      | Out today            | out-today          | 2026-10-12   | false |
      | Next Friday          | next-friday        | 2026-10-16   | false |
      | Monday after         | monday-after       | 2026-10-19   | false |
      | Still a draft        | still-a-draft      | 2026-10-16   | true  |
      | Not scheduled        | not-scheduled      |              | false |
    And today is "2026-10-12"

  Scenario: A successful deployment lists today's posts and the queue
    When the deployment summary is built for status "success"
    Then the summary is:
      """
      Blog deployed.

      Live today:
      - Out today: https://blog.nuphirho.dev/out-today

      Queued:
      - Fri 16 Oct: Next Friday
      - Mon 19 Oct: Monday after
      """

  Scenario: A day with nothing going live still shows the queue
    Given today is "2026-10-13"
    When the deployment summary is built for status "success"
    Then the summary is:
      """
      Blog deployed. Nothing new went live today.

      Queued:
      - Fri 16 Oct: Next Friday
      - Mon 19 Oct: Monday after
      """

  Scenario: An empty queue is said plainly
    Given today is "2026-10-19"
    When the deployment summary is built for status "success"
    Then the summary is:
      """
      Blog deployed.

      Live today:
      - Monday after: https://blog.nuphirho.dev/monday-after

      Nothing is queued.
      """

  Scenario: A failed deployment links the run so it can be investigated
    When the deployment summary is built for status "failure" with run URL "https://github.com/run/1"
    Then the summary is:
      """
      Blog deployment failure: https://github.com/run/1

      Queued:
      - Fri 16 Oct: Next Friday
      - Mon 19 Oct: Monday after
      """
