package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}

// Run builds the static site from templates into the output directory.
func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("site-build", flag.ContinueOnError)
	fs.SetOutput(stderr)

	templatesDir := fs.String("templates-dir", "site/templates", "Path to templates directory")
	staticDir := fs.String("static-dir", "site/static", "Path to static assets directory")
	cssDir := fs.String("css-dir", "site/css", "Path to CSS directory")
	jsDir := fs.String("js-dir", "site/js", "Path to JS directory")
	outputDir := fs.String("output-dir", "site/dist", "Path to output directory")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if err := build(*templatesDir, *staticDir, *cssDir, *jsDir, *outputDir, stdout); err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}

	return 0
}

func build(templatesDir, staticDir, cssDir, jsDir, outputDir string, stdout io.Writer) error {
	// Clean output directory
	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("cleaning output dir: %w", err)
	}

	// Parse base template and partials
	basePattern := filepath.Join(templatesDir, "base.html")
	partialsPattern := filepath.Join(templatesDir, "partials", "*.html")

	baseFiles, err := filepath.Glob(basePattern)
	if err != nil || len(baseFiles) == 0 {
		return fmt.Errorf("no base template found at %s", basePattern)
	}

	partialFiles, err := filepath.Glob(partialsPattern)
	if err != nil {
		return fmt.Errorf("reading partials: %w", err)
	}

	sharedFiles := append(baseFiles, partialFiles...)

	// Build each page
	for _, page := range pages() {
		pageTemplate := filepath.Join(templatesDir, "pages", page.TemplateName)

		allFiles := make([]string, len(sharedFiles)+1)
		copy(allFiles, sharedFiles)
		allFiles[len(sharedFiles)] = pageTemplate

		funcMap := template.FuncMap{
			"isCurrentNav": func(nav string) bool {
				return page.CurrentNav == nav
			},
			"joinCSS": func() template.HTML {
				var sb strings.Builder
				for _, css := range page.PageCSS {
					sb.WriteString(fmt.Sprintf("  <link rel=\"stylesheet\" href=\"/%s\">\n", css))
				}
				return template.HTML(sb.String())
			},
			"joinJS": func() template.HTML {
				var sb strings.Builder
				for _, js := range page.PageJS {
					sb.WriteString(fmt.Sprintf("  <script src=\"/%s\"></script>\n", js))
				}
				return template.HTML(sb.String())
			},
		}

		tmpl, err := template.New(filepath.Base(baseFiles[0])).Funcs(funcMap).ParseFiles(allFiles...)
		if err != nil {
			return fmt.Errorf("parsing templates for %s: %w", page.OutputPath, err)
		}

		outPath := filepath.Join(outputDir, page.OutputPath)
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", outPath, err)
		}

		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("creating %s: %w", outPath, err)
		}

		if err := tmpl.ExecuteTemplate(f, "base", page); err != nil {
			f.Close()
			return fmt.Errorf("executing template for %s: %w", page.OutputPath, err)
		}
		f.Close()

		fmt.Fprintf(stdout, "built %s\n", page.OutputPath)
	}

	// Copy static assets
	if err := copyDir(staticDir, outputDir); err != nil {
		return fmt.Errorf("copying static assets: %w", err)
	}

	// Copy CSS
	if err := copyDir(cssDir, filepath.Join(outputDir, "css")); err != nil {
		return fmt.Errorf("copying CSS: %w", err)
	}

	// Copy JS
	if err := copyDir(jsDir, filepath.Join(outputDir, "js")); err != nil {
		return fmt.Errorf("copying JS: %w", err)
	}

	fmt.Fprintf(stdout, "site build complete\n")
	return nil
}

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

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}

		return os.WriteFile(target, data, 0o644)
	})
}
