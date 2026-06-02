package squash

import (
	"path/filepath"
	"strings"
)

var defaultExcludedDirs = map[string]bool{
	".git":         true,
	".hg":          true,
	".svn":         true,
	".idea":        true,
	".vscode":      true,
	"node_modules": true,
	"target":       true,
	"dist":         true,
	"build":        true,
	"vendor":       true,
}

type filter struct {
	include map[string]bool
	exclude map[string]bool
}

func newFilter(options Options) filter {
	f := filter{
		include: map[string]bool{},
		exclude: map[string]bool{},
	}

	for name := range defaultExcludedDirs {
		f.exclude[name] = true
	}
	for _, name := range options.Exclude {
		f.exclude[name] = true
	}
	for _, ext := range options.Include {
		f.include[ext] = true
	}

	return f
}

func (f filter) skipDir(name string) bool {
	return f.exclude[name]
}

func (f filter) includeFile(path string) bool {
	if len(f.include) == 0 {
		return true
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	return f.include[ext]
}
