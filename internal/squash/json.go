package squash

import (
	"encoding/json"
)

type jsonDocument struct {
	Tree  []string   `json:"tree"`
	Files []jsonFile `json:"files,omitempty"`
}

type jsonFile struct {
	Path     string `json:"path"`
	Language string `json:"language"`
	Content  string `json:"content"`
}

func renderJSON(treePaths []string, files []File, treeOnly bool) (string, error) {
	document := jsonDocument{Tree: treePaths}
	if !treeOnly {
		for _, file := range files {
			document.Files = append(document.Files, jsonFile{
				Path:     file.Path,
				Language: LanguageForPath(file.Path),
				Content:  file.Content,
			})
		}
	}

	content, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return "", err
	}
	return string(content) + "\n", nil
}
