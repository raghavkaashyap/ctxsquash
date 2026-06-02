package squash

import (
	"strings"
)

func renderTree(paths []string) string {
	var b strings.Builder
	b.WriteString("## Project Tree\n\n")
	b.WriteString("```text\n")
	b.WriteString(".\n")
	for _, path := range paths {
		depth := strings.Count(strings.TrimSuffix(path, "/"), "/")
		b.WriteString(strings.Repeat("  ", depth))
		b.WriteString("- ")
		b.WriteString(path)
		b.WriteByte('\n')
	}
	b.WriteString("```\n")
	return b.String()
}
