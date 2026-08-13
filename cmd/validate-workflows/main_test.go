package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRun_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	var out, errOut bytes.Buffer
	code := Run([]string{dir}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, errOut.String())
	}
}

func TestRun_AllValidYAML(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "valid.yml", "name: CI\non:\n  push:\n    branches: [main]\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo hello\n")
	var out, errOut bytes.Buffer
	code := Run([]string{dir}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, errOut.String())
	}
}

func TestRun_InvalidYAML_ColonInRunCommand(t *testing.T) {
	dir := t.TempDir()
	// Reproduces the exact bug that broke blog deploys: a colon-space inside an unquoted run value.
	writeFile(t, dir, "broken.yml", "jobs:\n  deploy:\n    steps:\n      - run: ./notify \"status: ${URL}\"\n")
	var out, errOut bytes.Buffer
	code := Run([]string{dir}, &out, &errOut)
	if code == 0 {
		t.Fatal("expected non-zero exit for invalid YAML, got 0")
	}
	if errOut.Len() == 0 {
		t.Fatal("expected error output, got none")
	}
}

func TestRun_MixedFiles_FailsOnInvalid(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "good.yml", "name: OK\non:\n  push:\n    branches: [main]\n")
	writeFile(t, dir, "bad.yml", "run: ./cmd \"with colon: here\"\n")
	var out, errOut bytes.Buffer
	code := Run([]string{dir}, &out, &errOut)
	if code == 0 {
		t.Fatal("expected non-zero exit when any file is invalid")
	}
}

func TestRun_NonExistentDirectory(t *testing.T) {
	var out, errOut bytes.Buffer
	code := Run([]string{"/nonexistent/path"}, &out, &errOut)
	if code == 0 {
		t.Fatal("expected non-zero exit for non-existent directory")
	}
}

func TestRun_IgnoresNonYAMLFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "README.md", "# not yaml")
	writeFile(t, dir, "script.sh", "#!/bin/bash\necho hello: world")
	var out, errOut bytes.Buffer
	code := Run([]string{dir}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected exit 0 when no .yml files present, got %d", code, )
	}
}
