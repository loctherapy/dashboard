package view

import "fmt"

type PrintMode int

const (
	Console PrintMode = iota
	TView
)

type ToDoPrinterFactory struct{}

func (f ToDoPrinterFactory) CreatePrinter(mode PrintMode) (*ToDoPrinter, error) {
	switch mode {
	case Console:
		return NewToDoPrinter(ConsolePrintSettings{})
	case TView:
		return NewToDoPrinter(TViewPrintSettings{})
	default:
		return nil, fmt.Errorf("unknown print mode")
	}
}
