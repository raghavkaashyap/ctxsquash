package squash

import (
	"strings"
)

func Render(options Options) (string, error) {
	result, _, err := RenderWithWarnings(options)
	return result, err
}

func RenderWithWarnings(options Options) (string, []SecretWarning, error) {
	options, err := options.normalized()
	if err != nil {
		return "", nil, err
	}

	treePaths, files, err := collect(options)
	if err != nil {
		return "", nil, err
	}

	warnings := findSecretWarnings(files)
	switch options.Format {
	case FormatJSON:
		result, err := renderJSON(treePaths, files, options.TreeOnly)
		return result, warnings, err
	default:
		return renderMarkdown(treePaths, files, options.TreeOnly), warnings, nil
	}
}

func renderMarkdown(treePaths []string, files []File, treeOnly bool) string {
	var b strings.Builder
	b.WriteString("# ctxsquash Context\n\n")
	b.WriteString(renderTree(treePaths))

	if treeOnly {
		return b.String()
	}

	b.WriteString("\n## Files\n")
	for _, file := range files {
		b.WriteString("\n### ")
		b.WriteString(file.Path)
		b.WriteString("\n\n")
		fence := markdownFence(file.Content)
		b.WriteString(fence)
		b.WriteString(LanguageForPath(file.Path))
		b.WriteByte('\n')
		b.WriteString(strings.TrimRight(file.Content, "\n"))
		b.WriteByte('\n')
		b.WriteString(fence)
		b.WriteByte('\n')
	}

	return b.String()
}

func markdownFence(content string) string {
	longest := 0
	current := 0
	for _, char := range content {
		if char == '`' {
			current++
			if current > longest {
				longest = current
			}
			continue
		}
		current = 0
	}

	if longest < 3 {
		return "```"
	}
	return strings.Repeat("`", longest+1)
}
