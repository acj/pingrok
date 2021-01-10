package main

import (
	"log"
	"sync"
)

type circularBuffer struct {
	mux           sync.Mutex
	buf           []latencyDataPoint
	count         int
	currentOffset int
}

func newCircularBuffer(size int) *circularBuffer {
	return &circularBuffer{
		buf: make([]latencyDataPoint, size, size),
	}
}

func (b *circularBuffer) snapshot() []latencyDataPoint {
	snap := make([]latencyDataPoint, len(b.buf))

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

func (b *circularBuffer) insert(value latencyDataPoint) {
	b.mux.Lock()
	b.buf[b.currentOffset] = value

	if b.count < len(b.buf) {
		b.count++
	}
	b.currentOffset = (b.currentOffset + 1) % len(b.buf)
	b.mux.Unlock()
}

func (b *circularBuffer) size() int {
	b.mux.Lock()
	defer b.mux.Unlock()
	return len(b.buf)
}
