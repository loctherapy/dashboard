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

		// Define a function to fetch and update todos
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

		// Run the application with the main UI
		mainUI(app, flex)
	},
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

func mainUI(app *tview.Application, flex *tview.Flex) {
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
