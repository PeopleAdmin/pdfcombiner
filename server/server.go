// Package server contains handlers which respond to requests to combine
// pdf files and hands them off to the combiner package.
// TODO more informative validation errors
package server

import (
	"net/http"
	"pdfcombiner/combiner"
	"pdfcombiner/job"
)

type CombinerServer struct{}

var invalidMessage = "{\"response\":\"invalid params\"}"
var okMessage = []byte("{\"response\":\"ok\"}\n")

// Handler to recieve a POSTed JSON body encoding a Job and if it validates,
// send it along to be fulfilled.
func (c CombinerServer) ProcessJob(w http.ResponseWriter, r *http.Request) {
	job, err := job.NewFromJson(r.Body)
	if err != nil || !job.IsValid() {
		http.Error(w, invalidMessage, http.StatusBadRequest)
		return
	}
	w.Write(okMessage)
	go combiner.Combine(job)
}

// No-op handler for responding to things like health checks.
func (c CombinerServer) Ping(w http.ResponseWriter, r *http.Request) {}
