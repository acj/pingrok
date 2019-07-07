package main

import (
	"log"
	"sync"
)

type CircularBuffer struct {
	mux sync.Mutex
	buf []LatencyReport
	count int
	currentOffset int
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		buf: make([]LatencyReport, size, size),
	}
}

func (b *CircularBuffer) Snapshot() []LatencyReport {
	snap := make([]LatencyReport, len(b.buf))

	b.mux.Lock()
	n := copy(snap[0:b.count-b.currentOffset], b.buf[b.count-(b.count-b.currentOffset):])
	if n != (b.count - b.currentOffset) {
		log.Fatalf("unexpected short copy: %d bytes", n)
	}
	n = copy(snap[n:], b.buf[0:b.currentOffset])
	if n != b.currentOffset {
		log.Fatalf("unexpected short copy: %d bytes", n)
	}
	b.mux.Unlock()
	return snap
}

func (b *CircularBuffer) Insert(value LatencyReport) {
	b.mux.Lock()
	b.buf[b.currentOffset] = value

	if b.count < len(b.buf) {
		b.count++
	}
	b.currentOffset = (b.currentOffset + 1) % len(b.buf)
	b.mux.Unlock()
}

func (b *CircularBuffer) Size() int {
	b.mux.Lock()
	defer b.mux.Unlock()
	return len(b.buf)
}
