// Package job provides the Job type and some methods on it.
// TODO probably need a document type
package job

import(
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
  "log"
  "fmt"
  "strings"
  "encoding/json"
)

// A Job includes all the data necessary to execute a pdf combination.
// It is mainly constructed from a JSON string in a HTTP request, but the
// last two fields contain internal state.
type Job struct {
  BucketName string     `json:"bucket_name"`
  EmployerId int        `json:"employer_id"`
  DocList    []string   `json:"doc_list"`
  Downloaded []string   `json:"downloaded"`
  Callback   string     `json:"callback"`
  Errors     []error    `json:"errors"`
  bucket     *s3.Bucket
}

type JobResponse struct {
  Success bool `json:"success"`
  Job     Job  `json:"job"`
}

// Does the job contain all the fields necessary to start the combination?
func (j *Job) IsValid() bool {
 return (j.BucketName != "") &&
        (j.Callback   != "") &&
        (j.EmployerId > 0)   &&
        (j.DocCount() > 0)
}

// Retrieve the requested document from S3 as a byte slice.
func (j *Job) Get(docname string) (data []byte, err error) {
  //TODO j.connect() if necessary
  data, err = j.bucket.Get(j.s3Path(docname))
  return
}

// The number of documents requested.
func (j *Job) DocCount() int {
  return len(j.DocList)
}

// The number of documents actually completed.
func (j *Job) CompleteCount() int {
  return len(j.Downloaded)
}

// Add a document to the list of completed (downloaded) docs.
func (j *Job) MarkComplete(newdoc string) {
  j.Downloaded = append(j.Downloaded, newdoc)
}

// Have any documents been successfully downloaded?
// TODO is it appropriate to use this to determine success in ToJSON()?
func (j *Job) HasDownloadedDocs() bool {
  return len(j.Downloaded) > 0
}

func (j *Job) RecipientUrl() string{
  return j.Callback
}

func (j *Job) ToJson() (json_response []byte) {
  response := JobResponse{
    Job: *j,
    Success: j.HasDownloadedDocs() }
  json_response, _ = json.Marshal(response)
  fmt.Printf("%s\n",json_response)
  return
}

// Add to the list of encountered errors, translating obscure ones.
func (j *Job) AddError(newErr error) {
  log.Println(newErr)
  if strings.Contains(newErr.Error(), "Get : 301 response missing Location header") {
    newErr = fmt.Errorf("bucket %s not accessible from this account", j.BucketName)
  }
  j.Errors = append(j.Errors, newErr)
}

// Connect to AWS and add the handle to the Job object.
func (j *Job) Connect() {
  auth, err := aws.EnvAuth()
  if err != nil { panic(err) }
  s := s3.New(auth, aws.USEast)
  j.bucket = s.Bucket(j.BucketName)
}

// Construct an absolute path within a bucket to a given document.
func (j *Job) s3Path(docname string) string {
  return fmt.Sprintf("%d/docs/%s", j.EmployerId, docname)
}

