package main

import (
  "flag"
  "fmt"
  "io/ioutil"
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
  if err != nil { panic(err) }
  fmt.Println(len(data))

  path := "/tmp/"+name
  err = ioutil.WriteFile(path, data, 0644)
  if err != nil { panic(err) }
  c <- data

}


func main() {

  flag.Parse()
  c := make(chan []byte)
  go getFile("000.txt",c)
  go getFile("000.txt",c)
  fmt.Printf("got %d bytes\n",len(<-c))
  fmt.Printf("got %d bytes\n",len(<-c))

}
