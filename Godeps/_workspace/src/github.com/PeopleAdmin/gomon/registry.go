package gomon

import (
	"fmt"
	"time"
)

// The Aggregators which data is being collected from in the main loop.
var registry []*Aggregator

// Register builds an Aggregator instance and adds it to the registry.
func Register(name string, interval time.Duration, unit string, floatEmitter func() float64) {
	newAggregator := &Aggregator{
		metricName:      name,
		produceInterval: interval,
		get:             floatEmitter,
		metricStream:    make(chan float64),
		batchStream:     commonBatchStream,
		unitName:        unit,
	}
	registry = append(registry, newAggregator)
}

// RegisterInt is a convenience function that wraps int emitters with a float cast.
func RegisterInt(name string, interval time.Duration, unit string, intEmitter func() int) {
	wrapped := func() float64 { return float64(intEmitter()) }
	Register(name, interval, unit, wrapped)
}

// RegisterDelta wraps an emitter with DeltaSinceLastCall to turn it into a
// rate measurement instead of a point in time metric.
func RegisterDelta(name string, interval time.Duration, unit string, floatEmitter func() float64) {
	wrapped := DeltaSinceLastCall(floatEmitter)
	Register(name, interval, unit, wrapped)
}

// RegisterDeltaInt is RegisterDelta for ints.
func RegisterDeltaInt(name string, interval time.Duration, unit string, intEmitter func() int) {
	wrapped := func() float64 { return float64(intEmitter()) }
	RegisterDelta(name, interval, unit, wrapped)
}

func (a *Aggregator) String() string {
	return fmt.Sprintf("%s, emitting %d times per minute.", a.metricName, a.produceInterval)
}
