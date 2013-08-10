// Package Poster contains types and methods for sending notifications
// that a job has completed.
package poster

// To create notifications, an object has to have a destination in mind
// and be able to serialize itself into a JSON message.
type Notifiable interface {
	RecipientUrl() string
	ToJson() []byte
}

func SendNotification(n Notifiable) (success bool, err error) {
	success = true
	return
}
