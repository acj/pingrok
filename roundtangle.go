package main

import (
	"log"
)

type Buffer struct {
	buf []LatencyReport
	currentOffset int
	saturated bool
}

func NewCircularBuffer(size int) *Buffer {
	return &Buffer{
		buf: make([]LatencyReport, size, size),
	}
}

func (b *Buffer) Snapshot() []LatencyReport {
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

func (b *Buffer) Insert(value LatencyReport) {
	b.buf[b.currentOffset] = value
	//log.Printf("offset is now (%d + 1) %% %d = %d", b.currentOffset, len(b.buf), (b.currentOffset + 1) % len(b.buf))
	b.currentOffset = (b.currentOffset + 1) % len(b.buf)

	if b.currentOffset == 0 {
		b.saturated = true
	}
}

func (b *Buffer) Size() int {
	return len(b.buf)
}
