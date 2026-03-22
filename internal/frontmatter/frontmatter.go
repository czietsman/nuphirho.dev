// Package frontmatter parses YAML frontmatter from markdown post files.
package frontmatter

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Post holds the parsed frontmatter and body content of a markdown post file.
type Post struct {
	Title    string   `yaml:"title"`
	Slug     string   `yaml:"slug"`
	Subtitle string   `yaml:"subtitle"`
	Tags     []string `yaml:"tags"`
	Draft      bool     `yaml:"draft"`
	Series     string   `yaml:"series"`
	AllowEmdash bool   `yaml:"allow_emdash"`
	Content    string   `yaml:"-"`
}

// ValidationResult holds errors and warnings from frontmatter validation.
type ValidationResult struct {
	Errors   []string
	Warnings []string
}

// Passed returns true if there are no validation errors.
func (v *ValidationResult) Passed() bool {
	return len(v.Errors) == 0
}

var slugPattern = regexp.MustCompile(`^[a-z0-9-]+$`)

var secretPattern = regexp.MustCompile(`(?i)(api[_-]?key|secret|token|password|credential)\s*[:=]\s*["']?[A-Za-z0-9+/_-]{20,}`)

// Parse reads a markdown file's content and returns the parsed Post and
// a ValidationResult. The returned Post is always populated with whatever
// could be parsed; callers should check the ValidationResult before using it.
func Parse(raw string) (*Post, *ValidationResult) {
	result := &ValidationResult{}
	post := &Post{}

	// Split frontmatter from content
	fm, content, err := splitFrontmatter(raw)
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
		return post, result
	}

	// Parse YAML
	if err := yaml.Unmarshal([]byte(fm), post); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("invalid frontmatter YAML: %v", err))
		return post, result
	}

	// Strip leading H1 if it matches the title
	content = stripLeadingH1(content, post.Title)
	post.Content = content

	// Validate required fields
	validate(post, result)

	// Content checks
	if strings.Contains(post.Content, "\u2014") && !post.AllowEmdash {
		result.Errors = append(result.Errors, "em dash detected")
	}

	if secretPattern.MatchString(post.Content) {
		result.Errors = append(result.Errors, "potential secret detected in post content")
	}

	return post, result
}

func splitFrontmatter(raw string) (frontmatter, content string, err error) {
	// Frontmatter must start with ---
	trimmed := strings.TrimLeft(raw, "\n")
	if !strings.HasPrefix(trimmed, "---") {
		return "", "", fmt.Errorf("missing frontmatter")
	}

	// Find the closing ---
	rest := trimmed[3:]
	rest = strings.TrimLeft(rest, " \t")
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	} else if len(rest) > 1 && rest[0:2] == "\r\n" {
		rest = rest[2:]
	}

	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return "", "", fmt.Errorf("missing frontmatter")
	}

	fm := rest[:idx]
	remaining := rest[idx+4:] // skip \n---

	// Skip the rest of the closing --- line
	if nlIdx := strings.IndexByte(remaining, '\n'); nlIdx >= 0 {
		remaining = remaining[nlIdx+1:]
	} else {
		remaining = ""
	}

	// Trim leading blank lines from content
	remaining = strings.TrimLeft(remaining, "\n")

	return fm, strings.TrimRight(remaining, "\n \t"), nil
}

func stripLeadingH1(content, title string) string {
	lines := strings.SplitN(content, "\n", 2)
	if len(lines) == 0 {
		return content
	}

	firstLine := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(firstLine, "# ") {
		return content
	}

	h1Text := strings.TrimSpace(firstLine[2:])
	if h1Text != title {
		return content
	}

	// Strip the H1 and any immediately following blank lines
	if len(lines) == 1 {
		return ""
	}
	return strings.TrimLeft(lines[1], "\n")
}

func validate(post *Post, result *ValidationResult) {
	if post.Title == "" {
		result.Errors = append(result.Errors, "missing required field: title")
	}
	if post.Slug == "" {
		result.Errors = append(result.Errors, "missing required field: slug")
	} else if !slugPattern.MatchString(post.Slug) {
		result.Errors = append(result.Errors, "slug must be lowercase and hyphenated")
	}
	if len(post.Tags) == 0 {
		result.Errors = append(result.Errors, "missing required field: tags")
	}
}
