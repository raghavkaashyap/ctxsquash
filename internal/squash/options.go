package squash

import (
	"path/filepath"
	"strings"
)

type Options struct {
	Root     string
	Output   string
	Include  []string
	Exclude  []string
	TreeOnly bool
}

func (o Options) normalized() (Options, error) {
	root, err := filepath.Abs(o.Root)
	if err != nil {
		return Options{}, err
	}

	o.Root = root
	if o.Output != "" {
		output, err := filepath.Abs(o.Output)
		if err != nil {
			return Options{}, err
		}
		o.Output = output
	}
	o.Include = normalizeExtensions(o.Include)
	o.Exclude = normalizeNames(o.Exclude)
	return o, nil
}

func SplitCSV(value string) []string {
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func normalizeExtensions(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(strings.ToLower(value))
		value = strings.TrimPrefix(value, ".")
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func normalizeNames(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}
