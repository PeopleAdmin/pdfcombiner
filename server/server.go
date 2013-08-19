// Package server contains handlers which respond to requests to combine
// pdf files and hands them off to the combiner package.
// TODO more informative validation errors
package server

import (
	"net/http"
	"github.com/brasic/pdfcombiner/combiner"
	"github.com/brasic/pdfcombiner/job"
)

type CombinerServer struct{}

var invalidMessage = "{\"response\":\"invalid params\"}"
var okMessage = []byte("{\"response\":\"ok\"}\n")

// Start a HTTP server listening on `port` to respond
// to JSON-formatted combination requests.
func ListenOn(port string) {
	server := new(CombinerServer)
	http.HandleFunc("/health_check", server.Ping)
	http.HandleFunc("/", server.ProcessJob)
	println("Accepting connections on " + port)
	http.ListenAndServe(":"+port, nil)
}

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
