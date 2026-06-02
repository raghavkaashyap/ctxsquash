package squash

import (
	"strings"
)

func Render(options Options) (string, error) {
	treePaths, files, err := collect(options)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.WriteString("# ctxsquash Context\n\n")
	b.WriteString(renderTree(treePaths))

	if options.TreeOnly {
		return b.String(), nil
	}

	b.WriteString("\n## Files\n")
	for _, file := range files {
		b.WriteString("\n### ")
		b.WriteString(file.Path)
		b.WriteString("\n\n")
		b.WriteString("```")
		b.WriteString(LanguageForPath(file.Path))
		b.WriteByte('\n')
		b.WriteString(strings.TrimRight(file.Content, "\n"))
		b.WriteString("\n```\n")
	}

	return b.String(), nil
}
