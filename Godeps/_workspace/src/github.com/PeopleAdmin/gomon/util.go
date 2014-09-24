package gomon

import (
	"github.com/crowdmob/goamz/cloudwatch"
	"time"
)

// DeltaSinceLastCall transforms a function returning a metric count into a
// function that returns the difference since the last time that function was
// called.
func DeltaSinceLastCall(getter func() float64) (generator func() float64) {
	lastCount := getter()
	currentCount := lastCount
	generator = func() (delta float64) {
		currentCount = getter()
		delta = currentCount - lastCount
		lastCount = currentCount
		return
	}
	return
}

// Create a new cloudwatch.MetricDatum representing statistics about several observations.
func newStatMetricDatum(name string, observations []float64, unitName string) (stats cloudwatch.MetricDatum, allZero bool) {
	allZero, min, max, sum := calculateStats(observations)
	stats = cloudwatch.MetricDatum{
		MetricName: name,
		Dimensions: Dimensions,
		Timestamp:  time.Now(),
		Unit:       unitName,
		StatisticValues: []cloudwatch.StatisticSet{
			cloudwatch.StatisticSet{
				Maximum:     max,
				Minimum:     min,
				SampleCount: float64(len(observations)),
				Sum:         sum,
			},
		},
	}
	return
}

// A StatisticSet needs to know the min, max and sum of a slice of
// observations.  It's also useful to know whether all the members of the slice
// are zero.
func calculateStats(numbers []float64) (allZero bool, min, max, sum float64) {
	min = numbers[0]
	max = numbers[0]
	allZero = true
	for _, value := range numbers {
		if value != 0 {
			allZero = false
		}
		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
		sum += value
	}
	return
}

