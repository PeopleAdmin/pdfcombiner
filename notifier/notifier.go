// Package notifier contains types and methods for sending notifications
// that a job has completed.
package notifier

import (
	"io"
	"fmt"
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
	client := &http.Client{}
	req, _ := http.NewRequest("POST", destination, body)
	req.SetBasicAuth("admin","admin")
	req.Header.Set("Content-Type","application/json")
	response, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	return
}
