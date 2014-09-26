package gomon

import (
	"encoding/xml"
	"fmt"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	"log"
	"strconv"
	"time"
)

func startBatcher(incomingBatches <-chan cloudwatch.MetricDatum) {
	var buffer []cloudwatch.MetricDatum
	buffer = make([]cloudwatch.MetricDatum, 0)
	var newStat cloudwatch.MetricDatum
	nextSend := time.Now().Add(metricTransmitInterval * time.Second)
	for {
		select {
		case newStat = <-incomingBatches:
			buffer = append(buffer, newStat)
		case <-time.After(1 * time.Second):
		}
		if len(buffer) > 0 && time.Now().After(nextSend) {
			sendToAws(buffer)
			buffer = make([]cloudwatch.MetricDatum, 0)
			nextSend = time.Now().Add(metricTransmitInterval * time.Second)
		}
	}
}

func sendToAws(metrics []cloudwatch.MetricDatum) {
	debug("transmitting", len(metrics), "metrics to CW:", metrics)
	_, err := putData(cw, metrics)
	if err != nil {
		log.Println("ERROR sending metrics to AWS:" + err.Error())
	}
}

// Following two functions are adapted from
// github.com/crowdmob/goamz/cloudwatch, since the putData implementation of
// that package was broken.
func putData(c *cloudwatch.CloudWatch, metrics []cloudwatch.MetricDatum) (result *aws.BaseResponse, err error) {
	// Serialize the params
	params := aws.MakeParams("PutMetricData")
	for i, metric := range metrics {
		prefix := "MetricData.member." + strconv.Itoa(i+1)
		if metric.MetricName == "" {
			err = fmt.Errorf("no metric name supplied for metric: %d", i)
			return
		}
		params[prefix+".MetricName"] = metric.MetricName
		if metric.Unit != "" {
			params[prefix+".Unit"] = metric.Unit
		}
		if metric.Value != 0 {
			params[prefix+".Value"] = strconv.FormatFloat(metric.Value, 'f', -1, 64)
		}
		if !metric.Timestamp.IsZero() {
			params[prefix+".Timestamp"] = metric.Timestamp.UTC().Format(time.RFC3339)
		}
		for j, dim := range metric.Dimensions {
			dimprefix := prefix + ".Dimensions.member." + strconv.Itoa(j+1)
			params[dimprefix+".Name"] = dim.Name
			params[dimprefix+".Value"] = dim.Value
		}
		if len(metric.StatisticValues) == 1 {
			stat := metric.StatisticValues[0]
			statprefix := prefix + ".StatisticValues"
			params[statprefix+".Maximum"] = strconv.FormatFloat(stat.Maximum, 'f', -1, 64)
			params[statprefix+".Minimum"] = strconv.FormatFloat(stat.Minimum, 'f', -1, 64)
			params[statprefix+".SampleCount"] = strconv.FormatFloat(stat.SampleCount, 'f', -1, 64)
			params[statprefix+".Sum"] = strconv.FormatFloat(stat.Sum, 'f', -1, 64)
		}
	}
	result = new(aws.BaseResponse)
	err = query(c, "POST", "/", params, result)
	return
}

func query(c *cloudwatch.CloudWatch, method, path string, params map[string]string, resp interface{}) error {
	// Add basic Cloudwatch param
	params["Version"] = "2010-08-01"
	if Namespace != "" {
		params["Namespace"] = Namespace
	}

	r, err := c.Service.Query(method, path, params)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return c.Service.BuildError(r)
	}
	err = xml.NewDecoder(r.Body).Decode(resp)
	return err
}
