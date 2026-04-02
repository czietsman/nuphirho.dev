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

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}

// Run builds a daily publish summary message. It writes nothing when there is
// no queue or publish outcome to report.
func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("notify-summary", flag.ContinueOnError)
	fs.SetOutput(stderr)

	postsDir := fs.String("posts-dir", "posts", "Directory containing markdown posts")
	todayStr := fs.String("today", "", "Today in YYYY-MM-DD format")
	publishExitCode := fs.Int("publish-exit-code", 0, "Exit code from the publish step")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *todayStr == "" {
		fmt.Fprintln(stderr, "error: today is required (--today)")
		return 1
	}

	today, err := time.Parse("2006-01-02", *todayStr)
	if err != nil {
		fmt.Fprintf(stderr, "error: invalid today value: %s\n", err)
		return 1
	}

	summary, err := buildSummary(*postsDir, today, *publishExitCode)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}

	if summary == "" {
		return 0
	}

	fmt.Fprint(stdout, summary)
	return 0
}

func buildSummary(postsDir string, today time.Time, publishExitCode int) (string, error) {
	var todayTitles []string
	var tomorrowTitles []string

	tomorrow := today.AddDate(0, 0, 1).Format("2006-01-02")
	todayStr := today.Format("2006-01-02")

	err := filepath.WalkDir(postsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		post, validation := frontmatter.Parse(string(raw))
		if !validation.Passed() || post.Draft || post.PublishDate == nil {
			return nil
		}

		switch post.PublishDate.Format("2006-01-02") {
		case todayStr:
			todayTitles = append(todayTitles, post.Title)
		case tomorrow:
			tomorrowTitles = append(tomorrowTitles, post.Title)
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	sort.Strings(todayTitles)
	sort.Strings(tomorrowTitles)

	var lines []string
	if len(tomorrowTitles) > 0 {
		lines = append(lines, fmt.Sprintf("Queued for tomorrow: %s", strings.Join(tomorrowTitles, ", ")))
	}
	if len(todayTitles) > 0 {
		if publishExitCode == 0 {
			lines = append(lines, fmt.Sprintf("Published today: %s", strings.Join(todayTitles, ", ")))
		} else {
			lines = append(lines, fmt.Sprintf("Publish failed today: %s", strings.Join(todayTitles, ", ")))
		}
	}

	if len(lines) == 0 {
		return "", nil
	}

	return strings.Join(lines, "\n") + "\n", nil
}
