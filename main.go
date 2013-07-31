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
  auth, err := aws.EnvAuth()
  if err != nil {
    panic(err.Error())
  }
  s := s3.New(auth, aws.USEast)
  bucket := s.Bucket(bucketName)

  data, err := bucket.Get(name)
  fmt.Println(len(data))
  if err != nil {
    panic(err.Error())
  }
  c <- data

}


func main() {

  flag.Parse()
  c := make(chan []byte)
  go getFile("000.txt",c)
  go getFile("000.txt",c)
  fmt.Printf("%s\n",<-c)
  fmt.Printf("%s\n",<-c)


}
