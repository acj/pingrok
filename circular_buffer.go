package main

import (
	"log"
)

type CircularBuffer struct {
	buf []LatencyReport
	currentOffset int
	saturated bool
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		buf: make([]LatencyReport, size, size),
	}
}

func (b *CircularBuffer) Snapshot() []LatencyReport {
	snap := make([]LatencyReport, len(b.buf))

	if b.saturated {
		n := copy(snap[0:len(b.buf)-b.currentOffset], b.buf[len(b.buf)-(len(b.buf)-b.currentOffset):])
		if n != (len(b.buf) - b.currentOffset) {
			log.Fatalf("unexpected short copy: %d bytes", n)
		}
		n = copy(snap[n:], b.buf[0:b.currentOffset])
		if n != b.currentOffset {
			log.Fatalf("unexpected short copy: %d bytes", n)
		}
	} else {
		n := copy(snap[:], b.buf[:])
		if n != len(b.buf) {
			log.Fatalf("unexpected short copy: %d bytes", n)
		}
	}
	return snap
}

func (b *CircularBuffer) Insert(value LatencyReport) {
	b.buf[b.currentOffset] = value
	b.currentOffset = (b.currentOffset + 1) % len(b.buf)

	if b.currentOffset == 0 {
		b.saturated = true
	}
}

func (b *CircularBuffer) Size() int {
	return len(b.buf)
}
