// Package notifier contains types and methods for sending notifications
// that a job has completed.
package notifier

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Notifiable objects have a destination in mind and a way to serialize
// themselves into a Reader that yields JSON text.
type Notifiable interface {
	Recipient() string
	Content() io.Reader
}

// SendNotification sends an HTTP message to the recipient with the job status.
func SendNotification(n Notifiable) (response *http.Response, err error) {
	destination := n.Recipient()
	payload := n.Content()
	client := &http.Client{}
	req, _ := http.NewRequest("POST", destination, payload)
	req.SetBasicAuth("admin", "password")
	req.Header.Set("Content-Type", "application/json")
	response, err = client.Do(req)
	if err != nil {
		fmt.Println("Error posting notification:", err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	fmt.Println("Notification response body:", string(body))
	return
}
