package todo

import (
	"bufio"
	"os"
	"regexp"
)

type ToDoBuilder struct {
    ToDoPattern *regexp.Regexp
}

func NewToDoBuilder() *ToDoBuilder {
    pattern := regexp.MustCompile(`^\s*- \[ \] `)
    return &ToDoBuilder{ToDoPattern: pattern}
}

func (b *ToDoBuilder) Build(files []string) ([]FileToDos, error) {
    var results []FileToDos

    for _, file := range files {
        todos, err := b.extractToDos(file)
        if err != nil {
            return nil, err
        }
        if len(todos) > 0 {
            results = append(results, FileToDos{
                FilePath: file,
                ToDos:    todos,
            })
        }
    }

    return results, nil
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
