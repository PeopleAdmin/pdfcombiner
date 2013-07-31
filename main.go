package main

import (
  "fmt"
  "time"
  "io/ioutil"
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
)

var (
  bucketName = "pa-hrsuite-production"
  basedir = "/tmp/"
  keybase = "606/docs/"
  doclist = []string{"100001.pdf","100009.pdf","100030.pdf","100035.pdf"}
)

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


func main() {

  start := time.Now()
  bucket := connect()
  c := make(chan int)
  for _,doc := range doclist{
    go getFile(bucket,doc,c)
  }

  totalBytes := 0
  for _,doc := range doclist{
    recieved := <-c
    fmt.Printf("%s was %d bytes\n", doc,recieved)
    totalBytes += recieved
  }
  printSummary(start,totalBytes)

}
