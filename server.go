package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Server struct {
	formatter                   *Formatter
	router                      *mux.Router
	httpServer                  *http.Server
	latencyReportCircularBuffer *CircularBuffer
	timeWindow                  int
	samplesPerSecond            int
}

func NewServer(timeWindow int, samplesPerSecond int) *Server {
	router := mux.NewRouter()

	s := &Server{
		formatter: 	&Formatter{
			timeWindow:       timeWindow,
			samplesPerSecond: samplesPerSecond,
		},
		router: router,
		httpServer: &http.Server{
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
			Handler: router,
		},
		latencyReportCircularBuffer: NewCircularBuffer(timeWindow*samplesPerSecond),
		timeWindow:                  timeWindow,
		samplesPerSecond:            samplesPerSecond,
	}

	router.HandleFunc("/data.json", s.dataSnapshotHandler)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	return s
}

func (s *Server) Serve(address string, targetHost string) {
	replies := make(chan LatencyReport)
	discretizedReplies := make(chan []LatencyReport)

	pinger := NewPinger(targetHost, replies)
	pinger.Start()

	go s.discretizeReplies(replies, discretizedReplies)
	go s.addToCircularBuffer(discretizedReplies)

	s.httpServer.Addr = address
	if err := s.httpServer.ListenAndServe(); err != nil {
		log.Printf("failed to start server: %s", err.Error())
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	s.httpServer.Shutdown(ctx)
}

func (s *Server) discretizeReplies(in <-chan LatencyReport, out chan<- []LatencyReport) {
	// Assumption: inbound latencyReports are ordered by time
	currentAccumulatorSecondOffset := 0
	timeQuantum := 1.0 / float64(s.samplesPerSecond)
	currentSlice := make([]LatencyReport, s.samplesPerSecond, s.samplesPerSecond)
	for r := range in {
		currentSecond := int(r.TimeOffset)
		if currentAccumulatorSecondOffset != currentSecond {
			out<- currentSlice
			currentSlice = make([]LatencyReport, s.samplesPerSecond, s.samplesPerSecond)
			currentAccumulatorSecondOffset = currentSecond
		}

		currentSubsecondOffset := r.TimeOffset - float64(int(r.TimeOffset))
		currentSlice[int(currentSubsecondOffset / timeQuantum)] = r
	}
}

func (s *Server) addToCircularBuffer(replies chan []LatencyReport) {
	for oneSecondOfData := range replies {
		for _, r := range oneSecondOfData {
			s.latencyReportCircularBuffer.Insert(r)
		}
	}
}