package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TerminalFormatter struct {
	timeWindowSeconds int
	samplesPerSecond  int
	app               *tview.Application
	grid              *tview.Grid
	gridItems         [][]*tview.Box
}

func NewTerminalFormatter(timeWindowSeconds, samplesPerSecond int) *TerminalFormatter {
	grid := tview.NewGrid().
		SetSize(samplesPerSecond, timeWindowSeconds, 1, 1).
		SetBorders(true)

	gridItems := make([][]*tview.Box, timeWindowSeconds)
	for col := 0; col < timeWindowSeconds; col++ {
		gridItems[col] = make([]*tview.Box, samplesPerSecond)
	}

	for row := 0; row < samplesPerSecond; row++ {
		for col := 0; col < timeWindowSeconds; col++ {
			emptyCell := tview.NewBox().SetBackgroundColor(tcell.NewRGBColor(int32(row*col), 0, 0))

			grid.AddItem(emptyCell, row, col, 1, 1, 1, 1, false)

			gridItems[col][row] = emptyCell
		}
	}

	app := tview.NewApplication().SetRoot(grid, true).EnableMouse(true)

	return &TerminalFormatter{
		timeWindowSeconds: timeWindowSeconds,
		samplesPerSecond:  samplesPerSecond,
		app:               app,
		grid:              grid,
		gridItems:         gridItems,
	}
}

func (f TerminalFormatter) start() error {
	return f.app.Run()
}
