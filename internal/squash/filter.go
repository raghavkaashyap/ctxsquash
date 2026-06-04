package squash

import (
	"os"
	"path"
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
	include     map[string]bool
	exclude     map[string]bool
	ignoreRules []ignoreRule
}

type ignoreRule struct {
	pattern  string
	negate   bool
	dirOnly  bool
	hasSlash bool
}

func newFilter(options Options) (filter, error) {
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

	ignoreRules, err := loadGitIgnore(options.Root)
	if err != nil {
		return filter{}, err
	}
	f.ignoreRules = ignoreRules

	return f, nil
}

func loadGitIgnore(root string) ([]ignoreRule, error) {
	gitignorePath := filepath.Join(root, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var rules []ignoreRule
	for _, line := range strings.Split(string(content), "\n") {
		rule, ok := parseIgnoreRule(line)
		if ok {
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

func parseIgnoreRule(line string) (ignoreRule, bool) {
	line = trimUnescapedTrailingSpaces(strings.TrimSuffix(line, "\r"))
	if line == "" || strings.HasPrefix(line, "#") {
		return ignoreRule{}, false
	}

	rule := ignoreRule{}
	if strings.HasPrefix(line, "!") {
		rule.negate = true
		line = strings.TrimPrefix(line, "!")
	} else if strings.HasPrefix(line, `\#`) || strings.HasPrefix(line, `\!`) {
		line = strings.TrimPrefix(line, `\`)
	}

	line = strings.TrimPrefix(filepath.ToSlash(line), "/")
	line = strings.ReplaceAll(line, `\ `, " ")
	if strings.HasSuffix(line, "/") {
		rule.dirOnly = true
		line = strings.TrimSuffix(line, "/")
	}
	if line == "" {
		return ignoreRule{}, false
	}

	rule.pattern = line
	rule.hasSlash = strings.Contains(line, "/")
	return rule, true
}

func trimUnescapedTrailingSpaces(line string) string {
	end := len(line)
	for end > 0 && line[end-1] == ' ' && !escapedAt(line, end-1) {
		end--
	}
	return line[:end]
}

func escapedAt(value string, index int) bool {
	backslashes := 0
	for i := index - 1; i >= 0 && value[i] == '\\'; i-- {
		backslashes++
	}
	return backslashes%2 == 1
}

func (f filter) skipDir(name, rel string) bool {
	return f.exclude[name] || f.ignored(rel, true)
}

func (f filter) ignored(rel string, isDir bool) bool {
	ignored := false
	for _, rule := range f.ignoreRules {
		if rule.matches(rel, isDir) {
			ignored = !rule.negate
		}
	}
	return ignored
}

func (r ignoreRule) matches(rel string, isDir bool) bool {
	rel = path.Clean(filepath.ToSlash(rel))
	if rel == "." {
		return false
	}
	if r.dirOnly && !isDir && !strings.HasPrefix(rel, r.pattern+"/") {
		return false
	}

	if r.hasSlash {
		return matchIgnorePattern(r.pattern, rel) || strings.HasPrefix(rel, r.pattern+"/")
	}

	for {
		name := path.Base(rel)
		if matchIgnorePattern(r.pattern, name) {
			return true
		}
		parent := path.Dir(rel)
		if parent == "." || parent == rel {
			return false
		}
		rel = parent
	}
}

func matchIgnorePattern(pattern, value string) bool {
	matched, err := path.Match(pattern, value)
	return err == nil && matched
}

func (f filter) includeFile(path string) bool {
	if len(f.include) == 0 {
		return true
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	return f.include[ext]
}
