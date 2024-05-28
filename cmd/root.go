package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "dashboard",
    Short: "Dashboard is a CLI tool for managing your application",
    Long:  `A longer description of Dashboard with usage examples and details.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Root command logic
        fmt.Println("Welcome to the Dashboard CLI!")
        mainUI()
    },
}

type ToDo struct {
    Line string
}

type FileToDos struct {
    FilePath string
    ToDos    []ToDo
    Context  string
    Gravity  int
}

func fetchTodos() []FileToDos {
    // Dummy data, replace with your fetching logic
    return []FileToDos{
        {
            FilePath: "/path/to/file1.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 1"},
                {Line: "- [ ] Example TODO 2"},
            },
            Context: "GP",
            Gravity: 2,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
        {
            FilePath: "/path/to/file2.md",
            ToDos: []ToDo{
                {Line: "- [ ] Example TODO 3"},
                {Line: "- [ ] Example TODO 4"},
            },
            Context: "GP",
            Gravity: 3,
        },
    }
}

func mainUI() {
    app := tview.NewApplication()

    // Create header
    header := tview.NewTextView().
        SetText("TODO Application").
        SetTextAlign(tview.AlignCenter).
        SetTextColor(tview.Styles.PrimaryTextColor).
        SetDynamicColors(true).
        SetRegions(true)

    // Create text view for todo list
    todoTextView := tview.NewTextView().
        SetDynamicColors(true).
        SetRegions(true).
        SetWordWrap(true).
        SetChangedFunc(func() {
            app.Draw()
        })

    todos := fetchTodos()

    // Populate the text view with todos
    var sb strings.Builder
    for _, fileToDos := range todos {
        sb.WriteString(fmt.Sprintf("[yellow::b]CONTEXT: %s\n", strings.ToUpper(fileToDos.Context)))
        sb.WriteString(fmt.Sprintf("[cyan::b]  File: %s (Gravity: %d)\n", fileToDos.FilePath, fileToDos.Gravity))
        for _, todo := range fileToDos.ToDos {
            sb.WriteString(fmt.Sprintf("    %s\n", todo.Line))
        }
        sb.WriteString("\n")
    }
    todoTextView.SetText(sb.String())

    // Create flex layout for header and todo list
    flex := tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(header, 3, 1, false).
        AddItem(todoTextView, 0, 1, true)

    // Set keybinding to exit application
    app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyCtrlC {
            app.Stop()
        }
        return event
    })

    // Start application
    if err := app.SetRoot(flex, true).Run(); err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    // Define flags and configuration settings here if needed
}
