// Package gomon adds functionality to aggregate and batch groups of metrics
// for delivery to AWS Cloudwatch.
package gomon

import (
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	"log"
	"os"
)

const metricTransmitInterval = 60

var cw *cloudwatch.CloudWatch

// A set of AWS Dimensions that are added to all sent metrics.
var Dimensions = []cloudwatch.Dimension{}

// The Cloudwatch region to use.
var Region = aws.USEast

// Start begins collecting and batching metrics from all aggregators.
func Start() {
	connect()
	go startBatcher(commonBatchStream)

	for i := range registry {
		debug("starting", registry[i])
		go registry[i].StartLoop()
	}
}

// AddDimension adds a dimension to all delivered metrics.  A good thing to use
// would be something like StackName or InstanceId.
func AddDimension(name, value string) {
	Dimensions = append(Dimensions, cloudwatch.Dimension{Name: name, Value: value})
}

func connect() {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic("Couldn't connect to amazon:" + err.Error())
	}
	cw, err = cloudwatch.NewCloudWatch(auth, Region.CloudWatchServicepoint, Namespace)
	if err != nil {
		panic("Couldn't connect to amazon:" + err.Error())
	}
}

func debug(args ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		args = append([]interface{}{"MONITORING"}, args...)
		log.Println(args...)
	}
}
