// Package server contains handlers which respond to requests to combine
// pdf files and hands them off to the combiner package.
package server

import (
  "pdfcombiner/combiner"
  "pdfcombiner/job"
  "net/http"
  "encoding/json"
)

type CombinerServer struct {}

var InvalidParamsMessage = "{\"response\":\"invalid params\"}"

// Turn the JSON body into a *Job struct, setting err if necessary.
func decodeParams(r *http.Request) (j *job.Job, err error){
  j = &job.Job{}
  err = json.NewDecoder(r.Body).Decode(j)
  return
}

// Handler to recieve a POSTed JSON body encoding a Job and if it validates,
// send it along to be fulfilled.
func (c CombinerServer) ProcessJob(w http.ResponseWriter, r *http.Request) {
  job, err := decodeParams(r)
  if err != nil || !job.IsValid() {
    http.Error(w, InvalidParamsMessage, http.StatusBadRequest)
    return
  }
  go combiner.Combine(job)
}

// No-op handler for responding to things like health checks.
func (c CombinerServer) Ping(w http.ResponseWriter, r *http.Request) {}

