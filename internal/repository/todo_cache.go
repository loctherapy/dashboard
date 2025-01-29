package repository

import (
	"time"

	"github.com/loctherapy/dashboard/internal/model"
)

type FileInfoTodos struct {
	FileInfo FileInfo
	Todos []model.FileToDos
}

type ToDoCache struct {
	cache map[string]FileInfoTodos
}

func NewToDoCache() *ToDoCache {
	return &ToDoCache{
		cache: make(map[string]FileInfoTodos),
	}
}

func (c *ToDoCache) Push(filePath string, fileModTime time.Time, todos []model.FileToDos) {
	fileInfo := FileInfo{Path: filePath, ModTime: fileModTime}
	fileInfoTodos := FileInfoTodos{FileInfo: fileInfo, Todos: todos}
	c.cache[filePath] = fileInfoTodos
}

func (c *ToDoCache) Get(filePath string) (FileInfoTodos, bool) {
    fileInfoTodos, exists := c.cache[filePath]
    return fileInfoTodos, exists
}

func (c *ToDoCache) Delete(filePath string) {
    delete(c.cache, filePath)
}