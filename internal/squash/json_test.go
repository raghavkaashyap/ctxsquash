package squash

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestRenderJSON(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "simple-project")
	got, err := Render(Options{Root: root, Format: FormatJSON})
	if err != nil {
		t.Fatal(err)
	}

	var document jsonDocument
	if err := json.Unmarshal([]byte(got), &document); err != nil {
		t.Fatal(err)
	}
	if len(document.Tree) == 0 {
		t.Fatal("expected tree entries")
	}
	if len(document.Files) == 0 {
		t.Fatal("expected file entries")
	}
	if document.Files[0].Path != "README.md" || document.Files[0].Language != "markdown" {
		t.Fatalf("unexpected first file: %+v", document.Files[0])
	}
}

func TestRenderJSONTreeOnlyOmitsFiles(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "simple-project")
	got, err := Render(Options{Root: root, Format: FormatJSON, TreeOnly: true})
	if err != nil {
		t.Fatal(err)
	}

	var document jsonDocument
	if err := json.Unmarshal([]byte(got), &document); err != nil {
		t.Fatal(err)
	}
	if len(document.Tree) == 0 {
		t.Fatal("expected tree entries")
	}
	if len(document.Files) != 0 {
		t.Fatal("expected tree-only JSON to omit files")
	}
}

func TestRenderRejectsUnsupportedFormat(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "simple-project")
	if _, err := Render(Options{Root: root, Format: "xml"}); err == nil {
		t.Fatal("expected unsupported format error")
	}
}
