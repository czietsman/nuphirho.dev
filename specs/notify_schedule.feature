Feature: Daily publish notifications
  As the scheduled publish pipeline
  I want a daily Telegram notification only when there is blog activity to report
  So that I know what is queued for tomorrow and whether today's publish succeeded

  Scenario: Publish workflow sends a daily notification on the scheduled run
    Given the publish workflow configuration
    Then the publish workflow prepares a daily notification during scheduled runs

  Scenario: Publish workflow sends no message when there is nothing to report
    Given the publish workflow configuration
    Then the publish workflow only sends a notification when a message was generated

  Scenario: Daily notification can report queued posts for tomorrow
    Given the notify summary configuration
    Then queued posts for tomorrow are included in the daily notification

  Scenario: Daily notification can report today's successful publishes
    Given the notify summary configuration
    Then today's successful publishes are included in the daily notification

  Scenario: Daily notification can report today's publish failures
    Given the notify summary configuration
    Then today's publish failures are included in the daily notification
