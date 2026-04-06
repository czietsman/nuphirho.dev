package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

// --- Fake HTTP client ---

type fakeResponse struct {
	status int
	body   string
}

type fakeClient struct {
	responses map[string]fakeResponse
}

func (f *fakeClient) Do(req *http.Request) (*http.Response, error) {
	url := req.URL.String()
	if r, ok := f.responses[url]; ok {
		return &http.Response{
			StatusCode: r.status,
			Body:       io.NopCloser(strings.NewReader(r.body)),
		}, nil
	}
	return &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(strings.NewReader("not found")),
	}, nil
}

// --- Test context ---

type testContext struct {
	version   string
	stdout    bytes.Buffer
	stderr    bytes.Buffer
	exitCode  int
	tmpDir    string
	client    *fakeClient
}

func newTestContext() *testContext {
	tc := &testContext{
		client: &fakeClient{
			responses: make(map[string]fakeResponse),
		},
	}
	return tc
}

func (tc *testContext) setupFakeResponses() {
	sha := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
	for _, repo := range repos {
		// Content response for main branch, AGENTS.md path.
		url := fmt.Sprintf("https://raw.githubusercontent.com/%s/main/AGENTS.md", repo)
		tc.client.responses[url] = fakeResponse{
			status: 200,
			body:   fmt.Sprintf("# AGENTS.md for %s\n\nThis is governance content.\nLine 3.\nLine 4.\n", repo),
		}
		// Commit SHA response.
		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/commits?path=AGENTS.md&per_page=1", repo)
		tc.client.responses[apiURL] = fakeResponse{
			status: 200,
			body:   fmt.Sprintf(`[{"sha":"%s"}]`, sha),
		}
	}
}

// --- Step definitions ---

func (tc *testContext) githubTokenIsSet() error {
	// Token is provided via getenv in the Run call.
	return nil
}

func (tc *testContext) theCorpusVersionIs(version string) error {
	tc.version = version
	return nil
}

func (tc *testContext) theCaptureToolRuns() error {
	tc.tmpDir, _ = os.MkdirTemp("", "corpus-test-*")
	origDir, _ := os.Getwd()
	os.Chdir(tc.tmpDir)
	defer os.Chdir(origDir)

	tc.setupFakeResponses()
	getenv := func(key string) string {
		if key == "GITHUB_TOKEN" {
			return "fake-token"
		}
		return ""
	}
	tc.exitCode = Run([]string{tc.version}, &tc.stdout, &tc.stderr, getenv, tc.client)
	return nil
}

func (tc *testContext) agentFilesExist(count int, dir string) error {
	fullDir := filepath.Join(tc.tmpDir, dir)
	entries, err := os.ReadDir(fullDir)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", fullDir, err)
	}
	mdFiles := 0
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".md") {
			mdFiles++
		}
	}
	if mdFiles != count {
		return fmt.Errorf("expected %d agent files, got %d", count, mdFiles)
	}
	return nil
}

func (tc *testContext) theManifestContainsRows(count int) error {
	path := filepath.Join(tc.tmpDir, "experiments", fmt.Sprintf("governance-prompts-%s", tc.version), "manifest.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	rows := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "| ") && !strings.HasPrefix(line, "| #") && !strings.HasPrefix(line, "|---") {
			rows++
		}
	}
	if rows != count {
		return fmt.Errorf("expected %d manifest rows, got %d", count, rows)
	}
	return nil
}

func (tc *testContext) noRowHasStatus(status string) error {
	path := filepath.Join(tc.tmpDir, "experiments", fmt.Sprintf("governance-prompts-%s", tc.version), "manifest.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if strings.Contains(string(data), status) {
		return fmt.Errorf("manifest contains status %q", status)
	}
	return nil
}

func (tc *testContext) everyAgentFileHasCommitSHA() error {
	return tc.walkAgentFiles(func(content string) error {
		if !strings.Contains(content, "commit: ") {
			return fmt.Errorf("agent file missing commit SHA")
		}
		return nil
	})
}

func (tc *testContext) noCommitSHAEquals(val string) error {
	return tc.walkAgentFiles(func(content string) error {
		re := regexp.MustCompile(`commit: (.+)`)
		m := re.FindStringSubmatch(content)
		if len(m) > 1 && strings.TrimSpace(m[1]) == val {
			return fmt.Errorf("found commit SHA %q", val)
		}
		return nil
	})
}

func (tc *testContext) everySHAIsExactly40Chars() error {
	return tc.walkAgentFiles(func(content string) error {
		re := regexp.MustCompile(`commit: (.+)`)
		m := re.FindStringSubmatch(content)
		if len(m) > 1 {
			sha := strings.TrimSpace(m[1])
			if len(sha) != 40 {
				return fmt.Errorf("SHA length is %d, expected 40: %s", len(sha), sha)
			}
		}
		return nil
	})
}

func (tc *testContext) everyAgentFileHasContentBeyondFrontMatter() error {
	return tc.walkAgentFiles(func(content string) error {
		parts := strings.SplitN(content, "---\n", 3)
		if len(parts) < 3 || strings.TrimSpace(parts[2]) == "" {
			return fmt.Errorf("agent file has no content beyond front matter")
		}
		return nil
	})
}

func (tc *testContext) lineCountInFrontMatterMatchesActual() error {
	return tc.walkAgentFiles(func(content string) error {
		re := regexp.MustCompile(`lines: (\d+)`)
		m := re.FindStringSubmatch(content)
		if len(m) < 2 {
			return fmt.Errorf("no lines field in front matter")
		}
		expected, _ := strconv.Atoi(m[1])
		parts := strings.SplitN(content, "---\n", 3)
		if len(parts) < 3 {
			return fmt.Errorf("cannot split front matter")
		}
		actual := countLines(parts[2])
		if actual != expected {
			return fmt.Errorf("front matter says %d lines, actual is %d", expected, actual)
		}
		return nil
	})
}

func (tc *testContext) manifestExistsAt(path string) error {
	full := filepath.Join(tc.tmpDir, path)
	if _, err := os.Stat(full); err != nil {
		return fmt.Errorf("manifest not found at %s", full)
	}
	return nil
}

func (tc *testContext) manifestContainsHeader() error {
	path := filepath.Join(tc.tmpDir, "experiments", fmt.Sprintf("governance-prompts-%s", tc.version), "manifest.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(data), "Retrieval date:") {
		return fmt.Errorf("manifest missing retrieval date header")
	}
	return nil
}

func (tc *testContext) manifestTableHasColumns(cols string) error {
	path := filepath.Join(tc.tmpDir, "experiments", fmt.Sprintf("governance-prompts-%s", tc.version), "manifest.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)
	for _, col := range strings.Split(cols, ",") {
		col = strings.TrimSpace(col)
		if !strings.Contains(content, col) {
			return fmt.Errorf("manifest missing column %q", col)
		}
	}
	return nil
}

func (tc *testContext) readmeExistsAt(path string) error {
	full := filepath.Join(tc.tmpDir, path)
	if _, err := os.Stat(full); err != nil {
		return fmt.Errorf("README not found at %s", full)
	}
	return nil
}

func (tc *testContext) governancePromptsExistsAt(path string) error {
	full := filepath.Join(tc.tmpDir, path)
	if _, err := os.Stat(full); err != nil {
		return fmt.Errorf("GOVERNANCE_PROMPTS.md not found at %s", full)
	}
	return nil
}

func (tc *testContext) walkAgentFiles(fn func(string) error) error {
	dir := filepath.Join(tc.tmpDir, "experiments", fmt.Sprintf("governance-prompts-%s", tc.version), "agents")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return err
		}
		if err := fn(string(data)); err != nil {
			return fmt.Errorf("%s: %w", e.Name(), err)
		}
	}
	return nil
}

// --- Godog integration ---

func InitializeScenario(ctx *godog.ScenarioContext) {
	tc := newTestContext()

	ctx.Step(`^GITHUB_TOKEN is set in the environment$`, tc.githubTokenIsSet)
	ctx.Step(`^the corpus version is "([^"]*)"$`, tc.theCorpusVersionIs)
	ctx.Step(`^the capture tool runs$`, tc.theCaptureToolRuns)
	ctx.Step(`^(\d+) agent files exist in ([^\s]+)$`, tc.agentFilesExist)
	ctx.Step(`^the manifest contains (\d+) rows$`, tc.theManifestContainsRows)
	ctx.Step(`^no row in the manifest has status "([^"]*)"$`, tc.noRowHasStatus)
	ctx.Step(`^every agent file has a commit SHA in its front matter$`, tc.everyAgentFileHasCommitSHA)
	ctx.Step(`^no commit SHA equals "([^"]*)"$`, tc.noCommitSHAEquals)
	ctx.Step(`^every SHA is exactly 40 characters$`, tc.everySHAIsExactly40Chars)
	ctx.Step(`^every agent file has content beyond the front matter$`, tc.everyAgentFileHasContentBeyondFrontMatter)
	ctx.Step(`^the line count in the front matter matches the actual file content$`, tc.lineCountInFrontMatterMatchesActual)
	ctx.Step(`^manifest\.md exists at ([^\s]+)$`, tc.manifestExistsAt)
	ctx.Step(`^the manifest contains a header with retrieval date$`, tc.manifestContainsHeader)
	ctx.Step(`^the manifest table has columns: (.+)$`, tc.manifestTableHasColumns)
	ctx.Step(`^README\.md exists at ([^\s]+)$`, tc.readmeExistsAt)
	ctx.Step(`^GOVERNANCE_PROMPTS\.md exists at ([^\s]+)$`, tc.governancePromptsExistsAt)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if tc.tmpDir != "" {
			os.RemoveAll(tc.tmpDir)
		}
		return ctx, nil
	})
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../specs/corpus_capture.feature"},
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("BDD scenarios failed")
	}
}
