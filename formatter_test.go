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
			`{"rows": [],"columns": [],"values": []}`,
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
			`{"rows": [0,1,2,3,4,5,6,7,8,9],"columns": [0],"values": [[7.373000,5.189000,2.886000,4.348000,6.587000,2.953000,8.544000,3.726000,7.517000,8.353000]]}`,
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
			`{"rows": [0,1,2,3,4,5,6,7,8,9],"columns": [0,1],"values": [[7.373000,5.189000,2.886000,4.348000,6.587000,2.953000,8.544000,3.726000,7.517000,8.353000],[9.373000,9.189000,9.886000,9.348000,9.587000,9.953000,9.544000,9.726000,9.517000,9.353000]]}`,
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
