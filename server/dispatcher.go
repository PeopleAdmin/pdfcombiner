package server

import (
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/combiner"
	"github.com/PeopleAdmin/pdfcombiner/job"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/sqs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var maxQueue = 30
var sqsTransmissions chan sqs.Message
var queueName string
var queue *sqs.Queue

func init() {
	if sqsEnabled() {
		auth, err := aws.EnvAuth()
		if err != nil {
			fmt.Println("SQS connect failed:", err)
			os.Exit(1)
		}
		s := sqs.New(auth, aws.USEast)
		q, err := s.GetQueue(queueName)
		if err != nil {
			fmt.Printf("Couldn't use queue '%v': %v\n", queueName, err)
			os.Exit(1)
		} else {
			queue = q
			combiner.SetQueue(queue)
			go combiner.PollSqsForJobs()
			log.Println("Overflow messages will be sent to queue", queueName)
		}
		sqsTransmissions = make(chan sqs.Message, 1000)
		go startSqsSender(sqsTransmissions)
	}
}

// Decides whether to locally enqueue a job with Combiner or to send it to SQS
// to be handled by shared processing capacity.
func dispatchNewJob(r *http.Request, c CombinerServer) (j *job.Job, err error) {
	jobContent, err := ioutil.ReadAll(r.Body)
	j, err = job.NewFromJSON(jobContent)
	logJobReceipt(r, j)
	queueSize := combiner.CurrentWait()
	if queueSize > maxQueue {
		log.Printf("%v Local queue length above %v, sending to SQS", j.Id(), queueSize)
		sqsTransmissions <- sqs.Message{MessageId: j.Id(), Body: string(jobContent)}
	} else {
		executeNow(c, j)
	}
	return
}

func executeNow(c CombinerServer, j *job.Job) {
	j.Setup()
	go combiner.Combine(j)
}

func startSqsSender(incoming <-chan sqs.Message) {
	for {
		go sendToSqs(<-incoming)
	}
}

// Our message size is likely to be relatively large (>40 KB) because of
// embedded documents.  So we use the batch message sender for its higher
// limits (256KB per message) but don't actually batch any messages because the
// total batch size must also by less than 256KB, which we risk hitting with a
// few large docs.
func sendToSqs(message sqs.Message) {
	_, err := queue.SendMessageBatch([]sqs.Message{message})
	if err != nil {
		log.Println(message.MessageId, "Couldn't send to SQS, processing locally:", err)
		j, _ := job.NewFromJSON([]byte(message.Body))
		j.Setup()
		go combiner.Combine(j)
	}
}

func sqsEnabled() bool {
	queueName = os.Getenv("SQS_QUEUE_NAME")
	return os.Getenv("SQS_ENABLED") != "" && queueName != ""
}
