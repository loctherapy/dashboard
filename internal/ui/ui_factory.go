package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/loctherapy/dashboard/internal/todo"
	"github.com/rivo/tview"
)

type UIFactory struct {
	app           *tview.Application
	header        *tview.TextView
	buttonsFlex    *tview.Flex
	todoTextView  *tview.TextView
	mainContainer *tview.Flex
}

func NewUIFactory() *UIFactory {
	app := tview.NewApplication()
	header := createMainHeader()
	buttonFlex := createButtons(app)
	todoTextView := createTodoTextView(app)
	mainContainer := createLayout(header, buttonFlex, todoTextView)

	setupKeybindings(app)

	return &UIFactory{app, header, buttonFlex, todoTextView, mainContainer}
}

func createMainHeader() *tview.TextView {
	headerText := `
🔋 LOCTHERAPY DASHBOARD 🔋
────────────────────────────────────────`

	return tview.NewTextView().
		SetText(headerText).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorGreen).
		SetDynamicColors(true).
		SetRegions(true)
}

func createButtons(app *tview.Application) *tview.Flex {
	button1 := tview.NewButton("1 - All").SetSelectedFunc(func() {
		// TODO: Add logic to filter and display all todos
	})
	button2 := tview.NewButton("2 - CF").SetSelectedFunc(func() {
		// TODO: Add logic to filter and display CF todos
	})
	button3 := tview.NewButton("3 - GP").SetSelectedFunc(func() {
		// TODO: Add logic to filter and display GP todos
	})

	buttonsFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(button1, 0, 1, false).
		AddItem(button2, 0, 1, false).
		AddItem(button3, 0, 1, false)

	return buttonsFlex
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

func createLayout(header *tview.TextView, buttonFlex *tview.Flex, todoTextView *tview.TextView) *tview.Flex {
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(buttonFlex, 3, 1, false). // Adjusted size and weight for buttons
		AddItem(todoTextView, 0, 1, true)

	return mainFlex
}

func setupKeybindings(app *tview.Application) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			app.Stop()
		}
		return event
	})
}

func (f *UIFactory) startUpdatingTodos(app *tview.Application, todoTextView *tview.TextView) {
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

func (f *UIFactory) Run() {
	f.startUpdatingTodos(f.app, f.todoTextView)
	runUI(f.app, f.mainContainer)
}
