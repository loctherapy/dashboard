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
	buttons       []*tview.Button
	todoTextView  *tview.TextView
	mainContainer *tview.Flex
}


func NewUIFactory() *UIFactory {
	app := tview.NewApplication()
	header := createMainHeader()
	buttons := createButtons(app)
	todoTextView := createTodoTextView(app)
	mainContainer := createLayout(header, buttons, todoTextView)

	setupKeybindings(app)

	return &UIFactory{app, header, buttons, todoTextView, mainContainer}
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

func createButtons(app *tview.Application) []*tview.Button {
    button1 := tview.NewButton("1 - All").SetSelectedFunc(func() {
        // TODO: Add logic to filter and display all todos
    })
    button2 := tview.NewButton("2 - CF").SetSelectedFunc(func() {
        // TODO: Add logic to filter and display CF todos
    })
    button3 := tview.NewButton("3 - GP").SetSelectedFunc(func() {
        // TODO: Add logic to filter and display GP todos
    })

    buttons := make([]*tview.Button, 0)
	buttons = append(buttons, button1, button2, button3)

    return buttons
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

func createLayout(header *tview.TextView, buttons []*tview.Button, todoTextView *tview.TextView) *tview.Flex {
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false)

	for _, button := range buttons {
		flex.AddItem(button, 1, 1, false)
	
	}
		
	flex.AddItem(todoTextView, 0, 1, true)

	return flex
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

func (f *UIFactory) Run() {
	startUpdatingTodos(f.app, f.todoTextView)

	runUI(f.app, f.mainContainer)
}