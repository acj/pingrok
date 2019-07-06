package main

import (
	"context"
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
	// TODO:
	// - time window flag
	// - samples per second flag
	// - host/port flag

	bind := "0.0.0.0:8086"

	server := NewServer(30, 20)
	go server.Serve(bind)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	server.Shutdown(ctx)

	os.Exit(0)
}
