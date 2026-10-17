// Command deploy-summary writes the message the blog deployment sends to
// Telegram: which posts went live today and which are queued, using the same
// visibility rule as the blog (not a draft, publish_date on or before today).
// A failed deployment also carries the run URL so it can be investigated.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/czietsman/nuphirho.dev/internal/frontmatter"
)

const siteOrigin = "https://blog.nuphirho.dev"

type post struct {
	title string
	slug  string
	date  time.Time
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}

// Run reads the posts directory and prints the deployment summary.
func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("deploy-summary", flag.ContinueOnError)
	fs.SetOutput(stderr)
	status := fs.String("status", "", "Deployment job status, for example success or failure")
	runURL := fs.String("run-url", "", "GitHub Actions run URL, included when the deployment did not succeed")
	todayStr := fs.String("today", "", "Reference date YYYY-MM-DD (default: today in UTC)")
	if err := fs.Parse(args); err != nil {
		return 1
	}
	if fs.NArg() != 1 || *status == "" {
		fmt.Fprintln(stderr, "usage: deploy-summary -status <status> [-run-url <url>] [-today YYYY-MM-DD] <posts-dir>")
		return 1
	}

	today := time.Now().UTC()
	if *todayStr != "" {
		t, err := time.Parse("2006-01-02", *todayStr)
		if err != nil {
			fmt.Fprintf(stderr, "error: invalid -today %q: %s\n", *todayStr, err)
			return 1
		}
		today = t
	}
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

	posts, err := readPosts(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}
	fmt.Fprintln(stdout, summary(posts, today, *status, *runURL))
	return 0
}

func readPosts(dir string) ([]post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}
	var posts []post
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || strings.HasPrefix(e.Name(), "_") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		// Only the title, slug, draft flag, and date matter here; validation
		// is the pipeline's job elsewhere.
		p, _ := frontmatter.Parse(string(raw))
		if p.Draft || p.PublishDate == nil {
			continue
		}
		posts = append(posts, post{title: p.Title, slug: p.Slug, date: p.PublishDate.Time})
	}
	sort.SliceStable(posts, func(i, j int) bool { return posts[i].date.Before(posts[j].date) })
	return posts, nil
}

func summary(posts []post, today time.Time, status, runURL string) string {
	var live, queued []string
	for _, p := range posts {
		switch {
		case p.date.Equal(today):
			live = append(live, fmt.Sprintf("- %s: %s/%s", p.title, siteOrigin, p.slug))
		case p.date.After(today):
			queued = append(queued, fmt.Sprintf("- %s: %s", p.date.Format("Mon 2 Jan"), p.title))
		}
	}

	var sections []string
	if status != "success" {
		sections = append(sections, fmt.Sprintf("Blog deployment %s: %s", status, runURL))
	} else if len(live) == 0 {
		sections = append(sections, "Blog deployed. Nothing new went live today.")
	} else {
		sections = append(sections, "Blog deployed.", "Live today:\n"+strings.Join(live, "\n"))
	}
	if len(queued) == 0 {
		sections = append(sections, "Nothing is queued.")
	} else {
		sections = append(sections, "Queued:\n"+strings.Join(queued, "\n"))
	}
	return strings.Join(sections, "\n\n")
}
