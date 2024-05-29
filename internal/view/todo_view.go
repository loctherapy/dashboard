package view

import (
	"fmt"
	"os"
	"sort"
	"strings"

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
	buttonFlex := tview.NewFlex()
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

func getContexts(todos []model.FileToDos) []string {
	contexts := make(map[string]struct{})
	for _, todo := range todos {
		contexts[todo.Context] = struct{}{}
	}

	var contextList []string

	// Add a button for all contexts
	contextList = append(contextList, "All")

	for context := range contexts {
		contextList = append(contextList, context)
	}

	sort.Strings(contextList)

	return contextList
}

func (f *View) createButtons(contexts []string) {
	// Clear the buttonsFlex before adding new buttons
	f.buttonsFlex.Clear()
	
	// Create a button for each context
	for id, context := range contexts {
		newId := id + 1
		buttonName := fmt.Sprintf("%d - %s", newId, strings.ToUpper(context))
		button := tview.NewButton(buttonName).SetSelectedFunc(func() {})
		f.buttonsFlex.AddItem(button, 0, 1, false)
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

		// Get the unique contexts from the todos
		contexts := getContexts(todos)
		
		// Create buttons for each context
		f.createButtons(contexts)
	})
}
