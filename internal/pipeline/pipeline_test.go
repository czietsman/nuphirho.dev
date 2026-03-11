package pipeline_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/czietsman/nuphirho.dev/internal/devto"
	"github.com/czietsman/nuphirho.dev/internal/hashnode"
	"github.com/czietsman/nuphirho.dev/internal/pipeline"
	"github.com/czietsman/nuphirho.dev/internal/tags"
)

// --- Fakes ---

type fakeHN struct {
	publishCalls []string // slugs
	publishErr   error
	slugErrors   map[string]error
	lastInputs   map[string]hashnode.PostInput // slug -> last input
}

func (f *fakeHN) Publish(input hashnode.PostInput) (*hashnode.PublishResult, error) {
	f.publishCalls = append(f.publishCalls, input.Slug)
	if f.lastInputs == nil {
		f.lastInputs = make(map[string]hashnode.PostInput)
	}
	f.lastInputs[input.Slug] = input
	if err, ok := f.slugErrors[input.Slug]; ok {
		return nil, err
	}
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
	createCalls []string // slugs
	createErr   error
	slugErrors  map[string]error
	lastInputs  map[string]devto.ArticleInput // slug -> last input
}

func (f *fakeDT) CreateArticle(input devto.ArticleInput) (*devto.PublishResult, error) {
	f.createCalls = append(f.createCalls, input.Slug)
	if f.lastInputs == nil {
		f.lastInputs = make(map[string]devto.ArticleInput)
	}
	f.lastInputs[input.Slug] = input
	if err, ok := f.slugErrors[input.Slug]; ok {
		return nil, err
	}
	if f.createErr != nil {
		return nil, f.createErr
	}
	return &devto.PublishResult{
		Action:          "create",
		ArticleID:       12345,
		URL:             "https://dev.to/nuphirho/" + input.Slug,
		Published:       input.Published,
		EmbedsConverted: countEmbeds(input.Content),
	}, nil
}

func countEmbeds(content string) int {
	count := 0
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "{% embed ") || strings.HasPrefix(line, "%[") {
			count++
		}
	}
	return count
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

type fakeSeriesResolver struct {
	err   error
	calls []string // names resolved
}

func (f *fakeSeriesResolver) ResolveSeriesID(name string) (string, error) {
	f.calls = append(f.calls, name)
	if f.err != nil {
		return "", f.err
	}
	return "series-resolved", nil
}

// --- godog context ---

type pipelineCtx struct {
	glossary       tags.Glossary
	hn             *fakeHN
	dt             *fakeDT
	prober         *fakeProber
	seriesResolver *fakeSeriesResolver
	files          []pipeline.PostFile
	result         *pipeline.RunResult
	output         string
	dryRun         bool
	probeRan       bool
}

func (pc *pipelineCtx) reset() {
	pc.glossary = nil
	pc.hn = &fakeHN{slugErrors: make(map[string]error)}
	pc.dt = &fakeDT{slugErrors: make(map[string]error)}
	pc.prober = &fakeProber{}
	pc.seriesResolver = &fakeSeriesResolver{}
	pc.files = nil
	pc.result = nil
	pc.output = ""
	pc.dryRun = false
	pc.probeRan = false
}

// --- Given steps ---

func (pc *pipelineCtx) aTagGlossary(glossaryJSON *godog.DocString) error {
	g, err := tags.ParseGlossary([]byte(glossaryJSON.Content))
	if err != nil {
		return err
	}
	pc.glossary = g
	return nil
}

func (pc *pipelineCtx) aPostFileWith(path string, table *godog.Table) error {
	fields := tableToMap(table)
	pc.files = append(pc.files, pipeline.PostFile{
		Path:    path,
		Content: buildMarkdown(fields),
	})
	return nil
}

func (pc *pipelineCtx) hashnodePublishWillFailWith(errMsg string) error {
	pc.hn.publishErr = fmt.Errorf("%s", errMsg)
	return nil
}

func (pc *pipelineCtx) hashnodePublishWillFailForSlugWith(slug, errMsg string) error {
	pc.hn.slugErrors[slug] = fmt.Errorf("%s", errMsg)
	return nil
}

func (pc *pipelineCtx) devtoCrossPostWillFailWith(errMsg string) error {
	pc.dt.createErr = fmt.Errorf("%s", errMsg)
	return nil
}

func (pc *pipelineCtx) dryRunModeIsEnabled() error {
	pc.dryRun = true
	return nil
}

// --- When steps ---

func (pc *pipelineCtx) seriesResolutionWillFailWith(errMsg string) error {
	pc.seriesResolver.err = fmt.Errorf("%s", errMsg)
	return nil
}

func (pc *pipelineCtx) thePipelineRuns() error {
	cfg := pipeline.Config{
		Hashnode:       pc.hn,
		DevTo:          pc.dt,
		SeriesResolver: pc.seriesResolver,
		Glossary:       pc.glossary,
		DryRun:         pc.dryRun,
	}

	var buf bytes.Buffer
	pc.result = pipeline.Run(cfg, pc.files, &buf)
	pc.output = buf.String()
	return nil
}

func (pc *pipelineCtx) thePipelineRunsWithProbe() error {
	var buf bytes.Buffer
	exitCode := pipeline.Probe(pc.prober, &buf)
	pc.probeRan = true
	pc.output = buf.String()
	pc.result = &pipeline.RunResult{ExitCode: exitCode}
	return nil
}

// --- Then steps ---

func (pc *pipelineCtx) hashnodePublishIsCalledWithSlug(slug string) error {
	for _, s := range pc.hn.publishCalls {
		if s == slug {
			return nil
		}
	}
	return fmt.Errorf("expected Hashnode publish called with slug %q, calls: %v", slug, pc.hn.publishCalls)
}

func (pc *pipelineCtx) hashnodePublishIsNotCalled() error {
	if len(pc.hn.publishCalls) > 0 {
		return fmt.Errorf("expected no Hashnode publish calls, got: %v", pc.hn.publishCalls)
	}
	return nil
}

func (pc *pipelineCtx) devtoCrossPostIsCalledWithSlug(slug string) error {
	for _, s := range pc.dt.createCalls {
		if s == slug {
			return nil
		}
	}
	return fmt.Errorf("expected Dev.to cross-post called with slug %q, calls: %v", slug, pc.dt.createCalls)
}

func (pc *pipelineCtx) devtoCrossPostIsNotCalled() error {
	if len(pc.dt.createCalls) > 0 {
		return fmt.Errorf("expected no Dev.to cross-post calls, got: %v", pc.dt.createCalls)
	}
	return nil
}

func (pc *pipelineCtx) theSummaryContains(expected string) error {
	if !strings.Contains(pc.output, expected) {
		return fmt.Errorf("expected output to contain %q, got:\n%s", expected, pc.output)
	}
	return nil
}

func (pc *pipelineCtx) theExitCodeIs(code int) error {
	if pc.result.ExitCode != code {
		return fmt.Errorf("expected exit code %d, got %d\noutput:\n%s", code, pc.result.ExitCode, pc.output)
	}
	return nil
}

func (pc *pipelineCtx) theOutputIsValidJSON() error {
	var result []interface{}
	if err := json.Unmarshal([]byte(pc.output), &result); err != nil {
		return fmt.Errorf("output is not valid JSON: %v\noutput:\n%s", err, pc.output)
	}
	return nil
}

func (pc *pipelineCtx) theJSONResultHasHashnodeAction(file, action string) error {
	entry, err := pc.findJSONEntry(file)
	if err != nil {
		return err
	}
	hn, ok := entry["hashnode"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no hashnode entry for %s", file)
	}
	if a, _ := hn["action"].(string); a != action {
		return fmt.Errorf("expected hashnode action %q, got %q", action, a)
	}
	return nil
}

func (pc *pipelineCtx) theJSONResultHasDevtoAction(file, action string) error {
	entry, err := pc.findJSONEntry(file)
	if err != nil {
		return err
	}
	dt, ok := entry["devto"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no devto entry for %s", file)
	}
	if a, _ := dt["action"].(string); a != action {
		return fmt.Errorf("expected devto action %q, got %q", action, a)
	}
	return nil
}

func (pc *pipelineCtx) theJSONResultHasDevtoEmbedsConverted(file string, count int) error {
	entry, err := pc.findJSONEntry(file)
	if err != nil {
		return err
	}
	dt, ok := entry["devto"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no devto entry for %s", file)
	}
	ec, _ := dt["embeds_converted"].(float64)
	if int(ec) != count {
		return fmt.Errorf("expected embeds_converted %d, got %v", count, dt["embeds_converted"])
	}
	return nil
}

func (pc *pipelineCtx) theJSONResultHasDevtoPublished(file string, val string) error {
	entry, err := pc.findJSONEntry(file)
	if err != nil {
		return err
	}
	dt, ok := entry["devto"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no devto entry for %s", file)
	}
	expected := val == "true"
	published, _ := dt["published"].(bool)
	if published != expected {
		return fmt.Errorf("expected published %v, got %v", expected, published)
	}
	return nil
}

func (pc *pipelineCtx) theHashnodeCallIncludesSeriesID(slug, id string) error {
	input, ok := pc.hn.lastInputs[slug]
	if !ok {
		return fmt.Errorf("no Hashnode call for slug %q", slug)
	}
	if input.SeriesID != id {
		return fmt.Errorf("expected Hashnode seriesID %q for slug %q, got %q", id, slug, input.SeriesID)
	}
	return nil
}

func (pc *pipelineCtx) theHashnodeCallDoesNotIncludeSeries(slug string) error {
	input, ok := pc.hn.lastInputs[slug]
	if !ok {
		return fmt.Errorf("no Hashnode call for slug %q", slug)
	}
	if input.SeriesID != "" {
		return fmt.Errorf("expected no seriesID for slug %q, got %q", slug, input.SeriesID)
	}
	return nil
}

func (pc *pipelineCtx) theSeriesResolverIsCalledOnceFor(name string) error {
	count := 0
	for _, n := range pc.seriesResolver.calls {
		if n == name {
			count++
		}
	}
	if count != 1 {
		return fmt.Errorf("expected series resolver called once for %q, got %d calls", name, count)
	}
	return nil
}

func (pc *pipelineCtx) theDevtoCallIncludesSeries(slug, series string) error {
	input, ok := pc.dt.lastInputs[slug]
	if !ok {
		return fmt.Errorf("no Dev.to call for slug %q", slug)
	}
	if input.Series != series {
		return fmt.Errorf("expected Dev.to series %q for slug %q, got %q", series, slug, input.Series)
	}
	return nil
}

func (pc *pipelineCtx) probeAllIsCalled() error {
	if !pc.prober.called {
		return fmt.Errorf("expected ProbeAll to be called")
	}
	return nil
}

func (pc *pipelineCtx) findJSONEntry(file string) (map[string]interface{}, error) {
	var entries []map[string]interface{}
	if err := json.Unmarshal([]byte(pc.output), &entries); err != nil {
		return nil, fmt.Errorf("parsing JSON: %v", err)
	}
	for _, e := range entries {
		if f, _ := e["file"].(string); f == file {
			return e, nil
		}
	}
	return nil, fmt.Errorf("no JSON entry for file %q", file)
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
	if series, ok := fields["series"]; ok && series != "" {
		sb.WriteString(fmt.Sprintf("series: %s\n", series))
	}
	sb.WriteString("---\n")

	if content, ok := fields["content"]; ok {
		// Unescape \n
		content = strings.ReplaceAll(content, `\n`, "\n")
		sb.WriteString(content)
	} else {
		sb.WriteString("Test content.")
	}

	return sb.String()
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	pc := &pipelineCtx{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		pc.reset()
		return ctx, nil
	})

	// Given
	ctx.Step(`^a tag glossary:$`, pc.aTagGlossary)
	ctx.Step(`^a post file "([^"]*)" with:$`, pc.aPostFileWith)
	ctx.Step(`^Hashnode publish will fail with "([^"]*)"$`, pc.hashnodePublishWillFailWith)
	ctx.Step(`^Hashnode publish will fail for slug "([^"]*)" with "([^"]*)"$`, pc.hashnodePublishWillFailForSlugWith)
	ctx.Step(`^Dev\.to cross-post will fail with "([^"]*)"$`, pc.devtoCrossPostWillFailWith)
	ctx.Step(`^series resolution will fail with "([^"]*)"$`, pc.seriesResolutionWillFailWith)
	ctx.Step(`^dry-run mode is enabled$`, pc.dryRunModeIsEnabled)

	// When
	ctx.Step(`^the pipeline runs$`, pc.thePipelineRuns)
	ctx.Step(`^the pipeline runs with --probe$`, pc.thePipelineRunsWithProbe)

	// Then
	ctx.Step(`^Hashnode publish is called with slug "([^"]*)"$`, pc.hashnodePublishIsCalledWithSlug)
	ctx.Step(`^Hashnode publish is not called$`, pc.hashnodePublishIsNotCalled)
	ctx.Step(`^Dev\.to cross-post is called with slug "([^"]*)"$`, pc.devtoCrossPostIsCalledWithSlug)
	ctx.Step(`^Dev\.to cross-post is not called$`, pc.devtoCrossPostIsNotCalled)
	ctx.Step(`^the summary contains "([^"]*)"$`, pc.theSummaryContains)
	ctx.Step(`^the exit code is (\d+)$`, pc.theExitCodeIs)
	ctx.Step(`^the output is valid JSON$`, pc.theOutputIsValidJSON)
	ctx.Step(`^the JSON result for "([^"]*)" has hashnode action "([^"]*)"$`, pc.theJSONResultHasHashnodeAction)
	ctx.Step(`^the JSON result for "([^"]*)" has devto action "([^"]*)"$`, pc.theJSONResultHasDevtoAction)
	ctx.Step(`^the JSON result for "([^"]*)" has devto embeds_converted (\d+)$`, pc.theJSONResultHasDevtoEmbedsConverted)
	ctx.Step(`^the JSON result for "([^"]*)" has devto published (true|false)$`, pc.theJSONResultHasDevtoPublished)
	ctx.Step(`^ProbeAll is called$`, pc.probeAllIsCalled)
	ctx.Step(`^the Hashnode call for "([^"]*)" includes series ID "([^"]*)"$`, pc.theHashnodeCallIncludesSeriesID)
	ctx.Step(`^the Hashnode call for "([^"]*)" does not include series$`, pc.theHashnodeCallDoesNotIncludeSeries)
	ctx.Step(`^the Dev\.to call for "([^"]*)" includes series "([^"]*)"$`, pc.theDevtoCallIncludesSeries)
	ctx.Step(`^the series resolver is called once for "([^"]*)"$`, pc.theSeriesResolverIsCalledOnceFor)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../specs/pipeline.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
