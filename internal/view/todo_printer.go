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
	Contexts     []ContextData
	PrintContext func(string) string
	PrintFile    func(string, int) string
	PrintToDo    func(string) string
}

type ContextData struct {
	Name           string
	Files          []model.FileToDos
	ContextGravity int
}

func NewToDoPrinter(printSettings IPrintSettings) (*ToDoPrinter, error) {
	const tmpl = `{{- range .Contexts }}
{{ call $.PrintContext .Name }}
{{- range .Files }}
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
	contextGravityMap := make(map[string]int)
	for _, fileToDos := range todos {
		contextMap[fileToDos.Context] = append(contextMap[fileToDos.Context], fileToDos)
		contextGravityMap[fileToDos.Context] = fileToDos.ContextGravity
	}

	// Sort files within each context by gravity
	for _, files := range contextMap {
		sort.Slice(files, func(i, j int) bool {
			return files[i].Gravity > files[j].Gravity
		})
	}

	// Create a slice of contexts with their gravities
	var contextDataList []ContextData
	for context, files := range contextMap {
		contextDataList = append(contextDataList, ContextData{
			Name:           context,
			Files:          files,
			ContextGravity: contextGravityMap[context],
		})
	}

	// Sort the contexts by their gravity
	sort.Slice(contextDataList, func(i, j int) bool {
		return contextDataList[i].ContextGravity < contextDataList[j].ContextGravity
	})

	// Create template data
	data := TemplateData{
		Contexts:     contextDataList,
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
