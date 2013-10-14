package monitoring

import (
	"fmt"
	"github.com/crowdmob/goamz/cloudwatch"
	"time"
)

// If this is false, periods in which a given metric is alwys zero will not
// result in an API transmission.  This will save money and bandwidth, but any
// alarms that depend on a quiescent system (such as a a scale-down event when
// a metric is low) might not fire, since the alarm will go into
// INSUFFICENT_DATA without any metrics.
const shouldTransmitZero = true

// A metricAggregator takes a fast stream of data points and batches them into
// a slower stream of statistical descriptions of those points.
type metricAggregator struct {
	metricName      string
	produceInterval time.Duration
	get             func() float64
	metricStream    chan float64
	batchStream     chan cloudwatch.MetricDatum
	unitName        string
}

// NewAggregator returns a metricAggregator instance.
func NewAggregator(
	name string,
	sender chan cloudwatch.MetricDatum,
	interval time.Duration,
	unit string,
	getter func() float64) *metricAggregator {

	return &metricAggregator{
		metricName:      name,
		produceInterval: interval,
		get:             getter,
		metricStream:    make(chan float64),
		batchStream:     sender,
		unitName:        unit,
	}
}

// Start listening on the fast channel until one minute's worth of data is
// available, then send the stats for that duration on the batch channel and
// repeat.  If all metrics for the duration are zero, discard the data.
func (m *metricAggregator) Start() {
	bufferLen := metricTransmitInterval / int(m.produceInterval)
	buffer := make([]float64, bufferLen)
	go getMetricsForever(m)

	for {
		for i := 0; i < bufferLen; i++ {
			buffer[i] = <-m.metricStream
		}
		stats, allZero := newStatMetricDatum(m.metricName, buffer, m.unitName)

		if shouldTransmitZero || !allZero {
			fmt.Fprintln(logFile, "packaging", m.metricName, "batch:")
			m.batchStream <- stats
		}
	}
}

// Produce metrics continuously by calling the given function.
func getMetricsForever(m *metricAggregator) {
	for {
		m.metricStream <- m.get()
		time.Sleep(m.produceInterval * time.Second)
	}
}
