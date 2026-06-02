package squash

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderSimpleProjectDeterministic(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "simple-project")
	first, err := Render(Options{Root: root})
	if err != nil {
		t.Fatal(err)
	}
	second, err := Render(Options{Root: root})
	if err != nil {
		t.Fatal(err)
	}

	if first != second {
		t.Fatal("expected deterministic output")
	}
	want := `# ctxsquash Context

## Project Tree

` + "```text" + `
.
- README.md
- assets/
- config.yml
- docs/
  - docs/notes.md
- src/
  - src/main.go
` + "```" + `

## Files

### README.md

` + "```markdown" + `
# Simple Project

Small fixture for ctxsquash tests.
` + "```" + `

### config.yml

` + "```yaml" + `
name: simple-project
enabled: true
` + "```" + `

### docs/notes.md

` + "```markdown" + `
# Notes

These notes are safe fixture data.
` + "```" + `

### src/main.go

` + "```go" + `
package main

import "fmt"

func main() {
	fmt.Println("hello")
}
` + "```" + `
`
	if first != want {
		t.Fatalf("unexpected output:\n%s", first)
	}
	if strings.Contains(first, "node_modules") {
		t.Fatal("expected excluded directories to be skipped")
	}
	if strings.Contains(first, "logo.bin") {
		t.Fatal("expected binary files to be skipped")
	}
	if !strings.Contains(first, "```go\npackage main") {
		t.Fatal("expected go code fence")
	}
}

func TestRenderTreeOnlyExcludesFileContents(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "simple-project")
	got, err := Render(Options{Root: root, TreeOnly: true})
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(got, "## Files") || strings.Contains(got, "package main") {
		t.Fatal("expected tree-only output to omit file contents")
	}
	if !strings.Contains(got, "## Project Tree") {
		t.Fatal("expected project tree")
	}
}

func TestRenderIncludeFilter(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "simple-project")
	got, err := Render(Options{Root: root, Include: []string{"md"}})
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(got, "src/main.go") {
		t.Fatal("expected go file to be excluded")
	}
	if !strings.Contains(got, "README.md") {
		t.Fatal("expected markdown file to be included")
	}
}

func TestRenderRequiresDirectory(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(file, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := Render(Options{Root: file}); err == nil {
		t.Fatal("expected error for file root")
	}
}
