package todo

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ToDoBuilder struct {
    ToDoPattern   *regexp.Regexp
    FrontMatterRE *regexp.Regexp
}

func NewToDoBuilder() *ToDoBuilder {
    todoPattern := regexp.MustCompile(`^\s*- \[ \] `)
    frontMatterRE := regexp.MustCompile(`(?m)^---\s*$`)
    return &ToDoBuilder{
        ToDoPattern:   todoPattern,
        FrontMatterRE: frontMatterRE,
    }
}


func (b *ToDoBuilder) Build(files []string) ([]FileToDos, error) {
    var results []FileToDos

    for _, file := range files {
        context, gravity, err := b.parseFrontMatter(file)
        if err != nil {
            return nil, err
        }

        todos, err := b.extractToDos(file)
        if err != nil {
            return nil, err
        }

        if len(todos) > 0 {
            results = append(results, FileToDos{
                FilePath: file,
                ToDos:    todos,
                Context:  context,
                Gravity:  gravity,
            })
        }
    }

    return results, nil
}

func (b *ToDoBuilder) parseFrontMatter(filePath string) (string, int, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var context string
    var gravity int
    inFrontMatter := false

    for scanner.Scan() {
        line := scanner.Text()
        if b.FrontMatterRE.MatchString(line) {
            if inFrontMatter {
                // End of front matter
                break
            } else {
                // Start of front matter
                inFrontMatter = true
                continue
            }
        }

        if inFrontMatter {
            if strings.HasPrefix(line, "context:") {
                context = strings.TrimSpace(strings.TrimPrefix(line, "context:"))
            } else if strings.HasPrefix(line, "gravity:") {
                gravityStr := strings.TrimSpace(strings.TrimPrefix(line, "gravity:"))
                gravity, _ = strconv.Atoi(gravityStr)
            }
        }
    }

    if err := scanner.Err(); err != nil {
        return "", 0, err
    }

    return context, gravity, nil
}

func (b *ToDoBuilder) extractToDos(filePath string) ([]ToDo, error) {
    var todos []ToDo

    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if b.ToDoPattern.MatchString(line) {
            todos = append(todos, ToDo{Line: line})
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return todos, nil
}
