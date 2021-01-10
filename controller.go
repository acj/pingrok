package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"math"
	"time"
)

type controller struct {
	config   *config
	uiBundle *uIBundle
}

func newController(config *config, uiBundle *uIBundle) *controller {
	return &controller{
		config:   config,
		uiBundle: uiBundle,
	}
}

func (c *controller) run() error {
	dataPointBuffer := newCircularBuffer(c.config.timeWindowSeconds * c.config.samplesPerSecond)
	partitioner := newDataPointPartitioner(dataPointBuffer, c.config.timeWindowSeconds, c.config.samplesPerSecond)
	dataPoints := make(chan latencyDataPoint)

	partitioner.start(dataPoints)
	pinger := newPinger(c.config.targetHost, dataPoints)
	pinger.start()

	go c.updateUILoop(1*time.Second, dataPointBuffer)

	return c.uiBundle.app.Run()
}

func (c *controller) updateUILoop(interval time.Duration, dataPointBuffer *circularBuffer) {
	c.uiBundle.heatmap.SetSelectionChangedFunc(func(row, col int) {
		if row > c.config.samplesPerSecond || col > c.config.timeWindowSeconds {
			c.uiBundle.infoCenterLeftCell.SetText(placeholderSelectCellText)
			return
		}

		dataPoint := dataPointBuffer.snapshot()[row+col*c.config.samplesPerSecond]
		c.uiBundle.infoCenterLeftCell.SetText(fmt.Sprintf("latency: %.02f ms @ Time Offset: %.02f seconds", dataPoint.latency, dataPoint.timeOffset))
	})

	for {
		applySnapshotToUI(dataPointBuffer.snapshot(), c.uiBundle, c.config.samplesPerSecond, c.config.timeWindowSeconds, c.config.overlayLatenciesOnHeatmap)
		time.Sleep(interval)
	}
}

func applySnapshotToUI(currentSnapshot []latencyDataPoint, uiBundle *uIBundle, samplesPerSecond int, timeWindowSeconds int, overlayLatenciesOnHeatmap bool) {
	minLatency := math.MaxFloat64
	maxLatency := 0.0
	for _, dataPoint := range currentSnapshot {
		if dataPoint.timeOffset == 0.0 {
			break
		}

		if dataPoint.latency >= maxLatency {
			maxLatency = dataPoint.latency
		}
		if dataPoint.latency <= minLatency {
			minLatency = dataPoint.latency
		}
	}

	uiBundle.app.QueueUpdateDraw(func() {
		haveFullSnapshot := false
		selectedRow, selectedCol := uiBundle.heatmap.GetSelection()

		for idx, dataPoint := range currentSnapshot {
			row := idx % samplesPerSecond
			col := idx / samplesPerSecond

			if dataPoint.timeOffset == 0.0 {
				continue
			}

			latencyRange := maxLatency - minLatency
			scaledRedLevel := int32(((dataPoint.latency - minLatency) / latencyRange) * 255.0)
			color := tcell.NewRGBColor(scaledRedLevel, 0, 0)
			currentCell := uiBundle.heatmap.GetCell(row, col)
			currentCell.SetBackgroundColor(color)

			if overlayLatenciesOnHeatmap {
				currentCell.SetText(fmt.Sprintf("%.1f", dataPoint.latency))
			}

			haveFullSnapshot = row == samplesPerSecond-1 && col == timeWindowSeconds-1
		}

		// Track the currently selected cell, if any
		if haveFullSnapshot && uiBundle.heatmap.HasFocus() {
			if selectedCol > 0 {
				uiBundle.heatmap.Select(selectedRow, selectedCol-1)
			} else {
				uiBundle.heatmap.Select(selectedRow, selectedCol)
			}
		}

		if minLatency != math.MaxFloat64 {
			uiBundle.infoCenterRightCell.SetText(fmt.Sprintf("Min: %.02f ms / Max: %.02f ms", minLatency, maxLatency))
		}
	})
}
