// Package pipeline orchestrates the publish flow: validate frontmatter,
// map tags, publish to Hashnode, cross-post to Dev.to, generate summary.
package pipeline

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/czietsman/nuphirho.dev/internal/devto"
	"github.com/czietsman/nuphirho.dev/internal/frontmatter"
	"github.com/czietsman/nuphirho.dev/internal/hashnode"
	"github.com/czietsman/nuphirho.dev/internal/tags"
)

// HashnodePub is the interface for Hashnode publish operations.
type HashnodePub interface {
	Publish(input hashnode.PostInput) (*hashnode.PublishResult, error)
}

// DevtoPub is the interface for Dev.to cross-post operations.
type DevtoPub interface {
	CreateArticle(input devto.ArticleInput) (*devto.PublishResult, error)
}

// Prober is the interface for contract probes.
type Prober interface {
	Run(w io.Writer) int
}

// PostFile represents a markdown file to process.
type PostFile struct {
	Path    string
	Content string
}

// Config holds the pipeline configuration.
type Config struct {
	Hashnode HashnodePub
	DevTo    DevtoPub
	Glossary tags.Glossary
	DryRun   bool
}

// FileResult holds the outcome of processing a single file.
type FileResult struct {
	File        string
	Post        *frontmatter.Post
	Validation  *frontmatter.ValidationResult
	HashnodeRes *hashnode.PublishResult
	DevToRes    *devto.PublishResult
	DevToTags   []string
	TagWarnings []string
	Skipped     bool
	SkipReason  string
	Error       string
}

// RunResult holds the full pipeline outcome.
type RunResult struct {
	Files    []*FileResult
	ExitCode int
}

// Probe runs contract probes and returns the exit code.
func Probe(prober Prober, w io.Writer) int {
	return prober.Run(w)
}

// Run executes the pipeline for the given post files.
func Run(cfg Config, files []PostFile, w io.Writer) *RunResult {
	result := &RunResult{}

	// Phase 1: Parse and validate all files
	parsed := make([]*FileResult, len(files))
	validationFailed := false

	for i, f := range files {
		post, val := frontmatter.Parse(f.Content)
		fr := &FileResult{
			File:       f.Path,
			Post:       post,
			Validation: val,
		}
		parsed[i] = fr
		if !val.Passed() {
			validationFailed = true
		}
	}

	// If any file fails validation, exit 1 before any publishing
	if validationFailed {
		result.Files = parsed
		result.ExitCode = 1
		if cfg.DryRun {
			writeDryRunJSON(w, parsed)
		} else {
			writeSummary(w, parsed)
		}
		return result
	}

	// Phase 2: Process each file independently
	publishFailed := false

	for _, fr := range parsed {
		// Skip drafts
		if fr.Post.Draft {
			fr.Skipped = true
			fr.SkipReason = "draft"
			continue
		}

		// Map tags
		hnMapResult := cfg.Glossary.MapTags(fr.Post.Tags, tags.Hashnode)
		dtMapResult := cfg.Glossary.MapTags(fr.Post.Tags, tags.DevTo)
		fr.DevToTags = dtMapResult.Tags
		fr.TagWarnings = append(fr.TagWarnings, hnMapResult.Warnings...)
		fr.TagWarnings = append(fr.TagWarnings, dtMapResult.Warnings...)

		// Publish to Hashnode
		hnInput := hashnode.PostInput{
			Title:    fr.Post.Title,
			Slug:     fr.Post.Slug,
			Subtitle: fr.Post.Subtitle,
			Content:  fr.Post.Content,
			Tags:     hnMapResult.Tags,
		}

		hnResult, err := cfg.Hashnode.Publish(hnInput)
		if err != nil {
			fr.Error = fmt.Sprintf("hashnode: %s", err.Error())
			publishFailed = true
			continue // skip Dev.to for this file
		}
		fr.HashnodeRes = hnResult

		// Cross-post to Dev.to
		dtInput := devto.ArticleInput{
			Title:     fr.Post.Title,
			Slug:      fr.Post.Slug,
			Content:   fr.Post.Content,
			Tags:      dtMapResult.Tags,
			Published: true,
		}

		dtResult, err := cfg.DevTo.CreateArticle(dtInput)
		if err != nil {
			fr.Error = fmt.Sprintf("devto: %s", err.Error())
			publishFailed = true
			continue
		}
		fr.DevToRes = dtResult
	}

	result.Files = parsed
	if publishFailed {
		result.ExitCode = 2
	}

	// Phase 3: Output
	if cfg.DryRun {
		writeDryRunJSON(w, parsed)
	} else {
		writeSummary(w, parsed)
	}

	return result
}

// --- Dry-run JSON output ---

type dryRunEntry struct {
	File        string        `json:"file"`
	Frontmatter dryRunFM     `json:"frontmatter"`
	Validation  dryRunVal    `json:"validation"`
	Hashnode    *dryRunHN    `json:"hashnode"`
	DevTo       *dryRunDT    `json:"devto"`
}

type dryRunFM struct {
	Title string   `json:"title"`
	Slug  string   `json:"slug"`
	Tags  []string `json:"tags"`
	Draft bool     `json:"draft"`
}

type dryRunVal struct {
	Passed   bool     `json:"passed"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

type dryRunHN struct {
	Action     string  `json:"action"`
	Slug       string  `json:"slug"`
	ExistingID *string `json:"existing_id"`
	DeletedID  *string `json:"deleted_id"`
}

type dryRunDT struct {
	Action          string   `json:"action"`
	Published       bool     `json:"published"`
	Tags            []string `json:"tags"`
	EmbedsConverted int      `json:"embeds_converted"`
}

func writeDryRunJSON(w io.Writer, files []*FileResult) {
	entries := make([]dryRunEntry, len(files))
	for i, fr := range files {
		entry := dryRunEntry{
			File: fr.File,
			Frontmatter: dryRunFM{
				Title: fr.Post.Title,
				Slug:  fr.Post.Slug,
				Tags:  fr.Post.Tags,
				Draft: fr.Post.Draft,
			},
			Validation: dryRunVal{
				Passed:   fr.Validation.Passed(),
				Errors:   nonNilSlice(fr.Validation.Errors),
				Warnings: nonNilSlice(fr.Validation.Warnings),
			},
		}

		if fr.Skipped {
			entry.Hashnode = &dryRunHN{Action: "skip", Slug: fr.Post.Slug}
			entry.DevTo = &dryRunDT{Action: "skip"}
		} else if !fr.Validation.Passed() {
			// validation failed: no platform results
		} else {
			// Hashnode result
			if fr.HashnodeRes != nil {
				hn := &dryRunHN{
					Action: fr.HashnodeRes.Action,
					Slug:   fr.Post.Slug,
				}
				if fr.HashnodeRes.PostID != "" {
					id := fr.HashnodeRes.PostID
					hn.ExistingID = &id
				}
				if fr.HashnodeRes.DeletedID != "" {
					id := fr.HashnodeRes.DeletedID
					hn.DeletedID = &id
				}
				entry.Hashnode = hn
			}

			// Dev.to result
			if fr.DevToRes != nil {
				dt := &dryRunDT{
					Action:          fr.DevToRes.Action,
					Published:       fr.DevToRes.Published,
					Tags:            fr.DevToTags,
					EmbedsConverted: fr.DevToRes.EmbedsConverted,
				}
				entry.DevTo = dt
			}
		}

		entries[i] = entry
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(entries)
}

// --- Summary text output ---

func writeSummary(w io.Writer, files []*FileResult) {
	published := 0
	skipped := 0
	failed := 0

	for _, fr := range files {
		if fr.Skipped {
			fmt.Fprintf(w, "%s: skipped (%s)\n", fr.Post.Slug, fr.SkipReason)
			skipped++
			continue
		}

		if !fr.Validation.Passed() {
			fmt.Fprintf(w, "%s: validation failed\n", fr.File)
			for _, e := range fr.Validation.Errors {
				fmt.Fprintf(w, "  Error: %s\n", e)
			}
			failed++
			continue
		}

		if fr.Error != "" {
			fmt.Fprintf(w, "%s: failed\n", fr.Post.Slug)
			fmt.Fprintf(w, "  Error: %s\n", fr.Error)
			if fr.HashnodeRes != nil {
				fmt.Fprintf(w, "  Hashnode: %s (%s)\n", fr.HashnodeRes.URL, fr.HashnodeRes.Action)
			}
			failed++
			continue
		}

		fmt.Fprintf(w, "%s: published\n", fr.Post.Slug)
		if fr.HashnodeRes != nil {
			fmt.Fprintf(w, "  Hashnode: %s (%s)\n", fr.HashnodeRes.URL, fr.HashnodeRes.Action)
		}
		if fr.DevToRes != nil {
			fmt.Fprintf(w, "  Dev.to:   %s (%s)\n", fr.DevToRes.URL, fr.DevToRes.Action)
		}
		canonicalURL := "https://blog.nuphirho.dev/" + fr.Post.Slug
		fmt.Fprintf(w, "  Medium:   import manually from %s\n", canonicalURL)
		published++

		// Warnings
		for _, w2 := range fr.TagWarnings {
			fmt.Fprintf(w, "  Warning:  %s\n", w2)
		}
		for _, w2 := range fr.Validation.Warnings {
			fmt.Fprintf(w, "  Warning:  %s\n", w2)
		}
	}

	total := len(files)
	fmt.Fprintf(w, "\n%d files processed", total)
	parts := []string{}
	if published > 0 {
		parts = append(parts, fmt.Sprintf("%d published", published))
	}
	if skipped > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", skipped))
	}
	if failed > 0 {
		parts = append(parts, fmt.Sprintf("%d failed", failed))
	}
	if len(parts) > 0 {
		fmt.Fprintf(w, ", %s", strings.Join(parts, ", "))
	}
	fmt.Fprintln(w, ".")
}

func nonNilSlice(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
