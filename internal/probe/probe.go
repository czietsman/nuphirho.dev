// Package probe provides contract probes that validate credentials and API
// shapes for all configured platforms before the publish pipeline runs.
package probe

import (
	"fmt"
	"io"
	"strings"
)

// CheckStatus represents the outcome of a single probe check.
type CheckStatus int

const (
	StatusOK   CheckStatus = iota
	StatusFail             // hard failure, affects exit code
	StatusWarn             // soft failure (e.g. cleanup), does not affect exit code
	StatusSkip             // skipped because a prior check failed
)

func (s CheckStatus) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusFail:
		return "FAIL"
	case StatusWarn:
		return "WARN"
	case StatusSkip:
		return "SKIP"
	default:
		return "?"
	}
}

// CheckResult holds the outcome of a single probe check.
type CheckResult struct {
	Platform string // "hashnode" or "devto"
	Name     string // e.g. "credentials", "publication", "draft mutation", "cleanup"
	Status   CheckStatus
	Detail   string // optional error detail
}

// Result holds the full probe run outcome.
type Result struct {
	Checks   []CheckResult
	ExitCode int
}

// HashnodeProber is the interface for Hashnode probe operations.
type HashnodeProber interface {
	CheckCredentials() error
	CheckPublication() error
	CreateProbeDraft() (string, error) // returns draft ID
	DeleteDraft(id string) error
}

// DevtoProber is the interface for Dev.to probe operations.
type DevtoProber interface {
	CheckCredentials() error
	CreateProbeArticle() (int, error) // returns article ID
	DeleteArticle(id int) error
}

// ProbeAll runs all contract probes and writes results to the given writer.
func ProbeAll(w io.Writer, hn HashnodeProber, dt DevtoProber) *Result {
	r := &Result{}

	probeHashnode(w, hn, r)
	probeDevTo(w, dt, r)

	// Summary
	fmt.Fprintln(w)
	failures := 0
	for _, c := range r.Checks {
		if c.Status == StatusFail {
			failures++
		}
	}
	if failures == 0 {
		fmt.Fprintln(w, "All probes passed.")
		r.ExitCode = 0
	} else {
		fmt.Fprintf(w, "%d probes failed.\n", failures)
		r.ExitCode = 1
	}

	return r
}

func probeHashnode(w io.Writer, hn HashnodeProber, r *Result) {
	// 1. Credentials
	err := hn.CheckCredentials()
	if err != nil {
		r.Checks = append(r.Checks, check(w, "hashnode", "credentials", StatusFail, err.Error()))
		r.Checks = append(r.Checks, check(w, "hashnode", "publication", StatusSkip, ""))
		r.Checks = append(r.Checks, check(w, "hashnode", "draft mutation", StatusSkip, ""))
		r.Checks = append(r.Checks, check(w, "hashnode", "cleanup", StatusSkip, ""))
		return
	}
	r.Checks = append(r.Checks, check(w, "hashnode", "credentials", StatusOK, ""))

	// 2. Publication resolves
	err = hn.CheckPublication()
	if err != nil {
		r.Checks = append(r.Checks, check(w, "hashnode", "publication", StatusFail, err.Error()))
		r.Checks = append(r.Checks, check(w, "hashnode", "draft mutation", StatusSkip, ""))
		r.Checks = append(r.Checks, check(w, "hashnode", "cleanup", StatusSkip, ""))
		return
	}
	r.Checks = append(r.Checks, check(w, "hashnode", "publication", StatusOK, ""))

	// 3. Draft mutation shape
	draftID, err := hn.CreateProbeDraft()
	if err != nil {
		r.Checks = append(r.Checks, check(w, "hashnode", "draft mutation", StatusFail, err.Error()))
		r.Checks = append(r.Checks, check(w, "hashnode", "cleanup", StatusSkip, ""))
		return
	}
	r.Checks = append(r.Checks, check(w, "hashnode", "draft mutation", StatusOK, ""))

	// 4. Cleanup (best-effort)
	err = hn.DeleteDraft(draftID)
	if err != nil {
		r.Checks = append(r.Checks, check(w, "hashnode", "cleanup", StatusWarn, err.Error()))
	} else {
		r.Checks = append(r.Checks, check(w, "hashnode", "cleanup", StatusOK, ""))
	}
}

func probeDevTo(w io.Writer, dt DevtoProber, r *Result) {
	// 1. Credentials
	err := dt.CheckCredentials()
	if err != nil {
		r.Checks = append(r.Checks, check(w, "devto", "credentials", StatusFail, err.Error()))
		r.Checks = append(r.Checks, check(w, "devto", "article mutation", StatusSkip, ""))
		r.Checks = append(r.Checks, check(w, "devto", "cleanup", StatusSkip, ""))
		return
	}
	r.Checks = append(r.Checks, check(w, "devto", "credentials", StatusOK, ""))

	// 2. Article create shape
	articleID, err := dt.CreateProbeArticle()
	if err != nil {
		r.Checks = append(r.Checks, check(w, "devto", "article mutation", StatusFail, err.Error()))
		r.Checks = append(r.Checks, check(w, "devto", "cleanup", StatusSkip, ""))
		return
	}
	r.Checks = append(r.Checks, check(w, "devto", "article mutation", StatusOK, ""))

	// 3. Cleanup (best-effort)
	err = dt.DeleteArticle(articleID)
	if err != nil {
		r.Checks = append(r.Checks, check(w, "devto", "cleanup", StatusWarn, err.Error()))
	} else {
		r.Checks = append(r.Checks, check(w, "devto", "cleanup", StatusOK, ""))
	}
}

func check(w io.Writer, platform, name string, status CheckStatus, detail string) CheckResult {
	// Format: [platform]__name_padded__STATUS
	// Total platform column is 10 chars: "[hashnode]" (10) or "[devto]   " (10)
	bracket := "[" + platform + "]"
	platformFmt := padRight(bracket, 10)
	nameFmt := padRight(name, 18)
	line := fmt.Sprintf("%s %s %s", platformFmt, nameFmt, status)
	if detail != "" && status != StatusOK {
		line += "  " + detail
	}
	fmt.Fprintln(w, line)
	return CheckResult{Platform: platform, Name: name, Status: status, Detail: detail}
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
