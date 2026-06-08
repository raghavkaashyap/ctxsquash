package squash

import "testing"

func TestFindSecretWarningsRedactsValues(t *testing.T) {
	files := []File{
		{
			Path:    "config.env",
			Content: "API_KEY=fake-value-for-tests\nsafe=true\n",
		},
	}

	warnings := findSecretWarnings(files)
	if len(warnings) != 1 {
		t.Fatalf("got %d warnings, want 1", len(warnings))
	}
	if warnings[0].Path != "config.env" || warnings[0].Line != 1 || warnings[0].Kind != "sensitive assignment" {
		t.Fatalf("unexpected warning: %+v", warnings[0])
	}
}

func TestFindSecretWarningsDetectsPrivateKeyHeader(t *testing.T) {
	files := []File{
		{
			Path:    "fixture.pem",
			Content: "-----BEGIN PRIVATE KEY-----\nfake body\n",
		},
	}

	warnings := findSecretWarnings(files)
	if len(warnings) != 1 {
		t.Fatalf("got %d warnings, want 1", len(warnings))
	}
	if warnings[0].Kind != "private key" {
		t.Fatalf("got %q, want private key", warnings[0].Kind)
	}
}
