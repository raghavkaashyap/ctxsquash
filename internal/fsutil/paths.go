package fsutil

import (
	"path/filepath"
	"strings"
)

func RelSlash(root, path string) (string, error) {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(strings.TrimPrefix(rel, "./")), nil
}
