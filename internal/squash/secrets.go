package squash

import (
	"regexp"
	"strings"
)

type SecretWarning struct {
	Path string
	Line int
	Kind string
}

var secretPatterns = []struct {
	kind    string
	pattern *regexp.Regexp
}{
	{
		kind:    "private key",
		pattern: regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`),
	},
	{
		kind:    "sensitive assignment",
		pattern: regexp.MustCompile(`(?i)\b(api[_-]?key|token|secret|password)\b\s*[:=]\s*['"]?[^\s'"]+`),
	},
}

func findSecretWarnings(files []File) []SecretWarning {
	var warnings []SecretWarning
	for _, file := range files {
		lines := strings.Split(file.Content, "\n")
		for lineNumber, line := range lines {
			for _, secretPattern := range secretPatterns {
				if secretPattern.pattern.MatchString(line) {
					warnings = append(warnings, SecretWarning{
						Path: file.Path,
						Line: lineNumber + 1,
						Kind: secretPattern.kind,
					})
					break
				}
			}
		}
	}
	return warnings
}
