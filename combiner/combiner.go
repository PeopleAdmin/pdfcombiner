package pdfcombiner

import (
  "time"
  "log"
  "os/exec"
  "io/ioutil"
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
)

var (
  verbose bool
  bucketName = "pa-hrsuite-production"
  basedir = "/tmp/"
  keybase = "606/docs/"
  defdocs = []string{"100001.pdf","100009.pdf","100030.pdf","100035.pdf",
                     "100035.pdf","100037.pdf","100038.pdf","100093.pdf"}
)

func handleError(cb, doc string) {
  if err := recover(); err != nil {
    log.Printf("work failed on doc %s, posting error '%s' to callback",doc, err)
  }
}

func getFile(bucket *s3.Bucket, docname string, c chan int) {
  defer handleError("http://callback.com", docname)
  s3key := keybase + docname
  data, err := bucket.Get(s3key)
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

func mergeWithCpdf(doclist []string) {
  cpdf, err := exec.LookPath("cpdf")
  if err != nil { log.Fatal("no cpdf") }
  combine_cmd := exec.Command(cpdf)
  combine_cmd.Dir = "/tmp"
  combine_cmd.Args = cpdfArgs(doclist)
  out, err := combine_cmd.Output()
  log.Println(string(out))
}

func getAllFiles(doclist []string, callback string) {
  start := time.Now()
  bucket := connect()
  c := make(chan int)
  for _,doc := range doclist{
    go getFile(bucket,doc,c)
  }

  totalBytes := 0
  for _,doc := range doclist{
    recieved := <-c
    if verbose{
      log.Printf("%s was %d bytes\n", doc,recieved)
    }
    totalBytes += recieved
  }
  printSummary(start,totalBytes,len(doclist))
}

// doclist is the list of pdfs to get, callback is the url to POST
// the job status to when finished.
func Combine(doclist []string, callback string) {
  getAllFiles(doclist,callback)
  mergeWithCpdf(doclist)
  log.Println("work complete, posting success to callback:",callback)
}

