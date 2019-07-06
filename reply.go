package main

import "fmt"

type Reply struct {
	TimeOffset float64
	Latency    float64
}

func (r *Reply) string() string {
	return fmt.Sprintf("offset: %f; latency: %f", r.TimeOffset, r.Latency)
}