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
	Printer            *ToDoPrinter
	app                *tview.Application
	buttonsFlex        *tview.Flex
	todoTextView       *tview.TextView
	mainContainer      *tview.Flex
	frame              *tview.Frame
	buttons            []*tview.Button // Store references to buttons
	todos              []model.FileToDos
	contexts           []string
	filterByContext    bool
	selectedContextName string
	selectedContextID  int
}

func NewView(printer *ToDoPrinter) *View {
	app := tview.NewApplication()
	buttonFlex := tview.NewFlex()
	todoTextView := createTodoTextView(app)
	mainContainer := createLayout(buttonFlex, todoTextView)
	selectedContextID := 0
	selectedContextName := "All"
	todos := []model.FileToDos{}
	filterByContext := false
	frame := createNewFrame(mainContainer)

	view := &View{
		Printer:            printer,
		app:                app,
		buttonsFlex:        buttonFlex,
		todoTextView:       todoTextView,
		mainContainer:      mainContainer,
		frame: 				frame,
		buttons:            []*tview.Button{},
		selectedContextID:  selectedContextID,
		todos:              todos,
		filterByContext:    filterByContext,
		selectedContextName: selectedContextName,
	}

	setupInitialKeybindings(view)

	return view
}

func createNewFrame(flex *tview.Flex) *tview.Frame {
	return tview.NewFrame(flex).
		SetBorders(1, 1, 1, 1, 2, 2)//.
		// AddText("🔋 LOCTHERAPY DASHBOARD 🔋", true, tview.AlignCenter, tcell.ColorGreen).
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

func createLayout(buttonFlex *tview.Flex, todoTextView *tview.TextView) *tview.Flex {
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
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
			if contextIndex == -2 {
				return event
			} else if contextIndex == -1 {
				f.filterByContext = false
			} else if contextIndex < len(f.buttons)-1 {
				f.filterByContext = true
				f.selectedContextID = contextIndex
				f.selectedContextName = f.contexts[contextIndex]
			}

			f.resetButtonColors()
			f.redrawToDos()
		}

		return event
	})
}

func (f *View) RunUI() {
	// Start the application
	if err := f.app.SetRoot(f.frame, true).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func getContexts(todos []model.FileToDos) []string {
	type Context struct {
		Name    string
		Gravity int
	}

	contextMap := make(map[string]Context)
	for _, todo := range todos {
		if existingContext, exists := contextMap[todo.Context]; !exists || existingContext.Gravity > todo.ContextGravity {
			contextMap[todo.Context] = Context{Name: todo.Context, Gravity: todo.ContextGravity}
		}
	}

	var contextList []Context
	for _, context := range contextMap {
		contextList = append(contextList, context)
	}

	// Sort contexts by gravity
	sort.Slice(contextList, func(i, j int) bool {
		return contextList[i].Gravity < contextList[j].Gravity
	})

	var sortedContexts []string
	for _, context := range contextList {
		sortedContexts = append(sortedContexts, context.Name)
	}

	return sortedContexts
}

func (f *View) createButtons() {
	// Clear the buttonsFlex and button references before adding new buttons
	f.buttonsFlex.Clear()
	f.buttons = []*tview.Button{}

	createButton := func(id int, context string) {
		buttonName := fmt.Sprintf("%d - %s", id, strings.ToUpper(context))
		button := tview.NewButton(buttonName)
		f.buttonsFlex.AddItem(button, 0, 1, false)
		f.buttons = append(f.buttons, button)
	}

	createButton(1, "All")

	// Create a button for each context
	for id, context := range f.contexts {
		newId := id + 2
		createButton(newId, context)
	}

	// Setup key bindings for the new buttons
	f.setupDynamicKeybindings()
}

func (f *View) getFilteredByContextToDos() []model.FileToDos {
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

func (f *View) resetButtonColors() {

	isSelected := func(index int) bool {
		if index == 0 && !f.filterByContext {
			return true
		}
		if index > 0 && f.filterByContext {
			return f.selectedContextID == index-1
		}

		return false
	}

	for index, button := range f.buttons {
		if isSelected(index) {
			button.SetBorder(true)
			button.SetBorderColor(tcell.ColorAntiqueWhite)
		} else {
			button.SetBorder(false)
		}
	}
}

func (f *View) redrawToDos() {
	// Get the todos to display
	filteredToDos := f.getFilteredByContextToDos()

	// Use the printer to convert the todos to a string
	todosString, err := f.Printer.Print(filteredToDos)
	if err != nil {
		fmt.Println("Error printing todos:", err)
		return
	}

	f.todoTextView.SetText(todosString)
}

func (f *View) DisplayToDos(todos []model.FileToDos) {
	// Store the todos
	f.todos = todos
	f.contexts = getContexts(todos)

	// Create buttons for each context
	f.createButtons()

	// Reset button colors
	f.resetButtonColors()

	// Redraw the todos
	f.redrawToDos()
}
