package main

import "fmt"

type LatencyReport struct {
	TimeOffset float64
	Latency    float64
}

func (dp *LatencyReport) string() string {
	return fmt.Sprintf("offset: %f; latency: %f", dp.TimeOffset, dp.Latency)
}