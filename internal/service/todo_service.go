package service

import (
	"github.com/loctherapy/dashboard/internal/model"
	"github.com/loctherapy/dashboard/internal/repository"
)

type ToDoService struct {
    Repo *repository.ToDoRepository
}

func NewToDoService(repo *repository.ToDoRepository) *ToDoService {
	return &ToDoService{Repo: repo}
}

func (s *ToDoService) GetTodos() ([]model.FileToDos, error) {
    return s.Repo.GetAll()
}
