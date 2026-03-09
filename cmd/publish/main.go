// Command publish is the CLI entrypoint for the nuphirho.dev publish pipeline.
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/czietsman/nuphirho.dev/internal/devto"
	"github.com/czietsman/nuphirho.dev/internal/hashnode"
	"github.com/czietsman/nuphirho.dev/internal/pipeline"
	"github.com/czietsman/nuphirho.dev/internal/probe"
	"github.com/czietsman/nuphirho.dev/internal/tags"
)

// Deps holds optional dependency overrides for testing.
// When nil or when a field is nil, the real implementation is used.
type Deps struct {
	Hashnode pipeline.HashnodePub
	DevTo    pipeline.DevtoPub
	Prober   pipeline.Prober
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr, os.Getenv, nil))
}

// Run is the testable CLI entry point.
func Run(args []string, stdout, stderr io.Writer, getenv func(string) string, d *Deps) int {
	fs := flag.NewFlagSet("publish", flag.ContinueOnError)
	fs.SetOutput(stderr)

	probeFlag := fs.Bool("probe", false, "Run contract probes and exit")
	dryRun := fs.Bool("dry-run", false, "Validate and report without publishing (JSON output)")
	tagsFile := fs.String("tags-file", "tags.json", "Path to tag glossary")
	hnToken := fs.String("hashnode-token", "", "Hashnode API token (or HASHNODE_TOKEN env)")
	hnPubID := fs.String("hashnode-publication", "", "Hashnode publication ID (or HASHNODE_PUBLICATION_ID env)")
	dtKey := fs.String("devto-api-key", "", "Dev.to API key (or DEVTO_API_KEY env)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	// Environment variable fallback
	if *hnToken == "" {
		*hnToken = getenv("HASHNODE_TOKEN")
	}
	if *hnPubID == "" {
		*hnPubID = getenv("HASHNODE_PUBLICATION_ID")
	}
	if *dtKey == "" {
		*dtKey = getenv("DEVTO_API_KEY")
	}

	// Validate required credentials
	if *hnToken == "" {
		fmt.Fprintln(stderr, "error: Hashnode token is required (--hashnode-token or HASHNODE_TOKEN)")
		return 1
	}
	if *hnPubID == "" {
		fmt.Fprintln(stderr, "error: Hashnode publication ID is required (--hashnode-publication or HASHNODE_PUBLICATION_ID)")
		return 1
	}
	if *dtKey == "" {
		fmt.Fprintln(stderr, "error: Dev.to API key is required (--devto-api-key or DEVTO_API_KEY)")
		return 1
	}

	// Probe mode: run probes and exit
	if *probeFlag {
		var prober pipeline.Prober
		if d != nil && d.Prober != nil {
			prober = d.Prober
		} else {
			prober = newRealProber(*hnToken, *hnPubID, *dtKey)
		}
		return pipeline.Probe(prober, stdout)
	}

	// Collect post files from remaining args
	filePaths := fs.Args()

	if len(filePaths) == 0 {
		if *dryRun {
			fmt.Fprintln(stdout, "[]")
			return 0
		}
		fmt.Fprintln(stderr, "error: no post files specified")
		return 1
	}

	// Load tag glossary
	glossary, err := tags.LoadGlossary(*tagsFile)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}

	// Read post files
	var postFiles []pipeline.PostFile
	for _, path := range filePaths {
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(stderr, "error: reading %s: %s\n", path, err)
			return 1
		}
		postFiles = append(postFiles, pipeline.PostFile{Path: path, Content: string(content)})
	}

	// Build pipeline clients
	var hn pipeline.HashnodePub
	var dt pipeline.DevtoPub

	if d != nil && d.Hashnode != nil {
		hn = d.Hashnode
	} else {
		hnClient := hashnode.New(*hnToken, *hnPubID, http.DefaultClient)
		hnClient.DryRun = *dryRun
		hn = hnClient
	}

	if d != nil && d.DevTo != nil {
		dt = d.DevTo
	} else {
		dtClient := devto.New(*dtKey, http.DefaultClient)
		dtClient.DryRun = *dryRun
		dt = dtClient
	}

	cfg := pipeline.Config{
		Hashnode: hn,
		DevTo:    dt,
		Glossary: glossary,
		DryRun:   *dryRun,
	}

	result := pipeline.Run(cfg, postFiles, stdout)
	return result.ExitCode
}

// realProber wraps probe.ProbeAll with real API adapters.
type realProber struct {
	hn probe.HashnodeProber
	dt probe.DevtoProber
}

func (p *realProber) Run(w io.Writer) int {
	return probe.ProbeAll(w, p.hn, p.dt).ExitCode
}

func newRealProber(hnToken, hnPubID, dtKey string) *realProber {
	hnClient := hashnode.New(hnToken, hnPubID, http.DefaultClient)
	dtClient := devto.New(dtKey, http.DefaultClient)
	return &realProber{
		hn: &hashnodeProbeAdapter{client: hnClient},
		dt: &devtoProbeAdapter{client: dtClient},
	}
}

// hashnodeProbeAdapter adapts hashnode.Client for probing.
type hashnodeProbeAdapter struct {
	client *hashnode.Client
}

func (a *hashnodeProbeAdapter) CheckCredentials() error {
	_, err := a.client.CheckPostBySlug("__probe__")
	return err
}

func (a *hashnodeProbeAdapter) CheckPublication() error {
	_, err := a.client.CheckPostBySlug("__probe__")
	return err
}

func (a *hashnodeProbeAdapter) CreateProbeDraft() (string, error) {
	input := hashnode.PostInput{
		Title:   "[probe] contract check",
		Slug:    "probe-contract-check",
		Content: "Automated probe. Safe to delete.",
	}
	result, err := a.client.CreateDraft(input)
	if err != nil {
		return "", err
	}
	return result.DraftID, nil
}

func (a *hashnodeProbeAdapter) DeleteDraft(_ string) error {
	// Draft deletion requires a removeDraft mutation not yet implemented
	// on the hashnode client. Best-effort: return nil for now.
	return nil
}

// devtoProbeAdapter adapts devto.Client for probing.
type devtoProbeAdapter struct {
	client *devto.Client
}

func (a *devtoProbeAdapter) CheckCredentials() error {
	// CreateArticle with an invalid/empty payload would test credentials,
	// but the simplest check is to attempt the articles/me/all lookup.
	input := devto.ArticleInput{
		Title:     "[probe] contract check",
		Slug:      "probe-contract-check",
		Content:   "Automated probe. Safe to delete.",
		Published: false,
	}
	_, err := a.client.CreateArticle(input)
	return err
}

func (a *devtoProbeAdapter) CreateProbeArticle() (int, error) {
	input := devto.ArticleInput{
		Title:     "[probe] contract check",
		Slug:      "probe-contract-check",
		Content:   "Automated probe. Safe to delete.",
		Published: false,
	}
	result, err := a.client.CreateArticle(input)
	if err != nil {
		return 0, err
	}
	return result.ArticleID, nil
}

func (a *devtoProbeAdapter) DeleteArticle(_ int) error {
	// Article deletion via Dev.to API not yet implemented.
	// Best-effort: return nil for now.
	return nil
}
