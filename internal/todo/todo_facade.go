package todo

import (
	"fmt"
)

type ToDoFacade struct {
	PrintMode PrintMode
	Pattern string
	fetcher *FileFetcher
	builder *ToDoBuilder
	printer *ToDoPrinter
}

func NewToDoFacade(pattern string, printMode PrintMode) (*ToDoFacade, error) {
	if pattern == "" {
		pattern = `.*\.md$`
	}

	fetcher := NewFileFetcher()
	builder := NewToDoBuilder()
	factory := ToDoPrinterFactory{}
	printer, err := factory.CreatePrinter(printMode)
	if err != nil {
		return nil, err
	}
	return &ToDoFacade{Pattern: pattern, fetcher: fetcher, builder: builder, printer: printer}, nil
}

func (f* ToDoFacade) Build() ([]FileToDos, error) {
	// Fetch .md files
	files, err := f.fetcher.Fetch(f.Pattern)
	if err != nil {
		fmt.Println("Error fetching files:", err)
		return nil, err
	}

	// Build todos
	todos, err := f.builder.Build(files)

	if err != nil {
		fmt.Println("Error building todos:", err)
		return nil, err
	}

	return todos, nil
}

func (f* ToDoFacade) Print() (string, error) {
	
	todos, err := f.Build()

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