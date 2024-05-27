package todo

import (
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

type ToDoPrinter struct {
    Template *template.Template
}

type TemplateData struct {
    ContextMap map[string][]FileToDos
    PrintContext func(string) string
    PrintFile func(string, int) string
    PrintToDo func(string) string
}

func NewToDoPrinter() (*ToDoPrinter, error) {
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
    return &ToDoPrinter{Template: t}, nil
}

func printContext(context string) string {
    contextColor := color.New(color.FgHiYellow, color.Bold).SprintFunc()
    return contextColor(fmt.Sprintf("🌳 %s\n", strings.ToUpper(context)))
}

func printFile(filePath string, gravity int) string {
    fileColor := color.New(color.FgHiCyan).SprintFunc()
    return fileColor(fmt.Sprintf("  💎 %s (Gravity: %d)", filePath, gravity))
}

func printToDo(todo string) string {
    todoColor := color.New(color.FgWhite).SprintFunc()
    return todoColor(fmt.Sprintf("    %s", todo))
}

func (p *ToDoPrinter) Print(todos []FileToDos) error {
    // Group by context
    contextMap := make(map[string][]FileToDos)
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
        ContextMap: contextMap,
        PrintContext: printContext,
        PrintFile: printFile,
        PrintToDo: printToDo,
    }

    // Render the template
    var sb strings.Builder
    if err := p.Template.Execute(&sb, data); err != nil {
        return err
    }
    fmt.Print(sb.String())
    return nil
}
