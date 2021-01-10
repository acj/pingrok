package main

type dataPointPartitioner struct {
	dataPointBuffer  *CircularBuffer
	timeWindow       int
	samplesPerSecond int
}

func newDataPointPartitioner(buffer *CircularBuffer, timeWindow int, samplesPerSecond int) *dataPointPartitioner {
	return &dataPointPartitioner{
		dataPointBuffer:  buffer,
		timeWindow:       timeWindow,
		samplesPerSecond: samplesPerSecond,
	}
}

func (s *dataPointPartitioner) start(dataPoints <-chan LatencyDataPoint) {
	dataPointsBinnedBySecond := make(chan []LatencyDataPoint)

	go s.partitionRepliesBySecond(dataPoints, dataPointsBinnedBySecond)
	go s.addToCircularBuffer(dataPointsBinnedBySecond)
}

func (s *dataPointPartitioner) partitionRepliesBySecond(in <-chan LatencyDataPoint, out chan<- []LatencyDataPoint) {
	// Assumption: inbound data points are ordered by time
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
			s.dataPointBuffer.Insert(r)
		}
	}
}
