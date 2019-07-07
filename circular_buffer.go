package main

import (
	"log"
)

type CircularBuffer struct {
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

	n := copy(snap[0:b.count-b.currentOffset], b.buf[b.count-(b.count-b.currentOffset):])
	if n != (b.count - b.currentOffset) {
		log.Fatalf("unexpected short copy: %d bytes", n)
	}
	n = copy(snap[n:], b.buf[0:b.currentOffset])
	if n != b.currentOffset {
		log.Fatalf("unexpected short copy: %d bytes", n)
	}
	return snap
}

func (b *CircularBuffer) Insert(value LatencyReport) {
	b.buf[b.currentOffset] = value

	if b.count < len(b.buf) {
		b.count++
	}
	b.currentOffset = (b.currentOffset + 1) % len(b.buf)
}

func (b *CircularBuffer) Size() int {
	return len(b.buf)
}
