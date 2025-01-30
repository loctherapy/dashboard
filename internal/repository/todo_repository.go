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
	cacheInitialized  bool
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
		cacheInitialized:  false,
	}
}

func (r *ToDoRepository) initilizeCache() (error) {
	files, err := r.FileFetcher.Fetch()
	if err != nil {
		return err
	}

	// Initialize the cache
	for _, file := range files {

		err := r.parseExtractPushToCache(file)
		if err != nil {
			return err
		}
	}

	r.cacheInitialized = true

	return nil
}

func (r *ToDoRepository) GetAll() ([]model.FileToDos, error) {
	if !r.cacheInitialized {
		err := r.initilizeCache()
		if err != nil {
			return nil, err
		}
	}

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

		err := r.parseExtractPushToCache(file)
		if err != nil {
			return nil, err
		}
	}

	// Remove files from cache that are no longer present
	for filePath := range r.Cache.cache {
		if _, found := files[filePath]; !found {
			r.Cache.Delete(filePath)
		}
	}

	return r.Cache.Dump(), nil
}

func (r *ToDoRepository) parseExtractPushToCache(file FileInfo) error {
	context, contextGravity, gravity, err := r.FrontMatterParser.Parse(file.Path)
	if err != nil {
		return err
	}

	todos, err := r.ToDoExtractor.Extract(file.Path)
	if err != nil {
		return err
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

	return nil
}