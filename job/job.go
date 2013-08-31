// Package job provides the Job type and some methods on it.
// TODO caching s3 connection objects in a pool might speed things up.
// TODO use callback interface for different types of success
// notifications (e.g. command line and URL callbacks)
package job

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/stat"
	"github.com/PeopleAdmin/pdfcombiner/testmode"
	"io"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
	"strings"
)

// A Job includes all the data necessary to execute a pdf combination.
// It is mainly constructed from a JSON string in a HTTP request, but the
// last two fields contain internal state.
type Job struct {
	BucketName     string               `json:"bucket_name"`
	DocList        []Document           `json:"doc_list"`
	Downloaded     []string             `json:"downloaded"`
	CombinedKey    string               `json:"combined_key"`
	Title          string               `json:"title,omitempty"`
	Callback       string               `json:"callback"`
	Errors         map[string]string    `json:"errors"`
	PerfStats      map[string]stat.Stat `json:"perf_stats"`
	uploadComplete bool
	bucket         *s3.Bucket
}

// A Document has a (file) name and a human readable title, possibly
// used for watermarking prior to combination.
type Document struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

// A JobResponse is sent as a notification.  It includes the success
// status and a copy of the original job.
type JobResponse struct {
	Success bool `json:"success"`
	Job     Job  `json:"job"`
}

// Given a slice of document names, return a slice of Documents.
func docsFromStrings(names []string) (docs []Document) {
	docs = make([]Document, len(names))
	for i, name := range names {
		docs[i] = Document{Name: name}
	}
	return
}

// New is the default Job constructor.
func New(bucket string, docs []string) (newJob *Job, err error) {
	newJob = &Job{
		BucketName: bucket,
		DocList:    docsFromStrings(docs),
	}
	err = newJob.setup()
	return
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

// IsValid determines whether the job contains all the fields necessary
// to start the combination.
func (j *Job) IsValid() bool {
	return (j.BucketName != "") &&
		(j.CombinedKey != "") &&
		(j.Callback != "") &&
		(j.DocCount() > 0)
}

// Get retrieves the requested document from S3 as a byte slice.
func (j *Job) Get(docname string) (data []byte, err error) {
	if !testmode.IsEnabled() {
		data, err = j.bucket.Get(j.s3Path(docname))
	}
	return
}

// UploadCombinedFile sends a file to the job's CombinedKey on S3.
func (j *Job) UploadCombinedFile(localPath string) (err error) {
	if testmode.IsEnabled() {
		return
	}
	content, err := ioutil.ReadFile(localPath)
	if err != nil {
		j.AddError(j.CombinedKey, err)
		return
	}
	err = j.bucket.Put(j.CombinedKey, content, "application/pdf", s3.Private)
	if err != nil {
		j.AddError(j.CombinedKey, err)
	}
	j.uploadComplete = true
	return
}

// DocCount returns the number of documents requested.
func (j *Job) DocCount() int {
	return len(j.DocList)
}

// CompleteCount returns the number of documents actually completed.
func (j *Job) CompleteCount() int {
	return len(j.Downloaded)
}

// MarkComplete adds a document to the list of downloaded docs.
func (j *Job) MarkComplete(newdoc string, info stat.Stat) {
	j.Downloaded = append(j.Downloaded, newdoc)
	j.PerfStats[newdoc] = info
}

// HasDownloadedDocs determines whether any documents been successfully
// downloaded.
// TODO is it appropriate to use this to determine success in ToJSON()?
func (j *Job) HasDownloadedDocs() bool {
	return len(j.Downloaded) > 0
}

// Recipient returns the notification URL to send status updates to.
func (j *Job) Recipient() string {
	return j.Callback
}

// ToJSON serializes the Job into a JSON byte slice.
func (j *Job) ToJSON() (jsonResponse []byte) {
	response := JobResponse{
		Job:     *j,
		Success: j.isSuccessful(),
	}
	jsonResponse, _ = json.Marshal(response)
	fmt.Printf("%s\n", jsonResponse)
	return
}

// Content returns a Reader object that yields the job as JSON.
func (j *Job) Content() io.Reader {
	return bytes.NewReader(j.ToJSON())
}

// AddError adds to the list of encountered errors, translating obscure ones.
func (j *Job) AddError(sourceFile string, newErr error) {
	log.Println(newErr)
	if strings.Contains(newErr.Error(), "Get : 301 response missing Location header") {
		newErr = fmt.Errorf("bucket %s not accessible from this account", j.BucketName)
	}
	j.Errors[sourceFile] = newErr.Error()
}

// In test mode, randomly fail 10% of the time.
func (j *Job) isSuccessful() bool {
	if testmode.IsEnabledFast() {
		return testmode.RandomBoolean(0.1)
	}
	return j.uploadComplete
}

// Initialize the fields which don't have usable zero-values.
func (j *Job) setup() (err error) {
	err = j.connect()
	j.Downloaded = make([]string, 0, len(j.DocList))
	j.Errors = make(map[string]string)
	j.PerfStats = make(map[string]stat.Stat)
	return
}

// Connect to AWS and add the handle to the Job object.
func (j *Job) connect() (err error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		return
	}
	s := s3.New(auth, aws.USEast)
	j.bucket = s.Bucket(j.BucketName)
	return
}

// Construct an absolute path within a bucket to a given document.
func (j *Job) s3Path(docname string) string {
	return docname
}
