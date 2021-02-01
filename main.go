package main

import (
	"flag"
	"fmt"
	"log"
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
	var overlayLatenciesOnHeatmap = flag.Bool("o", false, "Overlay latency numbers on heatmap")
	var logFilePath = flag.String("l", "pingrok.log", "Log file path")
	flag.Parse()

	logFile, err := os.OpenFile(*logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error opening log file: %s", err)
		os.Exit(1)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(logFile)

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
