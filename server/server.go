// Package server contains handlers which respond to requests to combine
// pdf files and hands them off to the combiner package.
// TODO more informative validation errors
package server

import (
	"fmt"
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

type CombinerServer struct {
	port int
	wg   *sync.WaitGroup
}

var invalidMessage = `{"response":"invalid params"}`
var okMessage = []byte("{\"response\":\"ok\"}\n")

// Start a HTTP server listening on `Port` to respond
// to JSON-formatted combination requests.
func (c CombinerServer) Listen(listen_port int) {
	c.port = listen_port
	c.wg = new(sync.WaitGroup)
	listener, err := net.Listen("tcp", c.portString())
	if err != nil {
		panic(err)
	}
	println("Accepting connections on " + c.portString())
	c.registerHandlers(listener)
	http.Serve(listener, http.DefaultServeMux)
	println("Waiting for all jobs to finish...\n")
	c.wg.Wait()
}

// Handler to recieve a POSTed JSON body encoding a Job and if it validates,
// send it along to be fulfilled, keeping track of the in-flight job count.
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
	return ":" + strconv.Itoa(c.port)
}

func (c CombinerServer) registerHandlers(listener net.Listener) {
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
		for sig := range sigs {
			fmt.Printf("\ncaptured %v -- closing socket\n", sig)
			listener.Close()
		}
	}()
}
