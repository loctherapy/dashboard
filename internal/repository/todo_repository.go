package repository

import (
	"regexp"

	"github.com/loctherapy/dashboard/internal/model"
)

type ToDoRepository struct {
	FileFetcher       *FileFetcher
	FrontMatterParser *FrontMatterParser
	ToDoExtractor     *ToDoExtractor
	Cache             *ToDoCache
}

func NewToDoRepository(fileFetcher *FileFetcher) *ToDoRepository {
	frontMatterRE := regexp.MustCompile(`(?m)^---\s*$`)
	contextRE := regexp.MustCompile(`^(?P<context_name>[a-zA-Z0-9_-]+)-(?P<context_gravity>\d+)$`)
	frontMatterParser := NewFrontMatterParser(frontMatterRE, contextRE)

	todoPattern := regexp.MustCompile(`^\s*- \[ \] `)
	todoExtractor := NewToDoExtractor(todoPattern)
	return &ToDoRepository{
		FileFetcher:       fileFetcher,
		FrontMatterParser: frontMatterParser,
		ToDoExtractor:     todoExtractor,
		Cache:             NewToDoCache(),
	}
}

func (r *ToDoRepository) GetAll() ([]model.FileToDos, error) {
	files, err := r.FileFetcher.Fetch()
	if err != nil {
		return nil, err
	}

	// Update cache with new or modified files
	for _, file := range files {
		cachedFile, exists := r.Cache.Get(file.Path)
		if exists && file.ModTime.Equal(cachedFile.FileInfo.ModTime) {
			continue
		}

		context, contextGravity, gravity, err := r.FrontMatterParser.Parse(file.Path)
		if err != nil {
			return nil, err
		}

		todos, err := r.ToDoExtractor.Extract(file.Path)
		if err != nil {
			return nil, err
		}

		if len(todos) == 0 {
			continue
		}

		r.Cache.Push(file.Path, file.ModTime, []model.FileToDos{
			{
				FilePath:       file.Path,
				ToDos:          todos,
				Context:        context,
				ContextGravity: contextGravity,
				Gravity:        gravity,
			},
		})
	}

	// Remove files from cache that are no longer present
	for filePath := range r.Cache.cache {
		found := false
		for _, file := range files {
			if file.Path == filePath {
				found = true
				break
			}
		}
		if !found {
			r.Cache.Delete(filePath)
		}
	}

	return r.Cache.Dump(), nil
}
