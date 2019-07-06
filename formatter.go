package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Formatter struct {
	timeWindow       int
	samplesPerSecond int
}

func (f Formatter) writeDataAsJSON(replies []Reply, w io.Writer) {
	fmt.Fprintln(w, "{")

	rows := make([]string, f.samplesPerSecond)
	for i := 0; i < f.samplesPerSecond; i++ {
		rows[i] = strconv.Itoa(i)
	}

	rowsJson := strings.Join(rows, ",")
	fmt.Fprintf(w, "\t\"rows\": [%s],\n", rowsJson)

	columns := make([]string, f.timeWindow)
	for i := 0; i < f.timeWindow; i++ {
		columns[i] = strconv.Itoa(i)
	}
	columnsJson := strings.Join(columns, ",")
	fmt.Fprintf(w, "\t\"columns\": [%s],\n", columnsJson)

	fmt.Fprint(w, "\t\"values\": [")

	for i := 0; i < f.timeWindow; i++ {
		var vals = replies[i*f.samplesPerSecond : (i+1)*f.samplesPerSecond]
		latencies := make([]string, len(vals), len(vals))
		for j := 0; j < len(vals); j++ {
			latencies[j] = fmt.Sprintf("%f", replies[i*f.samplesPerSecond + j].Latency)
		}
		fmt.Fprintf(w, "\t[%s]", strings.Join(latencies, ","))

		if i != (f.timeWindow - 1) {
			fmt.Fprintln(w, ",")
		}

		fmt.Fprint(w, "\n")
	}

	fmt.Fprint(w, "\t]")

	fmt.Fprintln(w, "}")
}