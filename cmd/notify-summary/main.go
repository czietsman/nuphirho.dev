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
	publishOutputFile := fs.String("publish-output-file", "", "Path to publish command output for target-level change reporting")

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

	summary, err := buildSummary(*postsDir, today, *publishExitCode, *publishOutputFile)
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

func buildSummary(postsDir string, today time.Time, publishExitCode int, publishOutputFile string) (string, error) {
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

	changeLines, err := summarizePublishOutput(publishOutputFile)
	if err != nil {
		return "", err
	}
	lines = append(lines, changeLines...)

	if len(todayTitles) > 0 {
		if publishOutputFile == "" && publishExitCode == 0 {
			lines = append(lines, fmt.Sprintf("Published today: %s", strings.Join(todayTitles, ", ")))
		} else if publishExitCode != 0 {
			lines = append(lines, fmt.Sprintf("Publish failed today: %s", strings.Join(todayTitles, ", ")))
		}
	}

	if len(lines) == 0 {
		return "", nil
	}

	return strings.Join(lines, "\n") + "\n", nil
}

func summarizePublishOutput(path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var lines []string
	var currentSlug string

	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if !strings.HasPrefix(line, "  ") && strings.Contains(trimmed, ": ") {
			currentSlug = strings.SplitN(trimmed, ": ", 2)[0]
			continue
		}

		if currentSlug == "" {
			continue
		}

		switch {
		case strings.HasPrefix(trimmed, "Hashnode:"):
			action := extractAction(trimmed)
			if action != "" && action != "unchanged" {
				lines = append(lines, fmt.Sprintf("%s: Hashnode %s", currentSlug, action))
			}
		case strings.HasPrefix(trimmed, "Dev.to:"):
			action := extractAction(trimmed)
			if action != "" && action != "unchanged" {
				lines = append(lines, fmt.Sprintf("%s: Dev.to %s", currentSlug, action))
			}
		}
	}

	return lines, nil
}

func extractAction(line string) string {
	start := strings.LastIndex(line, "(")
	end := strings.LastIndex(line, ")")
	if start == -1 || end == -1 || end <= start+1 {
		return ""
	}
	return line[start+1 : end]
}
