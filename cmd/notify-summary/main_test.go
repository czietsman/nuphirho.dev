package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

type notifySummaryCtx struct {
	tmpDir        string
	today         string
	exitCode      int
	stdout        string
	stderr        string
	publishCode   int
	publishOutput string
}

func (nc *notifySummaryCtx) reset() {
	tmpDir, err := os.MkdirTemp("", "notify-summary-*")
	if err != nil {
		panic(err)
	}
	nc.tmpDir = tmpDir
	nc.today = "2026-04-02"
	nc.exitCode = -1
	nc.stdout = ""
	nc.stderr = ""
	nc.publishCode = 0
	nc.publishOutput = ""
}

func (nc *notifySummaryCtx) cleanup() {
	if nc.tmpDir != "" {
		os.RemoveAll(nc.tmpDir)
	}
}

func (nc *notifySummaryCtx) aPostFileWith(path string, table *godog.Table) error {
	fields := tableToMap(table)
	content := buildMarkdown(fields)
	absPath := filepath.Join(nc.tmpDir, path)
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(absPath, []byte(content), 0644)
}

func (nc *notifySummaryCtx) todaysDateIs(date string) error {
	nc.today = date
	return nil
}

func (nc *notifySummaryCtx) thePublishStepExitCodeIs(code int) error {
	nc.publishCode = code
	return nil
}

func (nc *notifySummaryCtx) thePublishOutputIs(output string) error {
	nc.publishOutput = strings.ReplaceAll(output, `\n`, "\n")
	return nil
}

func (nc *notifySummaryCtx) theNotifySummaryCLIRuns() error {
	publishOutputPath := filepath.Join(nc.tmpDir, "publish-output.txt")
	if nc.publishOutput != "" {
		if err := os.WriteFile(publishOutputPath, []byte(nc.publishOutput), 0644); err != nil {
			return err
		}
	}

	args := []string{
		"--posts-dir", filepath.Join(nc.tmpDir, "posts"),
		"--today", nc.today,
		"--publish-exit-code", intToString(nc.publishCode),
	}
	if nc.publishOutput != "" {
		args = append(args, "--publish-output-file", publishOutputPath)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	nc.exitCode = Run(args, &stdoutBuf, &stderrBuf)
	nc.stdout = stdoutBuf.String()
	nc.stderr = stderrBuf.String()
	return nil
}

func (nc *notifySummaryCtx) stdoutContains(expected string) error {
	if !strings.Contains(nc.stdout, expected) {
		return fmt.Errorf("expected stdout to contain %q, got %q", expected, nc.stdout)
	}
	return nil
}

func (nc *notifySummaryCtx) stdoutIsEmpty() error {
	if strings.TrimSpace(nc.stdout) != "" {
		return fmt.Errorf("expected stdout to be empty, got %q", nc.stdout)
	}
	return nil
}

func (nc *notifySummaryCtx) theExitCodeIs(code int) error {
	if nc.exitCode != code {
		return fmt.Errorf("expected exit code %d, got %d\nstdout: %s\nstderr: %s", code, nc.exitCode, nc.stdout, nc.stderr)
	}
	return nil
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	if n == 2 {
		return "2"
	}
	return "1"
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	nc := &notifySummaryCtx{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		nc.reset()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		nc.cleanup()
		return ctx, nil
	})

	ctx.Step(`^a post file "([^"]*)" with:$`, nc.aPostFileWith)
	ctx.Step(`^today's date is "([^"]*)"$`, nc.todaysDateIs)
	ctx.Step(`^the publish step exit code is (\d+)$`, nc.thePublishStepExitCodeIs)
	ctx.Step(`^the publish output is:$`, nc.thePublishOutputIs)
	ctx.Step(`^the notify summary CLI runs$`, nc.theNotifySummaryCLIRuns)
	ctx.Step(`^stdout contains "([^"]*)"$`, nc.stdoutContains)
	ctx.Step(`^stdout is empty$`, nc.stdoutIsEmpty)
	ctx.Step(`^the exit code is (\d+)$`, nc.theExitCodeIs)
}

func TestNotifySummaryFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../specs/notify_summary.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func tableToMap(table *godog.Table) map[string]string {
	m := make(map[string]string)
	for _, row := range table.Rows {
		key := strings.TrimSpace(row.Cells[0].Value)
		val := strings.TrimSpace(row.Cells[1].Value)
		m[key] = val
	}
	return m
}

func buildMarkdown(fields map[string]string) string {
	var sb strings.Builder
	sb.WriteString("---\n")
	if title, ok := fields["title"]; ok && title != "" {
		sb.WriteString(fmt.Sprintf("title: %s\n", title))
	}
	if slug, ok := fields["slug"]; ok && slug != "" {
		sb.WriteString(fmt.Sprintf("slug: %s\n", slug))
	}
	if t, ok := fields["tags"]; ok && t != "" {
		sb.WriteString(fmt.Sprintf("tags: [%s]\n", t))
	}
	if draft, ok := fields["draft"]; ok && draft == "true" {
		sb.WriteString("draft: true\n")
	}
	if publishDate, ok := fields["publish_date"]; ok && publishDate != "" {
		sb.WriteString(fmt.Sprintf("publish_date: %s\n", publishDate))
	}
	sb.WriteString("---\n")
	sb.WriteString("Test content.")
	return sb.String()
}
