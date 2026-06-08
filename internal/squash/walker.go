package squash

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/raghavkaashyap/ctxsquash/internal/fsutil"
)

type File struct {
	Path    string
	Content string
}

func collect(options Options) ([]string, []File, error) {
	info, err := os.Stat(options.Root)
	if err != nil {
		return nil, nil, err
	}
	if !info.IsDir() {
		return nil, nil, fmt.Errorf("root path must be a directory: %s", options.Root)
	}

	f, err := newFilter(options)
	if err != nil {
		return nil, nil, err
	}
	var treePaths []string
	var files []File

	err = filepath.WalkDir(options.Root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == options.Root {
			return nil
		}

		rel, err := fsutil.RelSlash(options.Root, path)
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if f.skipDir(entry.Name(), rel) {
				return filepath.SkipDir
			}
			treePaths = append(treePaths, rel+"/")
			return nil
		}

		if !entry.Type().IsRegular() || f.ignored(rel, false) || !f.includeFile(path) {
			return nil
		}
		if options.Output != "" && filepath.Clean(path) == filepath.Clean(options.Output) {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Size() > options.MaxFileSize {
			return nil
		}

		binary, err := fsutil.IsBinaryFile(path)
		if err != nil {
			return err
		}
		if binary {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		treePaths = append(treePaths, rel)
		files = append(files, File{Path: rel, Content: string(content)})
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	sort.Strings(treePaths)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	return treePaths, files, nil
}
