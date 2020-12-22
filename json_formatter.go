package main

import (
	"fmt"
	"io"
	"strings"
)

type JSONFormatter struct {
	timeWindow       int
	samplesPerSecond int
}

func (f JSONFormatter) writeDataAsJSON(dataPoints []LatencyDataPoint, w io.Writer) {
	fmt.Fprintln(w, "[")

	for i := 0; i < f.timeWindow; i++ {
		var vals = dataPoints[i*f.samplesPerSecond : (i+1)*f.samplesPerSecond]
		latencies := make([]string, len(vals), len(vals))
		for j := 0; j < len(vals); j++ {
			latencies[j] = latencyReportAsJSON(i, j*(1000/f.samplesPerSecond), dataPoints[i*f.samplesPerSecond+j].Latency)
		}
		fmt.Fprintf(w, "\t%s", strings.Join(latencies, ","))

		if i != (f.timeWindow - 1) {
			fmt.Fprintln(w, ",")
		}

		fmt.Fprint(w, "\n")
	}

	fmt.Fprintln(w, "]")
}

func latencyReportAsJSON(offset int, subsecondOffset int, latency float64) string {
	return fmt.Sprintf(`{"offset": %d, "subsecond-offset": %d, "latency": %f}`, offset, subsecondOffset, latency)
}
