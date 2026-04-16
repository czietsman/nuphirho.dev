// Package hashnode provides a client for the Hashnode GraphQL API.
package hashnode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// timeString formats a *time.Time for logging. Returns "nil" if t is nil.
func timeString(t *time.Time) string {
	if t == nil {
		return "nil"
	}
	return t.Format(time.RFC3339Nano)
}

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
	EditedAt *time.Time
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
	// Log receives diagnostic output when non-nil. Writes are best-effort.
	Log io.Writer
	// Now overrides the current time for testing. Defaults to time.Now.
	Now func() time.Time

	// mutationsSent tracks mutations for test assertions
	mutationsSent   []string
	requestCount    int
	publishedBySlug map[string]publishedPost
	deletedBySlug   map[string]string
}

func (c *Client) nowOrDefault() time.Time {
	if c.Now != nil {
		return c.Now()
	}
	return time.Now()
}

type publishedPost struct {
	ID          string
	PublishedAt *time.Time
	UpdatedAt   *time.Time
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

// Publish handles the full publish flow: check existing, check deleted, restore, or create.
func (c *Client) Publish(input PostInput) (*PublishResult, error) {
	canonicalURL := "https://blog.nuphirho.dev/" + input.Slug

	if err := c.loadPublishedInventory(); err != nil {
		return nil, err
	}
	if err := c.loadDeletedInventory(); err != nil {
		return nil, err
	}

	existing := c.publishedBySlug[input.Slug]
	if existing.ID != "" {
		if c.Log != nil {
			remote := latestRemoteTimestamp(existing.UpdatedAt, existing.PublishedAt)
			fmt.Fprintf(c.Log, "[hashnode] %s: existing post id=%s publishedAt=%s updatedAt=%s remote=%s editedAt=%s\n",
				input.Slug, existing.ID,
				timeString(existing.PublishedAt), timeString(existing.UpdatedAt),
				timeString(remote), timeString(input.EditedAt))
		}
		if !shouldUpdatePublishedPost(input.EditedAt, c.nowOrDefault()) {
			if c.Log != nil {
				fmt.Fprintf(c.Log, "[hashnode] %s: unchanged (editedAt not within 24h window)\n", input.Slug)
			}
			return &PublishResult{Action: "unchanged", PostID: existing.ID}, nil
		}
		if c.Log != nil {
			fmt.Fprintf(c.Log, "[hashnode] %s: updating post\n", input.Slug)
		}
		if c.DryRun {
			return &PublishResult{Action: "update", PostID: existing.ID, DryRun: true}, nil
		}
		url, err := c.updatePost(existing.ID, input)
		if err != nil {
			return nil, err
		}
		return &PublishResult{Action: "update", PostID: existing.ID, URL: url}, nil
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
	return c.publishedBySlug[slug].ID, nil
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

	query := `query($id: ObjectId!) { publication(id: $id) { posts(first: 50) { edges { node { id slug publishedAt updatedAt } } } } }`
	vars := map[string]string{"id": c.PublicationID}

	resp, err := c.doGraphQL(query, vars)
	if err != nil {
		return err
	}

	pub := resp.path("data", "publication")
	if pub.isNull() {
		return fmt.Errorf("publication not found")
	}

	c.publishedBySlug = edgesToPublishedPosts(resp.path("data", "publication", "posts", "edges"))
	if c.Log != nil {
		fmt.Fprintf(c.Log, "[hashnode] inventory: %d published posts\n", len(c.publishedBySlug))
		for slug, p := range c.publishedBySlug {
			fmt.Fprintf(c.Log, "[hashnode]   slug=%q id=%s publishedAt=%s updatedAt=%s\n",
				slug, p.ID, timeString(p.PublishedAt), timeString(p.UpdatedAt))
		}
	}
	return nil
}

func (c *Client) loadDeletedInventory() error {
	if c.deletedBySlug != nil {
		return nil
	}

	query := `query($id: ObjectId!) { publication(id: $id) { posts(first: 50, filter: { deletedOnly: true }) { edges { node { id slug } } } } }`
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
		c.publishedBySlug = make(map[string]publishedPost)
	}
	c.publishedBySlug[input.Slug] = publishedPost{ID: postID}
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
		c.publishedBySlug = make(map[string]publishedPost)
	}
	c.publishedBySlug[input.Slug] = publishedPost{ID: id}
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
					c.publishedBySlug = make(map[string]publishedPost)
				}
				c.publishedBySlug[slug] = publishedPost{ID: restoredID}
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

func edgesToPublishedPosts(edges *gqlNode) map[string]publishedPost {
	result := make(map[string]publishedPost)
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
			result[slug] = publishedPost{
				ID:          id,
				PublishedAt: parseGraphQLTime(node["publishedAt"]),
				UpdatedAt:   parseGraphQLTime(node["updatedAt"]),
			}
		}
	}

	return result
}

func parseGraphQLTime(v interface{}) *time.Time {
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

func shouldUpdatePublishedPost(editedAt *time.Time, now time.Time) bool {
	if editedAt == nil {
		return false
	}
	cutoff := now.UTC().Add(-24 * time.Hour)
	return editedAt.UTC().After(cutoff)
}

func latestRemoteTimestamp(values ...*time.Time) *time.Time {
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
