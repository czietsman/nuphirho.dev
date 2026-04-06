// Command score-governance-prompts scores AGENTS.md files against PromptQ
// principles by piping assembled prompts to LLM CLI tools.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var models = []string{"claude", "codex", "gemini"}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}

// Run is the testable entry point.
func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("score-governance-prompts", flag.ContinueOnError)
	fs.SetOutput(stderr)

	model := fs.String("model", "claude", "Model to use: claude, codex, gemini, or all")
	runs := fs.Int("runs", 1, "Number of runs per file")
	baseDir := fs.String("dir", "experiments/governance-prompts-v1", "Experiment directory")
	version := fs.String("version", "v1", "Scoring version (controls prompt and results directory)")
	files := fs.String("files", "", "Comma-separated basenames to score (e.g. PrefectHQ-fastmcp,BerriAI-litellm)")
	failFast := fs.Bool("fail-fast", false, "Stop on first scoring error")
	force := fs.Bool("force", false, "Re-score files that already have valid results")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	var filter map[string]bool
	if *files != "" {
		filter = make(map[string]bool)
		for _, f := range strings.Split(*files, ",") {
			filter[strings.TrimSpace(f)] = true
		}
	}

	opts := scoreOpts{filter: filter, failFast: *failFast, force: *force, version: *version}

	if *model == "all" {
		for _, m := range models {
			code := scoreModel(m, *runs, *baseDir, opts, stdout, stderr)
			if code != 0 {
				return code
			}
		}
		archive(*baseDir, opts.version, stdout, stderr)
		return 0
	}

	code := scoreModel(*model, *runs, *baseDir, opts, stdout, stderr)
	if code != 0 {
		return code
	}
	archive(*baseDir, opts.version, stdout, stderr)
	return 0
}

type scoreOpts struct {
	filter   map[string]bool
	failFast bool
	force    bool
	version  string
}

func scoreModel(model string, runs int, baseDir string, opts scoreOpts, stdout, stderr io.Writer) int {
	promptFile := fmt.Sprintf("score-%s.txt", opts.version)
	promptPath := filepath.Join(baseDir, "prompts", promptFile)
	agentsDir := filepath.Join(baseDir, "agents")
	scoresDir := fmt.Sprintf("scores-%s", opts.version)
	resultsDir := filepath.Join(baseDir, "results", scoresDir, model)

	template, err := os.ReadFile(promptPath)
	if err != nil {
		fmt.Fprintf(stderr, "cannot read prompt template: %v\n", err)
		return 1
	}

	if err := os.MkdirAll(resultsDir, 0o755); err != nil {
		fmt.Fprintf(stderr, "cannot create results directory: %v\n", err)
		return 1
	}

	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		fmt.Fprintf(stderr, "cannot read agents directory: %v\n", err)
		return 1
	}

	scored, skipped, failed := 0, 0, 0
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		basename := strings.TrimSuffix(entry.Name(), ".md")
		if opts.filter != nil && !opts.filter[basename] {
			continue
		}
		agentPath := filepath.Join(agentsDir, entry.Name())

		content, err := readAgentContent(agentPath)
		if err != nil {
			fmt.Fprintf(stderr, "cannot read %s: %v\n", agentPath, err)
			failed++
			if opts.failFast {
				break
			}
			continue
		}

		prompt := strings.Replace(string(template), "{agent_source}", content, 1)

		for run := 1; run <= runs; run++ {
			outputFile := filepath.Join(resultsDir, basename+".json")
			if runs > 1 {
				outputFile = filepath.Join(resultsDir, fmt.Sprintf("%s_run%d.json", basename, run))
			}

			if !opts.force && hasValidResult(outputFile) {
				fmt.Fprintf(stdout, "Skipping %s (%s, run %d) -- valid result exists\n", basename, model, run)
				skipped++
				continue
			}

			fmt.Fprintf(stdout, "Scoring %s (%s, run %d)\n", basename, model, run)

			out, err := runModel(model, prompt)
			if err != nil {
				fmt.Fprintf(stderr, "  error: %v\n", err)
				failed++
				if opts.failFast {
					break
				}
				continue
			}

			if err := os.WriteFile(outputFile, out, 0o644); err != nil {
				fmt.Fprintf(stderr, "  write error: %v\n", err)
				failed++
				if opts.failFast {
					break
				}
				continue
			}
			scored++
		}
		if opts.failFast && failed > 0 {
			break
		}
	}

	fmt.Fprintf(stdout, "%s: scored %d, skipped %d, failed %d\n", model, scored, skipped, failed)
	if failed > 0 {
		return 1
	}
	return 0
}

// hasValidResult returns true if the file exists and contains valid JSON.
func hasValidResult(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return json.Valid(data)
}

// readAgentContent reads an agent file and strips YAML front matter.
func readAgentContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	content := string(data)
	parts := strings.SplitN(content, "---", 3)
	if len(parts) >= 3 {
		return strings.TrimSpace(parts[2]), nil
	}
	return content, nil
}

// runModel executes the model CLI with the given prompt.
// Claude reads from stdin; codex and gemini take the prompt as an argument.
func runModel(model, prompt string) ([]byte, error) {
	var cmd *exec.Cmd
	switch model {
	case "claude":
		cmd = exec.Command("claude", "-p")
		cmd.Stdin = strings.NewReader(prompt)
	case "codex":
		cmd = exec.Command("codex", "exec", prompt)
	case "gemini":
		cmd = exec.Command("gemini", "-p", prompt)
	default:
		return nil, fmt.Errorf("unknown model: %s", model)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%v: %s", err, string(out))
	}
	return extractJSON(out), nil
}

// extractJSON extracts a valid JSON object from model output that may
// contain markdown code fences, banner text, or other noise.
func extractJSON(data []byte) []byte {
	s := string(data)

	// Strategy 1: extract from markdown code fences.
	if openIdx := strings.Index(s, "```"); openIdx != -1 {
		afterOpen := s[openIdx:]
		if nlIdx := strings.Index(afterOpen, "\n"); nlIdx != -1 {
			content := afterOpen[nlIdx+1:]
			if closeIdx := strings.LastIndex(content, "```"); closeIdx != -1 {
				content = content[:closeIdx]
			}
			trimmed := strings.TrimSpace(content)
			if json.Valid([]byte(trimmed)) {
				return []byte(trimmed)
			}
		}
	}

	// Strategy 2: find the last valid JSON object in the output.
	// Models sometimes duplicate output or prepend banners; the last
	// complete JSON object is typically the clean one.
	lastIdx := strings.LastIndex(s, "\n{")
	if lastIdx == -1 && strings.HasPrefix(strings.TrimSpace(s), "{") {
		lastIdx = strings.Index(s, "{") - 1
	}
	if lastIdx >= -1 {
		candidate := s[lastIdx+1:]
		// Find matching closing brace by scanning for the last '}'.
		if closeIdx := strings.LastIndex(candidate, "}"); closeIdx != -1 {
			candidate = strings.TrimSpace(candidate[:closeIdx+1])
			if json.Valid([]byte(candidate)) {
				return []byte(candidate)
			}
		}
	}

	// Strategy 3: the whole output might already be valid JSON.
	trimmed := strings.TrimSpace(s)
	if json.Valid([]byte(trimmed)) {
		return []byte(trimmed)
	}

	// Give up -- return trimmed output as-is.
	return []byte(trimmed)
}

// archive copies results to EXPERIMENT_ARCHIVE_DIR if set.
func archive(baseDir, version string, stdout, stderr io.Writer) {
	archiveRoot := os.Getenv("EXPERIMENT_ARCHIVE_DIR")
	if archiveRoot == "" {
		return
	}

	scoresDir := fmt.Sprintf("scores-%s", version)
	srcDir := filepath.Join(baseDir, "results", scoresDir)
	dstDir := filepath.Join(archiveRoot, "experiments", "governance-prompts-v1", scoresDir)

	if _, err := os.Stat(archiveRoot); os.IsNotExist(err) {
		fmt.Fprintf(stderr, "Warning: EXPERIMENT_ARCHIVE_DIR not found at %s. Results not archived.\n", archiveRoot)
		return
	}

	fmt.Fprintf(stdout, "Archiving results...\n")
	if err := copyDir(srcDir, dstDir); err != nil {
		fmt.Fprintf(stderr, "Archive error: %v\n", err)
		return
	}
	fmt.Fprintf(stdout, "Archived. Remember to commit in the archive repo.\n")
}

// copyDir recursively copies src to dst.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}
