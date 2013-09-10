// Package job provides the Job type and some methods on it.
// TODO caching s3 connection objects in a pool might speed things up.
// TODO use callback interface for different types of success
// notifications (e.g. command line and URL callbacks)
package job

import (
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/testmode"
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
	Downloaded     []*Document					`json:"downloaded"`
	CombinedKey    string               `json:"combined_key"`
	Title          string               `json:"title,omitempty"`
	Callback       string               `json:"callback"`
	Errors         []error              `json:"errors"`
	uploadComplete bool
	bucket         *s3.Bucket
}

// New is the default Job constructor.
func New(bucket string, docs []string) (newJob *Job, err error) {
	newJob = &Job{
		BucketName: bucket,
	}
	newJob.DocList = docsFromStrings(docs)
	err = newJob.setup()
	return
}

// IsValid determines whether the job contains all the fields necessary
// to start the combination.
func (j *Job) IsValid() bool {
	return (j.BucketName != "") &&
		(j.CombinedKey != "") &&
		(j.Callback != "") &&
		j.DocsAreValid()
}

func (j *Job) DocsAreValid() bool {
	if j.DocCount() <= 0 {
		return false
	}
	for _, doc := range j.DocList {
		if !doc.isValid() {
			return false
		}
	}
	return true
}

// DocCount returns the number of documents requested.
func (j *Job) DocCount() int {
	return len(j.DocList)
}

// CompleteCount returns the number of documents actually completed.
func (j *Job) CompleteCount() int {
	return len(j.Downloaded)
}

// AddError adds to the list of encountered errors, translating obscure ones.
func (j *Job) AddError(newErr error) {
	log.Println(newErr)
	if strings.Contains(newErr.Error(), "Get : 301 response missing Location header") {
		newErr = fmt.Errorf("bucket %s not accessible from this account", j.BucketName)
	}
	j.Errors = append(j.Errors, newErr)
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
	err = j.s3Connect()
	j.Downloaded = make([]*Document, 0, len(j.DocList))
	j.Errors = make([]error, 0, len(j.DocList))
	j.setDocParents()
	return
}

// Initialize the S3 connection used to download and upload docs.
func (j *Job) s3Connect() (err error) {
	err = j.connect()
	return
}

// Point each Document's parent pointer to the enclosing job.
func (j *Job) setDocParents() {
	for i, _ := range j.DocList {
		j.DocList[i].parent = j
	}
}
