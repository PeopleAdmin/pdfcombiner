package pdfcombiner

import (
  "time"
  "log"
  "fmt"
  "io/ioutil"
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
  cpdf "pdfcombiner/cpdf"
)

var (
  verbose bool
  bucketName = "pa-hrsuite-production"
  basedir = "/tmp/"
  keybase = "606/docs/"
  defdocs = []string{"100001.pdf","100009.pdf","100030.pdf","100035.pdf",
                     "100035.pdf","100037.pdf","100038.pdf","100093.pdf"}
)

type Job struct {
  BucketName string
  EmployerId string
  DocList    []string
  Callback   string
  Errors     []string
  Bucket     *s3.Bucket
}

func (j *Job) IsValid() bool {
 return (j.BucketName != "") &&
        (j.EmployerId != "") &&
        (j.Callback   != "") &&
        (j.DocCount() > 0)
}

func (j *Job) DocCount() int {
  return len(j.DocList)
}

func (j *Job) handleGenericError() {
  if err := recover(); err != nil {
    msg := fmt.Sprintf("failed: %s", err)
    j.Errors = append(j.Errors, msg)
    log.Println(msg)
  }
}

func handleGetError(req *Job, doc string) {
  if err := recover(); err != nil {
    msg := fmt.Sprintf("failed while getting doc %s: '%s'",doc, err)
    req.Errors = append(req.Errors, msg)
    log.Println(msg)
  }
}

func getFile(req *Job, docname string, c chan int) {
  defer handleGetError(req, docname)
  s3key := keybase + docname
  data, err := req.Bucket.Get(s3key)
  if err != nil { panic(err) }

  path := "/tmp/"+docname
  err = ioutil.WriteFile(path, data, 0644)
  if err != nil { panic(err) }
  c <- len(data)
}

func (j *Job) connect() {
  auth, err := aws.EnvAuth()
  if err != nil { panic(err) }
  s := s3.New(auth, aws.USEast)
  j.Bucket = s.Bucket(bucketName)
}

func printSummary(start time.Time, bytes int, count int){
  elapsed := time.Since(start)
  seconds := elapsed.Seconds()
  mbps := float64(bytes) / 1024 / 1024 / seconds
  log.Printf("got %d bytes over %d files in %f secs (%f MB/s)\n",
             bytes, count, seconds, mbps)
}

func (j *Job) getAllFiles() {
  defer j.handleGenericError()
  start := time.Now()
  j.connect()
  c := make(chan int)
  for _,doc := range j.DocList{
    go getFile(j,doc,c)
  }

  totalBytes := 0
  for _,doc := range j.DocList{
    recieved := <-c
    if verbose{
      log.Printf("%s was %d bytes\n", doc,recieved)
    }
    totalBytes += recieved
  }
  printSummary(start, totalBytes, j.DocCount())
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

