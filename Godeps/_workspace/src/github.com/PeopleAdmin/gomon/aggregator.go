package gomon

import (
	"github.com/crowdmob/goamz/cloudwatch"
	"time"
)

// If this is false, periods in which a given metric is always zero will not
// result in an API transmission.  This will save money and bandwidth, but any
// alarms that depend on a quiescent system (such as a a scale-down event when
// a metric is low) might not fire, since the alarm will go into
// INSUFFICENT_DATA without any metrics.
const ShouldTransmitZero = true

// Only one stream is needed to send batches of data to.
var commonBatchStream = make(chan cloudwatch.MetricDatum)

// A common prefix for the metrics gathered by this app.
var Namespace string

// An Aggregator takes a fast stream of data points and batches them into
// a slower stream of statistical descriptions of those points.
type Aggregator struct {
	metricName      string
	produceInterval time.Duration
	get             func() float64
	metricStream    chan float64
	batchStream     chan cloudwatch.MetricDatum
	unitName        string
}

// StartLoop begins listening on the fast channel until one minute's worth of
// data is available, then sends the stats for that duration on the batch
// channel and repeats.  If all metrics for the duration are zero, optionally
// discard the values.
func (m *Aggregator) StartLoop() {
	bufferLen := metricTransmitInterval / int(m.produceInterval)
	buffer := make([]float64, bufferLen)
	go getMetricsForever(m)

	for {
		for i := 0; i < bufferLen; i++ {
			buffer[i] = <-m.metricStream
		}
		stats, allZero := newStatMetricDatum(m.metricName, buffer, m.unitName)

		if ShouldTransmitZero || !allZero {
			debug("batching metric:", m.metricName, stats)
			m.batchStream <- stats
		}
	}
}

// Produce metrics continuously by calling the given function.
func getMetricsForever(m *Aggregator) {
	for {
		m.metricStream <- m.get()
		time.Sleep(m.produceInterval * time.Second)
	}
}
