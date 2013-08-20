// Package notifier contains types and methods for sending notifications
// that a job has completed.
package notifier

import (
	"io"
	"net/http"
)

// To create notifications, an object has to have a destination in mind
// and be able to serialize itself into a JSON message.
type Notifiable interface {
	Recipient() string
	Content() io.Reader
}

func SendNotification(n Notifiable) (response *http.Response, err error) {
	destination := n.Recipient()
	body := n.Content()
	response, err = http.Post(destination, "application/json", body)
	return
}
