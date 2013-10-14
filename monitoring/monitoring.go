// Package monitoring adds functionality to send metrics to AWS
package monitoring

import (
	"github.com/PeopleAdmin/pdfcombiner/combiner"
	"github.com/crowdmob/goamz/aws"
	. "github.com/crowdmob/goamz/cloudwatch"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const metricTransmitInterval = 60
const ubuntuIdFile = "/var/lib/cloud/data/previous-instance-id"
const namespace = "pdfcombiner"

var instanceId string
var dimensions = []Dimension{}
var cw *CloudWatch

func StartMetricSender() {
	connect()
	addStackToDimensions()
	log.Println("Sending metrics to AWS every", metricTransmitInterval, "seconds.")

	receiver := make(chan MetricDatum)
	go collectBatchesForDelivery(receiver)

	go NewAggregator("JobsWaiting", receiver, 1, "Count",
		func() float64 {
			return float64(combiner.CurrentWait())
		}).Start()

	go NewAggregator("JobsRunning", receiver, 1, "Count",
		func() float64 {
			return float64(combiner.CurrentJobs())
		}).Start()

	go NewAggregator("RequestsReceivedPerMinute", receiver, 60, "Count",
		DeltaSinceLastCall(func() float64 {
			return float64(combiner.TotalJobsReceived())
		})).Start()

	go NewAggregator("RequestsFinishedPerMinute", receiver, 60, "Count",
		DeltaSinceLastCall(func() float64 {
			return float64(combiner.TotalJobsFinished())
		})).Start()

	go NewAggregator("ExcessRequestsPerMinute", receiver, 60, "Count", excessJobsMetric()).Start()

	// System metrics
	go NewAggregator("DiskFree", receiver, 60, "Megabytes", RootPartitionFreeMb()).Start()
	go NewAggregator("MemFree", receiver, 60, "Megabytes", FreeMemoryMb()).Start()
}

// A generator that when called return the net number of jobs proccessed since
// last called -- that is, (jobs_finished - jobs_received)
func excessJobsMetric() func() float64 {
	receivedGenerator := DeltaSinceLastCall(func() float64 {
		return float64(combiner.TotalJobsReceived())
	})
	processedGenerator := DeltaSinceLastCall(func() float64 {
		return float64(combiner.TotalJobsFinished())
	})

	return func() float64 {
		return receivedGenerator() - processedGenerator()
	}
}

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

// Create a new MetricDatum representing statistics about several observations.
func newStatMetricDatum(name string, observations []float64, unitName string) (stats MetricDatum, allZero bool) {
	allZero, min, max, sum := calculateStats(observations)
	stats = MetricDatum{
		MetricName: name,
		Dimensions: dimensions,
		Timestamp:  time.Now(),
		Unit:       unitName,
		StatisticValues: []StatisticSet{
			StatisticSet{
				Maximum:     max,
				Minimum:     min,
				SampleCount: float64(len(observations)),
				Sum:         sum,
			},
		},
	}
	return
}

// Create a MetricDatum based on a single datapoint.
func newSimpleMetricDatum(name string, value int32) MetricDatum {
	return MetricDatum{
		Dimensions: dimensions,
		MetricName: name,
		Timestamp:  time.Now(),
		Unit:       "Count",
		Value:      float64(value) + 1,
	}
}

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

func sendMetrics() {
	sendCurrentCounts()
}

func sendCurrentCounts() {
	metrics := []MetricDatum{
		queueLengthMetric(),
		runningLengthMetric(),
	}
	_ , err := cw.PutMetricData(metrics)
	if err != nil {
		log.Println("ERROR sending metrics to AWS:" + err.Error())
	}
}

func queueLengthMetric() MetricDatum {
	return newSimpleMetricDatum("JobsWaiting", combiner.CurrentWait())
}

func runningLengthMetric() MetricDatum {
	return newSimpleMetricDatum("JobsRunning", combiner.CurrentJobs())
}

func addInstanceIdToDimensions() (err error) {
	content, err := ioutil.ReadFile(ubuntuIdFile)
	if err == nil {
		instanceId = strings.Trim(string(content), "\n")
		dimensions = append(dimensions, Dimension{Name: "InstanceId", Value: instanceId})
	} else {
		log.Println("ERROR: Not able to get instanceid for monitoring, ensure that", ubuntuIdFile, "exists")
	}
	return
}

func addStackToDimensions() {
	if os.Getenv("STACK_NAME") != "" {
		dimensions = append(dimensions,
			Dimension{Name: "StackName", Value: os.Getenv("STACK_NAME")})
	} else {
		log.Printf("ERROR: Not able to get stack name, ensure that $STACK_NAME is set.")
	}
}

func connect() {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic("Couldn't connect to amazon:" + err.Error())
	}
	cw, err = NewCloudWatch(auth, aws.USEast.CloudWatchServicepoint, namespace)
	if err != nil {
		panic("Couldn't connect to amazon:" + err.Error())
	}
}
