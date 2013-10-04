// Package notifier contains types and methods for sending notifications
// that a job has completed.
package notifier

import (
	"github.com/PeopleAdmin/pdfcombiner/config"
	"io"
	"log"
	"net/http"
)

// Notifiable objects have a destination in mind and a way to serialize
// themselves into a Reader that yields JSON text.
type Notifiable interface {
	Recipient() string
	Content() io.Reader
	IsSuccessful() bool
}

// SendNotification sends an HTTP message to the recipient with the job status.
func SendNotification(n Notifiable) (response *http.Response, err error) {
	response, err = deliver(n)
	if err != nil {
		log.Println("ERROR posting notification:", err)
		return
	}
	logResponse(n, response)
	return
}

func deliver(n Notifiable) (response *http.Response, err error) {
	destination := n.Recipient()
	payload := n.Content()
	client := &http.Client{}
	req, _ := http.NewRequest("POST", destination, payload)
	req.SetBasicAuth(config.TransmitUser(), config.TransmitPass())
	req.Header.Set("Content-Type", "application/json")
	response, err = client.Do(req)
	return
}

func logResponse(n Notifiable, response *http.Response) {
	log.Printf("%sNotified %s that job success was %v.  Got response %s\n",
		prefix(response), n.Recipient(), n.IsSuccessful(), response.Status)
}

// Prefix the log message with error if it's not a 'good' code.
func prefix(response *http.Response) string{
	switch response.StatusCode {
	case 200:
		return ""
	default:
		return "ERROR: "
	}
}
