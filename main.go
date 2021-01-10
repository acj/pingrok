package main

import (
	"flag"
	"fmt"
	"os"
)

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
	config := &config{
		timeWindowSeconds:         *timeWindowSeconds,
		samplesPerSecond:          *samplesPerSecond,
		targetHost:                *targetHost,
		overlayLatenciesOnHeatmap: *overlayLatenciesOnHeatmap,
		uiBundle:                  uiBundle,
	}

	controller := newController(config, uiBundle)
	if err := controller.run(); err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
