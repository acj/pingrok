package main

type dataPointPartitioner struct {
	latencyReportCircularBuffer *CircularBuffer
	timeWindow                  int
	samplesPerSecond            int
}

func newDataPointPartitioner(timeWindow int, samplesPerSecond int) *dataPointPartitioner {
	return &dataPointPartitioner{
		latencyReportCircularBuffer: NewCircularBuffer(timeWindow * samplesPerSecond),
		timeWindow:                  timeWindow,
		samplesPerSecond:            samplesPerSecond,
	}
}

func (s *dataPointPartitioner) start(targetHost string) {
	replies := make(chan LatencyDataPoint)
	discretizedReplies := make(chan []LatencyDataPoint)

	pinger := NewPinger(targetHost, replies)
	pinger.Start()

	go s.partitionRepliesBySecond(replies, discretizedReplies)
	go s.addToCircularBuffer(discretizedReplies)
}

func (s *dataPointPartitioner) snapshot() []LatencyDataPoint {
	return s.latencyReportCircularBuffer.Snapshot()
}

func (s *dataPointPartitioner) partitionRepliesBySecond(in <-chan LatencyDataPoint, out chan<- []LatencyDataPoint) {
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

func (s *dataPointPartitioner) addToCircularBuffer(replies chan []LatencyDataPoint) {
	for oneSecondOfData := range replies {
		for _, r := range oneSecondOfData {
			s.latencyReportCircularBuffer.Insert(r)
		}
	}
}
