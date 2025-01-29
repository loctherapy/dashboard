package repository

import (
	"regexp"

	"github.com/loctherapy/dashboard/internal/model"
)

type ToDoRepository struct {
	FileFetcher   *FileFetcher
	FrontMatterParser *FrontMatterParser
	ToDoExtractor *ToDoExtractor
}

func NewToDoRepository(fileFetcher *FileFetcher) *ToDoRepository {
	frontMatterRE := regexp.MustCompile(`(?m)^---\s*$`)
	contextRE := regexp.MustCompile(`^(?P<context_name>[a-zA-Z0-9_-]+)-(?P<context_gravity>\d+)$`)
	frontMatterParser := NewFrontMatterParser(frontMatterRE, contextRE)

	todoPattern := regexp.MustCompile(`^\s*- \[ \] `)
	todoExtractor := NewToDoExtractor(todoPattern)
	return &ToDoRepository{
		FileFetcher:   fileFetcher,
		FrontMatterParser: frontMatterParser,
		ToDoExtractor: todoExtractor,
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

		todos, err := r.ToDoExtractor.Extract(file.Path)
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
