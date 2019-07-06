package main

import "net/http"

func (s *Server) dataSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	s.formatter.writeDataAsJSON(s.latencyReportCircularBuffer.Snapshot(), w)
}