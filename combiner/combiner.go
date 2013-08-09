// Package combiner contains methods to consume Job objects, orchestrating
// the download, combination, and upload of a group of related PDF files.
package combiner

import (
  "time"
  "log"
  "fmt"
  "io/ioutil"
  "errors"
  "pdfcombiner/cpdf"
  "pdfcombiner/job"
  "pdfcombiner/poster"
  "runtime"
)

var (
  downloadTimeout = 3 * time.Minute
  verbose = true
  basedir = "/tmp/"
  MaxGoroutines = 30
)

// A stat represents statistics about a completed document transfer operation.
type stat struct {
  filename string
  size     int
  dlSecs   time.Duration
}

// Get an individual file from S3.  If successful, writes the file out to disk
// and sends a stat object back to the main channel.  If there are errors they
// are sent back through the error channel.
func getFile(j *job.Job, docname string, c chan<- stat, e chan<- error) {
  start := time.Now()
  data, err := j.Get(docname)
  if err != nil {
    e <- fmt.Errorf("%v: %v", docname, err)
    return
  }
  path := basedir+docname
  err = ioutil.WriteFile(path, data, 0644)
  if err != nil {
    e <- err
    return
  }
  c <- stat{ filename: docname,
             size: len(data),
             dlSecs: time.Since(start) }
}

// Fan out workers to download each document in parallel, then block
// until all downloads are complete.
func getAllFiles(j *job.Job) {
  start := time.Now()
  j.Connect()
  c := make(chan stat, j.DocCount())
  e := make(chan error, j.DocCount())
  for _,doc := range j.DocList{
    throttle()
    go getFile(j,doc,c,e)
  }

  totalBytes := waitForDownloads(j,c,e)
  printSummary(start, totalBytes, j.CompleteCount())
}

// Prevents the system from being overwhelmed with work.
// Blocks until the number of Goroutines is less than a preset threshold.
func throttle() {
  for {
    if runtime.NumGoroutine() < MaxGoroutines { break }
    time.Sleep( 50 * time.Millisecond )
  }
}

// Listen on several channels for information from background download
// tasks -- each task will either send a stat through c, an error through
// e, or timeout.  Once all docs are accounted for, return the total number
// of bytes recieved.
func waitForDownloads(j *job.Job, c <-chan stat, e <-chan error) (totalBytes int) {
  for _,_ = range j.DocList{
    select {
      case packet := <-c:
        if verbose { log.Printf("%s was %d bytes\n", packet.filename,packet.size) }
        totalBytes += packet.size
        j.MarkComplete(packet.filename)
      case err := <-e:
        j.AddError(err)
      case <-time.After(downloadTimeout):
        j.AddError(errors.New("Timed out while downloading"))
        return
    }
  }
  return
}

// Print a summary of the download activity.
func printSummary(start time.Time, bytes int, count int){
  elapsed := time.Since(start)
  seconds := elapsed.Seconds()
  mbps := float64(bytes) / 1024 / 1024 / seconds
  log.Printf("got %d bytes over %d files in %f secs (%f MB/s)\n",
             bytes, count, seconds, mbps)
}

// Send an update on the success or failure of the operation to the
// callback URL provided by the job originator.
func postToCallback(j *job.Job){
  log.Println("work complete, posting status to callback:",j.Callback)
  _ = j.ToJson()
  poster.SendNotification(j)
}

// The entry point to this package.  Given a Job, download all the files,
// combine them into a single one, upload it to AWS and post the status to
// a callback endpoint.
func Combine(j *job.Job) bool {
  defer postToCallback(j)
  getAllFiles(j)
  if j.HasDownloadedDocs() {
    cpdf.Merge(j.Downloaded)
  }
  return true
}

