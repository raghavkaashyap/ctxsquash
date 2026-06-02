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
