package todo

import (
	"fmt"
	"sort"
	"strings"
	"text/template"
)

type ToDoPrinter struct {
    Template *template.Template
}

func NewToDoPrinter() (*ToDoPrinter, error) {
    const tmpl = `{{- range $context, $files := . }}
Context: {{ $context }}
{{- range $files }}
  File: {{ .FilePath }} (Gravity: {{ .Gravity }})
  {{- range .ToDos }}
    {{ .Line }}
  {{- end }}
{{- end }}
{{ end }}`
    t, err := template.New("todo").Parse(tmpl)
    if err != nil {
        return nil, err
    }
    return &ToDoPrinter{Template: t}, nil
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

    // Render the template
    var sb strings.Builder
    if err := p.Template.Execute(&sb, contextMap); err != nil {
        return err
    }
    fmt.Print(sb.String())
    return nil
}
