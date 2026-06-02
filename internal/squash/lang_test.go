package squash

import "testing"

func TestLanguageForPath(t *testing.T) {
	tests := map[string]string{
		"main.go":        "go",
		"README.md":      "markdown",
		"config.yml":     "yaml",
		"unknown.custom": "text",
	}

	for path, want := range tests {
		if got := LanguageForPath(path); got != want {
			t.Fatalf("LanguageForPath(%q) = %q, want %q", path, got, want)
		}
	}
}
