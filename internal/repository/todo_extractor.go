package repository

import (
	"bufio"
	"os"
	"regexp"

	"github.com/loctherapy/dashboard/internal/model"
)

type ToDoExtractor struct {
	ToDoPattern   *regexp.Regexp
}

func NewToDoExtractor(todoPattern *regexp.Regexp) *ToDoExtractor {
	return &ToDoExtractor{
		ToDoPattern: todoPattern,
	}
}

func (r *ToDoExtractor) Extract(filePath string) ([]model.ToDo, error) {
	var todos []model.ToDo

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if r.ToDoPattern.MatchString(line) {
			todos = append(todos, model.ToDo{Line: line})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
