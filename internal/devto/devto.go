// Package devto provides a client for the Dev.to REST API.
package devto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// HTTPClient is the interface for making HTTP requests.
// Inject a fake in tests; use http.DefaultClient in production.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ArticleInput holds the data needed to create or update an article.
type ArticleInput struct {
	Title     string
	Slug      string
	Content   string
	Tags      []string
	Published bool
	Series    string
	EditedAt  *time.Time
}

// PublishResult holds the outcome of a create or update operation.
type PublishResult struct {
	Action          string // "create", "update", "unchanged"
	ArticleID       int
	URL             string
	Published       bool
	EmbedsConverted int
	DryRun          bool
}

// Client is the Dev.to REST API client.
type Client struct {
	APIKey   string
	Endpoint string
	HTTP     HTTPClient
	DryRun   bool

	// requestsMade tracks requests for test assertions
	requestsMade []requestRecord
	inventory    map[string]articleRecord
	sleepFn      func(time.Duration)
}

const maxRetries = 3

type requestRecord struct {
	Method string
	Path   string
	Body   map[string]interface{}
}

type articleRecord struct {
	ID           int
	URL          string
	Published    bool
	CanonicalURL string
	PublishedAt  *time.Time
	EditedAt     *time.Time
}

// New creates a new Dev.to client.
func New(apiKey string, httpClient HTTPClient) *Client {
	return &Client{
		APIKey:   apiKey,
		Endpoint: "https://dev.to",
		HTTP:     httpClient,
		sleepFn:  time.Sleep,
	}
}

// SetSleepFn overrides the sleep function (for testing).
func (c *Client) SetSleepFn(fn func(time.Duration)) {
	c.sleepFn = fn
}

var embedRe = regexp.MustCompile(`(?m)^%\[(.+?)\]$`)

// convertEmbeds replaces Hashnode %[url] embeds with Dev.to {% embed url %} Liquid tags.
// Returns the converted content and the number of conversions made.
func convertEmbeds(content string) (string, int) {
	count := 0
	result := embedRe.ReplaceAllStringFunc(content, func(match string) string {
		count++
		sub := embedRe.FindStringSubmatch(match)
		return "{% embed " + sub[1] + " %}"
	})
	return result, count
}

// CreateArticle handles the full create/update flow:
// 1. Convert embeds
// 2. Look up existing article by canonical URL
// 3. Update if found, create if not
func (c *Client) CreateArticle(input ArticleInput) (*PublishResult, error) {
	canonicalURL := "https://blog.nuphirho.dev/" + input.Slug
	content, embedCount := convertEmbeds(input.Content)

	// Look up existing article by canonical URL
	existing, err := c.findByCanonicalURL(canonicalURL)
	if err != nil {
		return nil, err
	}

	articleBody := map[string]interface{}{
		"title":         input.Title,
		"body_markdown": content,
		"tags":          input.Tags,
		"published":     input.Published,
		"canonical_url": canonicalURL,
	}
	if input.Series != "" {
		articleBody["series"] = input.Series
	}
	body := map[string]interface{}{"article": articleBody}

	if existing != nil {
		if existing.Published == input.Published && !shouldUpdateArticle(existing, input.EditedAt) {
			return &PublishResult{
				Action:          "unchanged",
				ArticleID:       existing.ID,
				URL:             existing.URL,
				Published:       input.Published,
				EmbedsConverted: embedCount,
			}, nil
		}
		if c.DryRun {
			return &PublishResult{
				Action:          "update",
				ArticleID:       existing.ID,
				Published:       input.Published,
				EmbedsConverted: embedCount,
				DryRun:          true,
			}, nil
		}
		return c.updateArticle(existing.ID, body, embedCount, input.Published)
	}

	if c.DryRun {
		return &PublishResult{
			Action:          "create",
			Published:       input.Published,
			EmbedsConverted: embedCount,
			DryRun:          true,
		}, nil
	}
	return c.createArticle(body, embedCount, input.Published)
}

func (c *Client) findByCanonicalURL(canonicalURL string) (*articleRecord, error) {
	if err := c.loadInventory(); err != nil {
		return nil, err
	}
	if article, ok := c.inventory[canonicalURL]; ok {
		copy := article
		return &copy, nil
	}
	return nil, nil
}

func (c *Client) loadInventory() error {
	if c.inventory != nil {
		return nil
	}

	req, err := http.NewRequest("GET", c.Endpoint+"/api/articles/me/all", nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("api-key", c.APIKey)

	resp, err := c.doWithRetry(req)
	if err != nil {
		if strings.Contains(err.Error(), "connection") {
			return fmt.Errorf("connection error: %w", err)
		}
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode == 401 {
		return fmt.Errorf("authentication failed")
	}
	if resp.StatusCode == 429 {
		return fmt.Errorf("rate limited")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var articles []map[string]interface{}
	if err := json.Unmarshal(respBody, &articles); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	c.inventory = make(map[string]articleRecord, len(articles))
	for _, article := range articles {
		cu, _ := article["canonical_url"].(string)
		if cu == "" {
			continue
		}
		record := articleRecord{
			URL:          stringValue(article["url"]),
			Published:    boolValue(article["published"]),
			CanonicalURL: cu,
			PublishedAt:  parseOptionalTime(article["published_at"]),
			EditedAt:     parseOptionalTime(article["edited_at"]),
		}
		if id, ok := article["id"].(float64); ok {
			record.ID = int(id)
		}
		c.inventory[cu] = record
	}

	return nil
}

func (c *Client) createArticle(body map[string]interface{}, embedCount int, published bool) (*PublishResult, error) {
	resp, err := c.doRequest("POST", "/api/articles", body)
	if err != nil {
		return nil, err
	}

	id, _ := resp["id"].(float64)
	url, _ := resp["url"].(string)
	return &PublishResult{
		Action:          "create",
		ArticleID:       int(id),
		URL:             url,
		Published:       published,
		EmbedsConverted: embedCount,
	}, nil
}

func (c *Client) updateArticle(id int, body map[string]interface{}, embedCount int, published bool) (*PublishResult, error) {
	path := fmt.Sprintf("/api/articles/%d", id)
	resp, err := c.doRequest("PUT", path, body)
	if err != nil {
		return nil, err
	}

	url, _ := resp["url"].(string)
	return &PublishResult{
		Action:          "update",
		ArticleID:       id,
		URL:             url,
		Published:       published,
		EmbedsConverted: embedCount,
	}, nil
}

func (c *Client) doRequest(method, path string, body map[string]interface{}) (map[string]interface{}, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	req, err := http.NewRequest(method, c.Endpoint+path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.APIKey)

	c.requestsMade = append(c.requestsMade, requestRecord{Method: method, Path: path, Body: body})

	resp, err := c.doWithRetry(req)
	if err != nil {
		if strings.Contains(err.Error(), "connection") {
			return nil, fmt.Errorf("connection error: %w", err)
		}
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("authentication failed")
	}
	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limited")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return result, nil
}

// RequestsMade returns the request records for test assertions.
func (c *Client) RequestsMade() []requestRecord {
	return c.requestsMade
}

func (c *Client) doWithRetry(req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("reading request body: %w", err)
		}
	}

	var resp *http.Response
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		var err error
		resp, err = c.HTTP.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 429 {
			return resp, nil
		}

		if attempt == maxRetries {
			return resp, nil
		}

		delay := c.retryDelay(resp, attempt)
		resp.Body.Close()
		c.sleepFn(delay)
	}

	return resp, nil
}

func (c *Client) retryDelay(resp *http.Response, attempt int) time.Duration {
	if ra := resp.Header.Get("Retry-After"); ra != "" {
		if seconds, err := strconv.Atoi(ra); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return time.Duration(1<<uint(attempt)) * time.Second
}

func shouldUpdateArticle(a *articleRecord, editedAt *time.Time) bool {
	if editedAt == nil {
		return false
	}
	remote := latestTime(a.EditedAt, a.PublishedAt)
	if remote == nil {
		return true
	}
	return editedAt.UTC().After(remote.UTC())
}

func stringValue(v interface{}) string {
	s, _ := v.(string)
	return s
}

func boolValue(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

func parseOptionalTime(v interface{}) *time.Time {
	s, _ := v.(string)
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	utc := t.UTC()
	return &utc
}

func latestTime(values ...*time.Time) *time.Time {
	var latest *time.Time
	for _, value := range values {
		if value == nil {
			continue
		}
		if latest == nil || value.After(*latest) {
			copy := value.UTC()
			latest = &copy
		}
	}
	return latest
}
