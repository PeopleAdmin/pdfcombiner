package server

import (
  "pdfcombiner/combiner"
  "pdfcombiner/job"
  "net/http"
  "fmt"
  "encoding/json"
)

type CombinerServer struct {}

var InvalidParamsMessage = "{\"response\":\"invalid params\"}"

func decodeParams(r *http.Request) (j *job.Job, err error){
  j = &job.Job{}
  decoder := json.NewDecoder(r.Body)
  err = decoder.Decode(j)
  fmt.Println(j)
  return
}

func (c CombinerServer) ProcessJob(w http.ResponseWriter, r *http.Request) {
  job, err := decodeParams(r)
  if err != nil || !job.IsValid() {
    http.Error(w, InvalidParamsMessage, http.StatusBadRequest)
    return
  }
  go combiner.Combine(job)
}

func (c CombinerServer) Ping(w http.ResponseWriter, r *http.Request) {}

