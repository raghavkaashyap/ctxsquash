package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCommandStdout(t *testing.T) {
	var out bytes.Buffer
	cmd := newRootCommand(&out)
	cmd.SetArgs([]string{filepath.Join("..", "..", "testdata", "simple-project"), "--stdout"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "# ctxsquash Context") {
		t.Fatal("expected markdown on stdout")
	}
}

func TestRootCommandOutputFile(t *testing.T) {
	var out bytes.Buffer
	output := filepath.Join(t.TempDir(), "context.md")
	cmd := newRootCommand(&out)
	cmd.SetArgs([]string{filepath.Join("..", "..", "testdata", "simple-project"), "--output", output})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Fatal("expected no stdout when output is set")
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "# ctxsquash Context") {
		t.Fatal("expected markdown output file")
	}
}

func TestRootCommandMaxFileSize(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "small.txt"), []byte("small\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "large.txt"), []byte("large file content\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	cmd := newRootCommand(&out)
	cmd.SetArgs([]string{root, "--max-file-size", "8", "--stdout"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "large.txt") {
		t.Fatal("expected oversized file to be skipped")
	}
	if !strings.Contains(out.String(), "small.txt") {
		t.Fatal("expected small file to be included")
	}
}

func TestRootCommandWarnsAboutPossibleSecrets(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "config.env"), []byte("API_KEY=fake-value-for-tests\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := newRootCommand(&out, &errOut)
	cmd.SetArgs([]string{root, "--stdout"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	warnings := errOut.String()
	if !strings.Contains(warnings, "warning: possible sensitive assignment in config.env:1") {
		t.Fatalf("expected warning on stderr, got %q", warnings)
	}
	if strings.Contains(warnings, "fake-value-for-tests") {
		t.Fatal("expected warning to redact matched value")
	}
	if !strings.Contains(out.String(), "# ctxsquash Context") {
		t.Fatal("expected markdown on stdout")
	}
}
