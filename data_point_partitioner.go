package main

type dataPointPartitioner struct {
	dataPointBuffer  *circularBuffer
	timeWindow       int
	samplesPerSecond int
}

func newDataPointPartitioner(buffer *circularBuffer, timeWindow int, samplesPerSecond int) *dataPointPartitioner {
	return &dataPointPartitioner{
		dataPointBuffer:  buffer,
		timeWindow:       timeWindow,
		samplesPerSecond: samplesPerSecond,
	}
}

func (s *dataPointPartitioner) start(dataPoints <-chan latencyDataPoint) {
	dataPointsBinnedBySecond := make(chan []latencyDataPoint)

	go s.partitionRepliesBySecond(dataPoints, dataPointsBinnedBySecond)
	go s.addToCircularBuffer(dataPointsBinnedBySecond)
}

func (s *dataPointPartitioner) partitionRepliesBySecond(in <-chan latencyDataPoint, out chan<- []latencyDataPoint) {
	// Assumption: inbound data points are ordered by time
	currentAccumulatorSecondOffset := 0
	timeQuantum := 1.0 / float64(s.samplesPerSecond)
	currentSlice := make([]latencyDataPoint, s.samplesPerSecond, s.samplesPerSecond)
	for r := range in {
		currentSecond := int(r.timeOffset)
		if currentAccumulatorSecondOffset != currentSecond {
			out <- currentSlice
			currentSlice = make([]latencyDataPoint, s.samplesPerSecond, s.samplesPerSecond)
			currentAccumulatorSecondOffset = currentSecond
		}

		currentSubsecondOffset := r.timeOffset - float64(int(r.timeOffset))
		currentSlice[int(currentSubsecondOffset/timeQuantum)] = r
	}
}

func (s *dataPointPartitioner) addToCircularBuffer(replies chan []latencyDataPoint) {
	for oneSecondOfData := range replies {
		for _, r := range oneSecondOfData {
			s.dataPointBuffer.insert(r)
		}
	}
}
