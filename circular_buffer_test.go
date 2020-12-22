package main

import (
	"testing"
)

func TestNewCircularBuffer_returnsUsableBuffer(t *testing.T) {
	b := NewCircularBuffer(10)

	if b == nil {
		t.Error("expected a buffer but got nil")
	}
}

func TestSnapshot_returnsCorrectlySizedBuffer(t *testing.T) {
	expected := 10
	b := NewCircularBuffer(expected)

	snap := b.Snapshot()
	if len(snap) != expected {
		t.Errorf("expected buffer of size %d, but it was %d", expected, len(snap))
	}
}

func TestSnapshot_returnsCorrectSnapshotForJustSaturatedBuffer(t *testing.T) {
	b := NewCircularBuffer(4)

	b.Insert(LatencyDataPoint{Latency: 1.0})
	b.Insert(LatencyDataPoint{Latency: 2.0})
	b.Insert(LatencyDataPoint{Latency: 3.0})
	b.Insert(LatencyDataPoint{Latency: 4.0})

	snap := b.Snapshot()

	expected := LatencyDataPoint{Latency: 1.0}
	if snap[0] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[0])
	}
	expected = LatencyDataPoint{Latency: 2.0}
	if snap[1] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[1])
	}
	expected = LatencyDataPoint{Latency: 3.0}
	if snap[2] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[2])
	}
	expected = LatencyDataPoint{Latency: 4.0}
	if snap[3] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[3])
	}
}

func TestSnapshot_returnsCorrectSnapshotForOverSaturatedBuffer(t *testing.T) {
	b := NewCircularBuffer(4)

	b.Insert(LatencyDataPoint{Latency: 1.0})
	b.Insert(LatencyDataPoint{Latency: 2.0})
	b.Insert(LatencyDataPoint{Latency: 3.0})
	b.Insert(LatencyDataPoint{Latency: 4.0})
	b.Insert(LatencyDataPoint{Latency: 5.0})

	snap := b.Snapshot()

	expected := LatencyDataPoint{Latency: 2.0}
	if snap[0] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[0])
	}
	expected = LatencyDataPoint{Latency: 3.0}
	if snap[1] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[1])
	}
	expected = LatencyDataPoint{Latency: 4.0}
	if snap[2] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[2])
	}
	expected = LatencyDataPoint{Latency: 5.0}
	if snap[3] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[3])
	}
}

func TestSnapshot_returnsCorrectSnapshotForNonSaturatedBuffer(t *testing.T) {
	b := NewCircularBuffer(4)

	b.Insert(LatencyDataPoint{Latency: 1.0})
	b.Insert(LatencyDataPoint{Latency: 2.0})
	b.Insert(LatencyDataPoint{Latency: 3.0})

	snap := b.Snapshot()

	expected := LatencyDataPoint{Latency: 1.0}
	if snap[0] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[0])
	}
	expected = LatencyDataPoint{Latency: 2.0}
	if snap[1] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[1])
	}
	expected = LatencyDataPoint{Latency: 3.0}
	if snap[2] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[2])
	}
}

func TestInsert_addsValueToBuffer(t *testing.T) {
	b := NewCircularBuffer(10)
	expected := LatencyDataPoint{Latency: 1.0}

	b.Insert(expected)

	if snap := b.Snapshot(); snap[0] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[0])
	}
}

func TestInsert_overwritesOldestValueWhenBufferIsFull(t *testing.T) {
	b := NewCircularBuffer(1)

	b.Insert(LatencyDataPoint{Latency: 1.0})
	b.Insert(LatencyDataPoint{Latency: 2.0})

	snap := b.Snapshot()
	expected := LatencyDataPoint{Latency: 2.0}
	if snap[0] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[0])
	}
}
