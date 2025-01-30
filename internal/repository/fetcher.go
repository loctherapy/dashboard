package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type FileFetcher struct {
	FileRegex *regexp.Regexp
}

type FileInfo struct {
	Path    string
	ModTime time.Time
}

func NewFileFetcher(pattern string) (*FileFetcher, error) {
	if pattern == "" {
		pattern = `.*\.md$`
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %v", err)
	}

	return &FileFetcher{FileRegex: re}, nil
}

func (f *FileFetcher) Fetch() (map[string]FileInfo, error) {
	files := make(map[string]FileInfo)

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && f.FileRegex.MatchString(info.Name()) {
			files[path] = FileInfo{Path: path, ModTime: info.ModTime()}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
