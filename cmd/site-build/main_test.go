package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := Run([]string{
		"--templates-dir", "../../site/templates",
		"--static-dir", "../../site/static",
		"--css-dir", "../../site/css",
		"--js-dir", "../../site/js",
		"--output-dir", t.TempDir(),
	}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("Run returned %d; stderr: %s", code, stderr.String())
	}
}

func TestAllPagesBuilt(t *testing.T) {
	outDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := Run([]string{
		"--templates-dir", "../../site/templates",
		"--static-dir", "../../site/static",
		"--css-dir", "../../site/css",
		"--js-dir", "../../site/js",
		"--output-dir", outDir,
	}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("Run returned %d; stderr: %s", code, stderr.String())
	}

	expected := []string{
		"index.html",
		"about/index.html",
		"words-of-meaning/index.html",
		"roadmap/index.html",
	}

	for _, path := range expected {
		full := filepath.Join(outDir, path)
		if _, err := os.Stat(full); os.IsNotExist(err) {
			t.Errorf("expected output file %s does not exist", path)
		}
	}
}

func TestSharedElementsPresent(t *testing.T) {
	outDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := Run([]string{
		"--templates-dir", "../../site/templates",
		"--static-dir", "../../site/static",
		"--css-dir", "../../site/css",
		"--js-dir", "../../site/js",
		"--output-dir", outDir,
	}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("Run returned %d; stderr: %s", code, stderr.String())
	}

	pages := []string{
		"index.html",
		"about/index.html",
		"words-of-meaning/index.html",
		"roadmap/index.html",
	}

	for _, page := range pages {
		data, err := os.ReadFile(filepath.Join(outDir, page))
		if err != nil {
			t.Fatalf("reading %s: %v", page, err)
		}
		content := string(data)

		// All pages should have skip-link
		if !strings.Contains(content, "skip-link") {
			t.Errorf("%s: missing skip-link", page)
		}

		// All pages should have site-header nav
		if !strings.Contains(content, "site-header") {
			t.Errorf("%s: missing site-header", page)
		}

		// All pages should have the three nav links
		if !strings.Contains(content, "blog.nuphirho.dev") {
			t.Errorf("%s: missing Blog nav link", page)
		}
		if !strings.Contains(content, `href="/about"`) {
			t.Errorf("%s: missing About nav link", page)
		}
		if !strings.Contains(content, `href="/words-of-meaning"`) {
			t.Errorf("%s: missing Words nav link", page)
		}

		// No inline <style> blocks
		if strings.Contains(content, "<style>") {
			t.Errorf("%s: contains inline <style> block", page)
		}
	}
}

func TestThemeTogglePresence(t *testing.T) {
	outDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	Run([]string{
		"--templates-dir", "../../site/templates",
		"--static-dir", "../../site/static",
		"--css-dir", "../../site/css",
		"--js-dir", "../../site/js",
		"--output-dir", outDir,
	}, &stdout, &stderr)

	// All pages should have theme toggle
	for _, page := range []string{"index.html", "about/index.html", "words-of-meaning/index.html", "roadmap/index.html"} {
		data, _ := os.ReadFile(filepath.Join(outDir, page))
		if !strings.Contains(string(data), "theme-toggle") {
			t.Errorf("%s: missing theme toggle", page)
		}
	}
}

func TestActiveNavState(t *testing.T) {
	outDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	Run([]string{
		"--templates-dir", "../../site/templates",
		"--static-dir", "../../site/static",
		"--css-dir", "../../site/css",
		"--js-dir", "../../site/js",
		"--output-dir", outDir,
	}, &stdout, &stderr)

	// Landing page: no aria-current
	landing, _ := os.ReadFile(filepath.Join(outDir, "index.html"))
	if strings.Contains(string(landing), `aria-current="page"`) {
		t.Error("landing page should not have aria-current on any nav link")
	}

	// About page: aria-current on About link
	about, _ := os.ReadFile(filepath.Join(outDir, "about/index.html"))
	if !strings.Contains(string(about), `href="/about" aria-current="page"`) {
		t.Error("about page: About link should have aria-current")
	}

	// Words page: aria-current on Words link
	words, _ := os.ReadFile(filepath.Join(outDir, "words-of-meaning/index.html"))
	if !strings.Contains(string(words), `href="/words-of-meaning" aria-current="page"`) {
		t.Error("words page: Words link should have aria-current")
	}
}

func TestStaticAssetsCopied(t *testing.T) {
	outDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	Run([]string{
		"--templates-dir", "../../site/templates",
		"--static-dir", "../../site/static",
		"--css-dir", "../../site/css",
		"--js-dir", "../../site/js",
		"--output-dir", outDir,
	}, &stdout, &stderr)

	assets := []string{
		"CNAME",
		"roadmap/calendar-data.json",
		"css/styles.css",
		"css/roadmap.css",
		"js/theme.js",
		"js/roadmap.js",
	}

	for _, asset := range assets {
		full := filepath.Join(outDir, asset)
		if _, err := os.Stat(full); os.IsNotExist(err) {
			t.Errorf("expected static asset %s not found in output", asset)
		}
	}
}

func TestFooterPresence(t *testing.T) {
	outDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	Run([]string{
		"--templates-dir", "../../site/templates",
		"--static-dir", "../../site/static",
		"--css-dir", "../../site/css",
		"--js-dir", "../../site/js",
		"--output-dir", outDir,
	}, &stdout, &stderr)

	// About page should have footer
	about, _ := os.ReadFile(filepath.Join(outDir, "about/index.html"))
	if !strings.Contains(string(about), "site-footer") {
		t.Error("about page should have footer")
	}

	// Landing page should have footer
	landing, _ := os.ReadFile(filepath.Join(outDir, "index.html"))
	if !strings.Contains(string(landing), "site-footer") {
		t.Error("landing page should have footer")
	}
}

func TestPageSpecificCSS(t *testing.T) {
	outDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	Run([]string{
		"--templates-dir", "../../site/templates",
		"--static-dir", "../../site/static",
		"--css-dir", "../../site/css",
		"--js-dir", "../../site/js",
		"--output-dir", outDir,
	}, &stdout, &stderr)

	// Words page should NOT have page-specific CSS (styles in shared styles.css)
	words, _ := os.ReadFile(filepath.Join(outDir, "words-of-meaning/index.html"))
	if strings.Contains(string(words), "words-of-meaning.css") {
		t.Error("words page should not link to words-of-meaning.css (merged into styles.css)")
	}

	// Roadmap page should link to roadmap.css and roadmap.js
	roadmap, _ := os.ReadFile(filepath.Join(outDir, "roadmap/index.html"))
	if !strings.Contains(string(roadmap), "roadmap.css") {
		t.Error("roadmap page should link to roadmap.css")
	}
	if !strings.Contains(string(roadmap), "roadmap.js") {
		t.Error("roadmap page should link to roadmap.js")
	}
}
