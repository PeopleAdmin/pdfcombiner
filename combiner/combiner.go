package pdfcombiner

import (
  "time"
  "log"
  "fmt"
  "os/exec"
  "io/ioutil"
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
)

type CombineRequest struct {
  BucketName string
  EmployerId string
  DocList    []string
  Callback   string
  Errors     []string
  Bucket     *s3.Bucket
}

var (
  verbose bool
  bucketName = "pa-hrsuite-production"
  basedir = "/tmp/"
  keybase = "606/docs/"
  defdocs = []string{"100001.pdf","100009.pdf","100030.pdf","100035.pdf",
                     "100035.pdf","100037.pdf","100038.pdf","100093.pdf"}
)

func handleGenericError(req *CombineRequest) {
  if err := recover(); err != nil {
    msg := fmt.Sprintf("failed: %s", err)
    req.Errors = append(req.Errors, msg)
    log.Println(msg)
  }
}
func handleGetError(req *CombineRequest, doc string) {
  if err := recover(); err != nil {
    msg := fmt.Sprintf("failed while getting doc %s: '%s'",doc, err)
    req.Errors = append(req.Errors, msg)
    log.Println(msg)
  }
}

func getFile(req *CombineRequest, docname string, c chan int) {
  defer handleGetError(req, docname)
  s3key := keybase + docname
  data, err := req.Bucket.Get(s3key)
  if err != nil { panic(err) }

  path := "/tmp/"+docname
  err = ioutil.WriteFile(path, data, 0644)
  if err != nil { panic(err) }
  c <- len(data)
}

func connect() *s3.Bucket {
  auth, err := aws.EnvAuth()
  if err != nil { panic(err) }
  s := s3.New(auth, aws.USEast)
  return s.Bucket(bucketName)
}

func printSummary(start time.Time, bytes int, count int){
  elapsed := time.Since(start)
  seconds := elapsed.Seconds()
  mbps := float64(bytes) / 1024 / 1024 / seconds
  log.Printf("got %d bytes over %d files in %f secs (%f MB/s)\n",
             bytes, count, seconds, mbps)
}

func cpdfArgs(doclist []string) (args []string) {
  args = []string{"merge"}
  for _,doc := range doclist{
    args = append(args, doc)
  }
  args = append(args, "-o", "combined.pdf")
  return
}

func mergeWithCpdf(req *CombineRequest) {
  cpdf, err := exec.LookPath("cpdf")
  if err != nil { log.Fatal("no cpdf") }
  combine_cmd := exec.Command(cpdf)
  combine_cmd.Dir = "/tmp"
  combine_cmd.Args = cpdfArgs(req.DocList)
  out, err := combine_cmd.Output()
  log.Println(string(out))
}

func getAllFiles(req *CombineRequest) {
  defer handleGenericError(req)
  start := time.Now()
  req.Bucket = connect()
  c := make(chan int)
  for _,doc := range req.DocList{
    go getFile(req,doc,c)
  }

  totalBytes := 0
  for _,doc := range req.DocList{
    recieved := <-c
    if verbose{
      log.Printf("%s was %d bytes\n", doc,recieved)
    }
    totalBytes += recieved
  }
  printSummary(start,totalBytes,len(req.DocList))
}

func postToCallback(req *CombineRequest) {
  log.Println("work complete, posting success to callback:",req.Callback)
  log.Println(req)
}

func Combine(req *CombineRequest) {
  defer postToCallback(req)
  getAllFiles(req)
  mergeWithCpdf(req)
}

