package todo

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type FileFetcher struct {

}

func NewFileFetcher() *FileFetcher {
    return &FileFetcher{}
}

func (f *FileFetcher) Fetch(pattern string) ([]string, error) {
    if pattern == "" {
        pattern = `.*\.md$`
    }

    var files []string
    re, err := regexp.Compile(pattern)
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
