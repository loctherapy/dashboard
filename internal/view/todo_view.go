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
	selectedContextID int
}

func NewView(printer *ToDoPrinter) *View {
	app := tview.NewApplication()
	header := createMainHeader()
	buttonFlex := tview.NewFlex()
	todoTextView := createTodoTextView(app)
	mainContainer := createLayout(header, buttonFlex, todoTextView)
	selectedContextID := 0

	view := &View{
		Printer:       printer,
		app:           app,
		header:        header,
		buttonsFlex:   buttonFlex,
		todoTextView:  todoTextView,
		mainContainer: mainContainer,
		buttons:       []*tview.Button{},
		selectedContextID: selectedContextID,
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
			index := int(event.Rune() - '1')
			if index < len(f.buttons) {
				f.selectedContextID = index
				f.setContextButtonFocus()
			}
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

func (f *View) setContextButtonFocus() {
	f.app.SetFocus(f.buttons[f.selectedContextID])
}

func (f *View) createButtons(contexts []string) {
	// Clear the buttonsFlex and button references before adding new buttons
	f.buttonsFlex.Clear()
	f.buttons = []*tview.Button{}

	// Create a button for each context
	for id, context := range contexts {
		newId := id + 1
		buttonName := fmt.Sprintf("%d - %s", newId, strings.ToUpper(context))
		button := tview.NewButton(buttonName).SetSelectedFunc(func() {
			// message := fmt.Sprintf("Button %d - %s pressed\n", newId, strings.ToUpper(context))
			// tview.Print(f.buttonFlex, message, 0, 0, 0, 0, tcell.ColorDefault)
			// tview.Print("Button %d - %s pressed\n", newId, strings.ToUpper(context))
			// fmt.Printf("Button %d - %s pressed\n", newId, strings.ToUpper(context))
		})
		f.buttonsFlex.AddItem(button, 0, 1, false)
		f.buttons = append(f.buttons, button)
	}

	f.setContextButtonFocus()

	// Setup key bindings for the new buttons
	f.setupDynamicKeybindings()
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
