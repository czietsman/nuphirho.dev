package devto_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/czietsman/nuphirho.dev/internal/devto"
)

// fakeHTTP implements devto.HTTPClient for testing.
type fakeHTTP struct {
	articles    map[string]fakeArticle // canonical_url -> article
	authError   bool
	rateLimited bool
	unreachable bool
	requests    []*http.Request
	lastBody    map[string]interface{}
}

type fakeArticle struct {
	ID  int
	URL string
}

func newFakeHTTP() *fakeHTTP {
	return &fakeHTTP{
		articles: make(map[string]fakeArticle),
	}
}

func (f *fakeHTTP) Do(req *http.Request) (*http.Response, error) {
	f.requests = append(f.requests, req)

	if f.unreachable {
		return nil, fmt.Errorf("connection refused")
	}

	// Capture request body for assertion
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(body))
		var parsed map[string]interface{}
		if json.Unmarshal(body, &parsed) == nil {
			f.lastBody = parsed
		}
	}

	if f.authError {
		return &http.Response{
			StatusCode: 401,
			Body:       io.NopCloser(strings.NewReader(`{"error":"unauthorized"}`)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	}

	if f.rateLimited {
		return &http.Response{
			StatusCode: 429,
			Body:       io.NopCloser(strings.NewReader(`{"error":"rate limited"}`)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	}

	// Route based on method and path
	switch {
	case req.Method == "GET" && strings.Contains(req.URL.Path, "/api/articles/me/all"):
		return f.handleListArticles(), nil
	case req.Method == "POST" && req.URL.Path == "/api/articles":
		return f.handleCreateArticle(), nil
	case req.Method == "PUT" && strings.HasPrefix(req.URL.Path, "/api/articles/"):
		return f.handleUpdateArticle(req.URL.Path), nil
	}

	return &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(strings.NewReader(`{"error":"not found"}`)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

func (f *fakeHTTP) handleListArticles() *http.Response {
	articles := make([]map[string]interface{}, 0)
	for canonicalURL, article := range f.articles {
		articles = append(articles, map[string]interface{}{
			"id":            float64(article.ID),
			"canonical_url": canonicalURL,
			"url":           article.URL,
		})
	}
	b, _ := json.Marshal(articles)
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(b)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func (f *fakeHTTP) handleCreateArticle() *http.Response {
	resp := map[string]interface{}{
		"id":  float64(12345),
		"url": "https://dev.to/nuphirho/my-new-post",
	}
	b, _ := json.Marshal(resp)
	return &http.Response{
		StatusCode: 201,
		Body:       io.NopCloser(bytes.NewReader(b)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func (f *fakeHTTP) handleUpdateArticle(path string) *http.Response {
	// Extract slug from the existing articles to build URL
	idStr := strings.TrimPrefix(path, "/api/articles/")
	id, _ := strconv.Atoi(idStr)

	url := ""
	for _, article := range f.articles {
		if article.ID == id {
			url = article.URL
			break
		}
	}

	resp := map[string]interface{}{
		"id":  float64(id),
		"url": url,
	}
	b, _ := json.Marshal(resp)
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(b)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

// --- godog context ---

type devtoContext struct {
	fake   *fakeHTTP
	client *devto.Client
	result *devto.PublishResult
	err    error
}

func (dc *devtoContext) reset() {
	dc.fake = newFakeHTTP()
	dc.client = devto.New("test-api-key", dc.fake)
	dc.result = nil
	dc.err = nil
}

// --- Given steps ---

func (dc *devtoContext) aDevToClientConfiguredWithAPIKey(key string) error {
	dc.reset()
	dc.client = devto.New(key, dc.fake)
	return nil
}

func (dc *devtoContext) noArticleExistsWithCanonicalURL(url string) error {
	delete(dc.fake.articles, url)
	return nil
}

func (dc *devtoContext) anArticleExistsWithCanonicalURLAndID(url string, id int) error {
	slug := strings.TrimPrefix(url, "https://blog.nuphirho.dev/")
	dc.fake.articles[url] = fakeArticle{
		ID:  id,
		URL: "https://dev.to/nuphirho/" + slug,
	}
	return nil
}

func (dc *devtoContext) theDevToAPIReturnsA401UnauthorizedError() error {
	dc.fake.authError = true
	return nil
}

func (dc *devtoContext) theDevToAPIReturnsA429RateLimitError() error {
	dc.fake.rateLimited = true
	return nil
}

func (dc *devtoContext) theDevToAPIIsUnreachable() error {
	dc.fake.unreachable = true
	return nil
}

func (dc *devtoContext) dryRunModeIsEnabled() error {
	dc.client.DryRun = true
	return nil
}

// --- When steps ---

func (dc *devtoContext) thePipelineCreatesAnArticle(table *godog.Table) error {
	input := tableToArticleInput(table)
	dc.result, dc.err = dc.client.CreateArticle(input)
	return nil
}

// --- Then steps ---

func (dc *devtoContext) aPOSTRequestIsSentTo(path string) error {
	for _, r := range dc.client.RequestsMade() {
		if r.Method == "POST" && r.Path == path {
			return nil
		}
	}
	return fmt.Errorf("expected POST to %s, requests: %+v", path, dc.client.RequestsMade())
}

func (dc *devtoContext) aPUTRequestIsSentTo(path string) error {
	for _, r := range dc.client.RequestsMade() {
		if r.Method == "PUT" && r.Path == path {
			return nil
		}
	}
	return fmt.Errorf("expected PUT to %s, requests: %+v", path, dc.client.RequestsMade())
}

func (dc *devtoContext) theResponseContainsArticleIDAndURL(id int, url string) error {
	if dc.result == nil {
		return fmt.Errorf("no result (error: %v)", dc.err)
	}
	if dc.result.ArticleID != id {
		return fmt.Errorf("expected article ID %d, got %d", id, dc.result.ArticleID)
	}
	if dc.result.URL != url {
		return fmt.Errorf("expected URL %q, got %q", url, dc.result.URL)
	}
	return nil
}

func (dc *devtoContext) theResponseContainsArticleURL(url string) error {
	if dc.result == nil {
		return fmt.Errorf("no result (error: %v)", dc.err)
	}
	if dc.result.URL != url {
		return fmt.Errorf("expected URL %q, got %q", url, dc.result.URL)
	}
	return nil
}

func (dc *devtoContext) theRequestBodyHasPublishedSetTo(val string) error {
	body := dc.fake.lastBody
	if body == nil {
		return fmt.Errorf("no request body captured")
	}
	article, ok := body["article"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no article in body: %+v", body)
	}
	published, ok := article["published"].(bool)
	if !ok {
		return fmt.Errorf("published not a bool in body: %+v", article)
	}
	expected := val == "true"
	if published != expected {
		return fmt.Errorf("expected published=%v, got %v", expected, published)
	}
	return nil
}

func (dc *devtoContext) theRequestBodyContentContains(expected string) error {
	body := dc.fake.lastBody
	if body == nil {
		return fmt.Errorf("no request body captured")
	}
	article, ok := body["article"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no article in body")
	}
	content, _ := article["body_markdown"].(string)
	if !strings.Contains(content, expected) {
		return fmt.Errorf("expected body_markdown to contain %q, got %q", expected, content)
	}
	return nil
}

func (dc *devtoContext) theRequestBodyContentDoesNotContain(unexpected string) error {
	body := dc.fake.lastBody
	if body == nil {
		return fmt.Errorf("no request body captured")
	}
	article, ok := body["article"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no article in body")
	}
	content, _ := article["body_markdown"].(string)
	if strings.Contains(content, unexpected) {
		return fmt.Errorf("expected body_markdown NOT to contain %q, but it does", unexpected)
	}
	return nil
}

func (dc *devtoContext) theRequestBodyTagsAre(expected string) error {
	body := dc.fake.lastBody
	if body == nil {
		return fmt.Errorf("no request body captured")
	}
	article, ok := body["article"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no article in body")
	}
	tags, ok := article["tags"].([]interface{})
	if !ok {
		return fmt.Errorf("tags not an array: %+v", article["tags"])
	}

	// Parse expected JSON array
	var expectedTags []string
	if err := json.Unmarshal([]byte(expected), &expectedTags); err != nil {
		return fmt.Errorf("parsing expected tags: %w", err)
	}

	if len(tags) != len(expectedTags) {
		return fmt.Errorf("expected %d tags, got %d: %v", len(expectedTags), len(tags), tags)
	}
	for i, t := range tags {
		s, _ := t.(string)
		if s != expectedTags[i] {
			return fmt.Errorf("tag %d: expected %q, got %q", i, expectedTags[i], s)
		}
	}
	return nil
}

func (dc *devtoContext) theErrorIs(expected string) error {
	if dc.err == nil {
		return fmt.Errorf("expected error %q, got nil", expected)
	}
	if !strings.Contains(dc.err.Error(), expected) {
		return fmt.Errorf("expected error containing %q, got %q", expected, dc.err.Error())
	}
	return nil
}

func (dc *devtoContext) theErrorContains(expected string) error {
	return dc.theErrorIs(expected)
}

func (dc *devtoContext) noPOSTOrPUTRequestIsMade() error {
	for _, r := range dc.client.RequestsMade() {
		if r.Method == "POST" || r.Method == "PUT" {
			return fmt.Errorf("expected no POST/PUT, got %s %s", r.Method, r.Path)
		}
	}
	return nil
}

func (dc *devtoContext) theDryRunResultActionIs(action string) error {
	if dc.result == nil {
		return fmt.Errorf("no result (error: %v)", dc.err)
	}
	if dc.result.Action != action {
		return fmt.Errorf("expected action %q, got %q", action, dc.result.Action)
	}
	return nil
}

func (dc *devtoContext) theDryRunResultPublishedIs(val string) error {
	if dc.result == nil {
		return fmt.Errorf("no result (error: %v)", dc.err)
	}
	expected := val == "true"
	if dc.result.Published != expected {
		return fmt.Errorf("expected published=%v, got %v", expected, dc.result.Published)
	}
	return nil
}

func (dc *devtoContext) theDryRunResultExistingIDIs(id int) error {
	if dc.result == nil {
		return fmt.Errorf("no result (error: %v)", dc.err)
	}
	if dc.result.ArticleID != id {
		return fmt.Errorf("expected existing ID %d, got %d", id, dc.result.ArticleID)
	}
	return nil
}

func (dc *devtoContext) theDryRunResultEmbedsConvertedIs(count int) error {
	if dc.result == nil {
		return fmt.Errorf("no result (error: %v)", dc.err)
	}
	if dc.result.EmbedsConverted != count {
		return fmt.Errorf("expected embeds converted %d, got %d", count, dc.result.EmbedsConverted)
	}
	return nil
}

func tableToArticleInput(table *godog.Table) devto.ArticleInput {
	input := devto.ArticleInput{}
	for _, row := range table.Rows {
		key := strings.TrimSpace(row.Cells[0].Value)
		val := strings.TrimSpace(row.Cells[1].Value)
		switch key {
		case "title":
			input.Title = val
		case "slug":
			input.Slug = val
		case "content":
			// Unescape \n in feature table values
			input.Content = strings.ReplaceAll(val, `\n`, "\n")
		case "tags":
			input.Tags = strings.Split(val, ",")
		case "published":
			input.Published = val == "true"
		}
	}
	return input
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	dc := &devtoContext{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		dc.reset()
		return ctx, nil
	})

	// Given
	ctx.Step(`^a Dev\.to client configured with API key "([^"]*)"$`, dc.aDevToClientConfiguredWithAPIKey)
	ctx.Step(`^no article exists with canonical URL "([^"]*)"$`, dc.noArticleExistsWithCanonicalURL)
	ctx.Step(`^an article exists with canonical URL "([^"]*)" and ID (\d+)$`, dc.anArticleExistsWithCanonicalURLAndID)
	ctx.Step(`^the Dev\.to API returns a 401 unauthorized error$`, dc.theDevToAPIReturnsA401UnauthorizedError)
	ctx.Step(`^the Dev\.to API returns a 429 rate limit error$`, dc.theDevToAPIReturnsA429RateLimitError)
	ctx.Step(`^the Dev\.to API is unreachable$`, dc.theDevToAPIIsUnreachable)
	ctx.Step(`^dry-run mode is enabled$`, dc.dryRunModeIsEnabled)

	// When
	ctx.Step(`^the pipeline creates an article:$`, dc.thePipelineCreatesAnArticle)

	// Then
	ctx.Step(`^a POST request is sent to "([^"]*)"$`, dc.aPOSTRequestIsSentTo)
	ctx.Step(`^a PUT request is sent to "([^"]*)"$`, dc.aPUTRequestIsSentTo)
	ctx.Step(`^the response contains article ID (\d+) and URL "([^"]*)"$`, dc.theResponseContainsArticleIDAndURL)
	ctx.Step(`^the response contains article URL "([^"]*)"$`, dc.theResponseContainsArticleURL)
	ctx.Step(`^the request body has "published" set to (true|false)$`, dc.theRequestBodyHasPublishedSetTo)
	ctx.Step(`^the request body content contains "([^"]*)"$`, dc.theRequestBodyContentContains)
	ctx.Step(`^the request body content does not contain "([^"]*)"$`, dc.theRequestBodyContentDoesNotContain)
	ctx.Step(`^the request body tags are (\[.*\])$`, dc.theRequestBodyTagsAre)
	ctx.Step(`^the error is "([^"]*)"$`, dc.theErrorIs)
	ctx.Step(`^the error contains "([^"]*)"$`, dc.theErrorContains)
	ctx.Step(`^no POST or PUT request is made$`, dc.noPOSTOrPUTRequestIsMade)
	ctx.Step(`^the dry-run result action is "([^"]*)"$`, dc.theDryRunResultActionIs)
	ctx.Step(`^the dry-run result published is (true|false)$`, dc.theDryRunResultPublishedIs)
	ctx.Step(`^the dry-run result existing ID is (\d+)$`, dc.theDryRunResultExistingIDIs)
	ctx.Step(`^the dry-run result embeds converted is (\d+)$`, dc.theDryRunResultEmbedsConvertedIs)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../specs/devto_client.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
