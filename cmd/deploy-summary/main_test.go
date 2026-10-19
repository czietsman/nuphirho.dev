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

type summaryContext struct {
	dir    string
	today  string
	stdout bytes.Buffer
	stderr bytes.Buffer
	code   int
}

func (sc *summaryContext) thesePosts(table *godog.Table) error {
	for _, row := range table.Rows[1:] {
		title, slug, date, draft := row.Cells[0].Value, row.Cells[1].Value, row.Cells[2].Value, row.Cells[3].Value
		fm := fmt.Sprintf("---\ntitle: %q\nslug: %q\ndraft: %s\n", title, slug, draft)
		if date != "" {
			fm += "publish_date: " + date + "\n"
		}
		fm += "---\n\nBody.\n"
		if err := os.WriteFile(filepath.Join(sc.dir, slug+".md"), []byte(fm), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func (sc *summaryContext) todayIs(date string) error {
	sc.today = date
	return nil
}

func (sc *summaryContext) run(args ...string) {
	sc.code = Run(append(args, "-today", sc.today, sc.dir), &sc.stdout, &sc.stderr)
}

func (sc *summaryContext) builtForStatus(status string) error {
	sc.run("-status", status)
	return nil
}

func (sc *summaryContext) builtForStatusWithRunURL(status, runURL string) error {
	sc.run("-status", status, "-run-url", runURL)
	return nil
}

func (sc *summaryContext) theSummaryIs(doc *godog.DocString) error {
	if sc.code != 0 {
		return fmt.Errorf("exit code %d, stderr: %s", sc.code, sc.stderr.String())
	}
	if got := strings.TrimRight(sc.stdout.String(), "\n"); got != doc.Content {
		return fmt.Errorf("summary mismatch\n--- got ---\n%s\n--- want ---\n%s", got, doc.Content)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	sc := &summaryContext{}
	ctx.Before(func(c context.Context, s *godog.Scenario) (context.Context, error) {
		dir, err := os.MkdirTemp("", "deploy-summary-*")
		*sc = summaryContext{dir: dir}
		return c, err
	})
	ctx.After(func(c context.Context, s *godog.Scenario, err error) (context.Context, error) {
		return c, os.RemoveAll(sc.dir)
	})
	ctx.Step(`^these posts:$`, sc.thesePosts)
	ctx.Step(`^today is "([^"]*)"$`, sc.todayIs)
	ctx.Step(`^the deployment summary is built for status "([^"]*)"$`, sc.builtForStatus)
	ctx.Step(`^the deployment summary is built for status "([^"]*)" with run URL "([^"]*)"$`, sc.builtForStatusWithRunURL)
	ctx.Step(`^the summary is:$`, sc.theSummaryIs)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../specs/deploy_summary.feature"},
			Strict:   true,
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func TestMissingDirectoryIsAnError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"-status", "success"}, &stdout, &stderr); code != 1 {
		t.Fatalf("exit code %d, want 1", code)
	}
}
