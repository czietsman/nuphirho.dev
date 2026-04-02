package frontmatter

import "testing"

func TestParseBlankLineAfterOpeningDelimiterParsesFields(t *testing.T) {
	raw := `---

title: "My Post"
slug: my-post
tags: [test]
---

Body.
`

	post, result := Parse(raw)

	if !result.Passed() {
		t.Fatalf("expected parse to succeed, got errors: %v", result.Errors)
	}
	if post.Title != "My Post" {
		t.Fatalf("expected title to be parsed after a leading blank line, got %q", post.Title)
	}
	if post.Content != "Body." {
		t.Fatalf("expected content %q, got %q", "Body.", post.Content)
	}
}

func TestParseWindowsLineEndingsNormalisesContent(t *testing.T) {
	raw := "---\r\n" +
		"title: \"Windows Post\"\r\n" +
		"slug: windows-post\r\n" +
		"tags: [test]\r\n" +
		"---\r\n" +
		"\r\n" +
		"# Windows Post\r\n" +
		"\r\n" +
		"Body on Windows.\r\n"

	post, result := Parse(raw)

	if !result.Passed() {
		t.Fatalf("expected parse to succeed, got errors: %v", result.Errors)
	}
	if post.Content != "Body on Windows." {
		t.Fatalf("expected normalised content %q, got %q", "Body on Windows.", post.Content)
	}
}

func TestParseDelimiterOnlyFrontmatterReturnsMissingFrontmatter(t *testing.T) {
	post, result := Parse("---")

	if post == nil {
		t.Fatal("expected a non-nil post result")
	}
	if len(result.Errors) != 1 || result.Errors[0] != "missing frontmatter" {
		t.Fatalf("expected missing frontmatter error, got %v", result.Errors)
	}
}

func TestParseOpeningDelimiterMayHaveTrailingSpaces(t *testing.T) {
	raw := `---   
title: "Spaced Delimiter"
slug: spaced-delimiter
tags: [test]
---

Content.
`

	post, result := Parse(raw)

	if !result.Passed() {
		t.Fatalf("expected parse to succeed, got errors: %v", result.Errors)
	}
	if post.Title != "Spaced Delimiter" {
		t.Fatalf("expected title %q, got %q", "Spaced Delimiter", post.Title)
	}
}

func TestParseNoContentAfterClosingDelimiterReturnsEmptyBody(t *testing.T) {
	raw := `---
title: "Header Only"
slug: header-only
tags: [test]
---`

	post, result := Parse(raw)

	if !result.Passed() {
		t.Fatalf("expected parse to succeed, got errors: %v", result.Errors)
	}
	if post.Content != "" {
		t.Fatalf("expected empty content, got %q", post.Content)
	}
}

func TestStripLeadingH1PreservesSingleLineBody(t *testing.T) {
	content := "Body only"

	if got := stripLeadingH1(content, "Any Title"); got != content {
		t.Fatalf("expected content %q to be preserved, got %q", content, got)
	}
}
