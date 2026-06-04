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
	f, err := newFilter(options)
	if err != nil {
		t.Fatal(err)
	}

	if !f.includeFile("main.go") {
		t.Fatal("expected go file to be included")
	}
	if f.includeFile("main.py") {
		t.Fatal("expected py file to be excluded")
	}
}

func TestParseIgnoreRulePreservesEscapedPrefixes(t *testing.T) {
	rule, ok := parseIgnoreRule(`\#secret.txt`)
	if !ok {
		t.Fatal("expected escaped comment pattern to parse")
	}
	if rule.pattern != "#secret.txt" || rule.negate {
		t.Fatalf("got %+v, want literal #secret.txt", rule)
	}

	rule, ok = parseIgnoreRule(`\!important.txt`)
	if !ok {
		t.Fatal("expected escaped negation pattern to parse")
	}
	if rule.pattern != "!important.txt" || rule.negate {
		t.Fatalf("got %+v, want literal !important.txt", rule)
	}
}

func TestParseIgnoreRuleHandlesEscapedSpaces(t *testing.T) {
	rule, ok := parseIgnoreRule(`name\ `)
	if !ok {
		t.Fatal("expected escaped trailing space pattern to parse")
	}
	if rule.pattern != "name " {
		t.Fatalf("got %q, want trailing space", rule.pattern)
	}

	rule, ok = parseIgnoreRule(`name   `)
	if !ok {
		t.Fatal("expected trailing space pattern to parse")
	}
	if rule.pattern != "name" {
		t.Fatalf("got %q, want unescaped trailing spaces trimmed", rule.pattern)
	}
}

func TestParseIgnoreRulePreservesLeadingSpaces(t *testing.T) {
	rule, ok := parseIgnoreRule(` leading.txt`)
	if !ok {
		t.Fatal("expected leading space pattern to parse")
	}
	if rule.pattern != " leading.txt" {
		t.Fatalf("got %q, want leading space preserved", rule.pattern)
	}
}
