package main

import (
	"os"
	"strings"
	"testing"
)

func TestPublishWorkflowIncludesDailyNotificationStep(t *testing.T) {
	t.Parallel()

	content := readWorkflowFile(t, ".github/workflows/publish.yml")

	fragments := []string{
		"id: publish_step",
		"continue-on-error: true",
		`./publish --tags-file tags.json --posts-dir posts | tee publish-output.txt`,
		`exit_code=${PIPESTATUS[0]}`,
		`./notify-summary --posts-dir posts --today "$(date -u +%F)" --publish-exit-code "${{ steps.publish_step.outputs.exit_code }}" --publish-output-file publish-output.txt > notification.txt`,
		`if [ -s notification.txt ]; then`,
		`./notify "$MESSAGE"`,
		"github.event_name == 'schedule'",
		"steps.notification.outputs.has_message == 'true'",
		"TELEGRAM_BOT_TOKEN",
		"TELEGRAM_CHAT_ID",
	}

	for _, fragment := range fragments {
		if !strings.Contains(content, fragment) {
			t.Fatalf("publish workflow missing %q", fragment)
		}
	}
}

func TestReadmeDocumentsDailyPublishNotifications(t *testing.T) {
	t.Parallel()

	content := readWorkflowFile(t, "README.md")

	if !strings.Contains(content, "The scheduled publish run also sends a daily Telegram notification when there is something to report.") {
		t.Fatalf("README does not document the daily publish notification")
	}
}

func readWorkflowFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	return string(content)
}
