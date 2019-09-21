package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestFormatter_writeDataAsJSON(t *testing.T) {
	type fields struct {
		timeWindow       int
		samplesPerSecond int
	}
	type args struct {
		latencyReports []LatencyReport
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		output string
	}{
		{
			"empty list",
			fields{
				0,
				0,
			},
			args{
				[]LatencyReport{},
			},
			`[]`,
		},
		{
			"one second",
			fields{
				1,
				10,
			},
			args{
				[]LatencyReport{
					{TimeOffset: 0.050760, Latency: 7.373},
					{TimeOffset: 0.125935, Latency: 5.189},
					{TimeOffset: 0.201725, Latency: 2.886},
					{TimeOffset: 0.304634, Latency: 4.348},
					{TimeOffset: 0.405814, Latency: 6.587},
					{TimeOffset: 0.506894, Latency: 2.953},
					{TimeOffset: 0.612703, Latency: 8.544},
					{TimeOffset: 0.718270, Latency: 3.726},
					{TimeOffset: 0.811362, Latency: 7.517},
					{TimeOffset: 0.912460, Latency: 8.353},
				},
			},
			`[{"offset": 0, "subsecond-offset": 0, "latency": 7.373000},{"offset": 0, "subsecond-offset": 100, "latency": 5.189000},{"offset": 0, "subsecond-offset": 200, "latency": 2.886000},{"offset": 0, "subsecond-offset": 300, "latency": 4.348000},{"offset": 0, "subsecond-offset": 400, "latency": 6.587000},{"offset": 0, "subsecond-offset": 500, "latency": 2.953000},{"offset": 0, "subsecond-offset": 600, "latency": 8.544000},{"offset": 0, "subsecond-offset": 700, "latency": 3.726000},{"offset": 0, "subsecond-offset": 800, "latency": 7.517000},{"offset": 0, "subsecond-offset": 900, "latency": 8.353000}]`,
		},
		{
			"two seconds",
			fields{
				2,
				10,
			},
			args{
				[]LatencyReport{
					{TimeOffset: 0.050760, Latency: 7.373},
					{TimeOffset: 0.125935, Latency: 5.189},
					{TimeOffset: 0.201725, Latency: 2.886},
					{TimeOffset: 0.304634, Latency: 4.348},
					{TimeOffset: 0.405814, Latency: 6.587},
					{TimeOffset: 0.506894, Latency: 2.953},
					{TimeOffset: 0.612703, Latency: 8.544},
					{TimeOffset: 0.718270, Latency: 3.726},
					{TimeOffset: 0.811362, Latency: 7.517},
					{TimeOffset: 0.912460, Latency: 8.353},
					{TimeOffset: 0.050760, Latency: 9.373},
					{TimeOffset: 0.125935, Latency: 9.189},
					{TimeOffset: 0.201725, Latency: 9.886},
					{TimeOffset: 0.304634, Latency: 9.348},
					{TimeOffset: 0.405814, Latency: 9.587},
					{TimeOffset: 0.506894, Latency: 9.953},
					{TimeOffset: 0.612703, Latency: 9.544},
					{TimeOffset: 0.718270, Latency: 9.726},
					{TimeOffset: 0.811362, Latency: 9.517},
					{TimeOffset: 0.912460, Latency: 9.353},
				},
			},
			`[{"offset": 0, "subsecond-offset": 0, "latency": 7.373000},{"offset": 0, "subsecond-offset": 100, "latency": 5.189000},{"offset": 0, "subsecond-offset": 200, "latency": 2.886000},{"offset": 0, "subsecond-offset": 300, "latency": 4.348000},{"offset": 0, "subsecond-offset": 400, "latency": 6.587000},{"offset": 0, "subsecond-offset": 500, "latency": 2.953000},{"offset": 0, "subsecond-offset": 600, "latency": 8.544000},{"offset": 0, "subsecond-offset": 700, "latency": 3.726000},{"offset": 0, "subsecond-offset": 800, "latency": 7.517000},{"offset": 0, "subsecond-offset": 900, "latency": 8.353000},{"offset": 1, "subsecond-offset": 0, "latency": 9.373000},{"offset": 1, "subsecond-offset": 100, "latency": 9.189000},{"offset": 1, "subsecond-offset": 200, "latency": 9.886000},{"offset": 1, "subsecond-offset": 300, "latency": 9.348000},{"offset": 1, "subsecond-offset": 400, "latency": 9.587000},{"offset": 1, "subsecond-offset": 500, "latency": 9.953000},{"offset": 1, "subsecond-offset": 600, "latency": 9.544000},{"offset": 1, "subsecond-offset": 700, "latency": 9.726000},{"offset": 1, "subsecond-offset": 800, "latency": 9.517000},{"offset": 1, "subsecond-offset": 900, "latency": 9.353000}]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			f := Formatter{
				timeWindow:       tt.fields.timeWindow,
				samplesPerSecond: tt.fields.samplesPerSecond,
			}
			f.writeDataAsJSON(tt.args.latencyReports, buf)

			rawFormattedOutput, _ := ioutil.ReadAll(buf)
			formattedOutput := string(rawFormattedOutput)
			formattedOutput = strings.ReplaceAll(formattedOutput, "\n", "")
			formattedOutput = strings.ReplaceAll(formattedOutput, "\t", "")
			if formattedOutput != tt.output {
				t.Errorf("incorrect output: wanted '%s', got '%s'", tt.output, formattedOutput)
			}
		})
	}
}
