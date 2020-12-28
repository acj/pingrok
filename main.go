package main

import (
	"flag"
	"fmt"
	"os"
)

// TODO:
// * If we miss a reply, mark it somehow. Black frame? X? Use the cutoff value?
// * Create an interface that decouples the display engine from ping
// * Create a generator that oscillates smoothly between 0 and some known latency -- good for end-to-end testing
// * Error out if the sample rate is greater than the height of the terminal
// * Add mechanism to stop Server

type config struct {
	timeWindowSeconds         int
	samplesPerSecond          int
	targetHost                string
	overlayLatenciesOnHeatmap bool
	uiBundle                  *uIBundle
}

func main() {
	var timeWindowSeconds = flag.Int("t", 30, "seconds of data to display")
	var samplesPerSecond = flag.Int("r", 10, "number of pings per second")
	var targetHost = flag.String("h", "192.168.1.1", "the host to ping")
	var overlayLatenciesOnHeatmap = flag.Bool("o", false, "Overlay latencies on heatmap")
	flag.Parse()

	uiBundle := prepareUI(*samplesPerSecond, *timeWindowSeconds)
	partitioner := newDataPointPartitioner(*timeWindowSeconds, *samplesPerSecond)
	config := &config{
		timeWindowSeconds:         *timeWindowSeconds,
		samplesPerSecond:          *samplesPerSecond,
		targetHost:                *targetHost,
		overlayLatenciesOnHeatmap: *overlayLatenciesOnHeatmap,
		uiBundle:                  uiBundle,
	}

	controller := newController(config, uiBundle, partitioner)
	if err := controller.Run(); err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
