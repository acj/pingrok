package main

import (
	"testing"
)

func TestNewCircularBuffer_returnsUsableBuffer(t *testing.T) {
	b := newCircularBuffer(10)

	if b == nil {
		t.Error("expected a buffer but got nil")
	}
}

func TestSnapshot_returnsCorrectlySizedBuffer(t *testing.T) {
	expected := 10
	b := newCircularBuffer(expected)

	snap := b.snapshot()
	if len(snap) != expected {
		t.Errorf("expected buffer of size %d, but it was %d", expected, len(snap))
	}
}

func TestSnapshot_returnsCorrectSnapshotForJustSaturatedBuffer(t *testing.T) {
	b := newCircularBuffer(4)

	b.insert(latencyDataPoint{latency: 1.0})
	b.insert(latencyDataPoint{latency: 2.0})
	b.insert(latencyDataPoint{latency: 3.0})
	b.insert(latencyDataPoint{latency: 4.0})

	snap := b.snapshot()

	expected := latencyDataPoint{latency: 1.0}
	if snap[0] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[0])
	}
	expected = latencyDataPoint{latency: 2.0}
	if snap[1] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[1])
	}
	expected = latencyDataPoint{latency: 3.0}
	if snap[2] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[2])
	}
	expected = latencyDataPoint{latency: 4.0}
	if snap[3] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[3])
	}
}

func TestSnapshot_returnsCorrectSnapshotForOverSaturatedBuffer(t *testing.T) {
	b := newCircularBuffer(4)

	b.insert(latencyDataPoint{latency: 1.0})
	b.insert(latencyDataPoint{latency: 2.0})
	b.insert(latencyDataPoint{latency: 3.0})
	b.insert(latencyDataPoint{latency: 4.0})
	b.insert(latencyDataPoint{latency: 5.0})

	snap := b.snapshot()

	expected := latencyDataPoint{latency: 2.0}
	if snap[0] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[0])
	}
	expected = latencyDataPoint{latency: 3.0}
	if snap[1] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[1])
	}
	expected = latencyDataPoint{latency: 4.0}
	if snap[2] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[2])
	}
	expected = latencyDataPoint{latency: 5.0}
	if snap[3] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[3])
	}
}

func TestSnapshot_returnsCorrectSnapshotForNonSaturatedBuffer(t *testing.T) {
	b := newCircularBuffer(4)

	b.insert(latencyDataPoint{latency: 1.0})
	b.insert(latencyDataPoint{latency: 2.0})
	b.insert(latencyDataPoint{latency: 3.0})

	snap := b.snapshot()

	expected := latencyDataPoint{latency: 1.0}
	if snap[0] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[0])
	}
	expected = latencyDataPoint{latency: 2.0}
	if snap[1] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[1])
	}
	expected = latencyDataPoint{latency: 3.0}
	if snap[2] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[2])
	}
}

func TestInsert_addsValueToBuffer(t *testing.T) {
	b := newCircularBuffer(10)
	expected := latencyDataPoint{latency: 1.0}

	b.insert(expected)

	if snap := b.snapshot(); snap[0] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[0])
	}
}

func TestInsert_overwritesOldestValueWhenBufferIsFull(t *testing.T) {
	b := newCircularBuffer(1)

	b.insert(latencyDataPoint{latency: 1.0})
	b.insert(latencyDataPoint{latency: 2.0})

	snap := b.snapshot()
	expected := latencyDataPoint{latency: 2.0}
	if snap[0] != expected {
		t.Errorf("expected '%v' but got '%v'", expected, snap[0])
	}
}
