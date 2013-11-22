// Package monitoring adds functionality to send metrics to AWS
package monitoring

import (
	"github.com/PeopleAdmin/gomon"
	"github.com/PeopleAdmin/pdfcombiner/combiner"
	"github.com/PeopleAdmin/pdfcombiner/server"
	"log"
	"os"
)

func StartMetricSender() {
	addStackToDimensions()
	gomon.Namespace = "pdfcombiner"

	// App metrics
	gomon.RegisterInt("JobsWaiting", 1, "Count", combiner.CurrentWait)
	gomon.RegisterInt("JobsRunning", 1, "Count", combiner.CurrentJobs)
	gomon.RegisterDeltaInt("RequestsReceivedPerMinute", 60, "Count", combiner.TotalJobsReceived)
	gomon.RegisterDeltaInt("RequestsFinishedPerMinute", 60, "Count", combiner.TotalJobsFinished)
	gomon.Register("ExcessRequestsPerMinute", 60, "Count", excessJobsMetric())
	if server.SqsEnabled() {
		gomon.RegisterInt("GlobalQueueLength", 60, "Count", server.CurrentQueueLength)
	}

	// System metrics
	gomon.Register("DiskFree", 60, "Megabytes", gomon.RootPartitionFree)
	gomon.Register("MemFree", 60, "Megabytes", gomon.MemoryFree)

	gomon.Start()
}

// A generator that when called return the net number of jobs proccessed since
// last called -- that is, (jobs_finished - jobs_received)
func excessJobsMetric() func() float64 {
	receivedGenerator := gomon.DeltaSinceLastCall(func() float64 {
		return float64(combiner.TotalJobsReceived())
	})
	processedGenerator := gomon.DeltaSinceLastCall(func() float64 {
		return float64(combiner.TotalJobsFinished())
	})

	return func() float64 {
		return receivedGenerator() - processedGenerator()
	}
}

func addStackToDimensions() {
	if os.Getenv("STACK_NAME") != "" {
		gomon.AddDimension("StackName", os.Getenv("STACK_NAME"))
	} else {
		log.Printf("ERROR: Not able to get stack name, ensure that $STACK_NAME is set.")
	}
}
