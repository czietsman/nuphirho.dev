package main

import (
	"os"
	"strings"
	"testing"
)

func TestDependencyReviewPromptTemplatesExist(t *testing.T) {
	t.Parallel()

	requiredFiles := []string{
		"prompts/dependency-review/README.md",
		"prompts/dependency-review/go-dependencies.md",
		"prompts/dependency-review/github-actions.md",
		"prompts/dependency-review/npm-dependencies.md",
		"prompts/dependency-review/terraform-providers.md",
		"prompts/dependency-review/tools.md",
	}

	for _, path := range requiredFiles {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
}

func TestReadmeDocumentsDependencyReviewPrompts(t *testing.T) {
	t.Parallel()

	content, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}

	required := []string{
		"**prompts/** -- Reviewed prompt material, including dependency review briefs",
		"`prompts/dependency-review/` contains reviewed research briefs",
	}

	for _, fragment := range required {
		if !strings.Contains(string(content), fragment) {
			t.Fatalf("README.md missing %q", fragment)
		}
	}
}
