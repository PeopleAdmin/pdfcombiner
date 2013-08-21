// Package server contains handlers which respond to requests to combine
// pdf files and hands them off to the combiner package.
// TODO more informative validation errors
package server

import (
	"github.com/brasic/pdfcombiner/combiner"
	"github.com/brasic/pdfcombiner/job"
	"net/http"
	"strconv"
	"sync"
)

type CombinerServer struct {
	Port int
	wg   sync.WaitGroup
}

var invalidMessage = `{"response":"invalid params"}`
var okMessage = []byte("{\"response\":\"ok\"}\n")

// Start a HTTP server listening on `Port` to respond
// to JSON-formatted combination requests.
func (c CombinerServer) Listen() {
	http.HandleFunc("/health_check", c.Ping)
	http.HandleFunc("/", c.ProcessJob)
	println("Accepting connections on " + c.portString())
	http.ListenAndServe(c.portString(), nil)
}

func (c CombinerServer) Shutdown() {
	println("Waiting for all connections to finish...")
	c.wg.Wait()
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
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		combiner.Combine(job)
	}()
}

// No-op handler for responding to things like health checks.
func (c CombinerServer) Ping(w http.ResponseWriter, r *http.Request) {}

// http.ListenAndServe needs a string for the port.
func (c CombinerServer) portString() string {
	return ":" + strconv.Itoa(c.Port)
}
