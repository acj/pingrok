package main

import (
	"fmt"
	"log"
	"os"
)

// Example usage:
// $ sudo ping -i 0.1 -W 100 192.168.1.1 | ts -s '[%.s]' | sed -l -n -e 's/\[\(.*\)\].*time=\(.*\) ms/\1 \2/p' | go run .

// TODO:
// - Expose as a self-contained webserver. Vendor the js dependencies, etc.
// - Automatically update the heatmap

// Maybe:
// * If we miss a reply, mark it somehow. Black frame? X? Use the cutoff value?

type Reply struct {
	TimeOffset float64
	Latency    float64
}

func (r *Reply) string() string {
	return fmt.Sprintf("offset: %f; latency: %f", r.TimeOffset, r.Latency)
}

func discretizeReplies(samplesPerSecond int, in <-chan Reply, out chan<- []Reply) {
	// Assumption: inbound replies are ordered by time
	currentAccumulatorSecondOffset := 0
	timeQuantum := 1.0 / float64(samplesPerSecond)
	currentSlice := make([]Reply, samplesPerSecond, samplesPerSecond)

	for r := range in {
		currentSecond := int(r.TimeOffset)
		if currentAccumulatorSecondOffset != currentSecond {
			out<- currentSlice
			currentSlice = make([]Reply, samplesPerSecond, samplesPerSecond)
			currentAccumulatorSecondOffset = currentSecond
		}

		currentSubsecondOffset := r.TimeOffset - float64(int(r.TimeOffset))
		currentSlice[int(currentSubsecondOffset / timeQuantum)] = r
	}
}

func main() {
	// TODO: Accept a duration flag

	timeWindow := 30
	samplesPerSecond := 10

	jsonFilename := "data.json"
	jsonFile, err := os.Create(jsonFilename)
	if err != nil {
		log.Fatalf("couldn't open file '%s': %v", jsonFilename, err)
	}

	replies := make(chan Reply, 2000)
	pinger := NewPinger(replies)
	discretizedReplies := make(chan []Reply)
	go discretizeReplies(samplesPerSecond, replies, discretizedReplies)
	pinger.Start()

	formatter := Formatter{
		w: jsonFile,
		timeWindow: timeWindow,
		samplesPerSecond: samplesPerSecond,
	}
	timeSeriesData := NewCircularBuffer(timeWindow*samplesPerSecond)

	for oneSecondOfData := range discretizedReplies {
		if _, err := jsonFile.Seek(0, 0); err != nil {
			log.Fatalf("seek error: %v", err)
		}

		for _, r := range oneSecondOfData {
			timeSeriesData.Insert(r)
		}

		formatter.formatDataAsJSON(timeSeriesData.Snapshot())
	}

	if err := os.Stdin.Close(); err != nil {
		log.Fatalf("stdin: %v", err)
	}
}
