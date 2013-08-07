package job

import(
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
  "log"
  "fmt"
  "strings"
)

type Job struct {
  BucketName string
  EmployerId int
  DocList    []string
  Downloaded []string
  Callback   string
  Errors     []error
  Bucket     *s3.Bucket
}

func (j *Job) IsValid() bool {
 return (j.BucketName != "") &&
        (j.Callback   != "") &&
        (j.EmployerId > 0)   &&
        (j.DocCount() > 0)
}

func (j *Job) Get(docname string) (data []byte, err error) {
  data, err = j.Bucket.Get(j.s3Path(docname))
  return
}

func (j *Job) DocCount() int {
  return len(j.DocList)
}

func (j *Job) CompleteCount() int {
  return len(j.Downloaded)
}

func (j *Job) MarkComplete(newdoc string) {
  j.Downloaded = append(j.Downloaded, newdoc)
}

func (j *Job) HasDownloadedDocs() bool {
  return len(j.Downloaded) > 0
}

// Add to the list of encountered errors, translating obscure ones.
func (j *Job) AddError(newErr error) {
  log.Println(newErr)
  if strings.Contains(newErr.Error(), "Get : 301 response missing Location header") {
    newErr = fmt.Errorf("bucket %s not accessible from this account", j.BucketName)
  }
  j.Errors = append(j.Errors, newErr)
}

func (j *Job) Connect() {
  auth, err := aws.EnvAuth()
  if err != nil { panic(err) }
  s := s3.New(auth, aws.USEast)
  j.Bucket = s.Bucket(j.BucketName)
}

func (j *Job) s3Path(docname string) string {
  return fmt.Sprintf("%d/docs/%s", j.EmployerId, docname)
}

