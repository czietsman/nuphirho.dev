// Package hashnode provides a client for the Hashnode GraphQL API.
package hashnode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// HTTPClient is the interface for making HTTP requests.
// Inject a fake in tests; use http.DefaultClient in production.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// PostInput holds the data needed to publish or update a post.
type PostInput struct {
	Title    string
	Slug     string
	Subtitle string
	Content  string
	Tags     []string
	SeriesID string
}

// PublishResult holds the outcome of a publish or draft operation.
type PublishResult struct {
	Action    string // "publish", "update", "restore_and_update", "create_draft", "draft_exists"
	PostID    string
	DraftID   string
	DeletedID string
	URL       string
	DryRun    bool
}

// Client is the Hashnode GraphQL API client.
type Client struct {
	Token         string
	PublicationID string
	Endpoint      string
	HTTP          HTTPClient
	DryRun        bool

	// mutationsSent tracks mutations for test assertions
	mutationsSent   []string
	requestCount    int
	publishedBySlug map[string]string
	deletedBySlug   map[string]string
}

// New creates a new Hashnode client.
func New(token, publicationID string, httpClient HTTPClient) *Client {
	return &Client{
		Token:         token,
		PublicationID: publicationID,
		Endpoint:      "https://gql.hashnode.com",
		HTTP:          httpClient,
	}
}

// existingPost holds the data fetched from Hashnode for content comparison.
type existingPost struct {
	Title    string
	Subtitle string
	Content  string
}

// fetchExistingPost retrieves the current content of a published post by slug.
func (c *Client) fetchExistingPost(slug string) (*existingPost, error) {
	query := `query($id: ObjectId!, $slug: String!) { publication(id: $id) { post(slug: $slug) { title subtitle content { markdown } } } }`
	vars := map[string]string{"id": c.PublicationID, "slug": slug}

	resp, err := c.doGraphQL(query, vars)
	if err != nil {
		return nil, err
	}

	post := resp.path("data", "publication", "post")
	if post.isNull() {
		return nil, nil
	}

	title := post.str("title")
	subtitle := post.str("subtitle")
	contentNode := post.path("content")
	content := ""
	if !contentNode.isNull() {
		content = contentNode.str("markdown")
	}

	return &existingPost{Title: title, Subtitle: subtitle, Content: content}, nil
}

// Publish handles the full publish flow: check existing, check deleted, restore, or create.
func (c *Client) Publish(input PostInput) (*PublishResult, error) {
	canonicalURL := "https://blog.nuphirho.dev/" + input.Slug

	if err := c.loadPublishedInventory(); err != nil {
		return nil, err
	}
	if err := c.loadDeletedInventory(); err != nil {
		return nil, err
	}

	existingID := c.publishedBySlug[input.Slug]
	if existingID != "" {
		if c.DryRun {
			return &PublishResult{Action: "update", PostID: existingID, DryRun: true}, nil
		}

		// Fetch existing content and compare before updating
		existing, err := c.fetchExistingPost(input.Slug)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.Title == input.Title && existing.Subtitle == input.Subtitle && existing.Content == input.Content {
			return &PublishResult{Action: "unchanged", PostID: existingID}, nil
		}

		url, err := c.updatePost(existingID, input)
		if err != nil {
			return nil, err
		}
		return &PublishResult{Action: "update", PostID: existingID, URL: url}, nil
	}

	deletedID := c.deletedBySlug[input.Slug]
	if deletedID != "" {
		if c.DryRun {
			return &PublishResult{Action: "restore_and_update", DeletedID: deletedID, DryRun: true}, nil
		}
		restoredID, err := c.restorePost(deletedID)
		if err != nil {
			return nil, err
		}
		url, err := c.updatePost(restoredID, input)
		if err != nil {
			return nil, err
		}
		return &PublishResult{Action: "restore_and_update", PostID: restoredID, DeletedID: deletedID, URL: url}, nil
	}

	// Create new post
	if c.DryRun {
		return &PublishResult{Action: "publish", DryRun: true}, nil
	}
	postID, url, err := c.publishPost(input, canonicalURL)
	if err != nil {
		return nil, err
	}
	return &PublishResult{Action: "publish", PostID: postID, URL: url}, nil
}

// CreateDraft handles the draft flow: check existing draft, create if not found.
func (c *Client) CreateDraft(input PostInput) (*PublishResult, error) {
	canonicalURL := "https://blog.nuphirho.dev/" + input.Slug

	// Check for existing draft
	draftID, err := c.checkDraftBySlug(input.Slug)
	if err != nil {
		return nil, err
	}

	if draftID != "" {
		return &PublishResult{Action: "draft_exists", DraftID: draftID}, nil
	}

	if c.DryRun {
		return &PublishResult{Action: "create_draft", DryRun: true}, nil
	}

	id, err := c.createDraft(input, canonicalURL)
	if err != nil {
		return nil, err
	}
	return &PublishResult{Action: "create_draft", DraftID: id}, nil
}

// CheckPostBySlug queries for a published post by slug. Returns the post ID or empty string.
func (c *Client) CheckPostBySlug(slug string) (string, error) {
	if err := c.loadPublishedInventory(); err != nil {
		return "", err
	}
	return c.publishedBySlug[slug], nil
}

// CheckDeletedBySlug queries for deleted posts matching the given slug.
func (c *Client) CheckDeletedBySlug(slug string) (string, error) {
	if err := c.loadDeletedInventory(); err != nil {
		return "", err
	}
	return c.deletedBySlug[slug], nil
}

func (c *Client) loadPublishedInventory() error {
	if c.publishedBySlug != nil {
		return nil
	}

	query := `query($id: ObjectId!) { publication(id: $id) { posts(first: 100) { edges { node { id slug } } } } }`
	vars := map[string]string{"id": c.PublicationID}

	resp, err := c.doGraphQL(query, vars)
	if err != nil {
		return err
	}

	pub := resp.path("data", "publication")
	if pub.isNull() {
		return fmt.Errorf("publication not found")
	}

	c.publishedBySlug = edgesToSlugMap(resp.path("data", "publication", "posts", "edges"))
	return nil
}

func (c *Client) loadDeletedInventory() error {
	if c.deletedBySlug != nil {
		return nil
	}

	query := `query($id: ObjectId!) { publication(id: $id) { posts(first: 100, filter: { deletedOnly: true }) { edges { node { id slug } } } } }`
	vars := map[string]string{"id": c.PublicationID}

	resp, err := c.doGraphQL(query, vars)
	if err != nil {
		return err
	}

	pub := resp.path("data", "publication")
	if pub.isNull() {
		return fmt.Errorf("publication not found")
	}

	c.deletedBySlug = edgesToSlugMap(resp.path("data", "publication", "posts", "edges"))
	return nil
}

func (c *Client) checkDraftBySlug(slug string) (string, error) {
	query := `query($id: ObjectId!) { publication(id: $id) { drafts(first: 50) { edges { node { id slug } } } } }`
	vars := map[string]string{"id": c.PublicationID}

	resp, err := c.doGraphQL(query, vars)
	if err != nil {
		return "", err
	}

	edges := resp.path("data", "publication", "drafts", "edges")
	if edges == nil {
		return "", nil
	}

	edgesList, ok := edges.raw.([]interface{})
	if !ok {
		return "", nil
	}

	for _, edge := range edgesList {
		edgeMap, ok := edge.(map[string]interface{})
		if !ok {
			continue
		}
		node, ok := edgeMap["node"].(map[string]interface{})
		if !ok {
			continue
		}
		if nodeSlug, _ := node["slug"].(string); nodeSlug == slug {
			if id, _ := node["id"].(string); id != "" {
				return id, nil
			}
		}
	}

	return "", nil
}

func (c *Client) publishPost(input PostInput, canonicalURL string) (string, string, error) {
	tagsJSON := buildTagsJSON(input.Tags)
	query := `mutation($input: PublishPostInput!) { publishPost(input: $input) { post { id url } } }`
	inputMap := map[string]interface{}{
		"publicationId":      c.PublicationID,
		"title":              input.Title,
		"subtitle":           input.Subtitle,
		"slug":               input.Slug,
		"contentMarkdown":    input.Content,
		"tags":               tagsJSON,
		"originalArticleURL": canonicalURL,
	}
	if input.SeriesID != "" {
		inputMap["seriesId"] = input.SeriesID
	}
	vars := map[string]interface{}{"input": inputMap}

	c.mutationsSent = append(c.mutationsSent, "publishPost")
	resp, err := c.doGraphQLRaw(query, vars)
	if err != nil {
		return "", "", err
	}

	postID := resp.path("data", "publishPost", "post").str("id")
	postURL := resp.path("data", "publishPost", "post").str("url")
	if postID == "" {
		return "", "", fmt.Errorf("publishPost returned no post ID: %v", resp.raw)
	}
	if c.publishedBySlug == nil {
		c.publishedBySlug = make(map[string]string)
	}
	c.publishedBySlug[input.Slug] = postID
	if c.deletedBySlug != nil {
		delete(c.deletedBySlug, input.Slug)
	}

	return postID, postURL, nil
}

func (c *Client) updatePost(id string, input PostInput) (string, error) {
	tagsJSON := buildTagsJSON(input.Tags)
	query := `mutation($input: UpdatePostInput!) { updatePost(input: $input) { post { url } } }`
	inputMap := map[string]interface{}{
		"id":              id,
		"title":           input.Title,
		"subtitle":        input.Subtitle,
		"slug":            input.Slug,
		"contentMarkdown": input.Content,
		"tags":            tagsJSON,
	}
	if input.SeriesID != "" {
		inputMap["seriesId"] = input.SeriesID
	}
	vars := map[string]interface{}{"input": inputMap}

	c.mutationsSent = append(c.mutationsSent, "updatePost:"+id)
	resp, err := c.doGraphQLRaw(query, vars)
	if err != nil {
		return "", err
	}

	url := resp.path("data", "updatePost", "post").str("url")
	if c.publishedBySlug == nil {
		c.publishedBySlug = make(map[string]string)
	}
	c.publishedBySlug[input.Slug] = id
	return url, nil
}

func (c *Client) restorePost(id string) (string, error) {
	query := `mutation($input: RestorePostInput!) { restorePost(input: $input) { post { id slug } } }`
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"id": id,
		},
	}

	c.mutationsSent = append(c.mutationsSent, "restorePost:"+id)
	resp, err := c.doGraphQLRaw(query, vars)
	if err != nil {
		return "", err
	}

	restoredID := resp.path("data", "restorePost", "post").str("id")
	if restoredID == "" {
		return "", fmt.Errorf("restorePost returned no post ID: %v", resp.raw)
	}
	if c.deletedBySlug != nil {
		for slug, deletedID := range c.deletedBySlug {
			if deletedID == id {
				delete(c.deletedBySlug, slug)
				if c.publishedBySlug == nil {
					c.publishedBySlug = make(map[string]string)
				}
				c.publishedBySlug[slug] = restoredID
				break
			}
		}
	}

	return restoredID, nil
}

func (c *Client) createDraft(input PostInput, canonicalURL string) (string, error) {
	tagsJSON := buildTagsJSON(input.Tags)
	query := `mutation($input: CreateDraftInput!) { createDraft(input: $input) { draft { id slug } } }`
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"publicationId":      c.PublicationID,
			"title":              input.Title,
			"subtitle":           input.Subtitle,
			"slug":               input.Slug,
			"contentMarkdown":    input.Content,
			"tags":               tagsJSON,
			"originalArticleURL": canonicalURL,
			"settings":           map[string]interface{}{"slugOverridden": true},
		},
	}

	c.mutationsSent = append(c.mutationsSent, "createDraft")
	resp, err := c.doGraphQLRaw(query, vars)
	if err != nil {
		return "", err
	}

	draftID := resp.path("data", "createDraft", "draft").str("id")
	if draftID == "" {
		return "", fmt.Errorf("createDraft returned no draft ID: %v", resp.raw)
	}

	return draftID, nil
}

// LookupSeries finds a series by name. Returns the series ID or empty string.
func (c *Client) LookupSeries(name string) (string, error) {
	query := `query($id: ObjectId!) { publication(id: $id) { seriesList(first: 20) { edges { node { id name } } } } }`
	vars := map[string]string{"id": c.PublicationID}

	resp, err := c.doGraphQL(query, vars)
	if err != nil {
		return "", err
	}

	edges := resp.path("data", "publication", "seriesList", "edges")
	if edges == nil {
		return "", nil
	}

	edgesList, ok := edges.raw.([]interface{})
	if !ok {
		return "", nil
	}

	for _, edge := range edgesList {
		edgeMap, ok := edge.(map[string]interface{})
		if !ok {
			continue
		}
		node, ok := edgeMap["node"].(map[string]interface{})
		if !ok {
			continue
		}
		if nodeName, _ := node["name"].(string); nodeName == name {
			if id, _ := node["id"].(string); id != "" {
				return id, nil
			}
		}
	}
	return "", nil
}

// CreateSeries creates a new series and returns its ID.
func (c *Client) CreateSeries(name string) (string, error) {
	query := `mutation($input: CreateSeriesInput!) { createSeries(input: $input) { series { id } } }`
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"publicationId": c.PublicationID,
			"name":          name,
			"slug":          slugify(name),
		},
	}

	c.mutationsSent = append(c.mutationsSent, "createSeries")
	resp, err := c.doGraphQLRaw(query, vars)
	if err != nil {
		return "", err
	}

	seriesID := resp.path("data", "createSeries", "series").str("id")
	if seriesID == "" {
		return "", fmt.Errorf("createSeries returned no series ID: %v", resp.raw)
	}
	return seriesID, nil
}

// ResolveSeriesID looks up a series by name, creating it if not found.
func (c *Client) ResolveSeriesID(name string) (string, error) {
	if name == "" {
		return "", nil
	}
	id, err := c.LookupSeries(name)
	if err != nil {
		return "", err
	}
	if id != "" {
		return id, nil
	}
	return c.CreateSeries(name)
}

func slugify(name string) string {
	s := strings.ToLower(name)
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		if r == ' ' {
			return '-'
		}
		return -1
	}, s)
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}

// MutationsSent returns the list of mutations sent (for test assertions).
func (c *Client) MutationsSent() []string {
	return c.mutationsSent
}

// RequestCount returns the number of HTTP requests made.
func (c *Client) RequestCount() int {
	return c.requestCount
}

func buildTagsJSON(tags []string) []map[string]string {
	result := make([]map[string]string, len(tags))
	for i, t := range tags {
		result[i] = map[string]string{"slug": t, "name": t}
	}
	return result
}

func edgesToSlugMap(edges *gqlNode) map[string]string {
	result := make(map[string]string)
	if edges == nil {
		return result
	}

	edgesList, ok := edges.raw.([]interface{})
	if !ok {
		return result
	}

	for _, edge := range edgesList {
		edgeMap, ok := edge.(map[string]interface{})
		if !ok {
			continue
		}
		node, ok := edgeMap["node"].(map[string]interface{})
		if !ok {
			continue
		}
		slug, _ := node["slug"].(string)
		id, _ := node["id"].(string)
		if slug != "" && id != "" {
			result[slug] = id
		}
	}

	return result
}

// graphQL response navigation helpers

type gqlNode struct {
	raw interface{}
}

func (n *gqlNode) path(keys ...string) *gqlNode {
	if n == nil || n.raw == nil {
		return &gqlNode{}
	}
	current := n.raw
	for _, key := range keys {
		m, ok := current.(map[string]interface{})
		if !ok {
			return &gqlNode{}
		}
		current = m[key]
	}
	return &gqlNode{raw: current}
}

func (n *gqlNode) str(key string) string {
	if n == nil || n.raw == nil {
		return ""
	}
	m, ok := n.raw.(map[string]interface{})
	if !ok {
		return ""
	}
	s, _ := m[key].(string)
	return s
}

func (n *gqlNode) isNull() bool {
	return n == nil || n.raw == nil
}

func (c *Client) doGraphQL(query string, vars map[string]string) (*gqlNode, error) {
	varsIface := make(map[string]interface{}, len(vars))
	for k, v := range vars {
		varsIface[k] = v
	}
	return c.doGraphQLRaw(query, varsIface)
}

func (c *Client) doGraphQLRaw(query string, vars map[string]interface{}) (*gqlNode, error) {
	payload := map[string]interface{}{
		"query":     query,
		"variables": vars,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshalling GraphQL request: %w", err)
	}

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.Token)

	c.requestCount++
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

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	// Check for GraphQL errors
	if errs, ok := result["errors"]; ok {
		errList, ok := errs.([]interface{})
		if ok && len(errList) > 0 {
			firstErr, _ := errList[0].(map[string]interface{})
			msg, _ := firstErr["message"].(string)
			if strings.Contains(strings.ToLower(msg), "authenticated") {
				return nil, fmt.Errorf("authentication failed")
			}
			return nil, fmt.Errorf("GraphQL error: %s", msg)
		}
	}

	return &gqlNode{raw: result}, nil
}
