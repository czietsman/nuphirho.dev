package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/czietsman/nuphirho.dev/internal/devto"
	"github.com/czietsman/nuphirho.dev/internal/hashnode"
)

// --- Fakes ---

type fakeHN struct {
	publishCalls []string
	publishErr   error
}

func (f *fakeHN) Publish(input hashnode.PostInput) (*hashnode.PublishResult, error) {
	f.publishCalls = append(f.publishCalls, input.Slug)
	if f.publishErr != nil {
		return nil, f.publishErr
	}
	return &hashnode.PublishResult{
		Action: "publish",
		PostID: "post-001",
		URL:    "https://blog.nuphirho.dev/" + input.Slug,
	}, nil
}

type fakeDT struct {
	createCalls []string
	createErr   error
}

func (f *fakeDT) CreateArticle(input devto.ArticleInput) (*devto.PublishResult, error) {
	f.createCalls = append(f.createCalls, input.Slug)
	if f.createErr != nil {
		return nil, f.createErr
	}
	return &devto.PublishResult{
		Action:    "create",
		ArticleID: 12345,
		URL:       "https://dev.to/nuphirho/" + input.Slug,
		Published: input.Published,
	}, nil
}

type fakeProber struct {
	called   bool
	exitCode int
}

func (f *fakeProber) Run(w io.Writer) int {
	f.called = true
	fmt.Fprintln(w, "All probes passed.")
	return f.exitCode
}

// --- godog context ---

type cliCtx struct {
	tmpDir     string
	tagsFile   string
	envVars    map[string]string
	postFiles  []string // absolute paths
	hn         *fakeHN
	dt         *fakeDT
	prober     *fakeProber
	exitCode   int
	stdout     string
	stderr     string
	extraFlags []string
	noTagsFile bool
}

func (cc *cliCtx) reset() {
	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "publish-test-*")
	if err != nil {
		panic(err)
	}
	cc.tmpDir = tmpDir

	// Create default tags.json
	tagsJSON := `{"go":{"devto":"go"},"testing":{"devto":"testing"}}`
	cc.tagsFile = filepath.Join(tmpDir, "tags.json")
	os.WriteFile(cc.tagsFile, []byte(tagsJSON), 0644)

	// Create posts subdir
	os.MkdirAll(filepath.Join(tmpDir, "posts"), 0755)

	cc.envVars = make(map[string]string)
	cc.postFiles = nil
	cc.hn = &fakeHN{}
	cc.dt = &fakeDT{}
	cc.prober = &fakeProber{}
	cc.exitCode = -1
	cc.stdout = ""
	cc.stderr = ""
	cc.extraFlags = nil
	cc.noTagsFile = false
}

func (cc *cliCtx) cleanup() {
	if cc.tmpDir != "" {
		os.RemoveAll(cc.tmpDir)
	}
}

// --- Given steps ---

func (cc *cliCtx) validCredentialsViaEnvironmentVariables() error {
	cc.envVars["HASHNODE_TOKEN"] = "test-token"
	cc.envVars["HASHNODE_PUBLICATION_ID"] = "test-pub-id"
	cc.envVars["DEVTO_API_KEY"] = "test-api-key"
	return nil
}

func (cc *cliCtx) envVarIsSetTo(name, value string) error {
	cc.envVars[name] = value
	return nil
}

func (cc *cliCtx) aPostFileWith(name string, table *godog.Table) error {
	fields := tableToMap(table)
	content := buildMarkdown(fields)
	path := filepath.Join(cc.tmpDir, "posts", name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}
	cc.postFiles = append(cc.postFiles, path)
	return nil
}

func (cc *cliCtx) hashnodePublishWillFailWith(errMsg string) error {
	cc.hn.publishErr = fmt.Errorf("%s", errMsg)
	return nil
}

// --- When steps ---

func (cc *cliCtx) theCLIRuns() error {
	return cc.runWith("")
}

func (cc *cliCtx) theCLIRunsWith(flags string) error {
	return cc.runWith(flags)
}

func (cc *cliCtx) theCLIRunsWithoutTagsFile() error {
	cc.noTagsFile = true
	// Change to tmpDir so the default "tags.json" is found
	origDir, _ := os.Getwd()
	if err := os.Chdir(cc.tmpDir); err != nil {
		return err
	}
	defer os.Chdir(origDir)
	return cc.runWith("")
}

func (cc *cliCtx) runWith(flags string) error {
	var args []string

	// Add extra flags
	if flags != "" {
		args = append(args, strings.Fields(flags)...)
	}

	// Add --tags-file unless explicitly omitted or --probe or --unknown
	isProbe := strings.Contains(flags, "--probe")
	isUnknown := strings.Contains(flags, "--unknown")
	isDryRunNoFiles := strings.Contains(flags, "--dry-run") && len(cc.postFiles) == 0

	if !cc.noTagsFile && !isProbe && !isUnknown && !isDryRunNoFiles {
		args = append(args, "--tags-file", cc.tagsFile)
	}

	// Add post files
	args = append(args, cc.postFiles...)

	getenv := func(name string) string {
		return cc.envVars[name]
	}

	d := &Deps{
		Hashnode: cc.hn,
		DevTo:    cc.dt,
		Prober:   cc.prober,
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cc.exitCode = Run(args, &stdoutBuf, &stderrBuf, getenv, d)
	cc.stdout = stdoutBuf.String()
	cc.stderr = stderrBuf.String()
	return nil
}

// --- Then steps ---

func (cc *cliCtx) probeAllIsCalled() error {
	if !cc.prober.called {
		return fmt.Errorf("expected ProbeAll to be called")
	}
	return nil
}

func (cc *cliCtx) theExitCodeIs(code int) error {
	if cc.exitCode != code {
		return fmt.Errorf("expected exit code %d, got %d\nstdout: %s\nstderr: %s", code, cc.exitCode, cc.stdout, cc.stderr)
	}
	return nil
}

func (cc *cliCtx) stdoutContainsValidJSON() error {
	var result interface{}
	if err := json.Unmarshal([]byte(cc.stdout), &result); err != nil {
		return fmt.Errorf("stdout is not valid JSON: %v\nstdout: %s", err, cc.stdout)
	}
	return nil
}

func (cc *cliCtx) stdoutContains(expected string) error {
	if !strings.Contains(cc.stdout, expected) {
		return fmt.Errorf("expected stdout to contain %q, got:\n%s", expected, cc.stdout)
	}
	return nil
}

func (cc *cliCtx) stderrContains(expected string) error {
	if !strings.Contains(cc.stderr, expected) {
		return fmt.Errorf("expected stderr to contain %q, got:\n%s", expected, cc.stderr)
	}
	return nil
}

// --- Helpers ---

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
	sb.WriteString("---\n")
	if content, ok := fields["content"]; ok {
		sb.WriteString(strings.ReplaceAll(content, `\n`, "\n"))
	} else {
		sb.WriteString("Test content.")
	}
	return sb.String()
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	cc := &cliCtx{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		cc.reset()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		cc.cleanup()
		return ctx, nil
	})

	// Given
	ctx.Step(`^valid credentials via environment variables$`, cc.validCredentialsViaEnvironmentVariables)
	ctx.Step(`^([A-Z_]+) is set to "([^"]*)"$`, cc.envVarIsSetTo)
	ctx.Step(`^a post file "([^"]*)" with:$`, cc.aPostFileWith)
	ctx.Step(`^Hashnode publish will fail with "([^"]*)"$`, cc.hashnodePublishWillFailWith)

	// When
	ctx.Step(`^the CLI runs$`, cc.theCLIRuns)
	ctx.Step(`^the CLI runs with "([^"]*)"$`, cc.theCLIRunsWith)
	ctx.Step(`^the CLI runs without --tags-file$`, cc.theCLIRunsWithoutTagsFile)

	// Then
	ctx.Step(`^ProbeAll is called$`, cc.probeAllIsCalled)
	ctx.Step(`^the exit code is (\d+)$`, cc.theExitCodeIs)
	ctx.Step(`^stdout contains valid JSON$`, cc.stdoutContainsValidJSON)
	ctx.Step(`^stdout contains "([^"]*)"$`, cc.stdoutContains)
	ctx.Step(`^stderr contains "([^"]*)"$`, cc.stderrContains)
}

// BDD: specs/cli.feature :: Feature: CLI wrapper
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../specs/cli.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
