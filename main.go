package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO:
// * If we miss a reply, mark it somehow. Black frame? X? Use the cutoff value?
// * Create an interface that decouples the display engine from ping
// * Create a generator that oscillates smoothly between 0 and some known latency -- good for end-to-end testing

func main() {
	var timeWindow = flag.Int("t", 30, "seconds of data to display")
	var samplesPerSecond = flag.Int("r", 10, "number of pings per second")
	var bindHost = flag.String("s", "0.0.0.0:8086", "IP and port for web server")
	var targetHost = flag.String("h", "192.168.1.1", "the host to ping")
	flag.Parse()

	server := NewServer(*timeWindow, *samplesPerSecond)
	go server.Serve(*bindHost, *targetHost)

	fmt.Printf("Up and running. Browse to http://%s\n", *bindHost)
	fmt.Println("Control-C to quit")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	server.Shutdown(ctx)

	os.Exit(0)
}
