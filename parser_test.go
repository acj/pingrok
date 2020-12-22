package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestParser_handlesEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	p := NewParser(r)

	for r := range p.Next {
		t.Errorf("emitted unexpected value: %v", r)
	}
}

func TestParser_handlesSingleLineWithoutLineBreak(t *testing.T) {
	expectedTime := 0.01234
	expectedLatency := 12345.0
	r := strings.NewReader(fmt.Sprintf("%f %f", expectedTime, expectedLatency))
	p := NewParser(r)

	haveReceivedReply := false

	for report := range p.Next {
		if haveReceivedReply {
			t.Errorf("had already received a report when I saw this: %v", report)
		}
		haveReceivedReply = true

		if report.TimeOffset != expectedTime {
			t.Errorf("incorrect time: wanted %f, got %f", expectedTime, report.TimeOffset)
		}
		if report.Latency != expectedLatency {
			t.Errorf("incorrect latency: wanted %f, got %f", expectedLatency, report.Latency)
		}
	}
}

func TestParser_handlesSingleLineWithLineBreak(t *testing.T) {
	expectedTime := 0.01234
	expectedLatency := 12345.0
	r := strings.NewReader(fmt.Sprintf("%f %f\n", expectedTime, expectedLatency))
	p := NewParser(r)

	haveReceivedReply := false

	for report := range p.Next {
		if haveReceivedReply {
			t.Errorf("had already received a report when I saw this: %v", report)
		}
		haveReceivedReply = true

		if report.TimeOffset != expectedTime {
			t.Errorf("incorrect time: wanted %f, got %f", expectedTime, report.TimeOffset)
		}
		if report.Latency != expectedLatency {
			t.Errorf("incorrect latency: wanted %f, got %f", expectedLatency, report.Latency)
		}
	}
}

func TestParser_handlesMultipleLines(t *testing.T) {
	input := "0.050760 7.373\n0.058287 4.277\n0.071051 5.429\n0.078683 3.043"
	reader := strings.NewReader(input)

	outputs := []LatencyDataPoint{
		{TimeOffset: 0.050760, Latency: 7.373},
		{TimeOffset: 0.058287, Latency: 4.277},
		{TimeOffset: 0.071051, Latency: 5.429},
		{TimeOffset: 0.078683, Latency: 3.043},
	}

	p := NewParser(reader)

	for index, output := range outputs {
		if index >= len(outputs) {
			t.Errorf("had already received all the expected latencyReports when I saw this: %v", output)
		}
		candidateReport := <-p.Next
		if candidateReport != output {
			t.Errorf("reply mismatch: expected %v, got %v", output, candidateReport)
		}
	}
}
