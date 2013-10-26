// Pulls messages off of SQS when excess capacity is detected.
package combiner

import (
	"github.com/PeopleAdmin/pdfcombiner/job"
	"github.com/crowdmob/goamz/sqs"
	"log"
	"os"
	"time"
)

// How long to wait after receiving from an empty queue.  If the queue is
// nonempty, it will be pulled from continuously.
const sqsReceiveInterval = 30 * time.Second
const minCountBeforeReceive = 20

var queue *sqs.Queue

func PollSqsForJobs() {
	for {
		if CurrentJobs() < minCountBeforeReceive {
			debug("Checking SQS for new messages")
			resp, err := queue.ReceiveMessage(10)
			if err != nil {
				log.Println("ERROR! Couldn't receive messages from SQS", err)
				sleep()
				continue
			}
			debug("got", len(resp.Messages), "messages")
			for _, message := range resp.Messages {
				locallyEnqueue(message)
			}
		}
	}
}

func locallyEnqueue(message sqs.Message) {
	j, err := job.NewFromJSON([]byte(message.Body))
	if err != nil {
		log.Printf("ERROR! can't deserialize %v from sqs: %v\n", message.MessageId, err)
	}
	j.Setup()
	debug("Trying to locally enqueue", j.Id(), "callback:", j.Callback)
	if deleteFromSqs(&message) {
		go Combine(j)
	}
}

// TODO deletions should be batched...
func deleteFromSqs(message *sqs.Message) bool {
	_, err := queue.DeleteMessage(message)
	if err != nil {
		log.Printf("ERROR! couldn't delete sqs message %v, not processing: %v\n", message.MessageId, err)
		return false
	}
	debug(message.MessageId, "deleted")
	return true
}

func sleep() {
	time.Sleep(sqsReceiveInterval)
}

func debug(args ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		args = append([]interface{}{"SQS"}, args...)
		log.Println(args...)
	}
}

func SetQueue(q *sqs.Queue) {
	queue = q
}
