package combiner

import (
  "time"
  "log"
  "fmt"
  "io/ioutil"
  "errors"
  "pdfcombiner/cpdf"
  "pdfcombiner/job"
)

var (
  verbose = true
  bucketName = "pa-hrsuite-production"
  basedir = "/tmp/"
  keybase = "606/docs/"
  defdocs = []string{"100001.pdf","100009.pdf","100030.pdf","100035.pdf",
                     "100035.pdf","100037.pdf","100038.pdf","100093.pdf"}
)

type stat struct {
  filename string
  size     int
  dlSecs   time.Duration
}

func getFile(j *job.Job, docname string, c chan stat, e chan error) {
  start := time.Now()
  data, err := j.Get(docname)
  if err != nil {
    e <- fmt.Errorf("%v: %v", docname, err)
    return
  }
  path := "/tmp/"+docname
  err = ioutil.WriteFile(path, data, 0644)
  if err != nil {
    e <- err; return
    return
  }
  c <- stat{ filename: docname,
             size: len(data),
             dlSecs: time.Since(start) }
}

func getAllFiles(j *job.Job) {
  start := time.Now()
  j.Connect()
  c := make(chan stat)
  e := make(chan error)
  for _,doc := range j.DocList{
    go getFile(j,doc,c,e)
  }

  totalBytes := waitForDownloads(j,c,e)
  printSummary(start, totalBytes, j.DocCount())
}

func waitForDownloads(j *job.Job, c chan stat, e chan error) (totalBytes int) {
  for _,_ = range j.DocList{
    select {
      case packet := <-c:
        if verbose { log.Printf("%s was %d bytes\n", packet.filename,packet.size) }
        totalBytes += packet.size
        j.MarkComplete(packet.filename)
      case err := <-e:
        j.AddError(err)
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

func postToCallback(j *job.Job){
  log.Println("work complete, posting success to callback:",j.Callback)
  log.Println(j)
}

func Combine(j *job.Job) bool {
  defer postToCallback(j)
  getAllFiles(j)
  if j.HasDownloadedDocs() {
    cpdf.Merge(j.Downloaded)
  }
  return true
}

