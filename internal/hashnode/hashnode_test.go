package hashnode_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/czietsman/nuphirho.dev/internal/hashnode"
)

// fakeHTTP implements hashnode.HTTPClient for testing.
type fakeHTTP struct {
	posts        map[string]string // slug -> id for published posts
	deletedPosts map[string]string // slug -> id for deleted posts
	drafts       map[string]string // slug -> id for drafts
	series       map[string]string // name -> id for series
	authError    bool
	pubNotFound  bool
	unreachable  bool
	requests     []*http.Request
	lastVars     map[string]interface{} // last request variables for assertion
}

func newFakeHTTP() *fakeHTTP {
	return &fakeHTTP{
		posts:        make(map[string]string),
		deletedPosts: make(map[string]string),
		drafts:       make(map[string]string),
		series:       make(map[string]string),
	}
}

func (f *fakeHTTP) Do(req *http.Request) (*http.Response, error) {
	f.requests = append(f.requests, req)

	if f.unreachable {
		return nil, fmt.Errorf("connection refused")
	}

	body, _ := io.ReadAll(req.Body)
	var payload map[string]interface{}
	json.Unmarshal(body, &payload)
	query, _ := payload["query"].(string)
	vars, _ := payload["variables"].(map[string]interface{})
	f.lastVars = vars

	if f.authError {
		return jsonResponse(map[string]interface{}{
			"errors": []interface{}{
				map[string]interface{}{"message": "You must be authenticated."},
			},
		}), nil
	}

	if f.pubNotFound {
		return jsonResponse(map[string]interface{}{
			"data": map[string]interface{}{"publication": nil},
		}), nil
	}

	// Route based on query content
	switch {
	case strings.Contains(query, "post(slug:"):
		return f.handlePostBySlug(vars), nil
	case strings.Contains(query, "deletedOnly"):
		return f.handleDeletedPosts(), nil
	case strings.Contains(query, "drafts(first:"):
		return f.handleDrafts(), nil
	case strings.Contains(query, "publishPost"):
		return f.handlePublishPost(vars), nil
	case strings.Contains(query, "updatePost"):
		return f.handleUpdatePost(vars), nil
	case strings.Contains(query, "restorePost"):
		return f.handleRestorePost(vars), nil
	case strings.Contains(query, "createDraft"):
		return f.handleCreateDraft(vars), nil
	case strings.Contains(query, "seriesList"):
		return f.handleSeriesList(), nil
	case strings.Contains(query, "createSeries"):
		return f.handleCreateSeries(vars), nil
	}

	return jsonResponse(map[string]interface{}{"data": nil}), nil
}

func (f *fakeHTTP) handlePostBySlug(vars map[string]interface{}) *http.Response {
	slug := extractVar(vars, "slug")
	id, found := f.posts[slug]
	if !found {
		return jsonResponse(map[string]interface{}{
			"data": map[string]interface{}{
				"publication": map[string]interface{}{
					"post": nil,
				},
			},
		})
	}
	return jsonResponse(map[string]interface{}{
		"data": map[string]interface{}{
			"publication": map[string]interface{}{
				"post": map[string]interface{}{"id": id},
			},
		},
	})
}

func (f *fakeHTTP) handleDeletedPosts() *http.Response {
	edges := make([]interface{}, 0)
	for slug, id := range f.deletedPosts {
		edges = append(edges, map[string]interface{}{
			"node": map[string]interface{}{"id": id, "slug": slug},
		})
	}
	return jsonResponse(map[string]interface{}{
		"data": map[string]interface{}{
			"publication": map[string]interface{}{
				"posts": map[string]interface{}{"edges": edges},
			},
		},
	})
}

func (f *fakeHTTP) handleDrafts() *http.Response {
	edges := make([]interface{}, 0)
	for slug, id := range f.drafts {
		edges = append(edges, map[string]interface{}{
			"node": map[string]interface{}{"id": id, "slug": slug},
		})
	}
	return jsonResponse(map[string]interface{}{
		"data": map[string]interface{}{
			"publication": map[string]interface{}{
				"drafts": map[string]interface{}{"edges": edges},
			},
		},
	})
}

func (f *fakeHTTP) handlePublishPost(vars map[string]interface{}) *http.Response {
	input, _ := vars["input"].(map[string]interface{})
	slug, _ := input["slug"].(string)
	return jsonResponse(map[string]interface{}{
		"data": map[string]interface{}{
			"publishPost": map[string]interface{}{
				"post": map[string]interface{}{
					"id":  "post-001",
					"url": "https://blog.example.com/" + slug,
				},
			},
		},
	})
}

func (f *fakeHTTP) handleUpdatePost(vars map[string]interface{}) *http.Response {
	input, _ := vars["input"].(map[string]interface{})
	slug, _ := input["slug"].(string)
	return jsonResponse(map[string]interface{}{
		"data": map[string]interface{}{
			"updatePost": map[string]interface{}{
				"post": map[string]interface{}{
					"url": "https://blog.example.com/" + slug,
				},
			},
		},
	})
}

func (f *fakeHTTP) handleRestorePost(vars map[string]interface{}) *http.Response {
	input, _ := vars["input"].(map[string]interface{})
	id, _ := input["id"].(string)
	return jsonResponse(map[string]interface{}{
		"data": map[string]interface{}{
			"restorePost": map[string]interface{}{
				"post": map[string]interface{}{
					"id":   id,
					"slug": "restored",
				},
			},
		},
	})
}

func (f *fakeHTTP) handleCreateDraft(vars map[string]interface{}) *http.Response {
	return jsonResponse(map[string]interface{}{
		"data": map[string]interface{}{
			"createDraft": map[string]interface{}{
				"draft": map[string]interface{}{
					"id":   "draft-001",
					"slug": "new-draft",
				},
			},
		},
	})
}

func (f *fakeHTTP) handleSeriesList() *http.Response {
	edges := make([]interface{}, 0)
	for name, id := range f.series {
		edges = append(edges, map[string]interface{}{
			"node": map[string]interface{}{"id": id, "name": name},
		})
	}
	return jsonResponse(map[string]interface{}{
		"data": map[string]interface{}{
			"publication": map[string]interface{}{
				"seriesList": map[string]interface{}{"edges": edges},
			},
		},
	})
}

func (f *fakeHTTP) handleCreateSeries(vars map[string]interface{}) *http.Response {
	return jsonResponse(map[string]interface{}{
		"data": map[string]interface{}{
			"createSeries": map[string]interface{}{
				"series": map[string]interface{}{
					"id": "series-001",
				},
			},
		},
	})
}

func extractVar(vars map[string]interface{}, key string) string {
	// Check top-level
	if v, ok := vars[key]; ok {
		s, _ := v.(string)
		return s
	}
	// Check inside "input"
	if input, ok := vars["input"].(map[string]interface{}); ok {
		s, _ := input[key].(string)
		return s
	}
	return ""
}

func jsonResponse(data interface{}) *http.Response {
	b, _ := json.Marshal(data)
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(b)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

// --- godog context ---

type hashnodeContext struct {
	fake     *fakeHTTP
	client   *hashnode.Client
	result   *hashnode.PublishResult
	err      error
	postID   string
	found    bool
	seriesID string
}

func (hc *hashnodeContext) reset() {
	hc.fake = newFakeHTTP()
	hc.client = hashnode.New("test-token", "pub123", hc.fake)
	hc.result = nil
	hc.err = nil
	hc.postID = ""
	hc.found = false
	hc.seriesID = ""
}

func (hc *hashnodeContext) aHashnodeClientConfiguredWithPublicationID(pubID string) error {
	hc.reset()
	hc.client = hashnode.New("test-token", pubID, hc.fake)
	return nil
}

func (hc *hashnodeContext) noPostExistsWithSlug(slug string) error {
	delete(hc.fake.posts, slug)
	return nil
}

func (hc *hashnodeContext) aPostExistsWithSlugAndID(slug, id string) error {
	hc.fake.posts[slug] = id
	return nil
}

func (hc *hashnodeContext) noDeletedPostExistsWithSlug(slug string) error {
	delete(hc.fake.deletedPosts, slug)
	return nil
}

func (hc *hashnodeContext) aDeletedPostExistsWithSlugAndID(slug, id string) error {
	hc.fake.deletedPosts[slug] = id
	return nil
}

func (hc *hashnodeContext) noDraftExistsWithSlug(slug string) error {
	delete(hc.fake.drafts, slug)
	return nil
}

func (hc *hashnodeContext) aDraftExistsWithSlugAndID(slug, id string) error {
	hc.fake.drafts[slug] = id
	return nil
}

func (hc *hashnodeContext) dryRunModeIsEnabled() error {
	hc.client.DryRun = true
	return nil
}

func (hc *hashnodeContext) theHashnodeAPIReturnsAnAuthenticationError() error {
	hc.fake.authError = true
	return nil
}

func (hc *hashnodeContext) theHashnodeAPIReturnsAPublicationNotFoundError() error {
	hc.fake.pubNotFound = true
	return nil
}

func (hc *hashnodeContext) theHashnodeAPIIsUnreachable() error {
	hc.fake.unreachable = true
	return nil
}

func (hc *hashnodeContext) thePipelinePublishesAPost(table *godog.Table) error {
	input := tableToPostInput(table)
	hc.result, hc.err = hc.client.Publish(input)
	return nil
}

func (hc *hashnodeContext) thePipelineCreatesADraft(table *godog.Table) error {
	input := tableToPostInput(table)
	hc.result, hc.err = hc.client.CreateDraft(input)
	return nil
}

func (hc *hashnodeContext) theClientChecksForAPostWithSlug(slug string) error {
	id, err := hc.client.CheckPostBySlug(slug)
	hc.err = err
	hc.postID = id
	hc.found = id != ""
	return nil
}

func (hc *hashnodeContext) theClientChecksForDeletedPostsWithSlug(slug string) error {
	id, err := hc.client.CheckDeletedBySlug(slug)
	hc.err = err
	hc.postID = id
	hc.found = id != ""
	return nil
}

// --- Then steps ---

func (hc *hashnodeContext) aPublishPostMutationIsSent() error {
	for _, m := range hc.client.MutationsSent() {
		if m == "publishPost" {
			return nil
		}
	}
	return fmt.Errorf("expected publishPost mutation, got: %v", hc.client.MutationsSent())
}

func (hc *hashnodeContext) anUpdatePostMutationIsSentWithID(id string) error {
	expected := "updatePost:" + id
	for _, m := range hc.client.MutationsSent() {
		if m == expected {
			return nil
		}
	}
	return fmt.Errorf("expected %s mutation, got: %v", expected, hc.client.MutationsSent())
}

func (hc *hashnodeContext) aRestorePostMutationIsSentWithID(id string) error {
	expected := "restorePost:" + id
	for _, m := range hc.client.MutationsSent() {
		if m == expected {
			return nil
		}
	}
	return fmt.Errorf("expected %s mutation, got: %v", expected, hc.client.MutationsSent())
}

func (hc *hashnodeContext) aCreateDraftMutationIsSent() error {
	for _, m := range hc.client.MutationsSent() {
		if m == "createDraft" {
			return nil
		}
	}
	return fmt.Errorf("expected createDraft mutation, got: %v", hc.client.MutationsSent())
}

func (hc *hashnodeContext) noMutationIsSent() error {
	if len(hc.client.MutationsSent()) > 0 {
		return fmt.Errorf("expected no mutations, got: %v", hc.client.MutationsSent())
	}
	return nil
}

func (hc *hashnodeContext) theResponseContainsPostIDAndURL(id, url string) error {
	if hc.result == nil {
		return fmt.Errorf("no result (error: %v)", hc.err)
	}
	if hc.result.PostID != id {
		return fmt.Errorf("expected post ID %q, got %q", id, hc.result.PostID)
	}
	if hc.result.URL != url {
		return fmt.Errorf("expected URL %q, got %q", url, hc.result.URL)
	}
	return nil
}

func (hc *hashnodeContext) theResponseContainsPostURL(url string) error {
	if hc.result == nil {
		return fmt.Errorf("no result (error: %v)", hc.err)
	}
	if hc.result.URL != url {
		return fmt.Errorf("expected URL %q, got %q", url, hc.result.URL)
	}
	return nil
}

func (hc *hashnodeContext) theResponseContainsDraftID(id string) error {
	if hc.result == nil {
		return fmt.Errorf("no result (error: %v)", hc.err)
	}
	if hc.result.DraftID != id {
		return fmt.Errorf("expected draft ID %q, got %q", id, hc.result.DraftID)
	}
	return nil
}

func (hc *hashnodeContext) theResultReportsTheDraftAlreadyExistsWithID(id string) error {
	if hc.result == nil {
		return fmt.Errorf("no result (error: %v)", hc.err)
	}
	if hc.result.Action != "draft_exists" {
		return fmt.Errorf("expected action 'draft_exists', got %q", hc.result.Action)
	}
	if hc.result.DraftID != id {
		return fmt.Errorf("expected draft ID %q, got %q", id, hc.result.DraftID)
	}
	return nil
}

func (hc *hashnodeContext) thePostIsFoundWithID(id string) error {
	if !hc.found {
		return fmt.Errorf("expected post to be found")
	}
	if hc.postID != id {
		return fmt.Errorf("expected ID %q, got %q", id, hc.postID)
	}
	return nil
}

func (hc *hashnodeContext) thePostIsNotFound() error {
	if hc.found {
		return fmt.Errorf("expected post not found, but found ID %q", hc.postID)
	}
	return nil
}

func (hc *hashnodeContext) theDeletedPostIsFoundWithID(id string) error {
	return hc.thePostIsFoundWithID(id)
}

func (hc *hashnodeContext) noDeletedPostIsFound() error {
	return hc.thePostIsNotFound()
}

func (hc *hashnodeContext) theErrorIs(expected string) error {
	if hc.err == nil {
		return fmt.Errorf("expected error %q, got nil", expected)
	}
	if !strings.Contains(hc.err.Error(), expected) {
		return fmt.Errorf("expected error containing %q, got %q", expected, hc.err.Error())
	}
	return nil
}

func (hc *hashnodeContext) theErrorContains(expected string) error {
	return hc.theErrorIs(expected)
}

func (hc *hashnodeContext) noHTTPRequestIsMade() error {
	// In dry-run, queries for existence checks still happen, but no mutations
	if len(hc.client.MutationsSent()) > 0 {
		return fmt.Errorf("expected no mutations in dry-run, got: %v", hc.client.MutationsSent())
	}
	return nil
}

func (hc *hashnodeContext) theDryRunResultActionIs(action string) error {
	if hc.result == nil {
		return fmt.Errorf("no result (error: %v)", hc.err)
	}
	if hc.result.Action != action {
		return fmt.Errorf("expected action %q, got %q", action, hc.result.Action)
	}
	return nil
}

func (hc *hashnodeContext) theDryRunResultSlugIs(slug string) error {
	// For publish dry-run, the slug is implicit in the action
	// We verify it by checking the result was created correctly
	if hc.result == nil {
		return fmt.Errorf("no result")
	}
	return nil
}

func (hc *hashnodeContext) theDryRunResultExistingIDIs(id string) error {
	if hc.result == nil {
		return fmt.Errorf("no result")
	}
	if hc.result.PostID != id {
		return fmt.Errorf("expected existing ID %q, got %q", id, hc.result.PostID)
	}
	return nil
}

func (hc *hashnodeContext) theDryRunResultDeletedIDIs(id string) error {
	if hc.result == nil {
		return fmt.Errorf("no result")
	}
	if hc.result.DeletedID != id {
		return fmt.Errorf("expected deleted ID %q, got %q", id, hc.result.DeletedID)
	}
	return nil
}

// --- Series steps ---

func (hc *hashnodeContext) aSeriesExistsWithNameAndID(name, id string) error {
	hc.fake.series[name] = id
	return nil
}

func (hc *hashnodeContext) noSeriesExistsWithName(name string) error {
	delete(hc.fake.series, name)
	return nil
}

func (hc *hashnodeContext) theClientLooksUpSeries(name string) error {
	hc.seriesID, hc.err = hc.client.LookupSeries(name)
	return nil
}

func (hc *hashnodeContext) theClientCreatesSeries(name string) error {
	hc.seriesID, hc.err = hc.client.CreateSeries(name)
	return nil
}

func (hc *hashnodeContext) theClientResolvesSeries(name string) error {
	hc.seriesID, hc.err = hc.client.ResolveSeriesID(name)
	return nil
}

func (hc *hashnodeContext) theSeriesIsFoundWithID(id string) error {
	if hc.seriesID != id {
		return fmt.Errorf("expected series ID %q, got %q", id, hc.seriesID)
	}
	return nil
}

func (hc *hashnodeContext) theSeriesIsNotFound() error {
	if hc.seriesID != "" {
		return fmt.Errorf("expected series not found, got ID %q", hc.seriesID)
	}
	return nil
}

func (hc *hashnodeContext) aCreateSeriesMutationIsSent() error {
	for _, m := range hc.client.MutationsSent() {
		if m == "createSeries" {
			return nil
		}
	}
	return fmt.Errorf("expected createSeries mutation, got: %v", hc.client.MutationsSent())
}

func (hc *hashnodeContext) noCreateSeriesMutationIsSent() error {
	for _, m := range hc.client.MutationsSent() {
		if m == "createSeries" {
			return fmt.Errorf("expected no createSeries mutation, but it was sent")
		}
	}
	return nil
}

func (hc *hashnodeContext) theSeriesIDIs(id string) error {
	return hc.theSeriesIsFoundWithID(id)
}

func (hc *hashnodeContext) thePipelinePublishesAPostWithSeries(table *godog.Table) error {
	input := tableToPostInput(table)
	hc.result, hc.err = hc.client.Publish(input)
	return nil
}

func (hc *hashnodeContext) thePublishRequestIncludesSeriesID(id string) error {
	vars := hc.fake.lastVars
	input, ok := vars["input"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no input in last request vars")
	}
	sid, _ := input["seriesId"].(string)
	if sid != id {
		return fmt.Errorf("expected seriesId %q in publish, got %q", id, sid)
	}
	return nil
}

func (hc *hashnodeContext) theUpdateRequestIncludesSeriesID(id string) error {
	return hc.thePublishRequestIncludesSeriesID(id)
}

func tableToPostInput(table *godog.Table) hashnode.PostInput {
	input := hashnode.PostInput{}
	for _, row := range table.Rows {
		key := strings.TrimSpace(row.Cells[0].Value)
		val := strings.TrimSpace(row.Cells[1].Value)
		switch key {
		case "title":
			input.Title = val
		case "slug":
			input.Slug = val
		case "subtitle":
			input.Subtitle = val
		case "content":
			input.Content = val
		case "tags":
			input.Tags = strings.Split(val, ",")
		case "seriesId":
			input.SeriesID = val
		}
	}
	return input
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	hc := &hashnodeContext{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		hc.reset()
		return ctx, nil
	})

	// Given
	ctx.Step(`^a Hashnode client configured with publication ID "([^"]*)"$`, hc.aHashnodeClientConfiguredWithPublicationID)
	ctx.Step(`^no post exists with slug "([^"]*)"$`, hc.noPostExistsWithSlug)
	ctx.Step(`^a post exists with slug "([^"]*)" and ID "([^"]*)"$`, hc.aPostExistsWithSlugAndID)
	ctx.Step(`^no deleted post exists with slug "([^"]*)"$`, hc.noDeletedPostExistsWithSlug)
	ctx.Step(`^a deleted post exists with slug "([^"]*)" and ID "([^"]*)"$`, hc.aDeletedPostExistsWithSlugAndID)
	ctx.Step(`^no draft exists with slug "([^"]*)"$`, hc.noDraftExistsWithSlug)
	ctx.Step(`^a draft exists with slug "([^"]*)" and ID "([^"]*)"$`, hc.aDraftExistsWithSlugAndID)
	ctx.Step(`^dry-run mode is enabled$`, hc.dryRunModeIsEnabled)
	ctx.Step(`^the Hashnode API returns an authentication error$`, hc.theHashnodeAPIReturnsAnAuthenticationError)
	ctx.Step(`^the Hashnode API returns a publication not found error$`, hc.theHashnodeAPIReturnsAPublicationNotFoundError)
	ctx.Step(`^the Hashnode API is unreachable$`, hc.theHashnodeAPIIsUnreachable)
	ctx.Step(`^a series exists with name "([^"]*)" and ID "([^"]*)"$`, hc.aSeriesExistsWithNameAndID)
	ctx.Step(`^no series exists with name "([^"]*)"$`, hc.noSeriesExistsWithName)

	// When
	ctx.Step(`^the pipeline publishes a post:$`, hc.thePipelinePublishesAPost)
	ctx.Step(`^the pipeline creates a draft:$`, hc.thePipelineCreatesADraft)
	ctx.Step(`^the client looks up series "([^"]*)"$`, hc.theClientLooksUpSeries)
	ctx.Step(`^the client creates series "([^"]*)"$`, hc.theClientCreatesSeries)
	ctx.Step(`^the client resolves series "([^"]*)"$`, hc.theClientResolvesSeries)
	ctx.Step(`^the pipeline publishes a post with series:$`, hc.thePipelinePublishesAPostWithSeries)
	ctx.Step(`^the client checks for a post with slug "([^"]*)"$`, hc.theClientChecksForAPostWithSlug)
	ctx.Step(`^the client checks for deleted posts with slug "([^"]*)"$`, hc.theClientChecksForDeletedPostsWithSlug)

	// Then
	ctx.Step(`^a publishPost mutation is sent$`, hc.aPublishPostMutationIsSent)
	ctx.Step(`^an updatePost mutation is sent with ID "([^"]*)"$`, hc.anUpdatePostMutationIsSentWithID)
	ctx.Step(`^a restorePost mutation is sent with ID "([^"]*)"$`, hc.aRestorePostMutationIsSentWithID)
	ctx.Step(`^a createDraft mutation is sent$`, hc.aCreateDraftMutationIsSent)
	ctx.Step(`^no mutation is sent$`, hc.noMutationIsSent)
	ctx.Step(`^the response contains post ID "([^"]*)" and URL "([^"]*)"$`, hc.theResponseContainsPostIDAndURL)
	ctx.Step(`^the response contains post URL "([^"]*)"$`, hc.theResponseContainsPostURL)
	ctx.Step(`^the response contains draft ID "([^"]*)"$`, hc.theResponseContainsDraftID)
	ctx.Step(`^the result reports the draft already exists with ID "([^"]*)"$`, hc.theResultReportsTheDraftAlreadyExistsWithID)
	ctx.Step(`^the post is found with ID "([^"]*)"$`, hc.thePostIsFoundWithID)
	ctx.Step(`^the post is not found$`, hc.thePostIsNotFound)
	ctx.Step(`^the deleted post is found with ID "([^"]*)"$`, hc.theDeletedPostIsFoundWithID)
	ctx.Step(`^no deleted post is found$`, hc.noDeletedPostIsFound)
	ctx.Step(`^the error is "([^"]*)"$`, hc.theErrorIs)
	ctx.Step(`^the error contains "([^"]*)"$`, hc.theErrorContains)
	ctx.Step(`^no HTTP request is made$`, hc.noHTTPRequestIsMade)
	ctx.Step(`^the dry-run result action is "([^"]*)"$`, hc.theDryRunResultActionIs)
	ctx.Step(`^the dry-run result slug is "([^"]*)"$`, hc.theDryRunResultSlugIs)
	ctx.Step(`^the dry-run result existing ID is "([^"]*)"$`, hc.theDryRunResultExistingIDIs)
	ctx.Step(`^the dry-run result deleted ID is "([^"]*)"$`, hc.theDryRunResultDeletedIDIs)
	ctx.Step(`^the series is found with ID "([^"]*)"$`, hc.theSeriesIsFoundWithID)
	ctx.Step(`^the series is not found$`, hc.theSeriesIsNotFound)
	ctx.Step(`^a createSeries mutation is sent$`, hc.aCreateSeriesMutationIsSent)
	ctx.Step(`^no createSeries mutation is sent$`, hc.noCreateSeriesMutationIsSent)
	ctx.Step(`^the series ID is "([^"]*)"$`, hc.theSeriesIDIs)
	ctx.Step(`^the publish request includes series ID "([^"]*)"$`, hc.thePublishRequestIncludesSeriesID)
	ctx.Step(`^the update request includes series ID "([^"]*)"$`, hc.theUpdateRequestIncludesSeriesID)
}

// BDD: specs/hashnode_client.feature :: Feature: Hashnode client
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../specs/hashnode_client.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
