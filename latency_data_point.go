package main

import "fmt"

type LatencyDataPoint struct {
	TimeOffset float64
	Latency    float64
}

func (dp *LatencyDataPoint) string() string {
	return fmt.Sprintf("offset: %f; latency: %f", dp.TimeOffset, dp.Latency)
}