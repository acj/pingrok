package main

import (
	"flag"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"math"
	"os"
	"time"
)

// TODO:
// * If we miss a reply, mark it somehow. Black frame? X? Use the cutoff value?
// * Create an interface that decouples the display engine from ping
// * Create a generator that oscillates smoothly between 0 and some known latency -- good for end-to-end testing

const placeholderSelectCellText = "Select a cell to view latency and time"

func main() {
	var timeWindowSeconds = flag.Int("t", 30, "seconds of data to display")
	var samplesPerSecond = flag.Int("r", 10, "number of pings per second")
	var targetHost = flag.String("h", "192.168.1.1", "the host to ping")
	flag.Parse()

	var currentSnapshot []LatencyDataPoint

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

	heatmap.SetSelectionChangedFunc(func(row, col int) {
		if row > *samplesPerSecond || col > *timeWindowSeconds {
			infoCenterLeftCell.SetText(placeholderSelectCellText)
			return
		}

		dataPoint := currentSnapshot[row + col * *samplesPerSecond]
		infoCenterLeftCell.SetText(fmt.Sprintf("Latency: %.02f ms @ Time Offset: %.02f seconds", dataPoint.Latency, dataPoint.TimeOffset))
	})

	for row := 0; row < *samplesPerSecond; row++ {
		offsetMs := int(1000.0 * float64(row) / float64(*samplesPerSecond))
		cell := tview.NewTableCell(fmt.Sprintf("%d ms", offsetMs)).
			SetAlign(tview.AlignRight).
			SetTextColor(tcell.ColorWhite).
			SetExpansion(1)
		yAxisLabels.SetCell(row, 0, cell)
	}

	for col := 0; col < *timeWindowSeconds; col++ {
		cell := tview.NewTableCell(fmt.Sprintf("%02d", col)).
			SetAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorWhite).
			SetExpansion(1)
		xAxisLabels.SetCell(0, col, cell)
	}

	for row := 0; row < *samplesPerSecond; row++ {
		for col := 0; col < *timeWindowSeconds; col++ {
			cell := tview.NewTableCell("").
				SetExpansion(1).
				SetAlign(tview.AlignCenter)
			heatmap.SetCell(row, col, cell)
		}
	}

	// Placeholder cell for nothing-is-selected state
	heatmap.SetCell(*samplesPerSecond + 1, 0, tview.NewTableCell(""))
	heatmap.Select(*samplesPerSecond + 1, 0)

	// TODO: Handle mouse events
	// TODO: Handle keyboard events
	// TODO: Error out if the sample rate is greater than the height of the terminal
	// TODO: Rename Server. Sampler?
	// TODO: Add mechanism to stop Server
	// TODO: Follow the selected data point as time passes
	// TODO: Show min/max latency?

	pinger := NewServer(*timeWindowSeconds, *samplesPerSecond)
	pinger.Start(*targetHost)

	go func() {
		for {
			currentSnapshot = pinger.latencyReportCircularBuffer.Snapshot()

			minLatency := math.MaxFloat64
			maxLatency := 0.0
			for _, dataPoint := range currentSnapshot {
				if dataPoint.Latency >= maxLatency {
					maxLatency = dataPoint.Latency
				}
				if dataPoint.Latency <= minLatency {
					minLatency = dataPoint.Latency
				}
			}

			app.QueueUpdateDraw(func() {
				haveFullSnapshot := false
				selectedRow, selectedCol := heatmap.GetSelection()

				for idx, dataPoint := range currentSnapshot {
					row := idx % *samplesPerSecond
					col := idx / *samplesPerSecond

					if dataPoint.TimeOffset == 0.0 {
						continue
					}

					latencyRange := maxLatency - minLatency
					scaledRedLevel := int32(((dataPoint.Latency - minLatency) / latencyRange) * 255.0)
					color := tcell.NewRGBColor(scaledRedLevel, 0, 0)
					heatmap.GetCell(row, col).SetBackgroundColor(color)

					haveFullSnapshot = row == *samplesPerSecond - 1 && col == *timeWindowSeconds - 1
				}

				// Track the currently selected cell, if any
				if haveFullSnapshot && heatmap.HasFocus() {
					if selectedCol > 0 {
						heatmap.Select(selectedRow, selectedCol - 1)
					} else {
						heatmap.Select(selectedRow, selectedCol)
					}
				}

				infoCenterRightCell.SetText(fmt.Sprintf("Min: %.02f ms / Max: %.02f ms", minLatency, maxLatency))
			})

			time.Sleep(1 * time.Second)
		}
	}()

	if err := app.Run(); err != nil {
		fmt.Printf("couldn't start app: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
