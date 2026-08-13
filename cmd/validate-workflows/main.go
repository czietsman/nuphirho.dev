package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}

func Run(args []string, stdout, stderr io.Writer) int {
	dir := ".github/workflows"
	if len(args) > 0 {
		dir = args[0]
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(stderr, "error: cannot read directory %s: %v\n", dir, err)
		return 1
	}

	failed := false
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".yml" {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(stderr, "error: cannot read %s: %v\n", path, err)
			failed = true
			continue
		}
		var node yaml.Node
		if err := yaml.Unmarshal(data, &node); err != nil {
			fmt.Fprintf(stderr, "invalid YAML in %s: %v\n", path, err)
			failed = true
		}
	}

	if failed {
		return 1
	}
	fmt.Fprintln(stdout, "all workflow files are valid YAML")
	return 0
}
