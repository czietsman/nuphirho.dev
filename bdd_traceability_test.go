package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type bddReference struct {
	path  string
	label string
}

// BDD: specs/repository_validation.feature :: Scenario: Repository tests declare BDD traceability annotations
// BDD: specs/repository_validation.feature :: Scenario: Repository test annotations reference existing BDD specifications
func TestBDDTraceability(t *testing.T) {
	t.Parallel()

	testFiles, err := filepath.Glob("*_test.go")
	if err != nil {
		t.Fatalf("glob root test files: %v", err)
	}

	err = filepath.Walk(".", func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			if strings.HasPrefix(path, ".git") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, "_test.go") && path != "./bdd_traceability_test.go" {
			testFiles = append(testFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk test files: %v", err)
	}

	fset := token.NewFileSet()
	for _, path := range testFiles {
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}

		for _, decl := range node.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}

			refs := parseBDDReferences(fn.Doc)
			if len(refs) == 0 {
				t.Fatalf("%s:%s is missing a BDD traceability comment", path, fn.Name.Name)
			}

			for _, ref := range refs {
				content, err := os.ReadFile(ref.path)
				if err != nil {
					t.Fatalf("%s:%s references unreadable spec %s: %v", path, fn.Name.Name, ref.path, err)
				}
				if !strings.Contains(string(content), ref.label) {
					t.Fatalf("%s:%s references missing BDD label %q in %s", path, fn.Name.Name, ref.label, ref.path)
				}
			}
		}
	}
}

func parseBDDReferences(group *ast.CommentGroup) []bddReference {
	if group == nil {
		return nil
	}

	var refs []bddReference
	for _, comment := range group.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if !strings.HasPrefix(text, "BDD: ") {
			continue
		}

		refText := strings.TrimPrefix(text, "BDD: ")
		parts := strings.SplitN(refText, " :: ", 2)
		if len(parts) != 2 {
			continue
		}

		refs = append(refs, bddReference{
			path:  parts[0],
			label: parts[1],
		})
	}

	return refs
}
