package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/loctherapy/dashboard/internal/todo"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Dashboard is a CLI tool for managing your application",
	Long:  `A longer description of Dashboard with usage examples and details.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := tview.NewApplication()
		header := createHeader()
		todoTextView := createTodoTextView(app)
		flex := createLayout(header, todoTextView)

		setupKeybindings(app)
		startUpdatingTodos(app, todoTextView)

		runUI(app, flex)
	},
}

func createHeader() *tview.TextView {
	return tview.NewTextView().
		SetText("Dashboard").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetDynamicColors(true).
		SetRegions(true)
}

func createTodoTextView(app *tview.Application) *tview.TextView {
	todoTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	return todoTextView
}

func createLayout(header, todoTextView *tview.TextView) *tview.Flex {
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(todoTextView, 0, 1, true)
}

func setupKeybindings(app *tview.Application) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			app.Stop()
		}
		return event
	})
}

func startUpdatingTodos(app *tview.Application, todoTextView *tview.TextView) {
	updateTodos := func() {
		todosString, err := getToDoString()
		if err != nil {
			fmt.Println("Error getting todos:", err)
			return
		}
		app.QueueUpdateDraw(func() {
			todoTextView.SetText(todosString)
		})
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

func getToDoString() (string, error) {
	facade, err := todo.NewToDoFacade(`.*\.md$`, todo.TView)
	if err != nil {
		return "", err
	}

	todosString, err := facade.Print()
	if err != nil {
		return "", err
	}

	return todosString, nil
}

func runUI(app *tview.Application, flex *tview.Flex) {
	// Start the application
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
