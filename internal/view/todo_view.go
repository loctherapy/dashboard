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
	buttonsFlex   *tview.Flex
	todoTextView  *tview.TextView
	mainContainer *tview.Flex
	buttons       []*tview.Button // Store references to buttons
	todos		  []model.FileToDos
	contexts 	  []string
	filterByContext bool
	selectedContextName string
	selectedContextID int
}

func NewView(printer *ToDoPrinter) *View {
	app := tview.NewApplication()
	header := createMainHeader()
	buttonFlex := tview.NewFlex()
	todoTextView := createTodoTextView(app)
	mainContainer := createLayout(header, buttonFlex, todoTextView)
	selectedContextID := 0
	selectedContextName := "All"
	todos := []model.FileToDos{}
	filterByContext := false

	view := &View{
		Printer:       printer,
		app:           app,
		header:        header,
		buttonsFlex:   buttonFlex,
		todoTextView:  todoTextView,
		mainContainer: mainContainer,
		buttons:       []*tview.Button{},
		selectedContextID: selectedContextID,
		todos: todos,
		filterByContext: filterByContext,
		selectedContextName: selectedContextName,
	}

	setupInitialKeybindings(view)

	return view
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

func setupInitialKeybindings(view *View) {
	view.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			view.app.Stop()
		}
		return event
	})
}

func (f *View) setupDynamicKeybindings() {
	f.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			f.app.Stop()
		}

		// Handle number keys for buttons
		if event.Rune() >= '0' && event.Rune() <= '9' {
			
			/* Buttons: 	1 - All, 2 - CF, 3 - GP
			   Runes:   	1        2       3
			   ButtonIndex: 0        1       2
			   ContextIndex:?        0       1

			*/
			buttonIndex := int(event.Rune() - '1')
			contextIndex := buttonIndex - 1

			if contextIndex == -1 {
				f.filterByContext = false
			} else if contextIndex < len(f.buttons) {
				f.filterByContext = true
				f.selectedContextID = contextIndex
				f.selectedContextName = f.contexts[contextIndex]
			}
			f.setContextButtonFocus()
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

	for context := range contexts {
		contextList = append(contextList, context)
	}

	sort.Strings(contextList)

	return contextList
}

func (f *View) setContextButtonFocus() {
	if !f.filterByContext {
		f.app.SetFocus(f.buttons[0])
	} else {
		f.app.SetFocus(f.buttons[f.selectedContextID + 1])
	}
	f.app.SetFocus(f.todoTextView)
}

func (f *View) createButtons() {
	// Clear the buttonsFlex and button references before adding new buttons
	f.buttonsFlex.Clear()
	f.buttons = []*tview.Button{}

	createButton := func(id int, context string, filterByContext bool) {
		buttonName := fmt.Sprintf("%d - %s", id, strings.ToUpper(context))
		button := tview.NewButton(buttonName).SetSelectedFunc(func() {
			
			f.filterByContext = filterByContext
			if filterByContext {
				f.selectedContextName = context
			}
			
		})
		f.buttonsFlex.AddItem(button, 0, 1, false)
		f.buttons = append(f.buttons, button)
	}

	createButton(1, "All", false)

	// Create a button for each context
	for id, context := range f.contexts {
		newId := id + 2
		createButton(newId, context, true)
	}

	f.setContextButtonFocus()

	// Setup key bindings for the new buttons
	f.setupDynamicKeybindings()
}

func (f* View) getFilteredByContextToDos() []model.FileToDos {
	if !f.filterByContext {
		return f.todos
	}

	filteredToDos := []model.FileToDos{}
	for _, todo := range f.todos {
		if todo.Context == f.selectedContextName {
			filteredToDos = append(filteredToDos, todo)
		}
	}

	return filteredToDos
}

func (f *View) DisplayToDos(todos []model.FileToDos) {
	// Store the todos
	f.todos = todos
	f.contexts = getContexts(todos)

	// Get the todos to display
	filteredToDos := f.getFilteredByContextToDos()

	// Use the printer to convert the todos to a string
	todosString, err := f.Printer.Print(filteredToDos)
	if err != nil {
		fmt.Println("Error printing todos:", err)
		return
	}

	// Queue the update to the todoTextView
	f.app.QueueUpdateDraw(func() {
		f.todoTextView.SetText(todosString)
		
		// Create buttons for each context
		f.createButtons()
	})
}
