package main

import (
	"flag"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"math"
	"os"
	"strconv"
	"time"
)

// TODO:
// * If we miss a reply, mark it somehow. Black frame? X? Use the cutoff value?
// * Create an interface that decouples the display engine from ping
// * Create a generator that oscillates smoothly between 0 and some known latency -- good for end-to-end testing

func main() {
	var timeWindowSeconds = flag.Int("t", 30, "seconds of data to display")
	var samplesPerSecond = flag.Int("r", 10, "number of pings per second")
	var targetHost = flag.String("h", "192.168.1.1", "the host to ping")
	flag.Parse()

	yAxisLabels := tview.NewGrid().
		SetSize(*samplesPerSecond, 1, 0, 0).
		SetBorders(false)

	xAxisLabels := tview.NewGrid().
		SetSize(1, *timeWindowSeconds, 0, 0).
		SetBorders(false)

	heatmap := tview.NewGrid().
		SetSize(*samplesPerSecond, *timeWindowSeconds, 0, 0).
		SetBorders(false)

	rootGrid := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(10, 0).
		SetBorders(true).
		AddItem(yAxisLabels, 0, 0, 1, 1, 0, 0, false).
		AddItem(xAxisLabels, 1, 1, 1, 1, 0, 0, false).
		AddItem(heatmap, 0, 1, 1, 1, 0, 0, false)

	app := tview.NewApplication().
		SetRoot(rootGrid, true).
		EnableMouse(true)

	gridItems := make([][]*tview.Box, *samplesPerSecond)
	for row := 0; row < *samplesPerSecond; row++ {
		gridItems[row] = make([]*tview.Box, *timeWindowSeconds)
	}

	for row := 0; row < *samplesPerSecond; row++ {
		offsetMs := int(1000.0 * float64(row) / float64(*samplesPerSecond))
		labelBox := tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			SetText(fmt.Sprintf("%d ms", offsetMs)).
			SetTextColor(tcell.ColorWhite)
		yAxisLabels.AddItem(labelBox, row, 0, 1, 1, 0, 0, false)
	}

	for col := 0; col < *timeWindowSeconds; col++ {
		labelBox := tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			SetText(strconv.Itoa(col)).
			SetTextColor(tcell.ColorWhite)
		xAxisLabels.AddItem(labelBox, 0, col, 1, 1, 0, 0, false)
	}

	for row := 0; row < *samplesPerSecond; row++ {
		for col := 0; col < *timeWindowSeconds; col++ {
			emptyCell := tview.NewBox().SetBackgroundColor(tcell.ColorGreen)
			heatmap.AddItem(emptyCell, row, col, 1, 1, 1, 1, false)
			gridItems[row][col] = emptyCell
		}
	}

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
			snapshot := pinger.latencyReportCircularBuffer.Snapshot()

			minLatency := math.MaxFloat64
			maxLatency := 0.0
			for _, dataPoint := range snapshot {
				if dataPoint.Latency >= maxLatency {
					maxLatency = dataPoint.Latency
				}
				if dataPoint.Latency <= minLatency {
					minLatency = dataPoint.Latency
				}
			}

			app.QueueUpdateDraw(func() {
				for idx, dataPoint := range snapshot {
					row := idx % *samplesPerSecond
					col := idx / *samplesPerSecond

					if dataPoint.TimeOffset == 0.0 {
						continue
					}

					latencyRange := maxLatency - minLatency
					scaledRedLevel := int32(((dataPoint.Latency - minLatency) / latencyRange) * 255.0)
					color := tcell.NewRGBColor(int32(latencyRange) - scaledRedLevel, 0, 0)
					gridItems[row][col].SetBackgroundColor(color)
				}
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
