package squash

import "testing"

func TestSplitCSVTrimsEmptyValues(t *testing.T) {
	got := SplitCSV(" go, py,,md ")
	want := []string{"go", "py", "md"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestFilterIncludesExtensions(t *testing.T) {
	options := Options{Include: []string{"go", ".md"}}
	options.Include = normalizeExtensions(options.Include)
	f := newFilter(options)

	if !f.includeFile("main.go") {
		t.Fatal("expected go file to be included")
	}
	if f.includeFile("main.py") {
		t.Fatal("expected py file to be excluded")
	}
}
