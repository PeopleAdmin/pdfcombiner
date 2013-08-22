// Package server contains handlers which respond to requests to combine
// pdf files and hands them off to the combiner package.
// TODO more informative validation errors
package server

import (
	"github.com/brasic/pdfcombiner/combiner"
	"github.com/brasic/pdfcombiner/job"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

type combinerServer struct {
	port int
	wg   *sync.WaitGroup
}

var invalidMessage = `{"response":"invalid params"}`
var okMessage = []byte("{\"response\":\"ok\"}\n")

// Start a HTTP server listening on `Port` to respond
// to JSON-formatted combination requests.
func (c combinerServer) Listen(listen_port int) {
	c.port = listen_port
	c.wg = new(sync.WaitGroup)
	listener, err := net.Listen("tcp", c.portString())
	if err != nil {
		panic(err)
	}
	println("Accepting connections on " + c.portString())
	c.registerHandlers(listener)
	http.Serve(listener, http.DefaultServeMux)
	println("Waiting for all jobs to finish...")
	c.wg.Wait()
}

// Handler to recieve a POSTed JSON body encoding a Job and if it validates,
// send it along to be fulfilled, keeping track of the in-flight job count.
func (c combinerServer) ProcessJob(w http.ResponseWriter, r *http.Request) {
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
func (c combinerServer) Ping(w http.ResponseWriter, r *http.Request) {}

// http.ListenAndServe needs a string for the port.
func (c combinerServer) portString() string {
	return ":" + strconv.Itoa(c.port)
}

func (c combinerServer) registerHandlers(listener net.Listener) {
	http.HandleFunc("/health_check", c.Ping)
	http.HandleFunc("/", c.ProcessJob)
	handleSignals(listener)
}

// On OS shutdown or Ctrl-C, immediately close the tcp listener,
// but let background jobs finish before actually exiting.
func handleSignals(listener net.Listener) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, syscall.SIGTERM)
	go func() {
		for _ = range sigs {
			listener.Close()
		}
	}()
}
