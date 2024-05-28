package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type FileFetcher struct {
    FileRegex *regexp.Regexp
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

func (f *FileFetcher) Fetch() ([]string, error) {

    var files []string

    err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && f.FileRegex.MatchString(info.Name()) {
            files = append(files, path)
        }
        return nil
    })

    if err != nil {
        return nil, err
    }

    return files, nil
}
