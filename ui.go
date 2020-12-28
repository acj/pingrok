package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UIBundle struct {
	app *tview.Application
	heatmap *tview.Table
	infoCenterLeftCell *tview.TableCell
	infoCenterRightCell *tview.TableCell
}

func prepareUI(samplesPerSecond, timeWindowSeconds int) *UIBundle {
	yAxisLabels := tview.NewTable()
	xAxisLabels := tview.NewTable()

	heatmap := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, true)

	infoCenter := tview.NewTable()
	infoCenterLeftCell := tview.NewTableCell(placeholderSelectCellText).SetExpansion(1)
	infoCenterRightCell := tview.NewTableCell("")
	infoCenter.SetCell(0, 0, infoCenterLeftCell).
		SetCell(0, 1, infoCenterRightCell)

	rootGrid := tview.NewGrid().
		SetRows(0, 1, 1).
		SetColumns(8, 0).
		SetBorders(true).
		AddItem(yAxisLabels, 0, 0, 1, 1, 0, 0, false).
		AddItem(xAxisLabels, 1, 1, 1, 1, 0, 0, false).
		AddItem(heatmap, 0, 1, 1, 1, 0, 0, true).
		AddItem(infoCenter, 2, 1, 1, 1, 0, 0, false)

	app := tview.NewApplication().
		SetRoot(rootGrid, true).
		EnableMouse(true)

	for row := 0; row < samplesPerSecond; row++ {
		offsetMs := int(1000.0 * float64(row) / float64(samplesPerSecond))
		cell := tview.NewTableCell(fmt.Sprintf("%d ms", offsetMs)).
			SetAlign(tview.AlignRight).
			SetTextColor(tcell.ColorWhite).
			SetExpansion(1)
		yAxisLabels.SetCell(row, 0, cell)
	}

	for col := 0; col < timeWindowSeconds; col++ {
		cell := tview.NewTableCell(fmt.Sprintf("%02d", col)).
			SetAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorWhite).
			SetExpansion(1)
		xAxisLabels.SetCell(0, col, cell)
	}

	for row := 0; row < samplesPerSecond; row++ {
		for col := 0; col < timeWindowSeconds; col++ {
			cell := tview.NewTableCell("").
				SetExpansion(1).
				SetAlign(tview.AlignCenter)
			heatmap.SetCell(row, col, cell)
		}
	}

	// Placeholder cell for nothing-is-selected state
	heatmap.SetCell(samplesPerSecond + 1, 0, tview.NewTableCell(""))
	heatmap.Select(samplesPerSecond + 1, 0)

	return &UIBundle{
		app: app,
		heatmap: heatmap,
		infoCenterLeftCell: infoCenterLeftCell,
		infoCenterRightCell: infoCenterRightCell,
	}
}