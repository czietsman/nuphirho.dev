// Command capture-governance-prompts-corpus fetches AGENTS.md files from a
// curated list of repositories and writes them with YAML front matter to
// experiments/governance-prompts-{version}/agents/.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var repos = []string{
	"czietsman/nuphirho.dev",
	"agentsmd/agents.md",
	"ansible/ansible",
	"apache/airflow",
	"BerriAI/litellm",
	"biomejs/biome",
	"coleam00/Archon",
	"elastic/elasticsearch",
	"giselles-ai/giselle",
	"google/adk-python",
	"grafana/grafana",
	"iannuttall/source-agents",
	"langchain-ai/langchain",
	"langflow-ai/langflow",
	"leonardsellem/codex-subagents-mcp",
	"mark3labs/mcp-go",
	"microsoft/vscode",
	"msitarzewski/AGENT-ZERO",
	"openai/codex",
	"openai/openai-agents-python",
	"originalankur/GenerateAgents.md",
	"Piebald-AI/tweakcc",
	"PrefectHQ/fastmcp",
	"prisma/prisma",
	"roboflow/supervision",
	"Significant-Gravitas/AutoGPT",
	"temporalio/sdk-java",
	"trevor-nichols/agentrules-architect",
	"twostraws/SwiftAgents",
	"vercel/next.js",
	"vercel-labs/agent-skills",
	"agent-sh/agnix",
	"adobe/aem-boilerplate",
	"yibie/SPEC-AGENTS.md",
}

var branches = []string{"main", "master", "canary"}
var paths = []string{"AGENTS.md", ".github/AGENTS.md"}

// redirect defines a corpus file whose AGENTS.md is a redirect or
// explicitly references a secondary governance file.
type redirect struct {
	Repo   string
	Path   string
	Reason string
}

var redirects = []redirect{
	{"PrefectHQ/fastmcp", "CLAUDE.md", "Pure redirect -- entire AGENTS.md is \"CLAUDE.md\""},
	{"microsoft/vscode", ".github/copilot-instructions.md", "Pure redirect -- 3-line pointer"},
	{"BerriAI/litellm", "CLAUDE.md", "Explicit: \"See CLAUDE.md for standard commands\""},
	{"agent-sh/agnix", "CLAUDE.md", "Project memory entrypoint; AGENTS.md is a copy"},
	{"apache/airflow", ".github/instructions/code-review.instructions.md", "Explicit: \"Read this file\""},
	{"mark3labs/mcp-go", "openspec/AGENTS.md", "@ syntax: \"Always open @/openspec/AGENTS.md\""},
}

// Result holds the capture outcome for a single repository.
type Result struct {
	Repo    string
	Source  string
	Commit  string
	Lines   int
	Content string
	Status  string
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr, os.Getenv, http.DefaultClient))
}

// HTTPClient abstracts HTTP requests for testing.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Run is the testable entry point.
func Run(args []string, stdout, stderr io.Writer, getenv func(string) string, client HTTPClient) int {
	version := "v1"
	resolveRedirects := false
	for i, arg := range args {
		if arg == "--resolve-redirects" {
			resolveRedirects = true
			args = append(args[:i], args[i+1:]...)
			break
		}
	}
	if len(args) > 0 && args[0] != "" {
		version = args[0]
	}

	if resolveRedirects {
		return runResolveRedirects(version, stdout, stderr, getenv, client)
	}

	token := getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprintln(stderr, "GITHUB_TOKEN is not set")
		return 1
	}

	baseDir := fmt.Sprintf("experiments/governance-prompts-%s", version)
	agentsDir := filepath.Join(baseDir, "agents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		fmt.Fprintf(stderr, "failed to create directory: %v\n", err)
		return 1
	}

	date := time.Now().UTC().Format("2006-01-02")
	var results []Result
	retrieved, failed, unknowns := 0, 0, 0

	for _, repo := range repos {
		r := capture(repo, token, date, client)
		results = append(results, r)

		if r.Status == "OK" {
			retrieved++
			fname := strings.ReplaceAll(repo, "/", "-") + ".md"
			fpath := filepath.Join(agentsDir, fname)
			front := fmt.Sprintf("---\nsource: %s\ncommit: %s\ncaptured: %s\nlines: %d\n---\n",
				r.Source, r.Commit, date, r.Lines)
			if err := os.WriteFile(fpath, []byte(front+r.Content), 0o644); err != nil {
				fmt.Fprintf(stderr, "write %s: %v\n", fpath, err)
			}
		} else {
			failed++
		}
		if r.Commit == "unknown" || r.Commit == "" {
			unknowns++
		}
	}

	// Write manifest.
	writeManifest(filepath.Join(baseDir, "manifest.md"), results, date)

	// Write README.
	writeReadme(filepath.Join(baseDir, "README.md"), version, date, len(results))

	// Write or append GOVERNANCE_PROMPTS.md.
	writeIndex("experiments/GOVERNANCE_PROMPTS.md", version, date, len(results))

	fmt.Fprintf(stdout, "Retrieved: %d  Failed: %d  Unknown SHAs: %d\n", retrieved, failed, unknowns)

	if unknowns > 0 {
		fmt.Fprintln(stderr, "ERROR: some commit SHAs are unknown or missing")
		return 1
	}
	return 0
}

func capture(repo, token, date string, client HTTPClient) Result {
	r := Result{Repo: repo, Status: "Not found", Commit: "unknown"}

	// Try each branch/path combination.
	for _, branch := range branches {
		for _, path := range paths {
			url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", repo, branch, path)
			req, _ := http.NewRequest("GET", url, nil)
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode != 200 {
				if resp != nil {
					resp.Body.Close()
				}
				continue
			}
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				continue
			}
			r.Content = string(body)
			r.Lines = countLines(r.Content)
			r.Source = url
			r.Status = "OK"

			// Fetch commit SHA via API.
			r.Commit = fetchCommitSHA(repo, path, token, client)
			return r
		}
	}
	return r
}

type commitEntry struct {
	SHA string `json:"sha"`
}

func fetchCommitSHA(repo, path, token string, client HTTPClient) string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/commits?path=%s&per_page=1", repo, path)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "unknown"
	}

	var entries []commitEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil || len(entries) == 0 {
		return "unknown"
	}
	return entries[0].SHA
}

func countLines(s string) int {
	if s == "" {
		return 0
	}
	n := strings.Count(s, "\n")
	if !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}

func writeManifest(path string, results []Result, date string) {
	var b strings.Builder
	b.WriteString("# Corpus Manifest: governance-prompts-v1\n")
	b.WriteString("# Study: Structural Quality Gaps in Practitioner AI Governance Prompts\n")
	b.WriteString("# Paper: RE@Next! 2026 submission\n")
	b.WriteString(fmt.Sprintf("# Retrieval date: %s\n", date))
	b.WriteString("# Retrieved by: capture-governance-prompts-corpus\n\n")
	b.WriteString("| # | Repository | Short SHA | Full SHA | Lines | Status |\n")
	b.WriteString("|---|---|---|---|---|---|\n")
	for i, r := range results {
		short := r.Commit
		if len(short) >= 7 {
			short = short[:7]
		}
		b.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %d | %s |\n",
			i+1, r.Repo, short, r.Commit, r.Lines, r.Status))
	}
	os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeReadme(path, version, date string, count int) {
	content := fmt.Sprintf(`# governance-prompts-%s

Corpus of AGENTS.md governance prompts captured from %d public repositories.

- **Version:** %s
- **Captured:** %s
- **Tool:** capture-governance-prompts-corpus

See manifest.md for the full list of repositories, commit SHAs, and retrieval status.

## Purpose

This corpus supports the RE@Next! 2026 paper submission on structural quality
gaps in practitioner AI governance prompts. Each file in agents/ contains the
verbatim AGENTS.md content with a YAML front matter header recording the source
URL, commit SHA, capture date, and line count.
`, version, count, version, date)
	os.WriteFile(path, []byte(content), 0o644)
}

// redirectResult holds the outcome of resolving one redirect.
type redirectResult struct {
	Repo    string
	Path    string
	Reason  string
	Status  string
	Lines   int
	Commit  string
	Source  string
}

func runResolveRedirects(version string, stdout, stderr io.Writer, getenv func(string) string, client HTTPClient) int {
	token := getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprintln(stderr, "GITHUB_TOKEN is not set")
		return 1
	}

	baseDir := fmt.Sprintf("experiments/governance-prompts-%s", version)
	agentsDir := filepath.Join(baseDir, "agents")
	date := time.Now().UTC().Format("2006-01-02")

	var results []redirectResult
	appended, failed := 0, 0

	for _, rd := range redirects {
		rr := redirectResult{Repo: rd.Repo, Path: rd.Path, Reason: rd.Reason, Status: "Not found", Commit: "unknown"}
		fname := strings.ReplaceAll(rd.Repo, "/", "-") + ".md"
		fpath := filepath.Join(agentsDir, fname)

		// Check corpus file exists.
		if _, err := os.Stat(fpath); os.IsNotExist(err) {
			fmt.Fprintf(stderr, "corpus file not found: %s\n", fpath)
			rr.Status = "Corpus file missing"
			results = append(results, rr)
			failed++
			continue
		}

		// Fetch the redirect target.
		var content string
		var source string
		fetched := false
		for _, branch := range branches {
			url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", rd.Repo, branch, rd.Path)
			req, _ := http.NewRequest("GET", url, nil)
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode != 200 {
				if resp != nil {
					resp.Body.Close()
				}
				continue
			}
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				continue
			}
			content = string(body)
			source = url
			fetched = true
			break
		}

		if !fetched {
			fmt.Fprintf(stderr, "could not fetch %s from %s\n", rd.Path, rd.Repo)
			results = append(results, rr)
			failed++
			continue
		}

		// Fetch commit SHA.
		sha := fetchCommitSHA(rd.Repo, rd.Path, token, client)
		lines := countLines(content)

		// Append to corpus file.
		separator := fmt.Sprintf("\n\n---\n## Appended governance file: %s\nsource: %s\ncommit: %s\ncaptured: %s\nlines: %d\nreason: %s\n---\n\n%s",
			rd.Path, source, sha, date, lines, rd.Reason, content)

		f, err := os.OpenFile(fpath, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			fmt.Fprintf(stderr, "cannot open %s for append: %v\n", fpath, err)
			rr.Status = "Write error"
			results = append(results, rr)
			failed++
			continue
		}
		_, err = f.WriteString(separator)
		f.Close()
		if err != nil {
			fmt.Fprintf(stderr, "write error for %s: %v\n", fpath, err)
			rr.Status = "Write error"
			results = append(results, rr)
			failed++
			continue
		}

		rr.Status = "Appended"
		rr.Lines = lines
		rr.Commit = sha
		rr.Source = source
		results = append(results, rr)
		appended++

		fmt.Fprintf(stdout, "Appended %s (%d lines) to %s\n", rd.Path, lines, fname)
	}

	// Write redirect resolution report.
	writeRedirectReport(filepath.Join(baseDir, "results", "redirect-resolution.md"), results, date)

	// Update manifest.
	updateManifestRedirects(filepath.Join(baseDir, "manifest.md"), results)

	fmt.Fprintf(stdout, "Appended: %d  Failed: %d\n", appended, failed)
	if failed > 0 {
		return 1
	}
	return 0
}

func writeRedirectReport(path string, results []redirectResult, date string) {
	var b strings.Builder
	b.WriteString("# Redirect Resolution Report\n")
	b.WriteString(fmt.Sprintf("# Date: %s\n\n", date))
	b.WriteString("| Repository | Redirect target | Status | Lines appended | Commit SHA |\n")
	b.WriteString("|---|---|---|---|---|\n")
	for _, r := range results {
		short := r.Commit
		if len(short) >= 7 {
			short = short[:7]
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %s |\n",
			r.Repo, r.Path, r.Status, r.Lines, short))
	}
	os.WriteFile(path, []byte(b.String()), 0o644)
}

func updateManifestRedirects(path string, results []redirectResult) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	content := string(data)

	for _, r := range results {
		if r.Status != "Appended" {
			continue
		}
		// Find the manifest row for this repo and update the status.
		old := fmt.Sprintf("| %s |", r.Repo)
		// Find the line containing this repo.
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if strings.Contains(line, old) && strings.HasSuffix(strings.TrimSpace(line), "| OK |") {
				lines[i] = strings.TrimSuffix(strings.TrimSpace(line), "| OK |") +
					fmt.Sprintf("| OK + redirect: %s |", r.Path)
				break
			}
		}
		content = strings.Join(lines, "\n")
	}

	os.WriteFile(path, []byte(content), 0o644)
}

func writeIndex(path, version, date string, count int) {
	entry := fmt.Sprintf("- **%s** (%s): %d repositories captured\n", version, date, count)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		header := "# Governance Prompts\n\nCorpus tracking index for the governance-prompts experiment.\n\n## Versions\n\n"
		os.WriteFile(path, []byte(header+entry), 0o644)
		return
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(entry)
}
