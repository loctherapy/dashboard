package view

import (
	"sort"
	"strings"
	"text/template"

	"github.com/loctherapy/dashboard/internal/model"
)

type ToDoPrinter struct {
	Template      *template.Template
	PrintSettings IPrintSettings
}

type TemplateData struct {
	ContextMap   map[string][]model.FileToDos
	PrintContext func(string) string
	PrintFile    func(string, int) string
	PrintToDo    func(string) string
}

func NewToDoPrinter(printSettings IPrintSettings) (*ToDoPrinter, error) {
	const tmpl = `{{- range $context, $files := .ContextMap }}
{{ call $.PrintContext $context }}
{{- range $files }}
  {{ call $.PrintFile .FilePath .Gravity }}
  {{- range .ToDos }}
    {{ call $.PrintToDo .Line }}
  {{- end }}
{{ end }}
{{ end }}`
	t, err := template.New("todo").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	return &ToDoPrinter{Template: t, PrintSettings: printSettings}, nil
}

func (p *ToDoPrinter) Print(todos []model.FileToDos) (string, error) {
	// Group by context
	contextMap := make(map[string][]model.FileToDos)
	for _, fileToDos := range todos {
		contextMap[fileToDos.Context] = append(contextMap[fileToDos.Context], fileToDos)
	}

	// Sort files within each context by gravity
	for _, files := range contextMap {
		sort.Slice(files, func(i, j int) bool {
			return files[i].Gravity > files[j].Gravity
		})
	}

	// Create template data
	data := TemplateData{
		ContextMap:   contextMap,
		PrintContext: p.PrintSettings.PrintContext,
		PrintFile:    p.PrintSettings.PrintFile,
		PrintToDo:    p.PrintSettings.PrintToDo,
	}

	// Render the template
	var sb strings.Builder
	if err := p.Template.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}
