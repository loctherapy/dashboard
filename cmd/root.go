package cmd

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"

	"github.com/loctherapy/dashboard/internal/todo"
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

func mainUI() {
    app := tview.NewApplication()

    // Create header
    header := tview.NewTextView().
        SetText("Dashboard").
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

    facade, err := todo.NewToDoFacade(`.*\.md$`)

    if err != nil {
        fmt.Println("Error printing todos:", err)
    }

    todosString, err := facade.Print()

    if err != nil {
        fmt.Println("Error printing todos:", err)
    }

    todoTextView.SetText(todosString)

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
