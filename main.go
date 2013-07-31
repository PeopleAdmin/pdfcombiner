package main

import (
  "flag"
  "fmt"
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
)

var (
  bucketName string
  fileName   string
)

func init() {
  flag.StringVar(&bucketName, "b", "", "Bucket Name")
}

func getFile(name string, c chan []byte) {

}



func main() {

  flag.Parse()

  // The AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables are used.
  auth, err := aws.EnvAuth()
  if err != nil {
    panic(err.Error())
  }

  // Open Bucket
  s := s3.New(auth, aws.USEast)
  bucket := s.Bucket(bucketName)

  //#var data  []byte
  data, err := bucket.Get("000.txt")
  fmt.Println(len(data))
  fmt.Printf("%s\n",data)
  if err != nil {
    panic(err.Error())
  }
}
