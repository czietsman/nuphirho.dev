Feature: Notification CLI
  As the notification command
  I want to send Telegram messages from the terminal or GitHub Actions
  So that publishing reminders reach Christo's phone without opening LinkedIn

  Scenario: Notification is sent with credentials from environment variables
    Given Telegram credentials via environment variables
    And a notification message "Post 4 is live. Monitor engagement."
    When the notification CLI runs
    Then the Telegram sender is called with chat ID "123456"
    And the Telegram sender is called with message "Post 4 is live. Monitor engagement."
    And the exit code is 0

  Scenario: Missing message exits with code 1
    Given Telegram credentials via environment variables
    When the notification CLI runs
    Then the Telegram sender is not called
    And stderr contains "error: notification message is required"
    And the exit code is 1

  Scenario: Missing Telegram token exits with code 1
    Given TELEGRAM_CHAT_ID is set to "123456"
    And a notification message "Post 4 is live. Monitor engagement."
    When the notification CLI runs
    Then the Telegram sender is not called
    And stderr contains "error: Telegram bot token is required"
    And the exit code is 1

  Scenario: Sender failure exits with code 2
    Given Telegram credentials via environment variables
    And a notification message "Post 4 is live. Monitor engagement."
    And the Telegram sender will fail with "telegram API error"
    When the notification CLI runs
    Then stderr contains "error: telegram API error"
    And the exit code is 2
