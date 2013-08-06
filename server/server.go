package server

import (
  "pdfcombiner/combiner"
  "net/http"
  "fmt"
  "encoding/json"
)

type CombinerServer struct {}

var InvalidParamsMessage = "{\"response\":\"invalid params\"}"

func decodeParams(r *http.Request) (job *combiner.Job, err error){
  job = &combiner.Job{}
  decoder := json.NewDecoder(r.Body)
  err = decoder.Decode(job)
  fmt.Println(job)
  return
}

func (c CombinerServer) ProcessJob(w http.ResponseWriter, r *http.Request) {
  job, err := decodeParams(r)
  if err != nil || !job.IsValid() {
    http.Error(w, InvalidParamsMessage, http.StatusBadRequest)
    return
  }
  go job.Combine()
}

func (c CombinerServer) Ping(w http.ResponseWriter, r *http.Request) {}

