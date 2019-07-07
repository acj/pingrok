package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO:
// - Automatically update the heatmap

// Maybe:
// * If we miss a reply, mark it somehow. Black frame? X? Use the cutoff value?

func main() {
	var timeWindow = flag.Int("t", 30, "seconds of data to display")
	var samplesPerSecond = flag.Int("r", 10, "number of pings per second")
	var bindHost = flag.String("s", "0.0.0.0:8086", "IP and port for web server")
	var targetHost = flag.String("h", "192.168.1.1", "the host to ping")
	flag.Parse()

	server := NewServer(*timeWindow, *samplesPerSecond)
	go server.Serve(*bindHost, *targetHost)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	server.Shutdown(ctx)

	os.Exit(0)
}
