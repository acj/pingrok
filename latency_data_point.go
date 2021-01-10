package main

import "fmt"

type latencyDataPoint struct {
	timeOffset float64
	latency    float64
}

func (dp *latencyDataPoint) string() string {
	return fmt.Sprintf("offset: %f; latency: %f", dp.timeOffset, dp.latency)
}
