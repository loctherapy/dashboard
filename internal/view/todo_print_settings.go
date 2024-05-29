package view

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type IPrintSettings interface {
	PrintContext(string) string
	PrintFile(string, int) string
	PrintToDo(string) string
}

type ConsolePrintSettings struct{}

func (c ConsolePrintSettings) PrintContext(context string) string {
	contextColor := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	return contextColor(fmt.Sprintf("🌳 %s\n", strings.ToUpper(context)))
}

func (c ConsolePrintSettings) PrintFile(filePath string, gravity int) string {
	fileColor := color.New(color.FgHiCyan).SprintFunc()
	return fileColor(fmt.Sprintf("  💎 %s (Gravity: %d)", filePath, gravity))
}

func (c ConsolePrintSettings) PrintToDo(todo string) string {
	todoColor := color.New(color.FgWhite).SprintFunc()
	return todoColor(fmt.Sprintf("    %s", todo))
}

type TViewPrintSettings struct{}

func (t TViewPrintSettings) PrintContext(context string) string {
	return fmt.Sprintf("[yellow::b]🌳 %s[-:-:-]\n", strings.ToUpper(context))
}

func (t TViewPrintSettings) PrintFile(filePath string, gravity int) string {
	return fmt.Sprintf("[cyan::b]  💎 %s (Gravity: %d)[-:-:-]", filePath, gravity)
}

func (t TViewPrintSettings) PrintToDo(todo string) string {
	return fmt.Sprintf("[white]    %s[-:-:-]", todo)
}
