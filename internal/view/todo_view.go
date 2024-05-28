package view

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/loctherapy/dashboard/internal/model"
	"github.com/rivo/tview"
)

type View struct {
	Printer       *ToDoPrinter
	app           *tview.Application
	header        *tview.TextView
	buttonsFlex    *tview.Flex
	todoTextView  *tview.TextView
	mainContainer *tview.Flex
}

func NewView(printer *ToDoPrinter) *View {
	app := tview.NewApplication()
	header := createMainHeader()
	buttonFlex := createButtons(app)
	todoTextView := createTodoTextView(app)
	mainContainer := createLayout(header, buttonFlex, todoTextView)

	setupKeybindings(app)

	return &View{printer, app, header, buttonFlex, todoTextView, mainContainer}
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

func (f *View) RunUI() {
	// Start the application
	if err := f.app.SetRoot(f.mainContainer, true).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func (f *View) DisplayToDos(todos []model.FileToDos) {
	// Use the printer to convert the todos to a string
	todosString, err := f.Printer.Print(todos)
	if err != nil {
		fmt.Println("Error printing todos:", err)
		return
	}

	// Queue the update to the todoTextView
	f.app.QueueUpdateDraw(func() {
		f.todoTextView.SetText(todosString)
	})
}
