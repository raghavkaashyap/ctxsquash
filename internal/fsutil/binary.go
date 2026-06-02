package fsutil

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const sniffSize = 8192

var binaryExtensions = map[string]bool{
	".bin":   true,
	".class": true,
	".dll":   true,
	".exe":   true,
	".gif":   true,
	".ico":   true,
	".jpg":   true,
	".jpeg":  true,
	".pdf":   true,
	".png":   true,
	".so":    true,
	".webp":  true,
	".zip":   true,
}

func IsBinaryFile(path string) (bool, error) {
	if binaryExtensions[strings.ToLower(filepath.Ext(path))] {
		return true, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	buffer := make([]byte, sniffSize)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}
	if n == 0 {
		return false, nil
	}

	sample := buffer[:n]
	if bytes.IndexByte(sample, 0) >= 0 {
		return true, nil
	}
	return !utf8.Valid(sample), nil
}
