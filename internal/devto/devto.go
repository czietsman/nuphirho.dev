// Package devto provides a client for the Dev.to REST API.
package devto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
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
}

// PublishResult holds the outcome of a create or update operation.
type PublishResult struct {
	Action          string // "create", "update"
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
}

type requestRecord struct {
	Method string
	Path   string
	Body   map[string]interface{}
}

// New creates a new Dev.to client.
func New(apiKey string, httpClient HTTPClient) *Client {
	return &Client{
		APIKey:   apiKey,
		Endpoint: "https://dev.to",
		HTTP:     httpClient,
	}
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
	existingID, err := c.findByCanonicalURL(canonicalURL)
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"article": map[string]interface{}{
			"title":         input.Title,
			"body_markdown": content,
			"tags":          input.Tags,
			"published":     input.Published,
			"canonical_url": canonicalURL,
		},
	}

	if existingID > 0 {
		if c.DryRun {
			return &PublishResult{
				Action:          "update",
				ArticleID:       existingID,
				Published:       input.Published,
				EmbedsConverted: embedCount,
				DryRun:          true,
			}, nil
		}
		return c.updateArticle(existingID, body, embedCount)
	}

	if c.DryRun {
		return &PublishResult{
			Action:          "create",
			Published:       input.Published,
			EmbedsConverted: embedCount,
			DryRun:          true,
		}, nil
	}
	return c.createArticle(body, embedCount)
}

func (c *Client) findByCanonicalURL(canonicalURL string) (int, error) {
	req, err := http.NewRequest("GET", c.Endpoint+"/api/articles/me/all", nil)
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("api-key", c.APIKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "connection") {
			return 0, fmt.Errorf("connection error: %w", err)
		}
		return 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode == 401 {
		return 0, fmt.Errorf("authentication failed")
	}
	if resp.StatusCode == 429 {
		return 0, fmt.Errorf("rate limited")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var articles []map[string]interface{}
	if err := json.Unmarshal(respBody, &articles); err != nil {
		return 0, fmt.Errorf("parsing response: %w", err)
	}

	for _, article := range articles {
		if cu, _ := article["canonical_url"].(string); cu == canonicalURL {
			if id, ok := article["id"].(float64); ok {
				return int(id), nil
			}
		}
	}

	return 0, nil
}

func (c *Client) createArticle(body map[string]interface{}, embedCount int) (*PublishResult, error) {
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
		Published:       true,
		EmbedsConverted: embedCount,
	}, nil
}

func (c *Client) updateArticle(id int, body map[string]interface{}, embedCount int) (*PublishResult, error) {
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
		Published:       true,
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

	resp, err := c.HTTP.Do(req)
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
