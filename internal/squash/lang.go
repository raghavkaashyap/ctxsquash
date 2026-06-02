package squash

import (
	"path/filepath"
	"strings"
)

var languageByExtension = map[string]string{
	".c":    "c",
	".cpp":  "cpp",
	".cs":   "csharp",
	".css":  "css",
	".go":   "go",
	".html": "html",
	".java": "java",
	".js":   "javascript",
	".json": "json",
	".jsx":  "jsx",
	".kt":   "kotlin",
	".md":   "markdown",
	".php":  "php",
	".py":   "python",
	".rb":   "ruby",
	".rs":   "rust",
	".sh":   "bash",
	".sql":  "sql",
	".ts":   "typescript",
	".tsx":  "tsx",
	".txt":  "text",
	".xml":  "xml",
	".yaml": "yaml",
	".yml":  "yaml",
}

func LanguageForPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if lang, ok := languageByExtension[ext]; ok {
		return lang
	}
	return "text"
}
