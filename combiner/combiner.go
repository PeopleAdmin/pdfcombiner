package combiner

import (
  "time"
  "log"
  "fmt"
  "io/ioutil"
  "errors"
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
  "pdfcombiner/cpdf"
)

var (
  verbose = true
  bucketName = "pa-hrsuite-production"
  basedir = "/tmp/"
  keybase = "606/docs/"
  defdocs = []string{"100001.pdf","100009.pdf","100030.pdf","100035.pdf",
                     "100035.pdf","100037.pdf","100038.pdf","100093.pdf"}
)

type Job struct {
  BucketName string
  EmployerId int
  DocList    []string
  Callback   string
  Errors     []error
  Bucket     *s3.Bucket
}

type stat struct {
  filename string
  size     int
  dlSecs   time.Duration
  err      error
}

func (j *Job) IsValid() bool {
 return (j.BucketName != "") &&
        (j.Callback   != "") &&
        (j.EmployerId > 0)   &&
        (j.DocCount() > 0)
}

func (j *Job) DocCount() int {
  return len(j.DocList)
}

func (j *Job) AddError(newErr error) {
  log.Println(newErr)
  j.Errors = append(j.Errors, newErr)
}

func (j *Job) handleGenericError(err error) {
  msg := fmt.Sprintf("failed: %s", err)
  j.AddError(errors.New(msg))
}

func (j *Job) handleGetError(err error, doc string) {
  msg := fmt.Sprintf("failed while getting doc %s: '%s'",doc, err)
  j.AddError(errors.New(msg))
}

func (j *Job) connect() {
  auth, err := aws.EnvAuth()
  if err != nil { panic(err) }
  s := s3.New(auth, aws.USEast)
  j.Bucket = s.Bucket(bucketName)
}

func (j *Job) getFile(docname string, c chan stat) {
  start := time.Now()
  data, err := j.Bucket.Get(j.s3Path(docname))
  if err != nil {
    j.handleGetError(err,docname)
    return
  }
  path := "/tmp/"+docname
  err = ioutil.WriteFile(path, data, 0644)
  if err != nil {
    j.handleGenericError(err)
  }
  c <- stat{ filename: docname,
             size: len(data),
             dlSecs: time.Since(start) }
}

func (j *Job) s3Path(docname string) string {
  return fmt.Sprintf("%d/docs/%s", j.EmployerId, docname)
}

func (j *Job) getAllFiles() {
  start := time.Now()
  j.connect()
  c := make(chan stat)
  for _,doc := range j.DocList{
    go j.getFile(doc,c)
  }

  totalBytes := j.waitForDownloads(c)
  printSummary(start, totalBytes, j.DocCount())
}

func (j *Job) waitForDownloads(c chan stat) (totalBytes int) {
  for _,_ = range j.DocList{
    select {
      case packet := <-c:
        if verbose { log.Printf("%s was %d bytes\n", packet.filename,packet.size) }
        totalBytes += packet.size
      case <-time.After(2 * time.Minute):
        j.AddError(errors.New("Timed out while downloading"))
        return
    }
  }
  return
}

func printSummary(start time.Time, bytes int, count int){
  elapsed := time.Since(start)
  seconds := elapsed.Seconds()
  mbps := float64(bytes) / 1024 / 1024 / seconds
  log.Printf("got %d bytes over %d files in %f secs (%f MB/s)\n",
             bytes, count, seconds, mbps)
}

func (j *Job) postToCallback(){
  log.Println("work complete, posting success to callback:",j.Callback)
  log.Println(j)
}

func (j *Job) Combine() bool {
  defer j.postToCallback()
  j.getAllFiles()
  cpdf.Merge(j.DocList)
  return true
}

