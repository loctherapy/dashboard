package controller

import (
	"fmt"
	"time"

	"github.com/loctherapy/dashboard/internal/service"
	"github.com/loctherapy/dashboard/internal/view"
)

type ToDoController struct {
    Service *service.ToDoService
	View *view.View
}

func NewToDoController(service *service.ToDoService, view *view.View) *ToDoController {
	return &ToDoController{
		Service: service,
		View: view,
	}
}

func (c *ToDoController) startGettingToDos() {
	updateTodos := func() {
		todos, err := c.Service.GetTodos()
		if err != nil {
			fmt.Println("Error getting todos:", err)
			return
		}
		c.View.DisplayToDos(todos)
	}

	// Initial fetch and update
	go updateTodos()

	// Create a ticker to update todos periodically
	ticker := time.NewTicker(1 * time.Second) // Adjust the interval as needed
	go func() {
		for {
			select {
			case <-ticker.C:
				updateTodos()
			}
		}
	}()
}

func (c *ToDoController) GetTodos() {
	c.startGettingToDos()
	c.View.RunUI()
}
