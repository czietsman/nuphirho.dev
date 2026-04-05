package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/czietsman/nuphirho.dev/internal/frontmatter"
	"github.com/czietsman/nuphirho.dev/internal/tags"
)

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}

// Run validates that all hyphenated tags in the given post files have
// Dev.to mappings in the tag glossary. Draft posts are skipped.
func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("validate-tags", flag.ContinueOnError)
	fs.SetOutput(stderr)

	tagsFile := fs.String("tags-file", "tags.json", "Path to tag glossary JSON file")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	postPaths := fs.Args()
	if len(postPaths) == 0 {
		return 0
	}

	glossary, err := tags.LoadGlossary(*tagsFile)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}

	failed := false
	for _, path := range postPaths {
		unmapped := validatePost(path, glossary, stderr)
		if unmapped {
			failed = true
		}
	}

	if failed {
		return 1
	}
	return 0
}

func validatePost(path string, glossary tags.Glossary, stderr io.Writer) bool {
	raw, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(stderr, "error: reading %s: %s\n", path, err)
		return true
	}

	post, result := frontmatter.Parse(string(raw))
	if !result.Passed() {
		return false
	}

	r := glossary.ValidatePostTags(post.Tags, post.Draft)
	if r.Skipped || r.Valid {
		return false
	}

	for _, tag := range r.Unmapped {
		fmt.Fprintf(stderr, "%s: unmapped Dev.to tag %q (add a mapping to tags.json)\n", path, tag)
	}
	return true
}
