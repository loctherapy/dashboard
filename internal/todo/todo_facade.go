package todo

import (
	"fmt"
)

type ToDoFacade struct {
	Pattern string
	fetcher *FileFetcher
	builder *ToDoBuilder
	printer *ToDoPrinter
}

func NewToDoFacade(pattern string) (*ToDoFacade, error) {
	if pattern == "" {
		pattern = `.*\.md$`
	}

	fetcher := NewFileFetcher()
	builder := NewToDoBuilder()
	printer, err := NewToDoPrinter()
	if err != nil {
		return nil, err
	}
	return &ToDoFacade{Pattern: pattern, fetcher: fetcher, builder: builder, printer: printer}, nil
}

func (f* ToDoFacade) Print() (string, error) {
	// Fetch .md files

	files, err := f.fetcher.Fetch(f.Pattern)
	if err != nil {
		fmt.Println("Error fetching files:", err)
		return "", err
	}

	// Build todos
	todos, err := f.builder.Build(files)
	if err != nil {
		fmt.Println("Error building todos:", err)
		return "", err
	}

	// Print results using ToDoPrinter
	todosString, err := f.printer.Print(todos)
	if err != nil {
		fmt.Println("Error printing todos", err)
		return "", err
	}

	return todosString, nil
}