package job

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"time"
)

// A JobResponse is sent as a notification -- it includes the success
// status as well as a subset of the job fields.
type JobResponse struct {
	Success    bool        `json:"success"`
	Errors     []string    `json:"errors"`
	Downloaded []*Document `json:"downloaded"`
	Times      metrics     `json:"times"`
}

type metrics struct {
	ReceivedAt   time.Time
	DownloadTime time.Duration
	TotalTime    time.Duration
}

// NewFromJSON constructs a Job from an io.Reader containing JSON
// conforming to the required portion of the Job schema.
func NewFromJSON(encoded io.Reader) (newJob *Job, err error) {
	newJob = &Job{}
	err = json.NewDecoder(encoded).Decode(newJob)
	if err != nil {
		return
	}
	if !newJob.IsValid() {
		err = errors.New("missing required fields")
		return
	}
	err = newJob.setup()
	return
}

// ToJSON serializes the Job into a JSON byte slice.
func (j *Job) ToJSON() (jsonResponse []byte) {
	response := JobResponse{
		Success:    j.IsSuccessful(),
		Errors:     errStrs(j.Errors),
		Downloaded: j.Downloaded,
		Times:      newMetrics(j),
	}
	jsonResponse, _ = json.Marshal(response)
	return
}

func errStrs(errors []error) (strings []string) {
	strings = make([]string, len(errors))
	for i, err := range errors {
		strings[i] = err.Error()
	}
	return
}

func newMetrics(j *Job) (m metrics) {
	m.ReceivedAt = j.receivedAt
	m.DownloadTime = j.DownloadsDoneAt.Sub(j.receivedAt)
	m.TotalTime = time.Since(j.receivedAt)
	return
}

// Recipient returns the notification URL to send status updates to.
func (j *Job) Recipient() string {
	return j.Callback
}

// Content returns a Reader object that yields the job as JSON.
func (j *Job) Content() io.Reader {
	return bytes.NewReader(j.ToJSON())
}
