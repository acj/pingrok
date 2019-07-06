package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Server struct {
	formatter *Formatter
	router *mux.Router
	httpServer *http.Server
	replyCircularBuffer *Buffer
	timeWindow int
	samplesPerSecond int
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
		replyCircularBuffer: NewCircularBuffer(timeWindow*samplesPerSecond),
		timeWindow: timeWindow,
		samplesPerSecond: samplesPerSecond,
	}

	router.HandleFunc("/data.json", s.dataSnapshotHandler)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	return s
}

func (s *Server) Serve(address string) {
	replies := make(chan Reply)
	pinger := NewPinger(replies)
	pinger.Start()

	discretizedReplies := make(chan []Reply)
	go discretizeReplies(s.samplesPerSecond, replies, discretizedReplies)
	go s.addToCircularBuffer(discretizedReplies)

	s.httpServer.Addr = address
	if err := s.httpServer.ListenAndServe(); err != nil {
		log.Printf("failed to start server: %s", err.Error())
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	s.httpServer.Shutdown(ctx)
}

func discretizeReplies(samplesPerSecond int, in <-chan Reply, out chan<- []Reply) {
	// Assumption: inbound replies are ordered by time
	currentAccumulatorSecondOffset := 0
	timeQuantum := 1.0 / float64(samplesPerSecond)
	currentSlice := make([]Reply, samplesPerSecond, samplesPerSecond)

	for r := range in {
		currentSecond := int(r.TimeOffset)
		if currentAccumulatorSecondOffset != currentSecond {
			out<- currentSlice
			currentSlice = make([]Reply, samplesPerSecond, samplesPerSecond)
			currentAccumulatorSecondOffset = currentSecond
		}

		currentSubsecondOffset := r.TimeOffset - float64(int(r.TimeOffset))
		currentSlice[int(currentSubsecondOffset / timeQuantum)] = r
	}
}

func (s *Server) addToCircularBuffer(replies chan []Reply) {
	for oneSecondOfData := range replies {
		for _, r := range oneSecondOfData {
			s.replyCircularBuffer.Insert(r)
		}
	}
}