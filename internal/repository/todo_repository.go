package repository

import (
	"bufio"
	"os"
	"regexp"

	"github.com/loctherapy/dashboard/internal/model"
)

type ToDoRepository struct {
	ToDoPattern   *regexp.Regexp
	FileFetcher   *FileFetcher
	FrontMatterParser *FrontMatterParser
}

func NewToDoRepository(fileFetcher *FileFetcher) *ToDoRepository {
	frontMatterRE := regexp.MustCompile(`(?m)^---\s*$`)
	contextRE := regexp.MustCompile(`^(?P<context_name>[a-zA-Z0-9_-]+)-(?P<context_gravity>\d+)$`)
	frontMatterParser := NewFrontMatterParser(frontMatterRE, contextRE)

	todoPattern := regexp.MustCompile(`^\s*- \[ \] `)
	
	return &ToDoRepository{
		ToDoPattern:   todoPattern,
		FrontMatterParser: frontMatterParser,
		FileFetcher:   fileFetcher,
	}
}

func (r *ToDoRepository) GetAll() ([]model.FileToDos, error) {
	var results []model.FileToDos

	files, err := r.FileFetcher.Fetch()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		context, contextGravity, gravity, err := r.FrontMatterParser.Parse(file.Path)
		if err != nil {
			return nil, err
		}

		todos, err := r.extractToDos(file.Path)
		if err != nil {
			return nil, err
		}

		if len(todos) > 0 {
			results = append(results, model.FileToDos{
				FilePath:       file.Path,
				ToDos:          todos,
				Context:        context,
				ContextGravity: contextGravity,
				Gravity:        gravity,
			})
		}
	}

	return results, nil
}

func (r *ToDoRepository) extractToDos(filePath string) ([]model.ToDo, error) {
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
