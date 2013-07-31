// Benchmark parallel downloading and combination of some pdf files.
package main

import (
  "flag"
  "fmt"
  "time"
  "net/http"
  "io/ioutil"
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
)

var (
  verbose bool
  bucketName = "pa-hrsuite-production"
  basedir = "/tmp/"
  keybase = "606/docs/"
  doclist = []string{"100001.pdf","100009.pdf","100030.pdf","100035.pdf",
                     "100035.pdf","100037.pdf","100038.pdf","100093.pdf",}
)

func init(){
  flag.BoolVar(&verbose, "v", false, "Display extra data")
}

func getFile(bucket *s3.Bucket, docname string, c chan int) {
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

func printSummary(start time.Time, bytes int){
  elapsed := time.Since(start)
  seconds := elapsed.Seconds()
  mbps := float64(bytes) / 1024 / 1024 / seconds
  fmt.Printf("got %d bytes over %d files in %f secs (%f MB/s)\n",
             bytes, len(doclist), seconds, mbps)
}

func combine() {

  flag.Parse()
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
      fmt.Printf("%s was %d bytes\n", doc,recieved)
    }
    totalBytes += recieved
  }
  printSummary(start,totalBytes)

}

func handle_req(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("got conn\n")
  str := fmt.Sprintf("<html><body><head></head><h1>Time is %s</h1></body></html>", time.Now())
  fmt.Fprintf(w, str)
  fmt.Println("wrote",str)
}
// Exists only to prevent calls to favicon when testing the browser
func noopConn(w http.ResponseWriter, r *http.Request) {}

func main() {
  http.HandleFunc("/favicon.ico", noopConn)
  http.HandleFunc("/", handle_req)
  http.ListenAndServe(":8080", nil)
}
