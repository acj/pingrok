package main

type Server struct {
	latencyReportCircularBuffer *CircularBuffer
	timeWindow                  int
	samplesPerSecond            int
}

func NewServer(timeWindow int, samplesPerSecond int) *Server {
	return &Server{
		latencyReportCircularBuffer: NewCircularBuffer(timeWindow * samplesPerSecond),
		timeWindow:                  timeWindow,
		samplesPerSecond:            samplesPerSecond,
	}
}

func (s *Server) Start(targetHost string) {
	replies := make(chan LatencyDataPoint)
	discretizedReplies := make(chan []LatencyDataPoint)

	pinger := NewPinger(targetHost, replies)
	pinger.Start()

	go s.partitionRepliesBySecond(replies, discretizedReplies)
	go s.addToCircularBuffer(discretizedReplies)
}

func (s *Server) Snapshot() []LatencyDataPoint {
	return s.latencyReportCircularBuffer.Snapshot()
}

func (s *Server) partitionRepliesBySecond(in <-chan LatencyDataPoint, out chan<- []LatencyDataPoint) {
	// Assumption: inbound latencyReports are ordered by time
	currentAccumulatorSecondOffset := 0
	timeQuantum := 1.0 / float64(s.samplesPerSecond)
	currentSlice := make([]LatencyDataPoint, s.samplesPerSecond, s.samplesPerSecond)
	for r := range in {
		currentSecond := int(r.TimeOffset)
		if currentAccumulatorSecondOffset != currentSecond {
			out <- currentSlice
			currentSlice = make([]LatencyDataPoint, s.samplesPerSecond, s.samplesPerSecond)
			currentAccumulatorSecondOffset = currentSecond
		}

		currentSubsecondOffset := r.TimeOffset - float64(int(r.TimeOffset))
		currentSlice[int(currentSubsecondOffset/timeQuantum)] = r
	}
}

func (s *Server) addToCircularBuffer(replies chan []LatencyDataPoint) {
	for oneSecondOfData := range replies {
		for _, r := range oneSecondOfData {
			s.latencyReportCircularBuffer.Insert(r)
		}
	}
}
