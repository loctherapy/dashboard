package fetcher

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type FileFetcher struct {
    Pattern string
}

func NewFileFetcher(pattern string) *FileFetcher {
    if pattern == "" {
        pattern = `.*\.md$`
    }
    return &FileFetcher{Pattern: pattern}
}

func (f *FileFetcher) Fetch() ([]string, error) {
    var files []string
    re, err := regexp.Compile(f.Pattern)
    if err != nil {
        return nil, fmt.Errorf("invalid regex pattern: %v", err)
    }

    err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && re.MatchString(info.Name()) {
            files = append(files, path)
        }
        return nil
    })

    if err != nil {
        return nil, err
    }

    return files, nil
}
